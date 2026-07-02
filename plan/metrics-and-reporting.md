# Metrics & Reporting

> The small set of metrics SysKit tracks to keep delivery honest, and how each is produced and read. Metrics inform conversation; they never replace judgment.

---

## Principles

- **Few metrics, well understood.** We track flow and quality, not activity.
- **Metrics are signals, not targets.** A gamed number is worse than no number (Goodhart's law). Coverage %, in particular, is a signal per `../specs/testing-strategy.md`.
- **Every metric has an owner and a home.** No orphan dashboards.

---

## Core metrics

| Metric | Question it answers | Source | Owner | Cadence |
|---|---|---|---|---|
| **Velocity** | How much *Done* work per sprint? | `estimation-and-velocity.md` table | Scrum Master | Per sprint |
| **Sprint burndown** | Are we on track *within* a sprint? | Sprint board (points remaining vs day) | Scrum Master | Daily |
| **Release burnup** | How close to the next milestone? | Backlog Done-points vs milestone scope | Product Owner | Per sprint |
| **Cycle time** | How long from *In Progress* → *Done*? | PR opened → squash-merged | Scrum Master | Per sprint |
| **Commitment reliability** | Do we finish what we commit? | Completed ÷ Committed points | Scrum Master | Per sprint |
| **Escaped defects** | Bugs found after a story is Done. | Issues labeled `bug` post-merge | Maintainer | Continuous |
| **CI health** | Is `main` trustworthy? | CI pass rate, time-to-green | Maintainer | Continuous |
| **Test coverage** | Is logic meaningfully exercised? | `go test -coverprofile` in CI | Developers | Per PR |

---

## Sprint burndown (within-sprint)

Track **story points remaining** each working day against an ideal line from committed total to zero.

```text
pts
 24 ●╲  ideal
    │ ╲╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
 18 │  ●───●            actual (flat early = work not closing)
    │        ╲
 12 │         ●──●
    │             ╲
  6 │              ●─●
    │                 ╲
  0 └────────────────────●──▶ day
    1  2  3  4  5  6  7  8  9 10
```

**Reading it:**
- Flat line early is normal (stories close near the end, not linearly).
- Flat line *late* (day 7+) → stories aren't reaching Done; raise at standup.
- A step down = a story hit Done (points only drop when Done, never partially).

---

## Release burnup (toward a milestone)

Plot cumulative **Done points within a milestone** against that milestone's total scope. Unlike burndown, a burnup also shows **scope changes** (the top line moves when stories are added/removed) — essential because SysKit re-plans at defined triggers.

```text
pts
 51 ┤························ scope (v0.1 = 51)
    │                  ╭─────
    │             ╭────╯  ← done
    │        ╭────╯
    │   ╭────╯
  0 ┼───╯
    S1   S2   S3
```

---

## Cycle time

Median time from a story's first commit/PR to its squash-merge. Short branches are a project standard (`../standards/branch-strategy.md`); rising cycle time is an early warning of oversized stories or review bottlenecks. Target: most stories merge within the sprint they start.

---

## Quality metrics (tied to Definition of Done)

Because DoD is strict, these mostly stay green — surfacing them makes regressions visible:

| Metric | Green condition |
|---|---|
| Unit + race suite | Passes on every PR. |
| Integration suite | Passes on Linux CI on every PR touching a collector. |
| Golden files | No unexplained diffs; intentional diffs reviewed. |
| Benchmarks | No unexplained regression on hot paths vs previous commit. |
| `govulncheck` | Clean; new deps have an ADR. |
| Coverage | Meaningfully exercised; guideline ≥ 80% overall, parsers/logic approaching complete. |

---

## Reporting rhythm

| Report | Audience | Contents | When |
|---|---|---|---|
| **Daily standup note** | Team | Yesterday / today / blockers; burndown glance. | Daily |
| **Sprint review summary** | Stakeholders | Demo notes, stories Done, velocity, next goal. | End of sprint (`templates/sprint-review.md`) |
| **Retro actions** | Team | 1–3 improvements + owners. | End of sprint (`templates/sprint-retrospective.md`) |
| **Milestone report** | Community | What shipped in `vX.Y`, CHANGELOG, known gaps. | On release tag |

---

## Anti-metrics (deliberately not tracked)

- **Lines of code / commit count** — measures typing, not value.
- **Individual velocity** — velocity is a team signal; per-person tracking harms collaboration.
- **Hours worked** — we plan by points and outcomes, not time-in-seat.

---

## Where the numbers live

- **Velocity + commitment reliability:** the table in `estimation-and-velocity.md`.
- **Burndown/burnup:** the sprint board tool (GitHub Projects) — this plan defines the method; the tool renders it.
- **Quality metrics:** CI output (`../.github/workflows/ci.yml`).
- **Risks influenced by metrics:** cross-referenced in `risk-register.md` (e.g., velocity drift → R-04, R-11).
