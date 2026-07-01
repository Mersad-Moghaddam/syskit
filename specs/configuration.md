# SysKit Configuration

> How SysKit is configured — precedence, file format, locations, and environment variables — while remaining fully usable with zero configuration.

---

## Philosophy

SysKit works out of the box. Configuration is entirely optional: a fresh install with no config file and no environment variables behaves sensibly, using built-in defaults. Configuration exists to let users encode their preferences once — a default output format, a preferred refresh interval, color settings — not to make the tool usable in the first place.

This "zero-config by default" stance keeps SysKit honest to the **Consistent CLI Experience** principle: every command has predictable defaults, and configuration only shifts those defaults, never introduces surprises.

---

## Precedence

When the same setting is specified in more than one place, SysKit resolves it by a fixed precedence, from highest to lowest:

```text
Command-line flags   (highest — most specific, per-invocation)
        ▲
Environment variables
        ▲
Configuration file
        ▲
Built-in defaults    (lowest — always present)
```

| Rank | Source | Scope | Example |
|---|---|---|---|
| 1 | Flags | Single invocation | `--format json` |
| 2 | Environment | Shell session | `SYSKIT_FORMAT=json` |
| 3 | Config file | User-persistent | `format = "json"` |
| 4 | Defaults | Built-in | `format = "table"` |

The rule is intuitive: the more specific and immediate the source, the higher it ranks. A flag on a single command overrides everything, because the user is expressing intent for *this run*. Defaults sit at the bottom and guarantee that every setting always has a value.

---

## Configuration File Format

SysKit uses **TOML** for its configuration file. TOML is unambiguous, comment-friendly, and maps cleanly onto typed configuration structs — a good fit for the **Clean Go** and **Minimal Dependencies** principles. It is easier to read and less error-prone by hand than JSON, and less whitespace-sensitive than YAML.

The file is optional. If it is absent, SysKit proceeds with defaults (adjusted by any environment variables and flags).

---

## File Locations

SysKit follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/) for locating its configuration file. It looks in this order and uses the first file found:

| Order | Path | Condition |
|---|---|---|
| 1 | `$XDG_CONFIG_HOME/syskit/config.toml` | If `XDG_CONFIG_HOME` is set |
| 2 | `~/.config/syskit/config.toml` | Fallback when `XDG_CONFIG_HOME` is unset |

An explicit path always wins over discovery:

```sh
syskit cpu --config /etc/syskit/config.toml
```

Following XDG keeps SysKit a well-behaved Linux citizen — consistent with the **Linux First** principle — placing user configuration exactly where Linux users expect it and keeping the home directory uncluttered.

---

## Environment Variables

Every configurable setting has a corresponding environment variable, prefixed with `SYSKIT_`. Environment variables override the config file but are overridden by flags.

| Variable | Equivalent setting | Example |
|---|---|---|
| `SYSKIT_FORMAT` | Default output format | `SYSKIT_FORMAT=json` |
| `SYSKIT_COLOR` | Color output (`auto`/`always`/`never`) | `SYSKIT_COLOR=never` |
| `SYSKIT_REFRESH` | Refresh interval for live commands | `SYSKIT_REFRESH=2s` |
| `SYSKIT_CONFIG` | Explicit config file path | `SYSKIT_CONFIG=/etc/syskit/config.toml` |

The naming is mechanical: a setting named `refresh_interval` maps to `SYSKIT_REFRESH_INTERVAL`. This predictability is part of the consistent CLI experience.

SysKit also honors the conventional `NO_COLOR` environment variable, which disables color regardless of other settings (see [cli-conventions.md](cli-conventions.md)).

---

## What Is Configurable

Configuration adjusts presentation and defaults — never behavior that would compromise SysKit's read-only, inspection-only nature.

| Setting | Values | Default | Purpose |
|---|---|---|---|
| `format` | `table` \| `json` \| `yaml` | `table` | Default output format |
| `color` | `auto` \| `always` \| `never` | `auto` | Colorized output control |
| `refresh_interval` | duration (e.g. `1s`, `500ms`) | `1s` | Refresh cadence for `watch`, `top`, `dashboard` |
| `no_header` | boolean | `false` | Suppress table headers |
| `verbosity` | `quiet` \| `normal` \| `verbose` \| `debug` | `normal` | Default diagnostic verbosity |

`color = "auto"` means SysKit enables color only when stdout is an interactive terminal and disables it when output is piped — TTY detection is described in [cli-conventions.md](cli-conventions.md).

---

## Per-Command Sections

Global settings sit at the top level of the file. Command-specific overrides live in named tables, allowing a user to, for example, default to table output everywhere but JSON specifically for `ports`. A command-level setting overrides the global setting for that command only.

```toml
# ~/.config/syskit/config.toml

# ── Global defaults ──────────────────────────────────
format           = "table"
color            = "auto"
refresh_interval = "1s"
verbosity        = "normal"

# ── Per-command overrides ────────────────────────────
[process]
# Default the process list to JSON for scripting.
format = "json"

[top]
# Refresh the interactive monitor faster.
refresh_interval = "500ms"

[dashboard]
refresh_interval = "2s"

[ports]
# Suppress headers for cleaner piping.
no_header = true
```

The effective value for a setting is resolved as: flag → environment → per-command section → global section → built-in default.

---

## Loading Model

Configuration maps onto a typed struct, loaded once at CLI startup and threaded down as plain values. The lower layers never read configuration files — configuration, like logging, is a CLI-layer concern that produces concrete parameters passed into services and formatters.

```go
type Config struct {
    Format          string            `toml:"format"`
    Color           string            `toml:"color"`
    RefreshInterval Duration          `toml:"refresh_interval"`
    Verbosity       string            `toml:"verbosity"`
    Commands        map[string]Config `toml:"-"` // per-command sections
}

// Load resolves the effective configuration: defaults, overlaid by the
// discovered TOML file, overlaid by SYSKIT_* environment variables.
// Flags are applied afterward by the command layer. A missing file is
// not an error — Load returns defaults.
func Load(path string) (*Config, error) {
    cfg := Defaults()
    if path == "" {
        path = discover() // XDG_CONFIG_HOME then ~/.config
    }
    if data, err := os.ReadFile(path); err == nil {
        if err := toml.Unmarshal(data, cfg); err != nil {
            return nil, fmt.Errorf("parsing config %s: %w", path, err)
        }
    } else if !errors.Is(err, fs.ErrNotExist) {
        return nil, fmt.Errorf("reading config %s: %w", path, err)
    }
    cfg.applyEnv() // SYSKIT_* overrides
    return cfg, nil
}
```

A missing config file is expected and silent. A *malformed* config file is a real error, surfaced to the user via the conventions in [error-handling.md](error-handling.md), because the user clearly intended to configure something and got it wrong.

---

## Summary

| Concern | Rule |
|---|---|
| Necessity | Optional — SysKit runs fully with zero config |
| Precedence | flags > env > file > defaults |
| Format | TOML |
| Location | `$XDG_CONFIG_HOME/syskit/config.toml`, fallback `~/.config/syskit/config.toml` |
| Environment | `SYSKIT_*` prefix; `NO_COLOR` honored |
| Structure | Global settings plus per-command `[section]` overrides |
| Loading | CLI-layer only; missing file is silent, malformed file is an error |

---

*Configuration should reward the user who wants to tune SysKit without ever punishing the user who does not. Sensible defaults first; preferences second.*
