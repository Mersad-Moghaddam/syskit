# Roadmap

> A user-facing summary of SysKit's development plan.

## Current Status

SysKit is in the implementation phase, building v0.1. The foundation is
complete: the Go module, CLI bootstrap, platform seam, collector contract,
renderers, configuration, logging, error handling, fixture tooling, and Go CI
are in place. The first user-facing command, `syskit system`, is available;
the CPU, memory, disk, and filesystem slices remain before v0.1 can release.

## Planned Milestones

| Milestone | Focus | Status |
|---|---|---|
| v0.1 | Foundation complete; system, CPU, memory, disk, basic output | In progress |
| v0.2 | Processes, networking, ports, filtering, sorting | Planned |
| v0.3 | Watch mode, terminal dashboard, live process monitor | Planned |
| v0.4 | Container-aware inspection and cgroup resource visibility | Planned |
| v0.5 | Plugin architecture and external collectors | Planned |
| v1.0 | Stable CLI contracts, packaging, complete documentation | Planned |

## Remaining Work for v0.1

- Implement and document the `system`, CPU, memory, disk, and filesystem
  vertical slices against their accepted specifications.
- Add fixture-backed parser, service, command, renderer, integration, and
  benchmark coverage for each command.
- Complete the v0.1 release checklist and publish the first release.

## Detailed Roadmap

See [specs/roadmap.md](../specs/roadmap.md) for the full development roadmap with technical breakdowns and future considerations.

## Release History

No releases have been published yet. The [CHANGELOG](../CHANGELOG.md) records
unreleased implementation work; versioned release notes will be added when
v0.1.0 is published.
