# EPIC-01 — Core System Inspection

> **Milestone:** v0.1 · **Sprints:** 1–3 · **Points:** 51
> The first four inspection commands, each delivered as a complete vertical slice.

---

## Goal

Ship `system`, `cpu`, `memory`, `disk`, and `filesystem` as working commands with table and JSON output, backed by fixture-tested collectors and golden-file output tests. This epic proves the architecture end-to-end and produces SysKit's first tagged release, **v0.1.0**.

## Why these, in this order

`system` is deliberately first (SYS-01): it is the simplest full slice (host/kernel/uptime/load) and validates that CLI → command → service → collector → platform → render works before the team invests in harder collectors. CPU utilization (CPU-02) introduces the two-sample/derived-rate pattern that later real-time features reuse.

## Success criteria

- Each command renders correct table and JSON output from fixtures and on a real Linux host.
- Static vs sampled data is cleanly separated (esp. CPU: identity vs utilization).
- Missing/optional data (e.g. cpufreq, `MemAvailable` on old kernels) is represented as *unavailable*, never as zero.
- `v0.1.0` is tagged on a green `main` with a complete CHANGELOG and getting-started build instructions.

## Stories

Authoritative list: [`../product-backlog.md`](../product-backlog.md#epic-01--core-system-inspection-v01).

| ID | Story | Pts | Sprint | Spec |
|---|---|---|---|---|
| SYS-01 | `syskit system`. | 8 | 1 | `../specs/features/system.md` |
| CPU-01 | `syskit cpu` static topology/model/cache. | 5 | 2 | `../specs/features/cpu.md` |
| CPU-02 | `syskit cpu` utilization + `--per-core`. | 8 | 2 | `../specs/features/cpu.md` |
| MEM-01 | `syskit memory`. | 8 | 2 | `../specs/features/memory.md` |
| DSK-01 | `syskit disk`. | 8 | 3 | `../specs/features/disk.md` |
| FS-01 | `syskit filesystem`. | 8 | 3 | `../specs/features/filesystem.md` |
| DOC-v01 | v0.1 user docs. | 3 | 3 | `../docs/getting-started.md` |
| REL-v01 | Release v0.1.0. | 3 | 3 | `../docs/release-process.md` |

## Per-command acceptance highlights

- **system:** host info, kernel version, OS release, uptime, load averages; correct on containers.
- **cpu:** logical/physical cores, sockets, model, flags summary; utilization only with two timestamped samples; per-core rows preserve `/proc/stat` CPU IDs.
- **memory:** physical/swap, buffers, caches, available, pressure; raw counters kept separate from derived values.
- **disk:** partition layout, filesystem usage, mount points, I/O stats.
- **filesystem:** inode usage, fs types, mount options.

Each story's full acceptance criteria live in its linked spec; the sprint file breaks them into tasks.

## Dependencies

Blocked by EPIC-00 (needs collector interface, render layer, test harness). Learning notes back each collector: `../learning/cpu.md`, `../learning/memory.md`, `../learning/disk.md`, `../learning/filesystem.md`.

## Definition of Done for the epic

All stories meet `../standards/definition-of-done.md`; every command has golden tests for table and JSON; every collector has fixtures for at least two hardware/VM variants and a Linux integration test; `v0.1.0` is tagged.

## Risks

- **R-03 fixture drift** — capture diverse fixtures (many-core, single-core VM, missing optional files).
- **R-09 golden churn** — stabilize output contracts early via FND-06 before piling on commands.
