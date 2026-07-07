# SysKit Error Handling

> How SysKit represents, propagates, and reports errors — from kernel-interface failures to user-facing diagnostics.

---

## Philosophy

Errors are values, and SysKit treats them with the same care as any other data. The constitution's **Clean Go** principle is explicit: "explicit error handling — no panics in library code." Every function that can fail returns an `error`, and every caller decides what to do with it.

Three commitments guide all error handling in SysKit:

1. **No panics in library code.** Collectors, services, and the platform layer never `panic` for expected conditions. A missing `/proc` file, a permission denial, or a malformed line is an `error` to be returned, not a crash. Panics are reserved exclusively for genuinely unrecoverable programmer errors (e.g. a nil dependency that indicates a wiring bug), and they never escape to the user as a stack trace.
2. **Errors carry context as they rise.** A failure deep in the platform layer is wrapped at each boundary so that by the time it reaches the CLI, it describes the full path of what went wrong.
3. **User-facing and internal errors are distinct.** Internal errors are for diagnosis; user-facing errors are for action. The CLI layer translates the former into the latter.

---

## Sentinel Errors

Well-known failure conditions are represented as exported sentinel errors, allowing callers to test for them with `errors.Is`:

```go
package platform

import "errors"

var (
    // ErrNotFound indicates a kernel interface file or directory is absent.
    ErrNotFound = errors.New("kernel interface not found")
    // ErrPermission indicates insufficient privilege to read an interface.
    ErrPermission = errors.New("permission denied reading kernel interface")
    // ErrUnsupported indicates the kernel does not provide this interface.
    ErrUnsupported = errors.New("kernel interface not supported")
)
```

Sentinels give the CLI layer a stable, testable vocabulary for deciding exit codes and messages without string-matching. Each collector may define its own domain sentinels (`ErrParse`, `ErrFieldMissing`) for the same purpose, as referenced in [testing-strategy.md](testing-strategy.md).

---

## Wrapping and Inspection

Errors are wrapped with `fmt.Errorf` and the `%w` verb as they cross each architectural boundary. Wrapping adds context — *what* was being attempted — while preserving the underlying error for inspection.

```go
func (c *Collector) Collect() (*CPUInfo, error) {
    data, err := c.fs.ReadFile("proc/stat")
    if err != nil {
        return nil, fmt.Errorf("reading /proc/stat: %w", err)
    }
    info, err := parseStat(data)
    if err != nil {
        return nil, fmt.Errorf("parsing /proc/stat: %w", err)
    }
    return info, nil
}
```

Callers inspect wrapped errors with `errors.Is` (for sentinels) and `errors.As` (for typed errors), never by comparing strings:

```go
info, err := collector.Collect()
if errors.Is(err, platform.ErrPermission) {
    // map to a permission-specific exit code and message
}

var perr *fs.PathError
if errors.As(err, &perr) {
    // inspect the underlying path operation
}
```

This layered wrapping produces chains like:

```text
parsing /proc/stat: reading /proc/stat: permission denied reading kernel interface
```

— which reads as a precise trace from symptom to root cause, without a single stack dump.

---

## Internal vs. User-Facing Errors

SysKit draws a firm line between the two.

| | Internal errors | User-facing errors |
|---|---|---|
| Audience | Developers, logs, diagnostics | End users at the terminal |
| Style | Lowercase fragments, wrapped | Full sentences, actionable |
| Destination | Returned up the stack | Printed to **stderr** |
| Created by | Any layer, via wrapping | The CLI layer, at the boundary |

Library and collector code produces **internal** errors: lowercase, no trailing punctuation, designed to be wrapped. Only at the CLI boundary are these translated into **user-facing** messages: complete sentences that tell the user what happened and what they can do about it.

```go
// CLI layer: translate an internal error into a user-facing diagnostic.
func present(err error) (message string, code int) {
    switch {
    case errors.Is(err, platform.ErrPermission):
        return "Permission denied. Try running with elevated privileges (sudo).", ExitPermission
    case errors.Is(err, platform.ErrUnsupported):
        return "This information is not available on your kernel.", ExitUnsupported
    default:
        return fmt.Sprintf("Error: %v", err), ExitGeneral
    }
}
```

The user sees a clear instruction; the underlying chain remains available for `--debug` diagnostics (see [logging-strategy.md](logging-strategy.md)).

---

## Message Format & Style

Consistent error wording is part of the **Consistent CLI Experience** principle.

**Internal / wrapped errors** (returned by library code):

- Lowercase first letter.
- No trailing punctuation.
- A short phrase naming the operation, followed by `: %w`.
- Example: `fmt.Errorf("reading /proc/meminfo: %w", err)`

**User-facing errors** (printed to stderr by the CLI):

- Full sentences with normal capitalization and terminal punctuation.
- State the problem and, where possible, the remedy.
- Example: `Permission denied. Try running with elevated privileges (sudo).`

This mirrors the Go standard library convention: wrapped errors compose into readable chains precisely because each fragment is a lowercase, punctuation-free clause.

---

## Exit Codes

SysKit uses a small, stable set of exit codes so that scripts and pipelines can react programmatically. This table is canonical; [cli-conventions.md](cli-conventions.md) mirrors it and the two must always match.

| Code | Name | Meaning |
|---|---|---|
| 0 | Success | Command completed successfully |
| 1 | General error | An unspecified runtime error occurred |
| 2 | Usage error | Invalid flags, arguments, or command usage |
| 3 | Permission | Insufficient privilege to read a kernel interface |
| 4 | Unsupported | Required kernel interface is missing or unsupported |
| 5 | Partial failure | Some data collected; one or more collectors failed |

```go
const (
    ExitSuccess     = 0
    ExitGeneral     = 1
    ExitUsage       = 2
    ExitPermission  = 3
    ExitUnsupported = 4
    ExitPartial     = 5
)
```

Exit codes are assigned exclusively at the CLI layer, derived from the sentinel errors that propagate up. Cobra emits exit code `2` for usage errors automatically; SysKit maps its own sentinels onto codes `3`–`5`.

---

## Partial-Failure Handling

Many SysKit commands aggregate data from several collectors (per the Service Layer's aggregation role). The failure of one collector must not blank out the entire result — the **Keep It Modular** principle means collectors are independent, and a fault in one is isolated from the others.

When a command aggregates multiple sources, it collects both results and errors, reports what succeeded, and surfaces what failed:

```go
func (s *SystemService) Collect() (*SystemInfo, error) {
    info := &SystemInfo{}
    var errs []error

    if host, err := s.host.Collect(); err != nil {
        errs = append(errs, fmt.Errorf("host info: %w", err))
    } else {
        info.Host = host
    }

    if load, err := s.load.Collect(); err != nil {
        errs = append(errs, fmt.Errorf("load average: %w", err))
    } else {
        info.Load = load
    }

    // Return the data we did gather, alongside a joined error describing gaps.
    return info, errors.Join(errs...)
}
```

The CLI layer prints the successfully collected data to **stdout**, prints the joined diagnostics to **stderr**, and exits with `ExitPartial` (5). The user gets everything the system could provide, plus an honest accounting of what it could not — never a silent omission and never a total failure over a single missing field.

`errors.Join` (Go 1.20+) preserves every underlying error, so each remains inspectable with `errors.Is`.

---

## Summary

| Concern | Rule |
|---|---|
| Panics | Never in library code; reserved for unrecoverable programmer bugs |
| Representation | Sentinel errors for known conditions; wrap with `%w` |
| Inspection | `errors.Is` / `errors.As`, never string comparison |
| Internal messages | lowercase, no trailing punctuation, operation-prefixed |
| User messages | full sentences, actionable, printed to stderr |
| Exit codes | Stable table (0–5), assigned at the CLI layer |
| Partial failures | Report what succeeded; join and surface what failed |

---

*Good error handling is invisible when things work and invaluable when they do not. SysKit's errors tell the user what to do and tell the developer what went wrong — never confusing the two.*
