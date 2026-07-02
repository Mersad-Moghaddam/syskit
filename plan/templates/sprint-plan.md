# Template: Sprint Plan

> Copy this to `../sprints/sprint-NN-<theme>.md` at Sprint Planning. Fill it during the ceremony (`../ceremonies/sprint-planning.md`).

---

```markdown
# Sprint NN — <theme>

**Dates:** <start> → <end> (2 weeks) · **Milestone:** <vX.Y> · **Committed points:** <sum>

## Sprint goal

<One sentence: the demonstrable value this sprint delivers.>

## Capacity

- Developers available: <n>
- Known reductions: <holidays / absence / conferences>
- Commitment ceiling (from rolling velocity): <pts>

## Committed backlog

| ID | Story | Pts | Owner | Status |
|---|---|---|---|---|
| <ID> | <title> | <pts> | <who> | Committed |

_All IDs and points are authoritative in ../product-backlog.md._

## Task breakdowns

<One task list per story, from ../templates/task-breakdown.md.>

## Definition of Ready (entry gate)

Every committed story passed ../standards/definition-of-ready.md at refinement.

## Definition of Done (exit gate)

No story counts until it meets ../standards/definition-of-done.md.

## Risks this sprint

<Reference risk IDs from ../risk-register.md that this sprint touches; note mitigations/spikes scheduled.>

## Dependencies

<Blocking stories/ADRs; note anything that must land before mid-sprint.>

## Notes

<Anything the team should remember: deferrable stories, release overhead, spikes.>
```

---

## Reminders

- Commit by **velocity**, not hope. Respect the **WIP limit**.
- If a High-exposure risk story is in scope, its spike must be scheduled here.
- In release sprints, leave slack for CHANGELOG, tag, and milestone exit-criteria checks (`../release-plan.md`).
