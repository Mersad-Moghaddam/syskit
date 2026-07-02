# Template: User Story

> Copy this to create a new backlog story. Register the story's ID, points, and status in [`../product-backlog.md`](../product-backlog.md) — that file is the source of truth; this template is the working detail.

---

```markdown
# <ID> — <short title>

**Epic:** <EPIC-NN> · **Points:** <1|2|3|5|8|13> · **Target sprint:** <N> · **Status:** <Backlog|Ready|Committed|In Progress|In Review|Done>
**Spec:** <../specs/features/xxx.md or other canonical spec>

## User story

As a <role>, I want <capability> so that <benefit>.

## Motivation

<Why this matters. Link to the spec's motivation rather than restating it.>

## Scope

- In: <what this story delivers>
- Out: <explicit non-goals for this story>

## CLI & flags

<Expected command(s) and flags, per the spec / cli-conventions.md.>

## Expected output

- Table: <sketch or link>
- Structured (JSON/YAML): <field names + units>

## Linux data sources

<e.g. /proc/stat, /sys/devices/system/cpu/, Netlink — per the spec.>

## Edge cases

- <e.g. missing cpufreq → unavailable, not zero>
- <e.g. process disappears mid-walk>

## Acceptance criteria (testable)

- [ ] <criterion 1>
- [ ] <criterion 2>
- [ ] <criterion 3>

## Fixtures required

<Which fixture sets, capturing which variation (many-core vs VM, cgroup v1 vs v2, missing files).>

## Dependencies

<Blocking story IDs / ADRs / new libraries needing an ADR.>

## Definition of Ready

Confirm all boxes in ../standards/definition-of-ready.md before this leaves refinement.

## Definition of Done

Confirm all boxes in ../standards/definition-of-done.md before merge.
```

---

## Notes

- **Vertical slice:** a feature story delivers CLI → command → service → collector → platform → render, with tests — not a single layer.
- **Testable acceptance:** each criterion must be checkable by a unit, integration, or golden test. "Works well" is not a criterion.
- **Points, not hours:** size relative to the reference stories in `../estimation-and-velocity.md`.
- Anything that would be ≥ 13 with hidden unknowns is split at refinement before becoming Ready.
