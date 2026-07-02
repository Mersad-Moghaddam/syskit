# Template: Sprint Review

> Copy this to record the outcome of a Sprint Review (`../ceremonies/sprint-review.md`). Feeds velocity into `../estimation-and-velocity.md`.

---

```markdown
# Sprint NN — Review

**Date:** <date> · **Sprint goal:** <goal> · **Goal met?** <yes/partial/no>

## Delivered (Done + accepted)

| ID | Story | Pts | Demo notes |
|---|---|---|---|
| <ID> | <title> | <pts> | <what was shown; edge case demonstrated> |

**Completed points (velocity this sprint):** <sum of Done+accepted>

## Not completed

| ID | Story | Pts | Why | Disposition |
|---|---|---|---|---|
| <ID> | <title> | <pts> | <blocker / underestimate / scope> | <returned to backlog / carry-over, re-pointed?> |

## Commitment reliability

Committed <X> pts → completed <Y> pts = <Y/X %>.

## Stakeholder feedback (→ new backlog items)

- <feedback> → <new backlog ID or "to be created">

## Milestone check (if a release sprint)

- [ ] Milestone exit criteria in ../release-plan.md met
- [ ] CHANGELOG finalized
- [ ] Tag <vX.Y.Z> created on green main

## Follow-ups

<Items to raise at retrospective or add to the backlog.>
```

---

## Reminder

Only **Done** work (per `../../standards/definition-of-done.md`) is demonstrated and counted. Partial work scores zero points and returns to the backlog — no partial credit.
