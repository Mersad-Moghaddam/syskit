package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds SysKit's resolved settings. It is loaded once at CLI startup and
// threaded down as plain values; lower layers never read configuration files
// (ARCHITECTURE.md §4, specs/configuration.md "Loading Model").
//
// SysKit is zero-config: with no file and no environment variables, Defaults
// applies and every command behaves sensibly. Configuration only shifts those
// defaults.
type Config struct {
	Format          string
	Color           string
	RefreshInterval time.Duration
	NoHeader        bool
	Verbosity       string

	// Commands holds per-command override sections (TOML `[process]`, `[top]`,
	// …). A per-command value overrides the global value for that command only,
	// but never outranks an environment variable or a flag (see resolveFormat).
	Commands map[string]commandConfig
}

// commandConfig is the subset of settings a per-command `[section]` may
// override. Only fields set in the section are non-nil, so resolution can tell
// "unset" from "set to the zero value".
type commandConfig struct {
	Format          *string
	Color           *string
	NoHeader        *bool
	Verbosity       *string
	RefreshInterval *time.Duration
}

// Defaults returns the built-in configuration used when nothing else sets a
// value (specs/configuration.md "What Is Configurable").
func Defaults() *Config {
	return &Config{
		Format:          formatTable,
		Color:           "auto",
		RefreshInterval: time.Second,
		NoHeader:        false,
		Verbosity:       "normal",
		Commands:        map[string]commandConfig{},
	}
}

// knownGlobalKeys are the top-level scalar settings; any other top-level table
// is treated as a per-command override section.
var knownGlobalKeys = map[string]bool{
	"format":           true,
	"color":            true,
	"refresh_interval": true,
	"no_header":        true,
	"verbosity":        true,
}

// Load resolves the effective configuration: defaults, overlaid by the
// discovered or specified TOML file, overlaid by SYSKIT_* environment
// variables. Command-line flags are applied afterward by the command layer.
//
// A missing file is not an error — Load returns defaults (adjusted by env). A
// malformed file is an error, because the user clearly intended to configure
// something and got it wrong (specs/configuration.md).
func Load(path string) (*Config, error) {
	cfg := Defaults()

	if path == "" {
		path = discoverConfigPath()
	}
	if path != "" {
		data, err := os.ReadFile(path)
		switch {
		case err == nil:
			if uerr := unmarshalConfig(data, cfg); uerr != nil {
				return nil, fmt.Errorf("parsing config %s: %w", path, uerr)
			}
		case errors.Is(err, fs.ErrNotExist):
			// Missing file is expected and silent.
		default:
			return nil, fmt.Errorf("reading config %s: %w", path, err)
		}
	}

	if err := cfg.applyEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// discoverConfigPath returns the XDG config path if a file exists there, else
// "". It checks $XDG_CONFIG_HOME/syskit/config.toml first, then
// ~/.config/syskit/config.toml (specs/configuration.md "File Locations").
func discoverConfigPath() string {
	if base := os.Getenv("XDG_CONFIG_HOME"); base != "" {
		return filepath.Join(base, "syskit", "config.toml")
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "syskit", "config.toml")
	}
	return ""
}

// unmarshalConfig decodes TOML into cfg. It uses a generic map decode so that
// arbitrary per-command `[section]` tables are captured without a named struct
// field per command.
func unmarshalConfig(data []byte, cfg *Config) error {
	var raw map[string]any
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return err
	}

	for key, val := range raw {
		if knownGlobalKeys[key] {
			if err := applyScalar(cfg, key, val); err != nil {
				return err
			}
			continue
		}
		// Any other top-level table is a per-command override section.
		section, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("unknown top-level key %q", key)
		}
		cc, err := commandConfigFromMap(section)
		if err != nil {
			return fmt.Errorf("in [%s]: %w", key, err)
		}
		cfg.Commands[key] = cc
	}
	return nil
}

// applyScalar sets a known global field on cfg from a decoded TOML value.
func applyScalar(cfg *Config, key string, val any) error {
	switch key {
	case "format":
		s, err := asString(val, key)
		if err != nil {
			return err
		}
		cfg.Format = s
	case "color":
		s, err := asString(val, key)
		if err != nil {
			return err
		}
		cfg.Color = s
	case "verbosity":
		s, err := asString(val, key)
		if err != nil {
			return err
		}
		cfg.Verbosity = s
	case "no_header":
		b, ok := val.(bool)
		if !ok {
			return fmt.Errorf("no_header must be a boolean, got %T", val)
		}
		cfg.NoHeader = b
	case "refresh_interval":
		s, err := asString(val, key)
		if err != nil {
			return err
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("refresh_interval %q: %w", s, err)
		}
		cfg.RefreshInterval = d
	}
	return nil
}

// commandConfigFromMap extracts the overridable subset from a per-command
// section table.
func commandConfigFromMap(section map[string]any) (commandConfig, error) {
	var cc commandConfig
	for key, val := range section {
		switch key {
		case "format":
			s, err := asString(val, key)
			if err != nil {
				return cc, err
			}
			cc.Format = &s
		case "color":
			s, err := asString(val, key)
			if err != nil {
				return cc, err
			}
			cc.Color = &s
		case "verbosity":
			s, err := asString(val, key)
			if err != nil {
				return cc, err
			}
			cc.Verbosity = &s
		case "no_header":
			b, ok := val.(bool)
			if !ok {
				return cc, fmt.Errorf("no_header must be a boolean, got %T", val)
			}
			cc.NoHeader = &b
		case "refresh_interval":
			s, err := asString(val, key)
			if err != nil {
				return cc, err
			}
			d, err := time.ParseDuration(s)
			if err != nil {
				return cc, fmt.Errorf("refresh_interval %q: %w", s, err)
			}
			cc.RefreshInterval = &d
		default:
			return cc, fmt.Errorf("unknown key %q", key)
		}
	}
	return cc, nil
}

func asString(val any, key string) (string, error) {
	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string, got %T", key, val)
	}
	return s, nil
}

// applyEnv overlays SYSKIT_* environment variables onto cfg. Environment
// variables outrank the config file (including per-command sections) and are
// themselves outranked by flags (specs/configuration.md "Precedence"). Every
// SYSKIT_* variable is global in scope.
func (c *Config) applyEnv() error {
	if v, ok := os.LookupEnv("SYSKIT_FORMAT"); ok {
		c.Format = v
	}
	if v, ok := os.LookupEnv("SYSKIT_COLOR"); ok {
		c.Color = v
	}
	if v, ok := os.LookupEnv("SYSKIT_VERBOSITY"); ok {
		c.Verbosity = v
	}
	if v, ok := os.LookupEnv("SYSKIT_NO_HEADER"); ok {
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return fmt.Errorf("SYSKIT_NO_HEADER %q: %w", v, err)
		}
		c.NoHeader = b
	}
	if v, ok := os.LookupEnv("SYSKIT_REFRESH_INTERVAL"); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("SYSKIT_REFRESH_INTERVAL %q: %w", v, err)
		}
		c.RefreshInterval = d
	}
	return nil
}

// resolveFormat returns the effective output format for the named command,
// applying the full precedence order (specs/configuration.md, worked example):
//
//	flag > env > per-command [section] > global (file/env) > built-in default
//
// The env tier is already folded into c.Format by applyEnv, and it sits above
// every file-based value — so a global SYSKIT_FORMAT outranks a per-command
// file section. To honor that, a per-command section is consulted only when the
// environment did not set the format.
//
//   - flagChanged: whether --format was given on this invocation.
//   - flagValue:   the --format value (only meaningful when flagChanged).
//   - envSet:      whether SYSKIT_FORMAT was set.
//   - command:     the subcommand name, for per-command section lookup.
func (c *Config) resolveFormat(flagChanged bool, flagValue string, envSet bool, command string) string {
	if flagChanged {
		return flagValue
	}
	if !envSet {
		if cc, ok := c.Commands[command]; ok && cc.Format != nil {
			return *cc.Format
		}
	}
	return c.Format
}

func (c *Config) resolveColor(flagChanged bool, flagValue string, envSet bool, command string) string {
	if flagChanged {
		return flagValue
	}
	if !envSet {
		if cc, ok := c.Commands[command]; ok && cc.Color != nil {
			return *cc.Color
		}
	}
	return c.Color
}

func (c *Config) resolveNoHeader(flagChanged, flagValue, envSet bool, command string) bool {
	if flagChanged {
		return flagValue
	}
	if !envSet {
		if cc, ok := c.Commands[command]; ok && cc.NoHeader != nil {
			return *cc.NoHeader
		}
	}
	return c.NoHeader
}

func (c *Config) resolveRefreshInterval(envSet bool, command string) time.Duration {
	if !envSet {
		if cc, ok := c.Commands[command]; ok && cc.RefreshInterval != nil {
			return *cc.RefreshInterval
		}
	}
	return c.RefreshInterval
}

func (c *Config) resolveConfiguredVerbosity(envSet bool, command string) string {
	if !envSet {
		if cc, ok := c.Commands[command]; ok && cc.Verbosity != nil {
			return *cc.Verbosity
		}
	}
	return c.Verbosity
}
