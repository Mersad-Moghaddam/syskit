# SysKit Delivery Plan

> The execution plan that takes SysKit from **0 (design/spec phase)** to **100 (stable v1.0 release)** using Scrum and Agile principles.

---

## Purpose

The `specs/`, `standards/`, and `decisions/` directories describe **what** SysKit is and **how** it must be built. This `plan/` directory describes **when and in what order** it gets built, and **how the team runs** while building it.

It is the bridge between the design-first documentation and a shipped tool:

- The **product backlog** decomposes the roadmap into epics, user stories, and story points.
- The **release plan** maps the roadmap milestones (`v0.1`…`v1.0`) onto time-boxed sprints.
- The **sprints** each carry a goal, a committed set of backlog items, and a task breakdown.
- The **ceremonies** and **templates** define how the team plans, tracks, reviews, and improves.

This plan does not restate specs. Where a story needs behavioral detail, it links to the canonical spec in `../specs/`. Where it needs a completion bar, it links to `../standards/`.

---

## How this plan is organized

```text
plan/
├── README.md                     # this file — index and reading order
├── agile-framework.md            # Scrum model: roles, cadence, ceremonies, artifacts
├── product-backlog.md            # SOURCE OF TRUTH: epics, stories, IDs, points, status
├── release-plan.md               # roadmap milestones → sprints → release tags
├── estimation-and-velocity.md    # points scale, planning poker, capacity, velocity
├── risk-register.md              # risks, likelihood/impact, mitigation, owners
├── metrics-and-reporting.md      # burndown, velocity, cycle time, dashboards
├── epics/                        # one file per epic, with its stories and acceptance
│   ├── EPIC-00-foundation.md
│   ├── EPIC-01-core-inspection.md
│   ├── EPIC-02-process-network.md
│   ├── EPIC-03-realtime-monitoring.md
│   ├── EPIC-04-containers.md
│   ├── EPIC-05-extensibility.md
│   ├── EPIC-06-stabilization-release.md
│   └── EPIC-07-cross-cutting-quality.md
├── sprints/                      # sprint-00 … sprint-14, each a committed plan
│   └── sprint-NN-*.md
├── ceremonies/                   # how each Scrum event is run for SysKit
│   ├── sprint-planning.md
│   ├── daily-standup.md
���   ├── backlog-refinement.md
│   ├── sprint-review.md
│   └── sprint-retrospective.md
└── templates/                    # copy-to-start artifacts
    ├── user-story.md
    ├── task-breakdown.md
    ├── sprint-plan.md
    ├── sprint-review.md
    └── sprint-retrospective.md
```

---

## Reading order

1. **New to the plan?** Read `agile-framework.md`, then `product-backlog.md`, then `release-plan.md`.
2. **Starting a sprint?** Open the sprint file in `sprints/`, then `ceremonies/sprint-planning.md`.
3. **Writing a new story?** Copy `templates/user-story.md` and register the ID in `product-backlog.md`.
4. **Reporting progress?** Use `metrics-and-reporting.md`.

---

## Guiding principles

SysKit's plan inherits the project's own posture (see `../specs/constitution.md`):

- **Documentation first, then code.** A story is not *Ready* until its spec exists and its acceptance criteria are testable (`../standards/definition-of-ready.md`).
- **Vertical slices over horizontal layers.** Each feature story is delivered end-to-end (CLI → command → service → collector → platform) so every sprint produces something demonstrable.
- **`main` is always releasable.** Trunk-based flow, squash-merge, linear history (`../standards/branch-strategy.md`).
- **Done means done.** No story is closed until it satisfies `../standards/definition-of-done.md`.
- **The plan is a living document.** Scope, order, and estimates change at refinement and retrospective — that is expected, not a failure.

---

## Current status snapshot

| Field | Value |
|---|---|
| Project phase | Design/Spec complete → entering implementation |
| First sprint | `sprint-00` (implementation transition) |
| Target of the plan | `v1.0.0` — stable release |
| Sprint length | 2 weeks |
| Planned sprints | 15 (Sprint 0 through Sprint 14) |
| Estimated horizon | ~30 weeks / ~7 months at planned velocity |

Live status of every epic and story is tracked in [`product-backlog.md`](product-backlog.md). Per-milestone timing is in [`release-plan.md`](release-plan.md).
