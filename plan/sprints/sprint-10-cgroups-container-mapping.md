# Sprint 10 — cgroups & Container Mapping

**Dates:** TBD → +2 weeks · **Milestone:** v0.4 (Containers) · **Committed points:** 21

## Sprint goal

Add cgroup v1/v2 parsing to the platform layer and map containers to their processes, laying the groundwork for the Docker commands.

## Capacity

- Nominal. CG-01 (13) dominates and carries the v1/v2 divergence risk.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| CG-01 | cgroup v1/v2 parsing in platform layer. | 13 | Committed |
| CNT-01 | Container-to-process mapping + container-aware process views. | 8 | Committed |

## Task breakdowns

**CG-01** — spec `../../specs/features/containers.md`
- [ ] platform: read cgroup v1 (`/sys/fs/cgroup/<controller>`) and v2 (unified) via `SysFS`.
- [ ] detect hierarchy version; normalize both into one model upward.
- [ ] fixtures: v1 host and v2 host captures; unit tests for both parse paths.
- [ ] integration: read real cgroup on Linux CI; benchmark.

**CNT-01** — spec `../../specs/features/containers.md`
- [ ] service: map cgroup/container IDs to owning PIDs; container-aware filter for the process view.
- [ ] extend process command/service with a container dimension (reuse PRC-01 + FLT-01).
- [ ] render: container column + JSON field (+ golden); tests; docs + CHANGELOG.

## Definition of Ready / Done

Standard gates. Acceptance: v1 and v2 both parse into the same normalized model; container mapping reuses existing process collection.

## Risks this sprint

- **R-08 (cgroup v1/v2 divergence) — MEDIUM.** Mitigation: contain the difference entirely in the platform layer; provide fixtures for both hierarchies; everything above the platform sees one model.
- **R-07 (kernel variants)** — CI may only expose one cgroup version; fixtures cover the other.

## Dependencies

Blocked by EPIC-02 (process/resource collectors that container views extend).

## Notes

Getting the normalized cgroup model right here keeps the Docker commands (Sprint 11) thin — they consume the model, not raw cgroup files.
