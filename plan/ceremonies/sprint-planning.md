# Ceremony: Sprint Planning

> **When:** Day 1 of the sprint · **Timebox:** 2 hours · **Output:** a committed sprint file in `../sprints/`

---

## Purpose

Decide **what** the team will deliver this sprint (the goal) and **how** it will start (the task breakdown). Planning turns *Ready* backlog items into a committed sprint backlog.

## Prerequisites

- The top of `../product-backlog.md` has enough *Ready* stories (they passed the [Definition of Ready](../../standards/definition-of-ready.md) at the previous refinement).
- The previous sprint's velocity is recorded in `../estimation-and-velocity.md`.
- Any deferred/carry-over stories are re-pointed if their scope changed.

## Attendees

Product Owner (required), Developers (required), Scrum Master (facilitates). Stakeholders optional.

## Agenda (2 hours)

| Time | Part | Activity |
|---|---|---|
| 0:00–0:15 | Context | PO states the milestone context and the candidate top-of-backlog stories. |
| 0:15–0:30 | Capacity | Team sets this sprint's capacity (holidays, absence) → commitment ceiling from rolling velocity. |
| 0:30–1:00 | Goal | Team agrees a single, testable **sprint goal** (one sentence). |
| 1:00–1:30 | Selection | Pull *Ready* stories until the point sum ≈ capacity-adjusted velocity. Confirm each still meets DoR. |
| 1:30–1:55 | Task breakdown | Decompose each committed story into tasks (see below). |
| 1:55–2:00 | Commit | Record the sprint file; set stories to `Committed` in the backlog. |

## The sprint goal

One sentence describing the increment's value, e.g. *"Ship `syskit cpu` with static info, utilization, and per-core output, fully tested."* If the team can't state a single goal, the sprint is trying to do too many unrelated things — narrow it.

## Task breakdown standard

For a vertical feature slice, decompose along the architecture layers so every task is small and independently reviewable:

```text
Story: CPU-02 — cpu utilization + --per-core
  □ platform: read /proc/stat samples via SysFS (+ fixtures)
  □ collector: parse counters, hold two timestamped samples
  □ service: derive per-core + aggregate utilization
  □ command: --per-core flag + validation
  □ render: table columns + JSON fields (golden)
  □ tests: unit (parse/derive/format), integration, benchmark
  □ docs: command help + CHANGELOG entry
```

Each task should be ≲ 1 day. A task nobody can size is a hidden unknown — surface it before committing.

## Commitment rules

- Commit by **velocity**, not optimism (`../estimation-and-velocity.md`).
- Respect the **WIP limit** (`../agile-framework.md`): don't commit more parallel stories than developers.
- If a High-exposure risk story is in scope (e.g. NET-01), confirm its **spike** is scheduled (`../risk-register.md`).
- Leave a little slack in release sprints for CHANGELOG/tag overhead.

## Output

A completed `../sprints/sprint-NN-*.md` from `../templates/sprint-plan.md`, with goal, committed stories (IDs + points), task breakdowns, capacity, and identified risks. Backlog statuses updated to `Committed`.
