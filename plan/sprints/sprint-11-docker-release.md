# Sprint 11 — Docker Commands → v0.4.0

**Dates:** TBD → +2 weeks · **Milestone:** v0.4 — **release sprint** · **Committed points:** 21

## Sprint goal

Ship `syskit docker` and `syskit docker inspect <id>`, complete container docs, and tag **v0.4.0**.

## Capacity

- Nominal, with release overhead absorbed.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| DOC-01 | `syskit docker` listing/usage/status. | 8 | Committed |
| DOC-02 | `syskit docker inspect <id>`. | 8 | Committed |
| DOC-v04 | v0.4 docs + container learning notes. | 2 | Committed |
| REL-v04 | Release v0.4.0. | 3 | Committed |

## Task breakdowns

**DOC-01 / DOC-02** — spec `../../specs/features/containers.md`
- [ ] platform/service: query the container runtime; combine with CG-01 cgroup model for resource usage.
- [ ] `docker` list: id, name, status, resource usage (reuse FLT-01 filters).
- [ ] `docker inspect <id>`: detailed view; graceful "not found".
- [ ] **Docker-absent behavior**: clear message + non-zero exit, never a crash (spec edge case).
- [ ] render: table + JSON (+ golden); integration where a runtime is available, else fixture-backed.
- [ ] new dependency (Docker client) requires an **ADR** per `../../standards/dependency-policy.md`.
- [ ] docs + CHANGELOG.

**REL-v04** — release
- [ ] Confirm v0.4 milestone exit criteria (`../release-plan.md`): cgroup v1+v2 covered; graceful Docker-absent behavior.
- [ ] CHANGELOG under v0.4.0; tag on green `main`.

## Definition of Ready / Done

Standard gates + v0.4 milestone gate.

## Risks this sprint

- **R-10 (dependency)** — a Docker API client is a new dependency; the ADR must be merged before DOC-01 code. Prefer the runtime's own interfaces where practical (`../../decisions/003-native-apis-over-shell.md`).
- **R-07 (kernel/runtime variants)** — CI may lack a Docker daemon; ensure fixture-backed tests plus a documented manual verification path.

## Dependencies

Blocked by Sprint 10 (CG-01 model, CNT-01 mapping).

## Notes

Container support is breadth, not new architecture — it reuses the process and cgroup layers. Keep the commands thin.
