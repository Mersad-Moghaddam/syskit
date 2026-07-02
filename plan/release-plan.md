# Release Plan

> How the roadmap milestones (`v0.1`…`v1.0`) map onto time-boxed sprints and release tags. Aligns [`../specs/roadmap.md`](../specs/roadmap.md) with [`product-backlog.md`](product-backlog.md).

---

## Approach

- **Milestones are outcomes; sprints are time-boxes.** A milestone releases when all its stories are Done — which may be mid-sprint or at a sprint boundary. We tag at the boundary for a clean, releasable `main`.
- **Sprint length:** 2 weeks. **Planned velocity:** ~23 pts/sprint after warm-up (see [`estimation-and-velocity.md`](estimation-and-velocity.md)).
- **Tags** are annotated SemVer on `main`, created only after CI is green (`../standards/branch-strategy.md`).

---

## Timeline overview

```text
Sprint:   0    1    2    3    4    5    6    7    8    9   10   11   12   13   14
         └─foundation─┘
              └────── v0.1 core ──────┘
                                └──── v0.2 proc/net ────┘
                                                  └──── v0.3 realtime ────┘
                                                                    └ v0.4 ┘
                                                                          └ v0.5 ┘
                                                                                └v1.0┘
Tag:                    v0.1.0            v0.2.0          v0.3.0     v0.4.0  v0.5.0  v1.0.0
```

---

## Sprint-by-sprint schedule

| Sprint | Theme | Epics | Committed stories | Pts | Release |
|---|---|---|---|---|---|
| **0** | Implementation transition | EPIC-00 | FND-01, FND-02, FND-03, FND-04, FND-10a | ~21 | — |
| **1** | Foundation complete + first command | EPIC-00, EPIC-01 | FND-05, FND-06, FND-07, FND-08, FND-09, FND-11, FND-10b, SYS-01 | ~24 | — |
| **2** | CPU + memory | EPIC-01 | CPU-01, CPU-02, MEM-01 | 21 | — |
| **3** | Disk + filesystem → **v0.1** | EPIC-01 | DSK-01, FS-01, DOC-v01, REL-v01 | 22 | **v0.1.0** |
| **4** | Filtering + processes | EPIC-02 | FLT-01, PRC-01 | 21 | — |
| **5** | Process tree + network | EPIC-02 | PRC-02, NET-01, NET-02 | 29* | — |
| **6** | Ports + YAML → **v0.2** | EPIC-02 | PRT-01, OUT-03, DOC-v02, REL-v02 | 18 | **v0.2.0** |
| **7** | TUI foundation + watch | EPIC-03 | TUI-01, RT-01, WCH-01 | 21 | — |
| **8** | Dashboard | EPIC-03 | DSH-01 | 13 | — |
| **9** | Top → **v0.3** | EPIC-03 | TOP-01, REL-v03 | 16 | **v0.3.0** |
| **10** | cgroups + mapping | EPIC-04 | CG-01, CNT-01 | 21 | — |
| **11** | Docker → **v0.4** | EPIC-04 | DOC-01, DOC-02, DOC-v04, REL-v04 | 21 | **v0.4.0** |
| **12** | Plugin core | EPIC-05 | PLG-01, PLG-02, PLG-03 | 26* | — |
| **13** | Plugin security + SDK → **v0.5** | EPIC-05 | PLG-04, PLG-05, PLG-06, REL-v05 | 19 | **v0.5.0** |
| **14** | Diagnostics + stabilize → **v1.0** | EPIC-06 | DIA-01, PERF-01, API-01, PKG-01, REL-01, DOCS-01, REL-v10 | 39* | **v1.0.0** |

\* Sprints marked with `*` exceed nominal velocity and carry **explicit overflow risk** — see notes below. They are scheduled tight on purpose so the risk is visible at planning, not discovered mid-sprint.

---

## Overflow-risk sprints (planned mitigations)

| Sprint | Load | Why | Mitigation |
|---|---|---|---|
| 5 | 29 pts | Netlink (NET-01, 13) is high-uncertainty. | If NET-01 slips, defer NET-02 to Sprint 6 and pull PRT-01 forward; v0.2 tag moves with it. |
| 12 | 26 pts | Plugin protocol (PLG-01, 13) is foundational and risky. | PLG-03 is the first to defer to Sprint 13. |
| 14 | 39 pts | Stabilization bundles many small independent items. | Split into a 14a/14b if velocity data says so at Sprint 13 refinement; DIA-01 is the only large item and can lead a dedicated sprint. |

**General rule:** the release tag follows the work. We ship a milestone when its stories are Done, never by forcing scope into a fixed date.

---

## Milestone exit criteria

Each milestone is released only when, in addition to its stories being Done:

| Milestone | Extra gate |
|---|---|
| **v0.1.0** | Vertical slice proven end-to-end; table+JSON stable; CI runs unit+integration+golden on Linux; getting-started build works from clean clone. |
| **v0.2.0** | Filtering/sorting consistent across list commands; YAML parity with JSON; Netlink integration tests green. |
| **v0.3.0** | TUI runs under `-race`; refresh pipeline has no leaks; `watch`/`dashboard`/`top` documented with keybindings. |
| **v0.4.0** | cgroup v1 and v2 both covered by fixtures; graceful behavior when Docker absent. |
| **v0.5.0** | Plugin isolation model reviewed against `../SECURITY.md`; example plugin builds against published SDK. |
| **v1.0.0** | CLI/output contracts frozen (API-01); packages install on target distros; release automation dry-run passes; full docs pass. See EPIC-06. |

---

## What moves the plan

Re-plan at these triggers, and only these:

1. **Velocity drift** > 20% over two sprints → re-forecast remaining sprints.
2. **A high-uncertainty story (13 pts) is discovered to be larger** → split it, re-point, and re-order at refinement.
3. **A dependency (ADR) changes** → update the affected epic and, if needed, resequence.

Every re-plan is recorded in the relevant sprint retrospective and reflected back into `product-backlog.md`.

---

## Post-1.0 (out of scope for this plan)

The roadmap's "Future Considerations" (remote monitoring, historical data, alerting, hardware info, Kubernetes) are **not** scheduled here. After v1.0 they enter the backlog and get their own release plan. This plan's job ends at a stable, documented, packaged v1.0.
