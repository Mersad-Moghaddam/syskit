# Sprint 08 — Dashboard

**Dates:** TBD → +2 weeks · **Milestone:** v0.3 (Real-Time Monitoring) · **Committed points:** 13

## Sprint goal

Ship `syskit dashboard` — an interactive terminal dashboard with a layout/widget system showing real-time CPU, memory, disk, and network metrics.

## Capacity

- Single large story (13). Under-committed on purpose: DSH-01 is complex and integrative, and it leaves slack to run **SPK-PKG** (packaging dry-run) to de-risk Sprint 14.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| DSH-01 | `syskit dashboard` layout/widget system. | 13 | Committed |
| SPK-PKG | Packaging spike — one working `.deb` build. | (spike) | Committed |

## Task breakdowns

**DSH-01** — spec `../../specs/features/dashboard.md`
- [ ] layout system: arrange widgets in a responsive terminal grid (Lip Gloss).
- [ ] widgets: CPU, memory, disk, network — each subscribing to the RT-01 refresh pipeline.
- [ ] wire widgets to existing services (no re-collection); handle terminal resize.
- [ ] keyboard navigation between widgets/panes.
- [ ] tests: update-function unit tests per widget; `-race` on the composed app; golden where output is deterministic.
- [ ] docs: dashboard keybindings + interval; CHANGELOG.

**SPK-PKG** — retires R-12
- [ ] Produce one working `.deb` from a built binary; note the toolchain and gaps.
- [ ] Write findings so PKG-01 (Sprint 14) is not the first package ever built.

## Definition of Ready / Done

Standard gates. Acceptance: dashboard refreshes live without races or goroutine leaks; reuses existing services; resizes cleanly.

## Risks this sprint

- **R-05 (TUI concurrency)** — the most integrative consumer of RT-01; `-race` is mandatory.
- **R-12 (packaging)** — actively retired this sprint via SPK-PKG.

## Dependencies

Blocked by Sprint 7 (TUI-01, RT-01). Needs all v0.1/v0.2 metric services.

## Notes

Deliberately light on story points so the dashboard gets the focus it needs and the packaging risk is retired early rather than discovered in the final sprint.
