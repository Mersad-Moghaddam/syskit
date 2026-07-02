# Ceremony: Sprint Review

> **When:** Day 10 · **Timebox:** 1 hour · **Output:** accepted increment, updated backlog, feedback captured

---

## Purpose

Inspect the **increment**: demonstrate what is *Done*, gather stakeholder feedback, and adapt the backlog. The review is about working software, not slides.

## Attendees

Product Owner, Developers, Scrum Master, and any stakeholders/users who want to see progress.

## The one rule

**Only Done work is demonstrated.** A story that does not meet `../../standards/definition-of-done.md` is not shown as complete and does not count toward velocity — it returns to the backlog. Half-working demos erode trust and hide the real state.

## Agenda (1 hour)

| Time | Activity |
|---|---|
| 0:00–0:05 | Scrum Master restates the sprint goal and what was committed. |
| 0:05–0:35 | Developers **demo** each Done story against its acceptance criteria — run the actual command / show the actual output. |
| 0:35–0:50 | PO accepts or rejects each story against its spec's acceptance criteria. |
| 0:50–1:00 | Stakeholder feedback captured as backlog items (not mid-sprint changes). |

## Demo standard for SysKit

Because SysKit is a CLI, demos are concrete and reproducible:

- Run the real command (`syskit cpu --per-core`, `--format json`) on a Linux host or fixture-backed test.
- Show table **and** structured output where the story delivered both.
- Show at least one edge case handled (missing data as *unavailable*, container behavior, etc.).
- For collectors: note fixtures used and that integration tests pass on Linux CI.

## Acceptance

The Product Owner accepts a story only if its spec's acceptance criteria are demonstrably met. Acceptance here + Done = the story counts toward velocity and moves to `Released` when the milestone tags.

## Output

- Backlog updated: accepted stories closed; rejected/partial ones returned with notes.
- New feedback captured as `Backlog` items (ordered later by the PO).
- A short review summary from `../templates/sprint-review.md`, including velocity for `../estimation-and-velocity.md`.
- At milestone boundaries, confirm the milestone exit criteria in `../release-plan.md` before tagging.
