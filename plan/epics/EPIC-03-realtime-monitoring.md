# EPIC-03 — Real-Time Monitoring

> **Milestone:** v0.3 · **Sprints:** 7–9 · **Points:** 50
> Interactive terminal UI, a race-safe live-refresh pipeline, and the `watch`, `dashboard`, and `top` commands — culminating in **v0.3.0**.

---

## Goal

Introduce SysKit's interactive experience: a Bubble Tea + Lip Gloss foundation, a concurrent refresh pipeline that never races, and three live commands that reuse existing collectors and services rather than duplicating collection logic.

## Why the pipeline is its own story

RT-01 (refresh pipeline) is separated from the TUI foundation because live refresh is where concurrency bugs live. Isolating it lets us design a single owner of state with channel-based updates and prove it under `-race` before three commands depend on it.

## Success criteria

- The TUI foundation follows the Elm-style model/update/view pattern (`../decisions/006-bubbletea-for-tui.md`).
- The refresh pipeline is race-free (`go test -race`), leak-free, and applies backpressure so a slow collector cannot flood the UI.
- `watch`, `dashboard`, and `top` reuse v0.1/v0.2 services — no re-collection code.
- Keybindings and refresh intervals are documented.
- `v0.3.0` is tagged on green `main`.

## Stories

Authoritative list: [`../product-backlog.md`](../product-backlog.md#epic-03--real-time-monitoring-v03).

| ID | Story | Pts | Sprint | Spec |
|---|---|---|---|---|
| TUI-01 | Bubble Tea + Lip Gloss foundation. | 8 | 7 | `../specs/features/dashboard.md` |
| RT-01 | Real-time refresh pipeline. | 8 | 7 | `../specs/features/dashboard.md` |
| WCH-01 | `syskit watch <command> --interval`. | 5 | 7 | `../specs/cli-conventions.md` |
| DSH-01 | `syskit dashboard` layout/widgets. | 13 | 8 | `../specs/features/dashboard.md` |
| TOP-01 | `syskit top` interactive process monitor. | 13 | 9 | `../specs/features/process.md` |
| REL-v03 | Release v0.3.0. | 3 | 9 | `../docs/release-process.md` |

## Dependencies & risk

- Blocked by EPIC-01 (all metric services) and EPIC-02 (process service feeds `top`).
- **R-05 (TUI concurrency races)** is the dominant risk — mitigated by the isolated pipeline story, single-owner-of-state design, and mandatory `-race` tests.
- The dependency chain TUI-01 → RT-01 → WCH-01 → DSH-01 → TOP-01 is the plan's longest; it cannot be compressed without additional parallel developers.

## Definition of Done for the epic

All stories meet the DoD; interactive components have unit-testable update logic; refresh pipeline passes `-race`; keybinding/interval behavior is documented and covered where testable; `v0.3.0` tagged.
