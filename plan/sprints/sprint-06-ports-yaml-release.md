# Sprint 06 — Ports & YAML → v0.2.0

**Dates:** TBD → +2 weeks · **Milestone:** v0.2 — **release sprint** · **Committed points:** 18

## Sprint goal

Ship `syskit ports`, add YAML output across all commands, and tag **v0.2.0**.

## Capacity

- Below nominal by design: release overhead + buffer to absorb any NET-02 carry-over from Sprint 5.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| PRT-01 | `syskit ports`. | 8 | Committed |
| OUT-03 | YAML formatter. | 5 | Committed |
| DOC-v02 | v0.2 docs + glossary updates. | 2 | Committed |
| REL-v02 | Release v0.2.0. | 3 | Committed |

_If NET-02 carried over from Sprint 5, it is pulled in here first and lower-priority items defer._

## Task breakdowns

**PRT-01** — spec `../../specs/features/ports.md`
- [ ] collector/service: listening ports, socket states, owning process (via Netlink from NET-01 + `/proc` socket inode mapping).
- [ ] command + render (table + JSON, + golden); reuse FLT-01.
- [ ] integration + benchmark; docs + CHANGELOG.

**OUT-03** — spec `../../specs/features/output-formats.md`
- [ ] YAML formatter behind existing `Formatter` interface; field/unit parity with JSON.
- [ ] golden tests for every existing command's YAML output.
- [ ] docs + CHANGELOG.

**REL-v02** — release
- [ ] Confirm v0.2 milestone exit criteria (`../release-plan.md`): filtering consistent across list commands; YAML parity; Netlink integration green.
- [ ] CHANGELOG under v0.2.0; tag on green `main`.

## Definition of Ready / Done

Standard gates + v0.2 milestone gate.

## Risks this sprint

- **R-01 spillover** — the deliberate under-commit absorbs NET carry-over. If none, the team pulls the next Sprint 7 story forward.
- **R-09 (golden churn)** — OUT-03 adds a golden file per command; land it before the release hardening pass.

## Dependencies

Blocked by Sprint 5 (NET-01 for PRT-01). REL-v02 depends on PRT-01 + OUT-03 + DOC-v02.

## Notes

After this sprint, list-command semantics (filter/sort) and all three output formats are stable — the base for real-time work.
