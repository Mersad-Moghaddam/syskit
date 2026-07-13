# SysKit Product Backlog

> **Source of truth** for all planned work. Every story has a stable ID, a size in story points, an owning epic, a target sprint, and a status. Sprint and epic files reference these IDs; they never redefine them.

---

## Legend

- **Status:** `Backlog` → `Ready` → `Committed` → `In Progress` → `In Review` → `Done`
- **Points:** modified Fibonacci (1, 2, 3, 5, 8, 13). See [`estimation-and-velocity.md`](estimation-and-velocity.md).
- **Spec:** canonical behavioral source in `../specs/`. A story is only *Ready* when its spec is accepted.
- IDs are stable. If a story is dropped, its ID is retired, never reused.

---

## Epic overview

| Epic | Title | Milestone | Points | Sprints |
|---|---|---|---|---|
| [EPIC-00](epics/EPIC-00-foundation.md) | Foundation & Delivery Infrastructure | v0.1 | 55 | 0–1 |
| [EPIC-01](epics/EPIC-01-core-inspection.md) | Core System Inspection | v0.1 | 51 | 1–3 |
| [EPIC-02](epics/EPIC-02-process-network.md) | Processes & Networking | v0.2 | 63 | 4–6 |
| [EPIC-03](epics/EPIC-03-realtime-monitoring.md) | Real-Time Monitoring | v0.3 | 50 | 7–9 |
| [EPIC-04](epics/EPIC-04-containers.md) | Containers | v0.4 | 42 | 10–11 |
| [EPIC-05](epics/EPIC-05-extensibility.md) | Extensibility (Plugins) | v0.5 | 45 | 12–13 |
| [EPIC-06](epics/EPIC-06-stabilization-release.md) | Stabilization & v1.0 Release | v1.0 | 39 | 14 |
| [EPIC-07](epics/EPIC-07-cross-cutting-quality.md) | Cross-Cutting Quality | all | ongoing | all |

Total scheduled: **345 points** across 15 sprints (Sprint 0–14). Cross-cutting quality work (EPIC-07) is embedded in every story's Definition of Done rather than separately pointed.

---

## EPIC-00 — Foundation & Delivery Infrastructure

Turns the planning repo into a compiling, tested Go project without shipping user features prematurely. Gated by `../docs/implementation-readiness.md`.

| ID | Story | Pts | Sprint | Status | Spec / Ref |
|---|---|---|---|---|---|
| FND-01 | As a maintainer, I complete the implementation-readiness sign-off so code may begin. | 3 | 0 | Done | `../docs/implementation-readiness.md` |
| FND-02 | Transition PR: create Go module, approved repo layout, update README status, CI, contributing. | 5 | 0 | Done | `../docs/project-structure.md` |
| FND-03 | CLI bootstrap with Cobra: root command, `--format`, `--help`, `version`. | 5 | 0 | Done | `../specs/cli-conventions.md`, `../decisions/005-cobra-for-cli.md` |
| FND-04 | Platform abstraction `SysFS` interface + `RealFS` + fixture-backed `TestFS`. | 8 | 0 | Done | `../specs/testing-strategy.md`, `../specs/architecture.md` |
| FND-05 | Collector interface + registration pattern (no cross-collector deps). | 5 | 1 | Done | `../specs/collectors.md` |
| FND-06 | Render layer skeleton: `Formatter` interface + table + JSON. | 8 | 1 | Done | `../specs/rendering.md`, `../specs/features/output-formats.md` |
| FND-07 | Error-handling patterns: sentinel errors, `%w` wrapping, exit codes. | 3 | 1 | Done | `../specs/error-handling.md` |
| FND-08 | Logging strategy scaffolding (structured, off by default). | 3 | 1 | Done | `../specs/logging-strategy.md` |
| FND-09 | Configuration loading (precedence: flags > env > file > default). | 5 | 1 | Done | `../specs/configuration.md` |
| FND-10 | Go CI pipeline: fmt, vet, `test -race`, integration tag, coverage, bench, govulncheck. | 8 | 0–1 | Done | `../specs/testing-strategy.md`, `../.github/workflows/ci.yml` |
| FND-11 | Test harness: golden helper, `testdata/` layout, `scripts/capture-fixtures.sh`. | 5 | 1 | Done | `../specs/testing-strategy.md` |

Subtotal: **58 pts** (Sprint 0: 21, Sprint 1: 37 — see release plan for exact split; FND-10 spans both).

---

## EPIC-01 — Core System Inspection (v0.1)

The first four inspection commands, each a full vertical slice: collector → service → command → render → tests.

| ID | Story | Pts | Sprint | Status | Spec / Ref |
|---|---|---|---|---|---|
| SYS-01 | `syskit system` — host, kernel, OS release, uptime, load averages. | 8 | 1 | Done | `../specs/features/system.md` |
| CPU-01 | `syskit cpu` — topology, model, cache, static identity. | 5 | 2 | Done | `../specs/features/cpu.md` |
| CPU-02 | `syskit cpu` utilization — two-sample derivation, `--per-core`. | 8 | 2 | Done | `../specs/features/cpu.md` |
| MEM-01 | `syskit memory` — physical/swap, buffers, caches, available, pressure. | 8 | 2 | Done | `../specs/features/memory.md` |
| DSK-01 | `syskit disk` — partition layout, usage, mount points, I/O stats. | 8 | 3 | Done | `../specs/features/disk.md` |
| FS-01 | `syskit filesystem` — inode usage, fs types, mount options. | 8 | 3 | Done | `../specs/features/filesystem.md` |
| REL-v01 | Release v0.1.0 — golden/integration hardening, CHANGELOG, tag. | 3 | 3 | Done | `../docs/release-process.md` |
| DOC-v01 | v0.1 user docs: getting-started build/install + command pages. | 3 | 3 | Done | `../docs/getting-started.md` |

Subtotal: **51 pts**.

---

## EPIC-02 — Processes & Networking (v0.2)

Process inspection, network visibility, and the filtering/sorting framework shared across commands.

| ID | Story | Pts | Sprint | Status | Spec / Ref |
|---|---|---|---|---|---|
| FLT-01 | Filtering & sorting framework (reusable across list commands). | 8 | 4 | Done | `../specs/cli-conventions.md` |
| PRC-01 | `syskit process` — listing, resource usage, filter by name/PID/user. | 13 | 4 | Done | `../specs/features/process.md` |
| PRC-02 | `syskit process tree` — tree view, handle disappearing PIDs. | 8 | 5 | Done | `../specs/features/process.md` |
| NET-01 | Netlink platform integration (sockets, message parse). | 13 | 5 | Done | `../specs/features/network.md`, `../decisions/003-native-apis-over-shell.md` |
| NET-02 | `syskit network` — interface stats, connections, routing. | 8 | 5 | Done | `../specs/features/network.md` |
| PRT-01 | `syskit ports` — listening ports, socket states, owning process. | 8 | 6 | Done | `../specs/features/ports.md` |
| OUT-03 | YAML output formatter. | 5 | 6 | Done | `../specs/features/output-formats.md` |
| REL-v02 | Release v0.2.0 — hardening, CHANGELOG, tag. | 3 | 6 | Done | `../docs/release-process.md` |
| DOC-v02 | v0.2 docs + glossary updates. | 2 | 6 | Done | `../docs/glossary.md` |

Subtotal: **68 pts**.

---

## EPIC-03 — Real-Time Monitoring (v0.3)

Interactive terminal UI, live refresh, and monitor commands.

| ID | Story | Pts | Sprint | Status | Spec / Ref |
|---|---|---|---|---|---|
| TUI-01 | Bubble Tea + Lip Gloss foundation (model/update/view, styling). | 8 | 7 | Done | `../specs/features/dashboard.md`, `../decisions/006-bubbletea-for-tui.md` |
| RT-01 | Real-time refresh pipeline (concurrent, race-safe, backpressure). | 8 | 7 | Done | `../specs/features/dashboard.md`, `../specs/testing-strategy.md` |
| WCH-01 | `syskit watch <command> --interval` — continuous refresh. | 5 | 7 | Done | `../specs/features/dashboard.md`, `../specs/cli-conventions.md` |
| DSH-01 | `syskit dashboard` — layout/widget system, real-time metrics. | 13 | 8 | Done | `../specs/features/dashboard.md` |
| TOP-01 | `syskit top` — interactive process monitor, sort/filter/keys. | 13 | 9 | Done | `../specs/features/process.md`, `../specs/features/dashboard.md` |
| REL-v03 | Release v0.3.0 — hardening, CHANGELOG, tag. | 3 | 9 | Done | `../docs/release-process.md` |

Subtotal: **50 pts**.

---

## EPIC-04 — Containers (v0.4)

Container runtime inspection and cgroup-based resource views.

| ID | Story | Pts | Sprint | Status | Spec / Ref |
|---|---|---|---|---|---|
| CG-01 | cgroup v1/v2 parsing in platform layer. | 13 | 10 | Done | `../specs/features/containers.md` |
| CNT-01 | Container-to-process mapping + container-aware process views. | 8 | 10 | Done | `../specs/features/containers.md` |
| DOC-01 | `syskit containers` — container listing, resource usage, status. | 8 | 11 | Done | `../specs/features/containers.md` |
| DOC-02 | `syskit containers inspect <id>` — detailed inspection. | 8 | 11 | Done | `../specs/features/containers.md` |
| REL-v04 | Release v0.4.0 — hardening, CHANGELOG, tag. | 3 | 11 | Done | `../docs/release-process.md` |
| DOC-v04 | v0.4 docs + container concepts in learning notes. | 2 | 11 | Done | `../learning/roadmap.md` |

Subtotal: **42 pts**.

---

## EPIC-05 — Extensibility / Plugins (v0.5)

Out-of-process plugin system per `../decisions/007-out-of-process-plugins.md`.

| ID | Story | Pts | Sprint | Status | Spec / Ref |
|---|---|---|---|---|---|
| PLG-01 | Plugin interface + out-of-process protocol definition. | 13 | 12 | In Progress | `../specs/plugin-architecture.md`, `../specs/features/plugins.md` |
| PLG-02 | Plugin discovery & loading mechanism. | 8 | 12 | In Progress | `../specs/plugin-architecture.md` |
| PLG-03 | Custom collector registration via plugins. | 5 | 12 | In Progress | `../specs/features/plugins.md` |
| PLG-04 | Plugin isolation & security model. | 8 | 13 | In Progress | `../specs/plugin-architecture.md`, `../SECURITY.md` |
| PLG-05 | Plugin configuration system. | 3 | 13 | Backlog | `../specs/configuration.md` |
| PLG-06 | Plugin SDK, example plugin, and authoring docs. | 5 | 13 | Done | `../specs/features/plugins.md` |
| REL-v05 | Release v0.5.0 — hardening, CHANGELOG, tag. | 3 | 13 | Backlog | `../docs/release-process.md` |

Subtotal: **45 pts**.

---

## EPIC-06 — Stabilization & v1.0 Release

Everything required to call SysKit stable and shippable to distros.

| ID | Story | Pts | Sprint | Status | Spec / Ref |
|---|---|---|---|---|---|
| DIA-01 | `syskit diagnostics` — health checks, bottleneck detection. | 13 | 14 | In Progress | `../specs/features/diagnostics.md` |
| PERF-01 | Benchmark sweep + hot-path optimization to baseline targets. | 5 | 14 | Backlog | `../specs/testing-strategy.md` |
| API-01 | CLI/output contract freeze + SemVer stability guarantees. | 5 | 14 | Backlog | `../standards/versioning.md` |
| PKG-01 | Packaging: deb, rpm, AUR, static binaries, checksums. | 8 | 14 | Backlog | `../docs/release-process.md` |
| REL-01 | Release automation + changelog generation. | 5 | 14 | In Progress | `../docs/release-process.md`, `../standards/commit-conventions.md` |
| DOCS-01 | Documentation completion pass (all commands, install, man pages). | 3 | 14 | Backlog | `../docs/documentation-standards.md` |
| REL-v10 | Release v1.0.0 — final gate, tag, announcement. | — | 14 | Backlog | `../docs/release-process.md` |

Subtotal: **39 pts**.

---

## EPIC-07 — Cross-Cutting Quality (ongoing)

Not separately scheduled; enforced inside every story's Definition of Done. Tracked so it is never invisible.

| ID | Concern | Applies to | Ref |
|---|---|---|---|
| XQ-01 | Unit tests (parse/transform/format/error paths). | every story | `../standards/definition-of-done.md` |
| XQ-02 | Fixture-backed collector tests. | collector stories | `../specs/testing-strategy.md` |
| XQ-03 | Linux integration tests (build-tagged). | collector stories | `../specs/testing-strategy.md` |
| XQ-04 | Golden-file output tests. | command stories | `../specs/testing-strategy.md` |
| XQ-05 | Benchmarks on hot paths. | parse/render stories | `../specs/testing-strategy.md` |
| XQ-06 | `gofmt`/`goimports`/`go vet`/`govulncheck` clean. | every story | `../standards/definition-of-done.md` |
| XQ-07 | Docs + CHANGELOG updated. | user-visible stories | `../standards/definition-of-done.md` |
| XQ-08 | Conventional Commits + squash-merge + rebase. | every PR | `../standards/branch-strategy.md`, `../standards/commit-conventions.md` |

---

## Dependency map (critical path)

```text
FND-02 ─▶ FND-03 ─▶ FND-05 ─┬─▶ (all collectors)
FND-04 ─────────────────────┘
FND-06 ─▶ (all command render output)
FND-10/FND-11 ─▶ (all tests / Done)

SYS-01 ─▶ CPU-01 ─▶ CPU-02      (SYS-01 proves the vertical slice first)
FLT-01 ─▶ PRC-01 ─▶ PRC-02 ─▶ TOP-01
NET-01 ─▶ NET-02 ─▶ PRT-01
TUI-01 ─▶ RT-01 ─▶ WCH-01 ─▶ DSH-01 ─▶ TOP-01
CG-01 ─▶ CNT-01 ─▶ DOC-01 ─▶ DOC-02
PLG-01 ─▶ PLG-02 ─▶ {PLG-03, PLG-04, PLG-05, PLG-06}
(most collectors) ─▶ DIA-01
```

**Longest chain:** Foundation → TUI-01 → RT-01 → WCH-01 → DSH-01 → TOP-01. This is why real-time work spans three sprints (7–9) and cannot be compressed without more parallel developers.

---

## Backlog hygiene rules

1. New work enters here first, at the bottom, in `Backlog` status.
2. The Product Owner orders; refinement pulls the top items to `Ready`.
3. A story is only assigned to a sprint once it is `Ready` (passes DoR).
4. Points are set by the Developers at refinement, never by the Product Owner.
5. Status changes here are mirrored by the sprint board; this file is authoritative on scope and points.
