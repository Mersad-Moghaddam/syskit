# Sprint 09 — Top → v0.3.0

**Dates:** TBD → +2 weeks · **Milestone:** v0.3 — **release sprint** · **Committed points:** 16

## Sprint goal

Ship `syskit top` — an interactive process monitor with live sorting and filtering — and tag **v0.3.0**, completing the real-time milestone.

## Capacity

- Below nominal: one large story plus release overhead.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| TOP-01 | `syskit top` interactive process monitor. | 13 | Committed |
| REL-v03 | Release v0.3.0. | 3 | Committed |

## Task breakdowns

**TOP-01** — spec `../../specs/features/process.md` + `../../specs/features/dashboard.md`
- [ ] reuse the process service (PRC-01) + FLT-01 filters + RT-01 refresh pipeline.
- [ ] interactive: live re-sort (CPU/mem/PID), filter input, scroll, kill-key confirm (read-only scope respected — see note).
- [ ] render: full-screen live table via Bubble Tea; handle process churn between ticks.
- [ ] tests: update-function unit tests (sort/filter state); `-race`; golden on a frozen frame where feasible.
- [ ] docs: keybindings; CHANGELOG.

**REL-v03** — release
- [ ] Confirm v0.3 milestone exit criteria (`../release-plan.md`): TUI under `-race`, no leaks, keybindings documented for `watch`/`dashboard`/`top`.
- [ ] CHANGELOG under v0.3.0; tag on green `main`.

## Definition of Ready / Done

Standard gates + v0.3 milestone gate.

## Risks this sprint

- **R-05 (TUI concurrency)** — `top` combines refresh + interactive sort/filter; the highest-churn TUI surface. Race tests non-negotiable.
- Scope note: SysKit is **read-only** by charter (`../../docs/implementation-readiness.md` non-readiness signals). A "kill" affordance requires explicit spec/security review before inclusion — if unresolved, ship `top` as inspect-only and defer kill to a future story.

## Dependencies

Blocked by Sprints 4 (PRC-01, FLT-01) and 7 (RT-01). This is the tail of the plan's longest dependency chain.

## Notes

Completing v0.3 means SysKit now has static, structured, and live interactive modes — the feature core is done; remaining milestones add breadth (containers, plugins) and stability.
