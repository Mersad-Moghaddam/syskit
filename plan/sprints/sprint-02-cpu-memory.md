# Sprint 02 — CPU & Memory

**Dates:** TBD → +2 weeks · **Milestone:** v0.1 (Core Inspection) · **Committed points:** 21

## Sprint goal

Ship `syskit cpu` (static identity + utilization + per-core) and `syskit memory`, introducing the two-sample derived-rate pattern that later real-time features reuse.

## Capacity

- Developers available: core team at nominal velocity.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| CPU-01 | `syskit cpu` static topology/model/cache. | 5 | Done |
| CPU-02 | `syskit cpu` utilization + `--per-core`. | 8 | Done |
| MEM-01 | `syskit memory`. | 8 | Done |

_Authoritative in ../product-backlog.md._

## Task breakdowns

**CPU-01 / CPU-02** — spec: `../../specs/features/cpu.md`
- [ ] platform: `/proc/cpuinfo`, `/proc/stat`, `/sys/devices/system/cpu/` via `SysFS` (+ fixtures: 8-core, 1-core VM, missing cpufreq).
- [ ] collector: parse identity; parse `/proc/stat` counters; hold two timestamped samples.
- [ ] service: derive aggregate + per-core utilization (rate from two samples); keep raw counters separate.
- [ ] command: `--per-core` flag; utilization requires two samples.
- [ ] render: table (per spec's columns) + JSON with `cpu_id`, `user`, `system`, `idle`, `iowait`, `steal`, `guest`, `total`, derived util (+ golden).
- [ ] tests: parse standard/virtualized/malformed fixtures; benchmark `ParseStat`; Linux integration (non-zero cores, monotonic counters).
- [ ] docs + CHANGELOG.

**MEM-01** — spec: `../../specs/features/memory.md`
- [ ] platform: `/proc/meminfo`, `/proc/vmstat`, pressure (PSI) if available (+ fixtures: standard, legacy no-`MemAvailable`, high-pressure).
- [ ] collector: parse fields; missing `MemAvailable` → `ErrFieldMissing`.
- [ ] service: derive used/available/pressure; raw vs derived separated.
- [ ] render: table + JSON with units (+ golden); integration; benchmark.
- [ ] docs + CHANGELOG.

## Definition of Ready / Done

`../../standards/definition-of-ready.md` / `../../standards/definition-of-done.md`. Acceptance highlights: utilization only with two samples; per-core rows preserve `/proc/stat` IDs; missing frequency = *unavailable*, not zero.

## Risks this sprint

- **R-03 (fixture drift)** — capture virtualized fixtures; VMs hide topology (CPU spec edge case).
- **R-07 (kernel variants)** — assert invariants (non-zero cores), not machine-specific values.

## Dependencies

Blocked by EPIC-00 (collector interface, render, harness). CPU-01 → CPU-02 (identity before utilization).

## Notes

CPU-02's two-sample pattern is deliberately built here so RT-01 (Sprint 7) inherits a proven approach to derived rates.
