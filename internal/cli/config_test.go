package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadMissingFileReturnsDefaults confirms a missing config file is not an
// error: Load returns the built-in defaults (specs/configuration.md).
func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir()) // empty dir → no config.toml
	clearSyskitEnv(t)

	cfg, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, Defaults(), cfg)
}

// TestLoadMalformedFileErrors confirms a malformed TOML file is surfaced as an
// error, because the user intended to configure something and got it wrong.
func TestLoadMalformedFileErrors(t *testing.T) {
	clearSyskitEnv(t)
	path := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(path, []byte("format = = ="), 0o644))

	_, err := Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing config")
}

// TestLoadFileValues confirms global settings decode from the file.
func TestLoadFileValues(t *testing.T) {
	clearSyskitEnv(t)
	path := filepath.Join(t.TempDir(), "config.toml")
	content := `
format = "json"
color = "never"
refresh_interval = "2s"
no_header = true
verbosity = "verbose"
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "json", cfg.Format)
	assert.Equal(t, "never", cfg.Color)
	assert.Equal(t, 2*time.Second, cfg.RefreshInterval)
	assert.True(t, cfg.NoHeader)
	assert.Equal(t, "verbose", cfg.Verbosity)
}

// TestLoadEnvOverridesFile confirms SYSKIT_* env overrides a file value.
func TestLoadEnvOverridesFile(t *testing.T) {
	clearSyskitEnv(t)
	path := filepath.Join(t.TempDir(), "config.toml")
	require.NoError(t, os.WriteFile(path, []byte(`format = "table"`), 0o644))
	t.Setenv("SYSKIT_FORMAT", "json")

	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "json", cfg.Format)
}

// TestResolveFormatPrecedence encodes the worked example from
// specs/configuration.md: for `syskit process`, with env SYSKIT_FORMAT=json and
// a per-command [process] format="yaml" and global format="table", the global
// env var wins over the per-command file section → "json". Removing the env
// falls to the per-command section → "yaml"; removing that → global "table".
func TestResolveFormatPrecedence(t *testing.T) {
	cfg := Defaults()
	cfg.Format = "table" // global (file/env-resolved) value
	cfg.Commands["process"] = commandConfig{Format: strptr("yaml")}

	// env set: global env beats per-command section.
	got := cfg.resolveFormat(false, "table", true, "process")
	assert.Equal(t, "table", got, "env-resolved global outranks per-command section")

	// Simulate env not set: per-command section beats global.
	got = cfg.resolveFormat(false, "table", false, "process")
	assert.Equal(t, "yaml", got, "per-command section beats global when no env")

	// No per-command section for this command: falls to global.
	got = cfg.resolveFormat(false, "table", false, "cpu")
	assert.Equal(t, "table", got)

	// Explicit flag always wins.
	got = cfg.resolveFormat(true, "json", false, "process")
	assert.Equal(t, "json", got, "explicit flag wins over everything")
}

// TestResolveFormatWorkedExampleEndToEnd runs the worked example through Load so
// the env-outranks-section behavior is proven against real decoding.
func TestResolveFormatWorkedExampleEndToEnd(t *testing.T) {
	clearSyskitEnv(t)
	path := filepath.Join(t.TempDir(), "config.toml")
	content := `
format = "table"

[process]
format = "yaml"
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	t.Setenv("SYSKIT_FORMAT", "json")

	cfg, err := Load(path)
	require.NoError(t, err)
	// applyEnv folded json into cfg.Format; resolveFormat must keep env winning.
	_, envSet := os.LookupEnv("SYSKIT_FORMAT")
	assert.Equal(t, "json", cfg.resolveFormat(false, "table", envSet, "process"))
}

func TestResolvePresentationAndLiveSettings(t *testing.T) {
	cfg := Defaults()
	cfg.Color = "never"
	cfg.NoHeader = false
	cfg.RefreshInterval = time.Second
	cfg.Verbosity = "normal"
	cfg.Commands["top"] = commandConfig{
		Color:           strptr("always"),
		NoHeader:        boolptr(true),
		RefreshInterval: durationptr(500 * time.Millisecond),
		Verbosity:       strptr("debug"),
	}

	assert.Equal(t, "always", cfg.resolveColor(false, "auto", false, "top"))
	assert.True(t, cfg.resolveNoHeader(false, false, false, "top"))
	assert.Equal(t, 500*time.Millisecond, cfg.resolveRefreshInterval(false, "top"))
	assert.Equal(t, "debug", cfg.resolveConfiguredVerbosity(false, "top"))

	assert.Equal(t, "never", cfg.resolveColor(false, "auto", true, "top"), "environment tier bypasses command section")
	assert.False(t, cfg.resolveNoHeader(false, false, true, "top"))
	assert.Equal(t, time.Second, cfg.resolveRefreshInterval(true, "top"))
	assert.Equal(t, "normal", cfg.resolveConfiguredVerbosity(true, "top"))

	assert.Equal(t, "auto", cfg.resolveColor(true, "auto", false, "top"))
	assert.False(t, cfg.resolveNoHeader(true, false, false, "top"))
}

// TestLoadPerCommandSection confirms an arbitrary [section] is captured.
func TestLoadPerCommandSection(t *testing.T) {
	clearSyskitEnv(t)
	path := filepath.Join(t.TempDir(), "config.toml")
	content := `
[top]
refresh_interval = "500ms"

[ports]
no_header = true
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	cfg, err := Load(path)
	require.NoError(t, err)
	require.Contains(t, cfg.Commands, "ports")
	require.NotNil(t, cfg.Commands["ports"].NoHeader)
	assert.True(t, *cfg.Commands["ports"].NoHeader)
}

func strptr(s string) *string { return &s }
func boolptr(v bool) *bool    { return &v }
func durationptr(v time.Duration) *time.Duration {
	return &v
}

// clearSyskitEnv removes SYSKIT_* variables so a test's environment is
// deterministic regardless of the developer's shell.
func clearSyskitEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"SYSKIT_FORMAT", "SYSKIT_COLOR", "SYSKIT_VERBOSITY",
		"SYSKIT_NO_HEADER", "SYSKIT_REFRESH_INTERVAL", "SYSKIT_CONFIG",
	} {
		t.Setenv(k, "")
		require.NoError(t, os.Unsetenv(k))
	}
}
