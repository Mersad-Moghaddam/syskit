# SysKit Architecture (Spec-Layer Anchor)

> **Canonical source:** [`ARCHITECTURE.md`](../ARCHITECTURE.md) (repository root) is the single canonical baseline architecture for SysKit. It supersedes the detailed layer descriptions, ASCII diagram, and data-flow narrative that previously lived in this file, which had drifted from the other architecture documents (tracked as W-8).

## Purpose of this file

This file is retained as the **spec-layer cross-reference target**: the ADRs and other specs link to `specs/architecture.md`, and the planning-phase CI requires it to exist. It no longer restates the architecture — it points to the canonical document so there is exactly one place the layer model, dependency rules, component responsibilities, and data flow are defined.

For everything about how SysKit is structured, read [`ARCHITECTURE.md`](../ARCHITECTURE.md):

- High-level architecture and the strict downward-dependency rule — §3
- Components & responsibilities per layer — §4
- Data model and snapshot semantics — §5
- Key design decisions (ADR summaries) — §6

## Related specifications

The cross-cutting specs remain the detailed, non-duplicative companions to the canonical architecture:

- [Collector architecture](./collectors.md)
- [Rendering architecture](./rendering.md)
- [Plugin architecture](./plugin-architecture.md)
- [CLI conventions](./cli-conventions.md)
- [Error handling](./error-handling.md)
- [Configuration](./configuration.md)
- [Engineering constitution](./constitution.md)

The binding architectural decision itself is recorded in [ADR 004 — Layered architecture](../decisions/004-layered-architecture.md).
