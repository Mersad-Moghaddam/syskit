# Agile Framework for SysKit

> How SysKit runs Scrum: the roles, the cadence, the events, and the artifacts — adapted to an open-source, documentation-first, single-to-small-team project.

---

## Why Scrum (with Kanban limits)

SysKit is built by a small team (and community contributors) against a well-specified backlog. Scrum gives us a predictable rhythm for turning specs into shipped commands: a fixed cadence, a committed goal per sprint, and regular inspection. We borrow one idea from Kanban — a **work-in-progress (WIP) limit** — because a small team burns itself out by starting more than it finishes.

The framework is deliberately lightweight. Ceremonies exist to serve delivery, not the other way around. If a ceremony stops adding value, we change it at the retrospective.

---

## Roles

Open-source projects rarely have a full, dedicated Scrum team. We map Scrum roles onto realistic hats. One person may wear several.

| Scrum role | SysKit hat | Responsibility |
|---|---|---|
| **Product Owner** | Maintainer / project lead | Owns the backlog order, accepts stories against acceptance criteria, defends scope and non-goals, decides release timing. |
| **Scrum Master** | Facilitator (rotating) | Runs ceremonies, removes blockers, guards the process and WIP limit, keeps metrics honest. |
| **Developers** | Contributors (core + community) | Refine, estimate, commit to, and deliver stories to the Definition of Done. |
| **Stakeholders** | Users, downstream engineers | Give feedback at reviews and through issues; do not direct sprints mid-flight. |

**Decision rule:** the Product Owner orders the backlog; the Developers decide how much they can pull into a sprint. Neither overrides the other.

---

## Cadence

| Element | Choice | Rationale |
|---|---|---|
| Sprint length | **2 weeks** | Long enough to deliver a vertical feature slice with full tests; short enough to correct course. |
| Sprint start | Planning on day 1 | |
| Sprint end | Review + retrospective on day 10 | |
| Refinement | Mid-sprint (day 6), ~1 hr | Keeps the top of the backlog always ~1.5 sprints deep and Ready. |
| Release cadence | At milestone boundaries, not every sprint | A milestone (`v0.1`…) may span 2–3 sprints; we tag when the milestone's stories are Done. |

---

## The Scrum events

Each event has a dedicated guide in [`ceremonies/`](ceremonies/). Summary:

| Event | When | Timebox | Output |
|---|---|---|---|
| **Sprint Planning** | Day 1 | 2 hrs | Sprint goal + committed backlog + task breakdown → `sprints/sprint-NN-*.md` |
| **Daily Standup** | Every working day | 15 min | Blockers surfaced; async-friendly for OSS (written check-in acceptable) |
| **Backlog Refinement** | Day 6 | 1 hr | Next sprint's candidate stories estimated and made *Ready* |
| **Sprint Review** | Day 10 | 1 hr | Demo of Done work; stakeholder feedback; backlog updated |
| **Sprint Retrospective** | Day 10 | 45 min | 1–3 concrete process improvements for next sprint |

---

## Artifacts

| Artifact | Where it lives | Owner |
|---|---|---|
| **Product Backlog** | [`product-backlog.md`](product-backlog.md) | Product Owner |
| **Sprint Backlog** | `sprints/sprint-NN-*.md` | Developers |
| **Increment** | A merged, releasable state of `main` | Whole team |
| **Definition of Ready** | [`../standards/definition-of-ready.md`](../standards/definition-of-ready.md) | Product Owner |
| **Definition of Done** | [`../standards/definition-of-done.md`](../standards/definition-of-done.md) | Whole team |
| **Impediment log** | Standup notes + [`risk-register.md`](risk-register.md) | Scrum Master |

---

## The two gates

Every story passes through two gates. They are non-negotiable and defined outside this plan so they apply to code and docs alike:

```text
        Definition of Ready                    Definition of Done
        (can we start?)                        (can we merge?)
   ┌──────────────────────┐              ┌──────────────────────────┐
   │ spec exists+reviewed  │   sprint     │ spec satisfied            │
   │ CLI + flags defined   │   ───────▶   │ unit+integration+golden   │
   │ data sources known    │              │ bench on hot paths        │
   │ edge cases listed     │              │ vet/fmt/vuln clean        │
   │ acceptance testable   │              │ docs + CHANGELOG updated  │
   │ fixtures identified   │              │ reviewed + CI green       │
   └──────────────────────┘              └──────────────────────────┘
```

- A story that fails the DoR stays in refinement — it does **not** enter a sprint.
- A story that fails the DoD stays in the sprint (or returns to the backlog) — it is **not** counted as delivered.

---

## Work-in-progress limit

To protect a small team from thrash:

- **WIP limit = number of active developers.** A developer finishes (to Done) or explicitly blocks their current story before pulling another.
- A blocked story is moved to a `Blocked` state with the blocker recorded, freeing the developer — but the blocker becomes the Scrum Master's priority.

---

## Sprint states of a story

```text
  Backlog ──refine──▶ Ready ──plan──▶ Committed ──start──▶ In Progress
                                                              │
                                        ┌─────────────────────┤
                                        ▼                     ▼
                                     Blocked              In Review
                                        │                     │
                                        └────────┬────────────┘
                                                 ▼
                                               Done  ──(accepted at review)──▶ Released
```

---

## How estimation works here

We estimate in **story points** using a modified Fibonacci scale and planning poker. The full method — scale meaning, reference stories, capacity, and velocity — is in [`estimation-and-velocity.md`](estimation-and-velocity.md). Planning uses **velocity**, not hours, to decide commitment.

---

## Adapting for open source

Because contributors are part-time and distributed:

- **Standups are asynchronous-friendly.** A written daily update in the sprint channel satisfies the standup.
- **Commitment is a forecast, not a promise.** Community capacity fluctuates; velocity smooths this over time.
- **Every story is independently mergeable.** No story depends on unmerged work from another in-flight story within the same sprint where avoidable (see dependency notes in the backlog).
- **The Product Owner keeps 1.5 sprints of Ready stories** so a contributor with spare time is never blocked for lack of a defined task.

---

*The framework serves the goal: a stable, well-tested SysKit v1.0. When a rule stops serving that goal, change it — and record why in a retrospective.*
