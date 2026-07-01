# Go Coding Conventions

> How Go code is written, formatted, and organized across the SysKit codebase.

---

## Purpose

This document defines the concrete rules for writing Go in SysKit. It exists to make the codebase uniform, reviewable, and idiomatic, so that any contributor can read any file and recognize the patterns.

These rules operationalize the **Clean Go** and **Minimal Dependencies** principles from `../specs/constitution.md`. They are enforced in review (`code-review.md`) and required by the `../standards/definition-of-done.md`.

---

## Formatting and Tooling

Formatting is not a matter of taste. The following tools run in CI and must pass with zero output.

| Tool | Command | Rule |
|---|---|---|
| gofmt | `gofmt -l .` | Must produce no output. All code is gofmt-formatted. |
| goimports | `goimports -l -local github.com/<org>/syskit .` | Imports are grouped: stdlib, third-party, local. |
| go vet | `go vet ./...` | Must pass with no findings. |
| govulncheck | `govulncheck ./...` | No known vulnerabilities (see `dependency-policy.md`). |

Rules:

- Do **not** hand-align struct fields or comments; let gofmt decide.
- Group imports into three blocks separated by blank lines. goimports enforces order:

```go
import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/<org>/syskit/internal/collector/cpu"
)
```

- Never commit code that gofmt would rewrite. Configure your editor to run gofmt and goimports on save.

---

## Naming

Naming communicates intent. Match Go community conventions, not conventions from other languages.

| Element | Rule | Good | Bad |
|---|---|---|---|
| Packages | short, lowercase, single word, no underscores/plurals | `cpu`, `netlink` | `cpuCollector`, `net_utils` |
| Exported identifiers | `PascalCase`, no stutter with package name | `cpu.Info` | `cpu.CPUInfo` |
| Unexported identifiers | `camelCase` | `readStat` | `read_stat` |
| Interfaces | behavior + `-er` suffix where natural | `Collector`, `Reader` | `ICollector` |
| Acronyms | keep one case | `parseCPUID`, `HTTPClient` | `parseCpuId`, `HttpClient` |
| Errors | sentinel: `Err` prefix; type: `Error` suffix | `ErrNotSupported`, `ParseError` | `NotSupportedErr` |

Detailed cross-project naming (files, CLI commands, flags, env vars) lives in `naming-conventions.md`.

---

## Package Design

- Keep packages **small and focused**: one domain per package (`cpu`, `memory`, `disk`). This mirrors the independent-collector rule in `../specs/architecture.md`.
- Put implementation detail under `internal/` so it cannot be imported by external code.
- A package's public surface is a deliberate decision. Export only what callers need.
- Avoid `util`, `common`, and `helpers` grab-bag packages. Name packages after what they provide.
- Avoid import cycles by depending downward through the architecture layers only (CLI → Command → Service → Collector → Platform).

---

## Error Handling

No panics in library code. Return errors; let `main` decide how to exit.

- Wrap errors with `%w` to preserve the chain, and add context at each layer:

```go
data, err := os.ReadFile("/proc/stat")
if err != nil {
	return Info{}, fmt.Errorf("reading /proc/stat: %w", err)
}
```

- Do **not** prefix wrapped messages with "failed to" or "error:"; the chain already reads as a sequence. Write lowercase, no trailing punctuation.
- Define **sentinel errors** for conditions callers branch on, and compare with `errors.Is`:

```go
var ErrNotSupported = errors.New("collector not supported on this kernel")

if errors.Is(err, cpu.ErrNotSupported) {
	// degrade gracefully
}
```

- Use `errors.As` for typed errors carrying data (e.g., a parse position).
- `panic` is permitted only for truly unrecoverable programmer errors (e.g., an impossible switch default). It must never be triggered by system state or user input.
- Never discard errors with `_` unless the call genuinely cannot fail and a comment says why.

---

## Struct and Interface Design

**Accept interfaces, return structs.** Functions take the narrow interface they need and return concrete types.

```go
// Reader is the minimal dependency the collector needs.
type Reader interface {
	Read(path string) ([]byte, error)
}

// New returns a concrete *Collector; callers get a real type, not an interface.
func New(r Reader) *Collector {
	return &Collector{r: r}
}
```

- Keep interfaces **small** — one or two methods is ideal. Define them at the point of use (the consumer), not alongside the implementation.
- Do not define an interface until there is a second implementation or a genuine need for a test seam.
- Prefer value receivers for small immutable structs; use pointer receivers when the method mutates or the struct is large. Be consistent across a type's method set.
- Zero values should be useful where practical.

---

## Context Usage

- Any function that performs I/O, blocks, or may need cancellation takes `ctx context.Context` as its **first** parameter.
- Never store a `context.Context` in a struct; pass it through the call chain.
- Do not pass `nil` contexts. Use `context.Background()` at the top of `main` and derive from there.
- Do not use context values to pass optional parameters; reserve them for request-scoped, cross-cutting data.

```go
func (c *Collector) Collect(ctx context.Context) (Info, error) {
	select {
	case <-ctx.Done():
		return Info{}, ctx.Err()
	default:
	}
	// ...
}
```

---

## Avoiding Globals

- No mutable package-level state. Globals defeat testing, concurrency safety, and modularity.
- Dependencies are passed explicitly through constructors, not reached for via package variables.
- Package-level `var` is allowed only for immutable values: sentinel errors, lookup tables, and compiled regexes.
- Do not use `init()` for anything with observable side effects.

---

## Test Organization

- Tests live beside the code in `_test.go` files, in the same package for white-box tests or `package foo_test` for black-box API tests.
- Use **table-driven tests** for functions with multiple input/output cases:

```go
func TestParseLoadAvg(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    LoadAvg
		wantErr bool
	}{
		{name: "typical", in: "0.50 0.40 0.30 1/234 5678", want: LoadAvg{One: 0.50, Five: 0.40, Fifteen: 0.30}},
		{name: "empty", in: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLoadAvg(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
```

- Use `testify/assert` and `testify/require` for readable assertions; use `require` when a failure should stop the test.
- Name subtests descriptively so `-run TestParseLoadAvg/empty` targets one case.
- Use fake/stub implementations of interfaces (e.g., a fake `Reader` over an in-memory filesystem) instead of touching real `/proc` in unit tests. Real-system checks belong in integration tests.
- Benchmarks (`BenchmarkXxx`) are required for hot paths per the constitution's **Performance Matters** principle.

---

*Idiomatic Go is not a constraint we tolerate — it is the baseline that lets the rest of these standards hold.*
