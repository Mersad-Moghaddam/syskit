# Architecture Decision Records

> A log of the significant, cross-cutting, and hard-to-reverse decisions that shape SysKit.

---

## What is an ADR?

An Architecture Decision Record (ADR) is a short document that captures a single
architecturally significant decision, the context that forced it, and the
consequences of choosing it. Each ADR is immutable once accepted: rather than
editing a decision after the fact, we supersede it with a new record and update
the status of the old one.

SysKit uses the [MADR](https://adr.github.io/madr/) (Markdown Architectural
Decision Records) structure — a lightweight, prose-oriented format that fits our
Specification-Driven Development workflow described in the
[constitution](../specs/constitution.md).

---

## Why SysKit uses ADRs

SysKit is a learning project as much as it is a tool (constitution principle 10,
*Learn Before Build*), and its foundational technical choices — the language, the
platform, the data-access strategy, the architecture — are decisions we want to
be able to revisit and *understand* months or years later.

ADRs give us:

- **A durable record of intent.** The *why* behind a decision is far more
  valuable than the *what*, and it is the first thing lost when only code
  survives.
- **Honest trade-off analysis.** Every ADR forces us to write down what we gave
  up, not just what we gained.
- **Onboarding context.** A new contributor can read the decision log and
  understand the shape of the project before touching a single line of Go.
- **A guard against churn.** When someone proposes reversing a decision, the ADR
  is the starting point for that conversation.

This directly supports the constitution's *Documentation First* principle:
decisions are documented before they are committed.

---

## When is a new ADR required?

Write an ADR when a decision is any of the following:

- **Significant** — it affects the structure, dependencies, or public interfaces
  of the project (e.g. adopting a CLI framework).
- **Cross-cutting** — it touches multiple layers or subsystems rather than a
  single collector or command (e.g. the layered architecture itself).
- **Hard to reverse** — undoing it later would be expensive or disruptive (e.g.
  the choice of implementation language or target platform).

Routine choices — a helper function's signature, a variable name, the internal
layout of a single collector — do **not** need an ADR. When in doubt, ask whether
a future maintainer would benefit from knowing *why* the choice was made. If yes,
write the record.

---

## Numbering scheme

ADRs are numbered sequentially with a zero-padded three-digit prefix, in the
order they are accepted:

```
decisions/NNN-short-kebab-title.md
```

For example, `001-use-go.md`. Numbers are never reused. If a decision is
superseded, the original file and number remain in place with an updated status;
a new record with the next available number captures the replacement.

---

## Lifecycle

Every ADR carries exactly one status at a time:

| Status         | Meaning                                                                 |
|----------------|-------------------------------------------------------------------------|
| **Proposed**   | Drafted and under discussion; not yet binding.                          |
| **Accepted**   | Agreed and in force. This is the default state of the records below.    |
| **Deprecated** | No longer recommended, but not yet replaced by a specific decision.     |
| **Superseded** | Replaced by a later ADR. The status line names the superseding record.  |

The normal progression is **Proposed → Accepted**. From there a decision may
later become **Deprecated** or **Superseded by NNN**. An accepted ADR is never
edited to change its decision — it is superseded instead, preserving the
historical record.

---

## How to create a new ADR

1. Copy [`template.md`](./template.md) to `decisions/NNN-short-title.md`, using
   the next available number.
2. Fill in every section. Replace the guidance comments with real content. An
   accepted record must not contain unresolved task markers.
3. Open it with status **Proposed** for review.
4. Once agreed, change the status to **Accepted** with the acceptance date and
   add a row to the Decision Log below.

---

## Decision Log

| ID  | Title                                                              | Status   | Date       |
|-----|--------------------------------------------------------------------|----------|------------|
| [001](./001-use-go.md) | Use Go as the implementation language                 | Accepted | 2026-07-01 |
| [002](./002-linux-only.md) | Target Linux exclusively                          | Accepted | 2026-07-01 |
| [003](./003-native-apis-over-shell.md) | Read native kernel interfaces instead of shelling out | Accepted | 2026-07-01 |
| [004](./004-layered-architecture.md) | Adopt a layered architecture with independent collectors | Accepted | 2026-07-01 |
| [005](./005-cobra-for-cli.md) | Use Cobra for the CLI framework                    | Accepted | 2026-07-01 |
| [006](./006-bubbletea-for-tui.md) | Use Bubble Tea for the interactive dashboard TUI | Accepted | 2026-07-01 |
| [007](./007-out-of-process-plugins.md) | Prefer out-of-process plugins                 | Accepted | 2026-07-01 |
| [008](./008-no-persistent-storage.md) | No persistent storage, cache, or queue in core scope | Accepted | 2026-07-07 |

---

*This log is a living index. Its structure — one record per significant
decision, numbered and never rewritten — does not change.*
