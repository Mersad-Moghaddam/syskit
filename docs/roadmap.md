# Roadmap

> A user-facing summary of SysKit's development plan.

## Current Status

SysKit v0.2.0 is released with process inspection, network interface/route/DNS
views, and port inspection alongside the v0.1 system commands. Table, JSON,
and YAML output are available for one-shot commands. v0.3 and v0.4 work is in
progress, including live monitoring and cgroup-derived container inspection.

## Planned Milestones

| Milestone | Focus | Status |
|---|---|---|
| v0.1 | Foundation, system, CPU, memory, disk, filesystem, table/JSON | Released (v0.1.0) |
| v0.2 | Processes, networking, ports, filtering, sorting, YAML | Released (v0.2.0) |
| v0.3 | Watch mode, terminal dashboard, live process monitor | Released (v0.3.0) |
| v0.4 | Container-aware inspection and cgroup resource visibility | Released (v0.4.0) |
| v0.5 | Plugin architecture and external collectors | Planned |
| v1.0 | Stable CLI contracts, packaging, complete documentation | Planned |

## v0.1.0 Contents

- Native Linux inspection for host, CPU, memory, mounted storage, and inode information.
- Deterministic table and JSON output contracts with fixture, integration, benchmark, and race coverage.
- A single Go-built Linux binary; distribution packages arrive in a later milestone.

## v0.2.0 Contents

- Native procfs process listing and tree views with filtering, sorting, limits,
  user identities, partial-data reporting, memory percentage, and sampled CPU
  percentage.
- Native Netlink interface-address and route views, plus procfs/sysfs counters
  and resolver configuration.
- TCP, UDP, IPv6, and Unix socket inspection with best-effort process ownership.

## Detailed Roadmap

See [specs/roadmap.md](../specs/roadmap.md) for the full development roadmap with technical breakdowns and future considerations.

## Release History

See the [CHANGELOG](../CHANGELOG.md) for v0.1.0 and v0.2.0 release notes.
