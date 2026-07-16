# SysKit Governance

> How decisions are made, who makes them, and how the project is steered.

---

## Purpose

This document describes how the SysKit project is governed: the roles people hold, how decisions are reached, how significant changes are proposed and recorded, and how the governance model itself evolves.

SysKit welcomes external contributions across code, documentation, specifications,
Linux research, and project process. This document describes the rules that keep
that collaboration transparent and sustainable.

The engineering principles that decisions are measured against live in the [Engineering Constitution](specs/constitution.md). Governance defines *how* decisions are made; the Constitution defines *what* a good decision looks like. Where the two intersect, the Constitution supplies the technical standard and this document supplies the process.

---

## Roles & Responsibilities

SysKit recognizes four roles. They form a ladder: people advance by sustained, high-quality participation, and each role carries the responsibilities of the roles below it.

| Role | Responsibilities | How one advances |
|---|---|---|
| **User** | Uses SysKit, reports bugs, requests features, asks and answers questions. | Anyone is a user. Filing a well-formed issue or a helpful discussion post is the first step toward contributing. |
| **Contributor** | Submits pull requests, specifications, documentation, or reviews. Follows the Constitution and the contribution guidelines. | A user becomes a contributor the moment they open their first accepted pull request or specification. |
| **Reviewer** | Reviews contributions for correctness, quality, and alignment with the Constitution. Has authority to approve changes within their area of familiarity. | A contributor with a track record of sound, consistent contributions and reviews is invited to become a reviewer by the maintainers. |
| **Maintainer** | Holds merge and release authority. Owns the roadmap, resolves disputes, enforces the Code of Conduct, and stewards the project's long-term direction. | A reviewer who has demonstrated sustained judgment, reliability, and care for the project as a whole is invited to become a maintainer by consensus of the existing maintainers. |

Advancement is by invitation, based on demonstrated merit rather than volume alone. Roles may also be stepped down — voluntarily, or by maintainer consensus in cases of inactivity or a serious breach of the Code of Conduct.

---

## Decision-Making Process

SysKit uses two decision-making modes, chosen according to the weight of the change.

### Lazy Consensus

Routine changes proceed by **lazy consensus**: silence is assent. A change may move forward once it has at least one reviewer approval and no unresolved objections. Most work — bug fixes, documentation updates, small features that fit an existing specification, and dependency bumps — travels this path.

Lazy consensus is appropriate when a change:

- Fits within an already-approved specification or the existing architecture
- Does not alter a public CLI interface, output format, or stability guarantee
- Does not introduce a new dependency of significant weight
- Raises no principled objection during review

### Formal Review

Significant changes require **formal review** and explicit consensus among the maintainers before they proceed. A change is significant when it:

- Alters the public CLI surface, output formats, or backward-compatibility guarantees
- Introduces a new subsystem, collector category, or architectural boundary
- Adds a non-trivial external dependency
- Changes an engineering principle or an established convention
- Affects security, privacy, or the exposure of sensitive system data

Significant changes are not merged on a single approval. They are discussed openly, and the reasoning behind the eventual decision is recorded (see below). When consensus cannot be reached, the maintainers decide; disputes that remain deadlocked are resolved by a simple majority of maintainers.

---

## RFC & Significant-Change Workflow

Significant changes follow a lightweight RFC process backed by **Architecture Decision Records (ADRs)** stored in [`decisions/`](decisions/).

An ADR is **required** whenever a change:

- Establishes or revises an architectural boundary or a core abstraction
- Selects a technology, dependency, or approach where reasonable alternatives exist
- Changes a public interface, output contract, or stability guarantee
- Reverses or supersedes a decision recorded in a previous ADR

The workflow is:

1. **Propose** — Open an issue or discussion describing the problem, the options considered, and a recommendation.
2. **Draft the ADR** — Capture the context, the decision, the alternatives weighed, and the consequences as a new record in `decisions/`, in `proposed` status.
3. **Review** — The maintainers and interested reviewers evaluate the ADR against the [Engineering Constitution](specs/constitution.md). Formal-review rules apply.
4. **Decide** — On consensus, the ADR is marked `accepted` and merged. If declined, it is marked `rejected` and retained for the historical record.
5. **Supersede** — When a later decision replaces an earlier one, the old ADR is marked `superseded` with a pointer to its replacement. ADRs are never silently deleted; the decision history is preserved.

Routine changes that travel by lazy consensus do not require an ADR.

---

## Code Review Requirements

All changes reach `main` through pull requests. Direct commits to `main` are reserved for maintainers and used only in exceptional circumstances (for example, correcting a broken build).

Every pull request must:

- Receive at least one approving review from a reviewer or maintainer. Significant changes require maintainer approval.
- Pass all automated checks — build, tests, `go vet`, `go fmt`, and static analysis — without warnings.
- Include tests for new behavior and documentation for user-visible changes, per the Constitution's *Test Everything* and *Documentation First* principles.
- Have unresolved review comments addressed or explicitly deferred before merge.

Authors do not approve or merge their own significant changes. Review focuses on correctness, alignment with the Constitution, and long-term maintainability — not merely on whether the code works.

---

## Roadmap & Priorities

The direction of SysKit is expressed in the [roadmap](specs/roadmap.md), which is owned by the maintainers and informed by the whole community.

Priorities are set by weighing:

- Alignment with the project's product vision and the Engineering Constitution
- User-reported needs, expressed through issues and discussions
- Technical dependencies between milestones — earlier foundations enable later features
- The capacity and interest of active contributors

The roadmap is a living document. Milestones and their ordering may be revised as the project learns and the community grows. Substantial changes to the roadmap follow the formal-review process.

---

## Code of Conduct Enforcement

SysKit adopts the [Code of Conduct](CODE_OF_CONDUCT.md). Its enforcement is owned by the maintainers, who act as the community leaders referenced in that document.

Reports are received at **conduct@syskit.dev** and handled confidentially, following the enforcement guidelines in the Code of Conduct. Maintainers who are themselves the subject of a report recuse themselves from its handling. Enforcement decisions are made by the remaining maintainers.

---

## Amending This Document

Governance is not fixed. This document may be amended through the **formal review** process:

1. A change is proposed via a pull request against this file, with its rationale.
2. The change is discussed openly and evaluated against the project's principles.
3. It is adopted only on consensus of the maintainers.

Amendments that materially change how power is distributed or how decisions are made are also recorded as an ADR in [`decisions/`](decisions/), so the evolution of the governance model itself remains part of the project's decision history.

---

*This governance model is intentionally lightweight. It may grow more formal as
the contributor base grows, but its intent—transparent and principled
stewardship—does not change.*
