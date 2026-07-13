# Performance Baseline

> Reproducible workloads and regression rules for SysKit's performance-sensitive paths.

## Scope

SysKit benchmarks deterministic parsing, collection, and rendering workloads. They are
not end-to-end latency service-level objectives: host kernel state, process count,
storage, CPU frequency, and CI virtualization all affect wall-clock command latency.
The benchmarks instead make like-for-like regressions visible while keeping allocation
costs explicit.

The required hot paths are:

- parsing aggregate and per-core `/proc/stat` counters;
- walking a fixture-backed `/proc` containing 1,000 processes;
- formatting a four-column table containing 1,000 rows;
- parsing representative CPU metadata, network counters, and socket tables.

Every benchmark calls `b.ReportAllocs()` and builds its fixture outside the timed loop.

## Running the Suite

Run the same command used by CI:

```sh
go test -run=^$ -bench=. -benchmem ./...
```

For a comparison suitable for review, capture at least five samples from the base and
candidate revisions on the same otherwise-idle Linux host, then compare them with
`benchstat`:

```sh
go test -run=^$ -bench=. -benchmem -count=5 ./... > old.txt
go test -run=^$ -bench=. -benchmem -count=5 ./... > new.txt
benchstat old.txt new.txt
```

## Initial Baseline

The initial v1 stabilization baseline was captured on 2026-07-13 with Go 1.26.3,
linux/amd64, Linux 6.17, and an Intel Core i5-6500. Values below are orientation
points from one run; repository and CI benchmark output remain the authoritative
per-commit record.

| Workload | Time | Bytes/op | Allocs/op |
|---|---:|---:|---:|
| Parse CPU metadata | 4.70 µs | 440 | 9 |
| Parse `/proc/stat` (aggregate + 4 cores) | 12.14 µs | 3,008 | 11 |
| Parse one process stat row | 1.85 µs | 560 | 3 |
| Walk 1,000 fixture processes | 203.26 ms | 2,154,282 | 34,827 |
| Parse 3 network interfaces | 7.95 µs | 2,464 | 8 |
| Parse 2 TCP sockets | 8.22 µs | 2,352 | 20 |
| Render a 1,000-row table | 443 µs | 154,416 | 22 |

The fixture filesystem deliberately favors determinism over simulating real procfs
latency. The process-walk number therefore measures SysKit plus `testing/fstest`
overhead and must not be presented as expected live-host command latency.

During the baseline sweep, table rendering was changed to write cells and padding
directly to one builder instead of allocating a field slice and padding strings per
row. On the baseline host, five post-change samples ranged from 417–524 µs with 22
allocations, versus the initial 1.34 ms and 4,077 allocations. Golden tests verify
that this optimization did not alter rendered output.

## Regression Policy

A change requires investigation when a same-host comparison shows either:

- a statistically meaningful time regression greater than 15%; or
- any unexplained increase greater than 10% in bytes or allocations per operation.

Reviewers should repeat noisy runs before drawing a conclusion. A justified tradeoff
must be recorded in the change description and, when it changes a lasting performance
expectation, in this document. Optimizations must preserve output contracts and may not
bypass the collector or platform boundaries.
