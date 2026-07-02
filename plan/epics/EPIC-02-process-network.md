# EPIC-02 — Processes & Networking

> **Milestone:** v0.2 · **Sprints:** 4–6 · **Points:** 68
> Process inspection, network visibility, and the reusable filtering/sorting framework — culminating in **v0.2.0**.

---

## Goal

Deliver `process`, `process tree`, `network`, and `ports`, add YAML output, and introduce a filtering/sorting framework shared by all list-style commands. This epic brings SysKit's first list-heavy and Netlink-backed features online.

## Why filtering comes first

FLT-01 is scheduled before PRC-01 because process, network, and ports all need consistent `--filter`/`--sort` semantics. Building it once, first, avoids three divergent implementations and keeps the CLI predictable (`../specs/cli-conventions.md`).

## Success criteria

- List commands share one filtering/sorting contract with consistent flags and behavior.
- Netlink is integrated in the platform layer with parsing tested against fixtures and validated by Linux integration tests.
- `process tree` handles processes that appear/disappear mid-walk without crashing.
- YAML output reaches parity with JSON (same fields, same units).
- `v0.2.0` is tagged on green `main`.

## Stories

Authoritative list: [`../product-backlog.md`](../product-backlog.md#epic-02--processes--networking-v02).

| ID | Story | Pts | Sprint | Spec |
|---|---|---|---|---|
| FLT-01 | Filtering & sorting framework. | 8 | 4 | `../specs/cli-conventions.md` |
| PRC-01 | `syskit process` listing + resource usage + filters. | 13 | 4 | `../specs/features/process.md` |
| PRC-02 | `syskit process tree`. | 8 | 5 | `../specs/features/process.md` |
| NET-01 | Netlink platform integration. | 13 | 5 | `../specs/features/network.md` |
| NET-02 | `syskit network` interfaces/connections/routing. | 8 | 5 | `../specs/features/network.md` |
| PRT-01 | `syskit ports`. | 8 | 6 | `../specs/features/ports.md` |
| OUT-03 | YAML formatter. | 5 | 6 | `../specs/features/output-formats.md` |
| DOC-v02 | v0.2 docs + glossary. | 2 | 6 | `../docs/glossary.md` |
| REL-v02 | Release v0.2.0. | 3 | 6 | `../docs/release-process.md` |

## Dependencies & risk

- Blocked by EPIC-01 (render, collectors, test harness proven).
- **R-01 (Netlink)** is the epic's dominant risk. Mitigation: run spike **SPK-NET** in Sprint 4 refinement; NET-02 and PRT-01 are deferrable if NET-01 slips. Netlink is required by `../decisions/003-native-apis-over-shell.md` — no shelling out to `ss`/`netstat`.
- Learning support: `../learning/process.md`, `../learning/network.md`.

## Definition of Done for the epic

All stories meet the DoD; `/proc/[pid]` parsing and Netlink parsing have fixture + integration coverage; list output has golden tests including filtered/sorted variants; `v0.2.0` tagged.
