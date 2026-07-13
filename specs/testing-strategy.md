# SysKit Testing Strategy

> How SysKit tests every component, from individual functions to full CLI output, on real Linux systems.

---

## Philosophy

Testing is a constraint, not a convenience. The constitution's **Test Everything** principle states plainly that "a feature without tests is not complete" — and this document defines what that means in practice.

A feature is complete when:

- Its logic is covered by unit tests, including edge cases and error paths.
- Its data collection is verified against realistic fixtures captured from actual Linux interfaces.
- Its performance-sensitive paths are benchmarked (per the **Performance Matters** principle).
- Its user-facing output is pinned by golden-file tests.
- All tests pass under the race detector.

Code that collects system data cannot be tested by running against the developer's own `/proc` — that output is non-deterministic and machine-specific. SysKit's testability rests on a single architectural decision: **collectors never read the filesystem directly.** They read through an interface, and tests supply that interface from fixtures. This is the practical realization of the architecture's "enable testing through interface-based design" concern in the Platform Abstraction Layer.

---

## Test Categories

SysKit uses four complementary categories of tests. Each answers a different question.

| Category | Question it answers | Where it runs | Speed |
|---|---|---|---|
| Unit | Does this function transform input to output correctly? | Everywhere | Fast |
| Integration | Do collectors produce correct data against real `/proc` & `/sys`? | Linux CI | Medium |
| Benchmark | Is this path fast enough, and has it regressed? | Linux CI | Slow |
| Golden-file / e2e | Does the CLI emit exactly the expected output? | Everywhere | Fast |

### Unit Tests

Unit tests cover parsing, transformation, aggregation, filtering, sorting, and formatting logic. They are the majority of the test suite and run on any platform because they depend only on fixtures and pure functions, never on the host kernel.

Every collector parser is unit-tested against captured fixture files. Every service transformation is unit-tested against synthetic input models. Every formatter is unit-tested against known data structures.

### Integration Tests

Integration tests exercise collectors against the **real** `/proc` and `/sys` of the CI runner. They confirm that the abstraction correctly reads live kernel interfaces — something fixtures alone cannot guarantee, because fixtures can drift from reality.

Because they require a Linux kernel, these tests are guarded with a build tag and run only on Linux CI:

```go
//go:build linux && integration

package cpu_test

func TestCollectorReadsRealProcStat(t *testing.T) {
    c := cpu.NewCollector(platform.RealFS())
    stat, err := c.Collect()
    require.NoError(t, err)
    require.NotZero(t, stat.LogicalCores)
}
```

Integration tests assert on invariants that must hold on any Linux host (non-zero core count, monotonic counters, present mount points) rather than exact values, since those vary per machine.

### Benchmarks

Benchmarks track the performance of hot paths — parsing `/proc/stat`, walking `/proc/[pid]`, formatting large tables. They uphold the **Performance Matters** principle by making regressions visible.

```go
func BenchmarkParseProcStat(b *testing.B) {
    data := testdata.Read(b, "proc/stat")
    b.ReportAllocs()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cpu.ParseStat(data)
    }
}
```

Benchmarks always call `b.ReportAllocs()` — allocation count is a first-class metric because the constitution requires minimizing allocations in hot paths.

### Golden-File & End-to-End Tests

Golden-file tests capture the exact CLI output for a command and compare future runs against the stored "golden" file. They protect the output contracts described in [cli-conventions.md](cli-conventions.md): once `--format json` emits a shape, that shape is pinned.

```go
func TestCPUCommandJSONOutput(t *testing.T) {
    out := runCommand(t, cpu.FS(testdata.Fixtures), "cpu", "--format", "json")
    golden.Assert(t, out, "cpu_json.golden")
}
```

Golden files are regenerated intentionally with `go test ./... -update`, and the diff is reviewed like any other change. An unexpected golden diff is a signal that a user-visible contract moved.

---

## Mocking `/proc` and `/sys`

Collectors depend on a filesystem interface, never on the concrete OS filesystem. This is the single most important testability decision in SysKit.

### The Filesystem Abstraction

The platform abstraction layer exposes a narrow, read-only interface built around the standard library's `fs.FS`, extended with the file-reading semantics that `/proc` and `/sys` require:

```go
// Package platform abstracts read access to procfs and sysfs so that
// collectors can be tested against fixtures. See constitution: "Test Everything".
type SysFS interface {
    // ReadFile reads the named file relative to the mount root
    // (e.g. "proc/stat", "sys/devices/system/cpu/present").
    ReadFile(name string) ([]byte, error)
    // Open opens the named file for streaming reads of large pseudo-files.
    Open(name string) (fs.File, error)
    // ReadDir lists a pseudo-directory (e.g. "proc" to enumerate PIDs).
    ReadDir(name string) ([]fs.DirEntry, error)
    // ReadLink reads a procfs symbolic link such as proc/<pid>/fd/<fd>.
    ReadLink(name string) (string, error)
}
```

Production code uses a real implementation rooted at `/`:

```go
func RealFS() SysFS { return osFS{root: "/"} }
```

Tests use an `fs.FS`-backed implementation rooted at a fixtures directory, so `ReadFile("proc/stat")` resolves to `testdata/proc/stat`:

```go
func TestParseLoadAverage(t *testing.T) {
    fsys := platform.TestFS(os.DirFS("testdata/fixtures/idle-host"))
    c := system.NewCollector(fsys)

    info, err := c.Collect()
    require.NoError(t, err)
    require.Equal(t, 0.12, info.LoadAvg1)
}
```

Because collectors receive their `SysFS` by injection, the same collector code runs unchanged in production (against `/`) and in tests (against `testdata/`). No `runtime.GOOS` branching, no conditional compilation — consistent with the constitution's **Linux First** and **Clean Go** principles.

### Capturing Fixtures

Fixtures are real files captured from actual Linux hosts, preserving the exact byte format the kernel emits — whitespace, ordering, and all. They live under each package's `testdata/` directory in a layout mirroring the real filesystem:

```text
collector/cpu/testdata/
├── fixtures/
│   ├── 8-core-xeon/
│   │   ├── proc/stat
│   │   ├── proc/cpuinfo
│   │   └── sys/devices/system/cpu/present
│   └── 1-core-vm/
│       ├── proc/stat
│       └── proc/cpuinfo
└── golden/
    ├── cpu_table.golden
    └── cpu_json.golden
```

Multiple fixture sets capture meaningful variation: many cores vs. one, cgroup v1 vs. v2, kernel version differences, missing optional files. Capturing a diverse corpus is how SysKit satisfies "provide consistent data shapes regardless of kernel version variations" (architecture, Collector Layer).

Fixture capture should be automated once implementation begins. The planned helper path is `scripts/capture-fixtures.sh`; when added, it must record provenance — kernel version, distribution, architecture, container status, and capture date — alongside each fixture set in a `SOURCE` file.

---

## Table-Driven Tests

SysKit favors table-driven tests, the idiomatic Go pattern, for exhaustively covering parsing and transformation branches:

```go
func TestParseMemAvailable(t *testing.T) {
    tests := []struct {
        name    string
        fixture string
        want    uint64
        wantErr error
    }{
        {"standard host", "meminfo_standard", 6_291_456, nil},
        {"no MemAvailable field", "meminfo_legacy", 0, memory.ErrFieldMissing},
        {"malformed value", "meminfo_corrupt", 0, memory.ErrParse},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            data := testdata.Read(t, "proc/"+tt.fixture)
            got, err := memory.ParseMemAvailable(data)
            require.ErrorIs(t, err, tt.wantErr)
            require.Equal(t, tt.want, got)
        })
    }
}
```

Each row names a scenario, and subtests (`t.Run`) make failures pinpoint the exact case. Error assertions use `errors.Is` against the sentinel errors defined in [error-handling.md](error-handling.md).

---

## Race Detector

Concurrency appears wherever SysKit parallelizes collection or refreshes live data (dashboard, `watch`, `top`). All tests run under the race detector in CI:

```sh
go test -race ./...
```

Any data race is a failing build. The **Clean Go** principle and Go's concurrency model give us the tools; the race detector is how we prove we used them correctly.

---

## Coverage Expectations

Coverage is a signal, not a target to be gamed. The expectation is that logic is meaningfully exercised, not that a percentage is hit.

| Layer | Expectation |
|---|---|
| Collectors (parsers) | Every parse branch and error path covered via fixtures |
| Services (logic) | All transformation, filtering, and derived-metric paths covered |
| Formatters | Every output format covered by golden files |
| Commands | Flag validation and error surfacing covered |
| CLI wiring | Covered by end-to-end golden tests |

As a guideline, the project targets **≥ 80% statement coverage** overall, with the understanding that the parsing and logic layers should approach complete coverage while thin glue code may fall below. Coverage is reported in CI but a raw number never overrides review judgment about whether the *right* things are tested.

---

## Continuous Integration

Testing is enforced in CI, defined in [`.github/workflows/ci.yml`](../.github/workflows/ci.yml). Every push and pull request runs the full suite on Linux runners. The pipeline enforces, at minimum:

| Stage | Command | Purpose |
|---|---|---|
| Format | `gofmt -l .` | No unformatted files |
| Vet | `go vet ./...` | Static analysis clean |
| Unit + race | `go test -race ./...` | Correctness and race-freedom |
| Integration | `go test -tags=integration ./...` | Real `/proc` & `/sys` validation |
| Coverage | `go test -coverprofile=...` | Coverage reporting |
| Benchmarks | `go test -run=^$ -bench=.` | Regression tracking |

Because CI runs on Linux, integration tests execute against a genuine kernel — the environment SysKit is built for. A pull request that reduces coverage, breaks a golden file without justification, or introduces a race does not merge.

### Benchmark Tracking

Benchmark results are recorded per commit so that performance trends are visible over time. A significant regression in a hot-path benchmark (execution time or allocation count) is treated as a defect under the **Performance Matters** principle and is investigated before release, not after.

---

## What "Complete" Means

Restating the constitution in operational terms — a change is not complete until:

1. New logic has unit tests covering success, edge, and error cases.
2. New collectors have fixtures and, where a kernel is required, integration tests.
3. New or changed CLI output has golden files.
4. Performance-sensitive additions have benchmarks.
5. The full suite passes under `-race` in CI.

Until all five hold, the feature is in progress — not done.

---

*Tests are how SysKit earns trust: in its correctness, in its performance, and in the stability of its contracts. They are written with the feature, never after it.*
