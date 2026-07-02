# EPIC-04 — Containers

> **Milestone:** v0.4 · **Sprints:** 10–11 · **Points:** 42
> Container runtime inspection and cgroup-based resource views — culminating in **v0.4.0**.

---

## Goal

Add cgroup v1/v2 parsing to the platform layer, map containers to their processes, and deliver `docker` listing and `docker inspect`. Container-aware views build on the existing process and resource collectors rather than replacing them.

## Success criteria

- cgroup **v1 and v2** are both parsed, isolated inside the platform layer, and normalized into one model upward (`../specs/features/containers.md`).
- Containers map to their owning processes; process views can be filtered by container.
- `docker` and `docker inspect` degrade gracefully when Docker is absent or the socket is unreachable — a clear message, not a crash.
- `v0.4.0` is tagged on green `main`.

## Stories

Authoritative list: [`../product-backlog.md`](../product-backlog.md#epic-04--containers-v04).

| ID | Story | Pts | Sprint | Spec |
|---|---|---|---|---|
| CG-01 | cgroup v1/v2 parsing in platform layer. | 13 | 10 | `../specs/features/containers.md` |
| CNT-01 | Container-to-process mapping + container-aware process views. | 8 | 10 | `../specs/features/containers.md` |
| DOC-01 | `syskit docker` listing/usage/status. | 8 | 11 | `../specs/features/containers.md` |
| DOC-02 | `syskit docker inspect <id>`. | 8 | 11 | `../specs/features/containers.md` |
| DOC-v04 | v0.4 docs + container learning notes. | 2 | 11 | `../learning/roadmap.md` |
| REL-v04 | Release v0.4.0. | 3 | 11 | `../docs/release-process.md` |

## Dependencies & risk

- Blocked by EPIC-02 (process/resource collectors that container views extend).
- **R-08 (cgroup v1/v2 divergence)** — mitigated by containing the difference in the platform layer and providing fixtures for both hierarchies.
- **R-10 (dependency)** — a Docker API client is a new dependency and requires an ADR per `../standards/dependency-policy.md`; prefer reading the runtime's own interfaces where practical, consistent with `../decisions/003-native-apis-over-shell.md`.

## Definition of Done for the epic

All stories meet the DoD; cgroup parsing has fixtures for both v1 and v2 plus integration coverage; Docker-absent behavior is tested; `v0.4.0` tagged.
