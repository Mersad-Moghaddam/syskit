# Roadmap

> A user-facing summary of SysKit's development plan.

## Current Status

SysKit is in the design and specification phase. No features have been released yet, and no production Go code is present. The immediate objective is to finish the planning foundation so implementation can begin deliberately.

## Planned Milestones

| Milestone | Focus | Status |
|---|---|---|
| v0.1 | Foundation: system, CPU, memory, disk, basic output | Planned |
| v0.2 | Processes, networking, ports, filtering, sorting | Planned |
| v0.3 | Watch mode, terminal dashboard, live process monitor | Planned |
| v0.4 | Container-aware inspection and cgroup resource visibility | Planned |
| v0.5 | Plugin architecture and external collectors | Planned |
| v1.0 | Stable CLI contracts, packaging, complete documentation | Planned |

## What Must Happen Before v0.1

- Review the architecture and feature specs.
- Confirm command names, flag conventions, and output formats.
- Create the initial Go module and repository layout.
- Add collector interfaces and platform abstractions.
- Add fixture strategy for procfs and sysfs data.
- Replace planning-phase CI with implementation CI.

## Detailed Roadmap

See [specs/roadmap.md](../specs/roadmap.md) for the full development roadmap with technical breakdowns and future considerations.

## Release History

No releases have been published yet. Release notes will be maintained in [CHANGELOG](../CHANGELOG.md) once implementation begins.
