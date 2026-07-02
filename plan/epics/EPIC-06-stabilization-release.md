# EPIC-06 — Stabilization & v1.0 Release

> **Milestone:** v1.0 · **Sprint:** 14 · **Points:** 39
> Everything required to call SysKit stable, documented, and installable — culminating in **v1.0.0**.

---

## Goal

Deliver the diagnostics command, optimize hot paths to baseline targets, freeze the CLI/output contracts under SemVer, package for major distros, automate releases, and complete the documentation. When this epic is Done, SysKit is a production-ready 1.0 that downstream users can install and depend on.

## Success criteria

- `syskit diagnostics` runs health checks and flags resource bottlenecks (`../specs/features/diagnostics.md`).
- Hot-path benchmarks meet baseline targets with no unexplained regressions (`../specs/testing-strategy.md`).
- Public CLI commands, flags, and structured-output schemas are **frozen** and covered by SemVer stability guarantees (`../standards/versioning.md`).
- Installable packages exist and are verified: deb, rpm, AUR, static binaries with checksums.
- Release automation generates the changelog and artifacts from Conventional Commits.
- Documentation is complete: every command, install paths, and man pages.
- `v1.0.0` is tagged and announced.

## Stories

Authoritative list: [`../product-backlog.md`](../product-backlog.md#epic-06--stabilization--v10-release).

| ID | Story | Pts | Sprint | Spec |
|---|---|---|---|---|
| DIA-01 | `syskit diagnostics`. | 13 | 14 | `../specs/features/diagnostics.md` |
| PERF-01 | Benchmark sweep + hot-path optimization. | 5 | 14 | `../specs/testing-strategy.md` |
| API-01 | CLI/output contract freeze + SemVer guarantees. | 5 | 14 | `../standards/versioning.md` |
| PKG-01 | Packaging (deb, rpm, AUR, binaries, checksums). | 8 | 14 | `../docs/release-process.md` |
| REL-01 | Release automation + changelog generation. | 5 | 14 | `../standards/commit-conventions.md` |
| DOCS-01 | Documentation completion pass + man pages. | 3 | 14 | `../docs/documentation-standards.md` |
| REL-v10 | Release v1.0.0 (final gate + tag + announce). | — | 14 | `../docs/release-process.md` |

## Dependencies & risk

- Depends on all prior epics: diagnostics aggregates signals from CPU, memory, disk, process, and network services.
- **R-12 (packaging surprises)** — mitigated by the SPK-PKG dry-run scheduled back in Sprint 9 slack, so PKG-01 is not the first `.deb` ever built.
- This sprint carries 39 pts — above nominal velocity. Per the release plan it **may be split into 14a/14b** at Sprint 13 refinement based on measured velocity; DIA-01 can lead its own sub-sprint since it is the only large item.

## Definition of Done for the epic

Every story meets `../standards/definition-of-done.md`; contracts are frozen and documented; packages install and run on target distros in CI or a verified manual matrix; a release dry-run produces correct artifacts and changelog; `v1.0.0` tagged on green `main`.

## The v1.0 gate

`v1.0.0` releases only when **all** milestone exit criteria in `../release-plan.md` are met and the diagnostics, packaging, and documentation stories are Done. Reaching this tag is "100" — the completion target of this entire plan.
