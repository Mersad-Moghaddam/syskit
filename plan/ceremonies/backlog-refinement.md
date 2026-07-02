# Ceremony: Backlog Refinement

> **When:** ~Day 6 (mid-sprint) · **Timebox:** 1 hour · **Output:** next sprint's candidate stories estimated and *Ready*

---

## Purpose

Keep the top of `../product-backlog.md` always ~1.5 sprints deep in *Ready* work, so planning is fast and no contributor is ever blocked for lack of a defined task. Refinement is where stories are clarified, estimated, split, and ordered — never in the middle of Sprint Planning.

## Attendees

Product Owner (required — owns order and clarifies intent), Developers (required — estimate and split), Scrum Master (facilitates).

## Agenda (1 hour)

| Time | Activity |
|---|---|
| 0:00–0:10 | PO presents candidate stories from the backlog top, in priority order. |
| 0:10–0:35 | For each: clarify intent, confirm the spec link, check against the [Definition of Ready](../../standards/definition-of-ready.md). |
| 0:35–0:55 | **Planning poker** estimation (`../estimation-and-velocity.md`); split anything ≥ 13 with hidden unknowns. |
| 0:55–1:00 | Update `../product-backlog.md`: points, status → `Ready`, order. Note any spikes needed. |

## The Ready checklist (gate to leave refinement)

A story becomes `Ready` only when all of `../../standards/definition-of-ready.md` holds:

- [ ] Feature spec exists and is reviewed.
- [ ] User story and motivation clear.
- [ ] CLI, flags, and output (table + structured) documented.
- [ ] Linux data sources identified.
- [ ] Edge cases listed.
- [ ] Acceptance criteria testable.
- [ ] Fixtures identified.
- [ ] New dependencies reviewed (ADR if needed).
- [ ] Security / permission / partial-data behavior understood.

If any box fails, the story stays in refinement (or spawns a design-proposal issue) — it does **not** enter a sprint.

## Splitting stories

Prefer **vertical** splits (still end-to-end, smaller scope) over horizontal (layer-only) splits:

- CPU-01 (static info) vs CPU-02 (utilization) — split by *behavior*, each shippable.
- NET-01 (Netlink platform) vs NET-02 (network command) — split by *risk isolation*.

Avoid splits that leave a non-demonstrable fragment ("just the parser, no command").

## Spikes

If a story carries a High-exposure risk (`../risk-register.md`), schedule its spike here so the investigation lands *before* the story is committed (e.g. SPK-NET before NET-01).

## Output

`../product-backlog.md` updated: candidate stories pointed, ordered, and marked `Ready`; spikes noted; nothing larger than 13 left unsplit.
