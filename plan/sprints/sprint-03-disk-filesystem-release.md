# Sprint 03 — Disk & Filesystem → v0.1.0

**Dates:** TBD → +2 weeks · **Milestone:** v0.1 (Core Inspection) — **release sprint** · **Committed points:** 22

## Sprint goal

Ship `syskit disk` and `syskit filesystem`, complete v0.1 user docs, and tag **v0.1.0** — SysKit's first release, proving the full architecture with five inspection commands in table and JSON.

## Capacity

- Nominal velocity, with slack reserved for release overhead (CHANGELOG, tag, exit-criteria check).

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| DSK-01 | `syskit disk`. | 8 | Done |
| FS-01 | `syskit filesystem`. | 8 | Done |
| DOC-v01 | v0.1 user docs (build/install + command pages). | 3 | Committed |
| REL-v01 | Release v0.1.0 (hardening, CHANGELOG, tag). | 3 | Committed |

## Task breakdowns

**DSK-01** — spec `../../specs/features/disk.md`
- [ ] platform: `/proc/mounts`, `/proc/diskstats`, statfs via `SysFS` (+ fixtures).
- [ ] collector/service: partition layout, usage, mount points, I/O stats.
- [ ] render: table + JSON (+ golden); integration; benchmark on `diskstats` parse.
- [ ] docs + CHANGELOG.

**FS-01** — spec `../../specs/features/filesystem.md`
- [ ] platform/collector: inode usage, fs types, mount options (+ fixtures).
- [ ] service/render: table + JSON (+ golden); integration.
- [ ] docs + CHANGELOG.

**REL-v01** — release
- [ ] Confirm v0.1 milestone exit criteria in `../release-plan.md`.
- [ ] Golden/integration hardening pass across all five commands.
- [ ] Finalize CHANGELOG under v0.1.0; verify getting-started build from a clean clone.
- [ ] Tag `v0.1.0` on green `main` (`../../standards/branch-strategy.md`).

## Definition of Ready / Done

Standard gates. **Milestone gate:** vertical slice proven end-to-end; table+JSON stable; CI runs unit+integration+golden on Linux; clean-clone build works.

## Risks this sprint

- **R-09 (golden churn)** — the release hardening pass is where cosmetic drift is caught; regenerate goldens intentionally with reviewed diffs.
- **R-12 (packaging)** — note: packaging is *not* in v0.1; only binary build. SPK-PKG is scheduled for Sprint 9.

## Dependencies

Blocked by Sprints 1–2 (render, collectors). REL-v01 depends on DSK-01 + FS-01 + DOC-v01 being Done.

## Notes

This is the first **release sprint** — rehearse the review's milestone-check and the tag process; later releases reuse it.
