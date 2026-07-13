# Sprint 11 — Container Commands → v0.4.0

**Dates:** TBD → +2 weeks · **Milestone:** v0.4 — **release sprint** · **Committed points:** 21

## Sprint goal

Ship `syskit containers` and `syskit containers inspect <id>`, complete container docs, and tag **v0.4.0**.

## Capacity

- Nominal, with release overhead absorbed.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| DOC-01 | `syskit containers` listing/usage/status. | 8 | Committed |
| DOC-02 | `syskit containers inspect <id>`. | 8 | Committed |
| DOC-v04 | v0.4 docs + container learning notes. | 2 | Committed |
| REL-v04 | Release v0.4.0. | 3 | Committed |

## Task breakdowns

**DOC-01 / DOC-02** — spec `../../specs/features/containers.md`
- [ ] platform/service: read cgroup resource data and optionally enrich with runtime metadata.
- [ ] `containers` list: id, optional runtime hint, resource usage (reuse FLT-01 filters).
- [ ] `containers inspect <id>`: detailed view; graceful "not found".
- [ ] **Runtime-absent behavior**: cgroup-derived output remains available without a runtime socket.
- [ ] render: table + JSON (+ golden); integration where a runtime is available, else fixture-backed.
- [ ] any optional runtime-client dependency requires an **ADR** per `../../standards/dependency-policy.md`.
- [ ] docs + CHANGELOG.

**REL-v04** — release
- [ ] Confirm v0.4 milestone exit criteria (`../release-plan.md`): cgroup v1+v2 covered; graceful runtime-absent behavior.
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
