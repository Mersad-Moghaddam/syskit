# Architecture Overview (User-Facing Entry Point)

> **Canonical source:** [`ARCHITECTURE.md`](../ARCHITECTURE.md) (repository root) is the single canonical baseline architecture for SysKit. This file previously carried its own Mermaid diagram and component table, which drifted from the spec and ADR versions (tracked as W-8). That duplicated content has been removed in favor of a pointer.

## Purpose of this file

This file is retained as the **user-facing doorway** into the architecture from the `docs/` set (getting-started, onboarding, and the documentation map link here). It gives a newcomer the one-paragraph orientation and then sends them to the canonical document for detail.

## In one paragraph

SysKit is a Linux-only, read-only command-line tool: a single Go binary that reads native kernel interfaces (`/proc`, `/sys`, Netlink, cgroups) directly, transforms the data, and renders it as table, JSON, YAML, or an interactive terminal interface. Its hierarchical control center is a CLI-only discoverability layer that delegates every action to the existing command tree. The project is organized as a strict, one-way layered pipeline — CLI → Command → Service → Collector → Platform → kernel — so that operating-system access stays away from presentation code, collectors can be tested against fixtures, and output contracts stay stable. It holds no persistent storage, cache, or queue.

## Read next

- [`ARCHITECTURE.md`](../ARCHITECTURE.md) — the canonical baseline: layer diagram (§3), component responsibilities (§4), data model (§5), and key decisions (§6).
- [Detailed architecture specification](../specs/architecture.md) — spec-layer anchor (also points to the canonical doc).
- [Collector architecture](../specs/collectors.md) · [Rendering architecture](../specs/rendering.md) · [Plugin architecture](../specs/plugin-architecture.md)
- [Engineering constitution](../specs/constitution.md)

## Operating boundaries

SysKit is read-only. It inspects, correlates, and reports, but does not change system configuration, stop services, kill processes, edit kernel parameters, or manage containers in the core command set. See [`ARCHITECTURE.md`](../ARCHITECTURE.md) §7 for the full non-functional and safety considerations.
