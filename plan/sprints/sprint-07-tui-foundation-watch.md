# Sprint 07 — TUI Foundation & Watch

**Dates:** TBD → +2 weeks · **Milestone:** v0.3 (Real-Time Monitoring) · **Committed points:** 21

## Sprint goal

Stand up the interactive layer: a Bubble Tea + Lip Gloss foundation, a race-safe real-time refresh pipeline, and the first live command `syskit watch`.

## Capacity

- Nominal, with awareness that TUI is a new paradigm (warm-up within the sprint).

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| TUI-01 | Bubble Tea + Lip Gloss foundation. | 8 | Committed |
| RT-01 | Real-time refresh pipeline. | 8 | Committed |
| WCH-01 | `syskit watch <command> --interval`. | 5 | Committed |

## Task breakdowns

**TUI-01** — spec `../../specs/features/dashboard.md`, ADR `../../decisions/006-bubbletea-for-tui.md`
- [ ] model/update/view scaffold; Lip Gloss styles; keybinding + quit handling.
- [ ] unit-test the update function (pure state transitions).

**RT-01** — spec `../../specs/features/dashboard.md`, `../../specs/testing-strategy.md`
- [ ] single-owner-of-state refresh loop; channel-based updates; backpressure so a slow collector can't flood the UI.
- [ ] `go test -race` coverage; leak check (goroutines stop on exit).

**WCH-01** — ref `../../specs/cli-conventions.md`
- [ ] `syskit watch <command> --interval <dur>` reusing existing services (no re-collection).
- [ ] render existing table output on each tick; interval validation.
- [ ] tests (update logic + interval parsing); docs + CHANGELOG.

## Definition of Ready / Done

Standard gates. Acceptance: refresh pipeline passes `-race`; `watch` reuses v0.1/v0.2 services unchanged.

## Risks this sprint

- **R-05 (TUI concurrency) — HIGH.** RT-01 is isolated precisely so races are designed out before three commands depend on it. All tests under `-race`; no shared mutable state without a single owner.

## Dependencies

Blocked by EPIC-01/02 (services that feed live views). TUI-01 → RT-01 → WCH-01 (the epic's dependency chain begins here).

## Notes

Getting RT-01 right this sprint is what makes DSH-01 (Sprint 8) and TOP-01 (Sprint 9) tractable. Do not shortcut the pipeline to fit `watch` in.
