# Ceremony: Sprint Retrospective

> **When:** Day 10 (after the review) · **Timebox:** 45 minutes · **Output:** 1–3 concrete, owned improvements

---

## Purpose

Inspect **how the team works** and commit to a small number of concrete improvements. The retrospective is about process and collaboration, not the product. Its output is action, not venting.

## Attendees

Developers, Scrum Master (facilitates), Product Owner. Psychological safety is the precondition — this is blameless.

## Format (45 minutes)

| Time | Activity |
|---|---|
| 0:00–0:05 | Set the tone: blameless, focused on the system not individuals. |
| 0:05–0:20 | Gather data — each person notes: **Keep** / **Drop** / **Try**. |
| 0:20–0:35 | Group themes; discuss the highest-signal ones. |
| 0:35–0:45 | Decide **1–3** improvements, each with an **owner** and a way to tell if it worked. |

## Prompts tuned for SysKit

- Did any story miss the [Definition of Done](../../standards/definition-of-done.md)? Why — estimate, unknown, or review bottleneck?
- Were estimates accurate? Which story surprised us, and does a **reference story** in `../estimation-and-velocity.md` need updating?
- Did fixtures/integration tests catch (or miss) anything real? (Feeds R-03.)
- Did the WIP limit help or get in the way?
- Did any risk in `../risk-register.md` materialize? Re-score it.
- Did async standups keep us synced, or did a blocker sit too long?

## Rules for good actions

- **Few and concrete.** Three real changes beat ten aspirations.
- **Owned.** Every action has a name attached.
- **Checkable.** "Reduce cycle time" is a wish; "split any story >8 into vertical slices at refinement" is an action.
- **Revisited.** Start the next retro by reviewing last time's actions — did they work?

## Feedback loops this ceremony drives

- Estimates/reference stories → `../estimation-and-velocity.md`
- Risk scores → `../risk-register.md`
- Scope/order re-plans → `../product-backlog.md` and `../release-plan.md`
- Process rule changes → `../agile-framework.md`

## Output

A short retro record from `../templates/sprint-retrospective.md`: what we keep, what we drop, and the 1–3 owned actions for next sprint.
