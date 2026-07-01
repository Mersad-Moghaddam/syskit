# Collector Architecture

> How SysKit should gather Linux system data without mixing kernel access, business logic, and presentation.

Collectors are the boundary between SysKit's domain model and Linux kernel interfaces. A collector reads one domain, parses native data sources, and returns normalized structured data.

## Goals

- Keep `/proc`, `/sys`, Netlink, and cgroup access out of CLI code.
- Make data parsing testable with fixtures.
- Keep collectors independent so features can evolve separately.
- Preserve raw Linux semantics where they matter.
- Return partial data carefully when a host lacks optional interfaces.

## Collector Responsibilities

Collectors should:

- Read from platform adapters, not directly from the host filesystem.
- Parse raw bytes into typed domain records.
- Normalize units to explicit base units such as bytes, nanoseconds, or counters.
- Preserve source metadata when useful for diagnostics.
- Return domain-specific errors that can be classified by the service layer.

Collectors should not:

- Parse command-line flags.
- Render output.
- Read configuration files.
- Decide terminal colors.
- Shell out to traditional Linux utilities.

## Snapshot Model

Most collector calls should return a point-in-time snapshot. Rate-based metrics, such as CPU utilization or network throughput, require two snapshots and a time delta. The service layer should own those calculations so collectors remain simple and reusable.

## Data Source Priority

Use native Linux interfaces in this order unless a feature spec says otherwise:

1. procfs and sysfs files with stable kernel documentation.
2. Netlink for networking and routing data.
3. cgroup files for container and resource accounting.
4. Runtime APIs for optional container features.

Shelling out to tools such as `ps`, `df`, `ss`, `ip`, or `free` is out of scope for core collectors.

## Error Classification

Collectors should distinguish:

- Missing optional data.
- Missing required data.
- Permission denied.
- Malformed kernel data.
- Unsupported kernel capability.
- Race with disappearing resources, especially processes.

The service layer decides whether an error is fatal or can be represented as partial data.

## Fixtures

Every collector must define fixture needs before implementation. Fixture sets should capture:

- Small virtual machine.
- Multi-core host.
- Containerized environment.
- Restricted permissions.
- Older kernel behavior where practical.
- Malformed or truncated data for parser tests.

## Concurrency

Collectors may be called concurrently by services, dashboards, and watch mode. They should avoid package-level mutable state. Caching should be explicit, bounded, and owned by services unless a collector has a documented reason to cache.

## Acceptance Criteria

- Each collector has a clear domain boundary.
- Each collector can run against fixture data.
- Each collector exposes enough source context for useful errors.
- No collector imports CLI or rendering concerns.
- No collector shells out to external utilities for core data.
