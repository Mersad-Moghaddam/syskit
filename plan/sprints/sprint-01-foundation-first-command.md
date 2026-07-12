# Sprint 01 — Foundation Complete + First Command

**Dates:** TBD → +2 weeks · **Milestone:** v0.1 (Foundation) · **Committed points:** 24

## Sprint goal

Complete the delivery scaffolding — collector interface, render layer, error/logging/config patterns, test harness, full CI — and prove it end-to-end by shipping `syskit system` as the first real vertical slice.

## Capacity

- Developers available: core team.
- Still slightly below nominal; foundation stories carry setup overhead.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| FND-05 | Collector interface + registration. | 5 | Done |
| FND-06 | Render layer: `Formatter` + table + JSON. | 8 | Done |
| FND-07 | Error-handling patterns (sentinels, `%w`, exit codes). | 3 | Done |
| FND-08 | Logging scaffolding (structured, off by default). | 3 | Done |
| FND-09 | Configuration loading (flags > env > file > default). | 5 | Done |
| FND-11 | Test harness (golden helper, `testdata/`, capture script). | 5 | Done |
| FND-10b | Go CI pipeline — second half: `-race`, integration tag, coverage, bench, govulncheck. | (part of FND-10) | Done |
| SYS-01 | `syskit system` — first vertical slice. | 8 | Committed |

_Total exceeds a single dev's capacity — parallelized across the team within the WIP limit. Authoritative points in ../product-backlog.md._

## Task breakdowns

**FND-06 — render layer** (unblocks every command's output)
- [x] `Formatter` interface; table formatter (alignment, headers); JSON formatter.
- [x] Golden-test helper integration; unit tests per formatter.

**SYS-01 — system command** (the proof of the architecture)
- [ ] platform: read `/proc/uptime`, `/proc/loadavg`, `/etc/os-release`, kernel version via `SysFS` (+ fixtures).
- [ ] collector: parse into a `SystemInfo` struct (+ unit tests, error paths).
- [ ] service: assemble host summary.
- [ ] command: `syskit system` wiring + `--format`.
- [ ] render: table + JSON (+ golden files).
- [ ] integration: `//go:build linux && integration` — non-empty kernel version, sane uptime.
- [ ] docs: command help + getting-started example; CHANGELOG entry.

_Other stories decomposed per `../templates/task-breakdown.md` at planning._

## Definition of Ready / Done

Entry: `../../standards/definition-of-ready.md`. Exit: `../../standards/definition-of-done.md`. By sprint end, CI must run the full matrix (fmt/vet/race/integration/coverage/bench/vuln).

## Risks this sprint

- **R-09 (golden churn)** — stabilize the output contract in FND-06 *before* piling commands on it; SYS-01 is the first consumer and validates it.
- **R-03 (fixture drift)** — FND-11 delivers `scripts/capture-fixtures.sh` with provenance recording.

## Dependencies

Blocked by Sprint 0 (FND-04 for collectors, FND-03 for CLI). FND-05 + FND-06 must land before SYS-01's collector/render tasks.

## Notes

At the end of this sprint the **foundation epic (EPIC-00) is Done** and the team's true velocity starts to emerge. Begin using rolling velocity from Sprint 3 onward.
