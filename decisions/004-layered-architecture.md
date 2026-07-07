# 004. Adopt a layered architecture with independent collectors

**Status:** Accepted, 2026-07-01

> **Note (added 2026-07-07):** This ADR is preserved as the original decision record and is intentionally left unedited below. For the up-to-date, canonical description of how this layered architecture is ultimately documented — current diagrams, component responsibilities, and data flow — see [`ARCHITECTURE.md`](../ARCHITECTURE.md) at the repository root. This ADR records *why* the decision was made; `ARCHITECTURE.md` records *how it now stands*.

---

## Context

SysKit needs to grow from a handful of inspection commands (v0.1) to processes
and networking (v0.2), a real-time dashboard (v0.3), containers (v0.4), and a
plugin system (v0.5) — all while remaining testable, modular, and understandable
(constitution principles 4, *Keep It Modular*, and 5, *Test Everything*).

Two structural questions drive this decision:

1. **How do we separate user interaction, business logic, and data access** so
   that each can evolve independently — a new output format without touching
   collectors, a rewritten collector without touching commands?
2. **How do the per-domain data gatherers relate to each other** so that adding a
   network collector never risks breaking the CPU collector, and so both can be
   developed and tested in isolation?

A naive flat structure — commands reading `/proc` directly and formatting output
inline — is simplest to start but couples presentation, logic, and I/O together.
It becomes progressively harder to test (no seam to inject fake data), harder to
extend (every new format or command touches everything), and impossible to reuse
across the CLI and the TUI dashboard, which need the same data through different
presentations.

The [architecture spec](../specs/architecture.md) already lays out the intended
layered design and its rationale. This ADR ratifies that design as a binding
decision and records the trade-off explicitly.

---

## Decision

We will adopt a **six-layer architecture with a strict downward dependency rule**
and **independent per-domain collectors**.

The layers, from top to bottom:

1. **CLI Layer** — user interaction: flag/argument parsing, help text, output
   formatting (table/JSON/YAML), and terminal rendering (including the Bubble Tea
   dashboard, see [ADR 006](./006-bubbletea-for-tui.md)).
2. **Command Layer** — thin command definitions that validate input and translate
   user intent into service calls. Commands contain no business logic or I/O.
3. **Service Layer** — business logic: aggregating data from one or more
   collectors, filtering, sorting, and computing derived metrics (rates, deltas,
   percentages). Independent of output format.
4. **Collector Layer** — per-domain data gathering (CPU, memory, disk, network,
   process, …). Each collector reads through the platform abstraction and returns
   normalized, typed Go structs.
5. **Platform Abstraction Layer** — the only layer that touches the OS. It wraps
   `/proc`/`/sys` reads, Netlink sockets, and cgroup v1/v2 differences behind
   stable Go interfaces, containing kernel-version variation
   (see [ADR 003](./003-native-apis-over-shell.md)).
6. **Linux Kernel Interfaces** — the kernel-provided data sources themselves
   (`/proc`, `/sys`, Netlink, cgroups). Not code we write.

**Rules:**

- **Dependencies flow strictly downward.** Each layer may call only the layer
  immediately below it; no layer reaches upward or skips levels.
- **Layer boundaries are Go interfaces.** In particular, collectors depend on a
  platform-abstraction *interface*, not concrete file readers, so tests can inject
  fakes and stubs. Services depend on collector interfaces likewise.
- **Collectors are independent of one another.** The CPU collector has no
  knowledge of the network collector. New collectors are added without modifying
  existing ones — the foundation for the [v0.5 plugin system](../specs/roadmap.md).

---

## Consequences

### Positive

- **Testability.** Interface boundaries create seams for mocks and stubs. The
  platform layer can be faked so collectors, services, and commands are unit-
  tested without a real `/proc` — enabling deterministic tests in CI
  (*Test Everything*).
- **Modularity.** Presentation, logic, and I/O evolve independently. A new output
  formatter touches only the CLI layer; a rewritten collector touches only itself.
- **Reuse across presentations.** The CLI and the TUI dashboard consume the same
  services and collectors through different top layers — no duplicated data logic.
- **Parallel development.** Independent collectors can be built and reviewed by
  different contributors without coordination or merge contention.
- **Extensibility.** Adding a collector or command is additive, which is what
  makes the planned plugin architecture feasible.

### Negative

- **Indirection cost.** More packages, interfaces, and hand-offs than a flat
  design. A trivial command still passes through command → service → collector →
  platform, which is more ceremony for simple cases.
- **Boilerplate.** Each layer boundary needs an interface and, for tests, a fake
  implementation. Small features carry a fixed structural overhead.
- **Discipline required.** The downward-only rule must be enforced in review;
  it is easy to "just read a file" from a command under time pressure, which would
  erode the architecture.

### Neutral

- Thin commands and services mean logic lives one layer deeper than a newcomer
  might first expect — a learning curve, not a defect.
- The abstraction assumes a single OS target; because SysKit is Linux-only
  (see [ADR 002](./002-linux-only.md)), the platform layer abstracts kernel- and
  cgroup-version variation rather than operating systems.

---

## Alternatives Considered

- **Flat architecture (commands read `/proc` directly and format inline).**
  Simplest to start and requires the least code up front. Rejected because it
  couples presentation, business logic, and I/O, leaving no seam for testing with
  fake data and forcing wide-reaching changes for new formats or the dashboard.
  The architecture spec explicitly weighs and rejects this: layering "adds a small
  amount of indirection in exchange for significant gains in testability,
  modularity, and long-term maintainability."
- **Two layers (CLI + everything else).** A middle ground that folds services,
  collectors, and platform access together. Rejected because it still entangles
  data access with business logic, blocks independent collector testing, and does
  not support the plugin boundary needed in v0.5.
- **A shared global registry that collectors read from each other.** Would let
  collectors reuse one another's data. Rejected because it creates cross-collector
  coupling, breaking the independence that enables isolated testing and parallel
  development; cross-collector aggregation belongs in the service layer instead.

---

## References

- [Architecture](../specs/architecture.md) — full layer descriptions, data flow,
  and design-decision rationale
- [Constitution](../specs/constitution.md) — principle 4 (Keep It Modular),
  principle 5 (Test Everything)
- [ADR 002](./002-linux-only.md) — Linux-only scope of the platform layer
- [ADR 003](./003-native-apis-over-shell.md) — native interfaces the platform layer wraps
- [ADR 006](./006-bubbletea-for-tui.md) — the TUI as a second presentation over shared services
- [Roadmap](../specs/roadmap.md) — v0.5 plugin system built on independent collectors
