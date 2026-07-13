# Roadmap

> A user-facing summary of SysKit's development plan.

## Current Status

SysKit v0.1.0 is released with system, CPU, memory, disk, and filesystem
inspection commands in table and JSON formats. The next milestone is v0.2:
processes, networking, ports, filtering, and YAML output. v0.2 implementation
is in progress: filtering, process listing/tree with user identities, memory
percentages, and sampled CPU percentages, interface counters plus state, MTU, MAC, IPv4/IPv6 addresses, route, and DNS
views, TCP/UDP/IPv6/Unix port listing with best-effort process ownership, and
YAML output are available.

## Planned Milestones

| Milestone | Focus | Status |
|---|---|---|
| v0.1 | Foundation, system, CPU, memory, disk, filesystem, table/JSON | Released (v0.1.0) |
| v0.2 | Processes, networking, ports, filtering, sorting, YAML | In progress |
| v0.3 | Watch mode, terminal dashboard, live process monitor | In progress (dashboard summaries, generic watch, and interactive top available) |
| v0.4 | Container-aware inspection and cgroup resource visibility | In progress (cgroup v1/v2 detection and normalized metrics available) |
| v0.5 | Plugin architecture and external collectors | Planned |
| v1.0 | Stable CLI contracts, packaging, complete documentation | Planned |

## v0.1.0 Contents

- Native Linux inspection for host, CPU, memory, mounted storage, and inode information.
- Deterministic table and JSON output contracts with fixture, integration, benchmark, and race coverage.
- A single Go-built Linux binary; distribution packages arrive in a later milestone.

## Detailed Roadmap

See [specs/roadmap.md](../specs/roadmap.md) for the full development roadmap with technical breakdowns and future considerations.

## Release History

v0.1.0 is the first preview release. See the [CHANGELOG](../CHANGELOG.md) for its release notes.
