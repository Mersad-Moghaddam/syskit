# Sprint 04 — Filtering Framework & Processes

**Dates:** TBD → +2 weeks · **Milestone:** v0.2 (Processes & Networking) · **Committed points:** 21

## Sprint goal

Build the reusable filtering/sorting framework and ship `syskit process` on top of it, establishing consistent list-command semantics for every later list feature (network, ports, top).

## Capacity

- Nominal velocity. PRC-01 is a 13 — the sprint's dominant story.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| FLT-01 | Filtering & sorting framework. | 8 | Committed |
| PRC-01 | `syskit process` listing + resource usage + filters. | 13 | Committed |

## Task breakdowns

**FLT-01** — ref `../../specs/cli-conventions.md`
- [ ] Define `--filter`/`--sort` flag grammar and precedence; reusable across list commands.
- [ ] service-level filter/sort primitives (typed, testable, no CLI coupling).
- [ ] unit tests: filter predicates, multi-key sort, invalid-flag errors.
- [ ] docs: document the shared flags once, referenced by each command.

**PRC-01** — spec `../../specs/features/process.md`
- [ ] platform: enumerate `/proc/[pid]`, read `stat`/`status`/`cmdline` via `SysFS` (+ fixtures with several processes).
- [ ] collector: parse per-process fields; tolerate a PID disappearing mid-read.
- [ ] service: resource usage; apply FLT-01 filter/sort (by name/PID/user).
- [ ] command: `syskit process` + shared filter/sort flags.
- [ ] render: table + JSON (+ golden, including a filtered+sorted variant).
- [ ] integration: real `/proc` walk returns the current process; benchmark the walk.
- [ ] docs + CHANGELOG.

## Definition of Ready / Done

Standard gates. Acceptance: filtering/sorting behaves identically to how later commands will use it; disappearing PIDs never crash the walk.

## Risks this sprint

- **R-05 precursor** — the `/proc` walk is perf-sensitive; benchmark it now (feeds later `top`).
- **SPK-NET** — run the Netlink spike during **this sprint's refinement** to de-risk NET-01 in Sprint 5 (`../risk-register.md`).

## Dependencies

Blocked by EPIC-01 (render/harness). FLT-01 must land before PRC-01's service task. FLT-01 is a hard dependency for NET-02, PRT-01, TOP-01.

## Notes

Investing in FLT-01 first is deliberate: three later commands consume it. Building it once prevents divergent filter semantics.
