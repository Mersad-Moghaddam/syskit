# Estimation & Velocity

> How SysKit sizes work, forecasts capacity, and measures throughput. Estimation is in **story points**, planning is by **velocity**, never raw hours.

---

## Why story points

Points measure *relative effort, complexity, and uncertainty* together — not clock time. This suits SysKit because:

- Contributor availability varies (open source), but relative sizing is stable.
- Kernel-interface work carries real uncertainty; points capture "we're not sure" that hours pretend away.
- Velocity (points completed per sprint) self-corrects for a team's true pace over time.

---

## The scale (modified Fibonacci)

| Points | Meaning | Typical SysKit example |
|---|---|---|
| **1** | Trivial, no unknowns, <½ day. | Add a CHANGELOG entry; wire one flag. |
| **2** | Small, well-understood. | YAML formatter given the JSON one exists (OUT-03). |
| **3** | Straightforward feature slice, minor edges. | Error-handling patterns (FND-07); a release story. |
| **5** | Standard feature slice with tests + fixtures. | `syskit cpu` static info (CPU-01); config loading (FND-09). |
| **8** | Substantial: multiple data sources or derived metrics. | `syskit memory` (MEM-01); render skeleton (FND-06); dashboard-less collectors. |
| **13** | Large/uncertain: new subsystem or protocol. | Netlink integration (NET-01); plugin protocol (PLG-01); `top` (TOP-01). |
| **(21)** | Too big — **must be split** before it enters a sprint. | Any "implement network support" lump. |

**Rule:** nothing larger than 13 enters a sprint. A 21 is a refinement failure and must be decomposed.

---

## Reference stories (calibration anchors)

Use these to keep estimates consistent across sprints. When unsure, compare to the anchor.

| Anchor | Points | Why it's the anchor |
|---|---|---|
| OUT-03 (YAML formatter) | 2 | Pure addition on an existing interface. |
| CPU-01 (static CPU info) | 5 | One collector, several files, clear parsing, full tests. |
| MEM-01 (`memory`) | 8 | Multiple `/proc` sources, derived available/pressure, edge cases. |
| NET-01 (Netlink) | 13 | New platform capability, protocol parsing, high uncertainty. |

---

## Planning poker

Estimation happens at **backlog refinement**, by the Developers:

1. Product Owner reads the story and its spec link.
2. Each developer privately selects a card (1/2/3/5/8/13).
3. Reveal simultaneously.
4. If aligned → record the point value.
5. If spread (e.g. 3 vs 8) → the high and low voters explain; re-vote. Spread usually means a hidden unknown — surface it into the acceptance criteria or split the story.
6. Persistent 13s that hide unknowns get split before being accepted as *Ready*.

Points are recorded only in [`product-backlog.md`](product-backlog.md).

---

## Capacity & commitment

We commit by **velocity**, adjusted for known capacity changes:

```text
Sprint commitment ≈ rolling-average velocity × (available capacity this sprint ÷ normal capacity)
```

- **Normal capacity** = the team's typical developer-days in a 2-week sprint.
- Reduce commitment for holidays, conferences, or known contributor absence.
- Foundation sprints (0–1) run **below** nominal because tooling and unknowns dominate.

### Assumed baseline for this plan

This plan forecasts with a **planned velocity of ~23 pts/sprint** once warmed up, and a **warm-up of ~21–24 pts** in Sprints 0–1. These are *forecasts*; they are replaced by **actual measured velocity** after Sprint 1.

---

## Velocity tracking

After each sprint, record actuals here (Scrum Master):

| Sprint | Committed | Completed (Done) | Velocity | Rolling avg (3) | Notes |
|---|---|---|---|---|---|
| 0 | 21 | — | — | — | Foundation; expect variance. |
| 1 | 24 | — | — | — | First feature slice. |
| 2 | 21 | — | — | — | |
| 3 | 22 | — | — | — | v0.1 release overhead. |
| 4 | 21 | — | — | — | |
| 5 | 29 | — | — | — | Netlink risk (see release plan). |
| 6 | 18 | — | — | — | v0.2 release. |
| 7 | 21 | — | — | — | TUI warm-up. |
| 8 | 13 | — | — | — | Single large story. |
| 9 | 16 | — | — | — | v0.3 release. |
| 10 | 21 | — | — | — | |
| 11 | 21 | — | — | — | v0.4 release. |
| 12 | 26 | — | — | — | Plugin protocol risk. |
| 13 | 19 | — | — | — | v0.5 release. |
| 14 | 39 | — | — | — | Stabilization; may split. |

**Rolling average** uses the last 3 completed sprints and is what the next commitment is based on. Ignore Sprint 0–1 in the average until Sprint 3.

---

## Definition: "completed"

A story counts toward velocity **only when it meets the [Definition of Done](../standards/definition-of-done.md)** — merged to `main`, tests green, docs and CHANGELOG updated. Partially finished stories score **zero** points that sprint; they do not carry partial credit. Unfinished work returns to the backlog and is re-estimated if scope changed.

---

## Common estimation traps (avoid)

- **Estimating in disguised hours.** "8 = one week" defeats the purpose; estimate relative effort.
- **Padding for testing.** Tests are part of Done, so they are *inside* the estimate, not added on top.
- **Anchoring on the loudest voice.** That's why voting is simultaneous.
- **Never re-estimating.** If a story is discovered to be bigger mid-sprint, note it, finish or return it, and re-point at refinement.
