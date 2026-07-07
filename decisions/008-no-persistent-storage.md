# ADR-008: No Persistent Storage, Cache, or Queue in Core Scope

**Status:** Accepted
**Date:** 2026-07-07
**Deciders:** SysKit maintainers

## Context

SysKit's domain is transient system-telemetry snapshots (host, CPU, memory, disk, filesystem, process, network, ports, containers). Every invocation reads live kernel interfaces (`/proc`, `/sys`, Netlink, cgroups) into typed structs, renders them, and exits. The product is explicitly scoped as a read-only, single-static-binary CLI (ADR 001, ADR 003; `specs/product.md` Non-Goals).

Three mechanisms are commonly reached for by default in tools that collect and present data, and each needs an explicit decision rather than silent omission:

- **Cache** — could avoid re-reading kernel interfaces across calls.
- **Persistent storage** — could retain historical snapshots across invocations.
- **Queue / event broker** — could support asynchronous or fan-out processing.

Without a recorded decision, a future contributor could reasonably add one of these "by default," which would conflict with the single-binary, read-only identity the rest of the architecture depends on.

## Decision

SysKit will hold **no persistent storage, no cache subsystem, and no queue/broker** in core scope.

The only exception is bounded, in-memory, per-session state already implied by the live modes (`watch`, `top`, `dashboard`): these hold the *previous snapshot* in memory to compute deltas (e.g. CPU %, network throughput) between refresh ticks. This state lives inside the service/TUI model and is discarded when the process exits. It is not a cache subsystem and is not persisted.

If a future need for bounded, explicit, service-owned caching arises (e.g. repeated reads within one refresh cycle becoming a measured cost), it must be added deliberately per `specs/collectors.md`'s existing guidance — not adopted as a default.

## Alternatives Considered

| Option | Assessment |
|---|---|
| Add a cache layer (e.g. read-through cache in the platform layer) | Rejected: no invocation re-reads the same interface today; would add machinery with no measured benefit. |
| Add an embedded store (e.g. SQLite/BoltDB) for historical data | Rejected for core scope: conflicts with the single-static-binary, read-only, per-invocation identity (ADR 001). Historical data is listed only as a possible post-1.0 "Future Consideration," not a current requirement. |
| Add a queue/event bus for async collection or fan-out | Rejected: there is no asynchronous work or inter-process fan-out in the product; collection is synchronous read → render → exit. |

## Trade-offs

- **Accepted cost:** none identified for current scope — this decision adds nothing and removes a class of future missteps.
- **Deferred cost:** if historical data is adopted post-1.0, this ADR will need to be superseded, and that change will be the project's first departure from "no persistent state," carrying real design weight (storage format, retention, migration).

## Consequences

- Simplest possible runtime footprint: no datastore, no background service, no serialization format to version for storage.
- Collectors and services remain stateless and safe to run concurrently (supports the concurrency approach in the baseline architecture and `-race`-clean tests).
- The read-only boundary (`specs/product.md` Non-Goals; `docs/architecture.md` Operating Boundaries) is reinforced at the architecture level, not just documented as intent.
- Any future proposal to add storage, caching, or a queue must go through a new ADR that explicitly supersedes this one, rather than being introduced incrementally.

## Action Items

1. [x] Add this file to `decisions/`. Filed as `decisions/008-no-persistent-storage.md` to match the established `NNN-short-title.md` convention used by ADR 001–007 (`decisions/README.md`, `docs/project-structure.md`), rather than an `ADR-008-` prefix that would be inconsistent with the existing records.
2. [x] Cross-reference this ADR from `ARCHITECTURE.md` §6 (where the no-storage/cache/queue decision is asserted) so the rationale has one canonical source. `specs/architecture.md` and `docs/architecture.md` do not themselves assert "no storage/cache/queue"; both now defer to `ARCHITECTURE.md`, which links here.
3. [x] Update `specs/collectors.md` to reference this ADR where it discusses bounded, service-owned caching, so the "explicit, bounded, service-owned" guidance is traceable to a decision record.
