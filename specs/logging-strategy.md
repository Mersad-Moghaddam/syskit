# SysKit Logging Strategy

> How SysKit emits diagnostics — structured, leveled, on stderr, and strictly separated from data output.

---

## Philosophy

SysKit is a data tool first. Its primary output is the system information a user asked for, and that output must never be polluted by diagnostic noise. Two rules follow directly:

1. **Data goes to stdout; diagnostics go to stderr.** They are never mixed. A user piping `syskit cpu --format json` into `jq` must receive clean JSON on stdout, with any warnings routed to stderr where they cannot corrupt the machine-readable stream. This is the same output contract described in [cli-conventions.md](cli-conventions.md).
2. **Silent success is the default.** With no verbosity flags, a successful command logs nothing. The user sees their data and nothing else. Logging is opt-in diagnostics, not ambient chatter.

Logging is a **CLI-layer concern**. It exists to help the user and the developer understand what the program is doing — it is not part of how data is collected.

---

## Logging Belongs to the CLI Layer

The most important rule in this document: **library and collector code does not log.**

Collectors, services, and the platform abstraction return errors — they never write to a log. This is a direct consequence of the constitution's **Keep It Modular** and **Clean Go** principles, and it mirrors the error philosophy in [error-handling.md](error-handling.md): the lower layers *report* what happened by returning values; the CLI layer *decides* what to do with that report, including whether and how to log it.

```go
// WRONG — a collector must never log.
func (c *Collector) Collect() (*CPUInfo, error) {
    data, err := c.fs.ReadFile("proc/stat")
    if err != nil {
        slog.Error("failed to read /proc/stat", "err", err) // ✗ library logging
        return nil, err
    }
    ...
}

// RIGHT — the collector returns; the CLI layer logs at its discretion.
func (c *Collector) Collect() (*CPUInfo, error) {
    data, err := c.fs.ReadFile("proc/stat")
    if err != nil {
        return nil, fmt.Errorf("reading /proc/stat: %w", err) // ✓ return, don't log
    }
    ...
}
```

Keeping logging out of the lower layers means collectors remain pure and testable (they have no logging side effects to assert or suppress), and the CLI retains full control over verbosity and destination.

---

## Levels

SysKit uses four levels, matching `log/slog`'s standard set:

| Level | When to use |
|---|---|
| `error` | A failure that affects the result — a collector failed, a file could not be read |
| `warn` | A recoverable anomaly — a missing optional field, a fallback path taken |
| `info` | High-level progress — which collectors ran, how long collection took |
| `debug` | Fine-grained detail — file paths read, raw values parsed, timing per step |

The default logging level is effectively *off*: nothing is emitted on a successful run. Verbosity flags raise the threshold of what is shown.

---

## Structured Logging with `log/slog`

SysKit uses the standard library's `log/slog` — no third-party logging dependency, consistent with the **Minimal Dependencies** principle. Structured, key-value logging makes diagnostics machine-parseable and greppable.

The logger is constructed once at the CLI layer, writes to **stderr**, and its level is set from the verbosity flags:

```go
// Configured once during CLI startup, based on parsed flags.
func newLogger(verbosity Verbosity) *slog.Logger {
    var level slog.Level
    switch verbosity {
    case Quiet:
        level = slog.LevelError + 4 // effectively silence everything
    case Debug:
        level = slog.LevelDebug
    case Verbose:
        level = slog.LevelInfo
    default:
        level = slog.LevelError + 4 // default: silent success
    }
    handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
    return slog.New(handler)
}
```

Diagnostics are emitted with structured attributes so they carry context without prose:

```go
log.Debug("collecting", "domain", "cpu", "source", "/proc/stat")
log.Info("collection complete", "domain", "cpu", "elapsed", elapsed)
log.Warn("optional field absent", "field", "MemAvailable", "fallback", "computed")
log.Error("collector failed", "domain", "network", "err", err)
```

Note that every logger writes to `os.Stderr`. There is no code path in which the logger writes to stdout.

---

## Verbosity Flags

Verbosity is controlled entirely by global flags, defined once and applied consistently across every command (per the **Consistent CLI Experience** principle). These flags are part of the CLI contract in [cli-conventions.md](cli-conventions.md).

| Flag | Short | Effect |
|---|---|---|
| *(none)* | | Silent success — no diagnostics emitted |
| `--verbose` | `-v` | Enable `info`-level diagnostics on stderr |
| `--debug` | | Enable `debug`-level diagnostics on stderr |
| `--quiet` | `-q` | Suppress all diagnostics, including errors, on stderr |

Precedence when multiple are supplied: `--quiet` wins over `--debug`, which wins over `--verbose`. `--quiet` is the strongest signal — the user has asked for silence, and SysKit honors it even for errors (the exit code still communicates failure to scripts; see [error-handling.md](error-handling.md)).

```sh
# Silent success: only the CPU data on stdout.
syskit cpu

# Info diagnostics on stderr, data still clean on stdout.
syskit cpu --verbose

# Full debug trace on stderr — every file read and value parsed.
syskit cpu --debug

# Absolute silence on stderr; exit code still reflects success/failure.
syskit cpu --quiet
```

---

## Stdout / Stderr Separation in Practice

Because the streams are strictly separated, SysKit composes cleanly in pipelines and scripts:

```sh
# stdout carries clean JSON to jq; stderr diagnostics stay on the terminal.
syskit network --format json --debug | jq '.interfaces[].name'

# Capture data and diagnostics independently.
syskit disk --format json --verbose > data.json 2> run.log
```

The formatter writes exclusively to stdout. The logger writes exclusively to stderr. User-facing error messages (from [error-handling.md](error-handling.md)) also go to stderr. There is no configuration in which these cross.

---

## Summary

| Concern | Rule |
|---|---|
| Destination | Diagnostics → **stderr**; data → **stdout**; never mixed |
| Default | Silent success — no logging without a flag |
| Where | CLI layer only; libraries return errors, never log |
| Library | `log/slog` from the standard library, text handler on stderr |
| Levels | `error`, `warn`, `info`, `debug` |
| Flags | `--verbose`/`-v` (info), `--debug` (debug), `--quiet`/`-q` (silence) |
| Precedence | `--quiet` > `--debug` > `--verbose` |

---

*A good diagnostic stream is one the user forgets exists until the moment they need it — and finds exactly where they expect it, never in the way of their data.*
