# Architecture Overview

> A user-facing summary of how SysKit is expected to work once implementation begins.

SysKit is planned as a layered Linux CLI. Each layer has one job: receive user intent, coordinate a command, collect system data, transform it, and render it consistently. The design keeps operating-system access away from presentation code so that collectors can be tested with fixtures and output contracts can stay stable.

## System Architecture

```mermaid
flowchart TD
    User[User]
    CLI[CLI layer]
    Commands[Command layer]
    Services[Service layer]
    Collectors[Collector layer]
    Platform[Linux platform adapters]
    Kernel[/proc, /sys, Netlink, cgroups]
    Renderers[Renderers: table, JSON, YAML, TUI]

    User --> CLI
    CLI --> Commands
    Commands --> Services
    Services --> Collectors
    Collectors --> Platform
    Platform --> Kernel
    Services --> Renderers
    Renderers --> User
```

The dependency direction is intentionally one-way. CLI code may depend on services, services may depend on collectors, collectors may depend on platform adapters, and platform adapters may read kernel interfaces. Lower layers never import higher layers.

## Component Overview

| Component | Responsibility | Must not do |
|---|---|---|
| CLI layer | Parse arguments, load configuration, choose output format, display errors | Read `/proc` or `/sys` directly |
| Command layer | Validate command-specific flags and call services | Perform business logic |
| Service layer | Aggregate collector data, compute derived metrics, apply filters | Render terminal output |
| Collector layer | Parse raw Linux data into typed domain models | Know about Cobra, tables, or terminal colors |
| Platform adapters | Read procfs, sysfs, Netlink, and cgroup data | Interpret user intent |
| Renderers | Convert structured results into table, JSON, YAML, or dashboard views | Collect system data |

## Data Flow

A command such as `syskit cpu --format json` is expected to flow like this:

1. The CLI parses `cpu` and `--format json`.
2. The command validates flag combinations and calls the CPU service.
3. The service requests CPU snapshots from the CPU collector.
4. The collector reads `/proc/stat`, `/proc/cpuinfo`, and relevant sysfs paths through platform adapters.
5. Parsed data returns as structured models.
6. The service computes derived values such as utilization percentages.
7. The renderer emits the selected format.

The same service result should be usable by table output, JSON output, tests, and the future terminal dashboard.

## Extension Points

SysKit's first implementation should prioritize the core CLI and built-in collectors. The architecture still reserves several future extension points:

| Extension point | Planned use |
|---|---|
| Collectors | Add new system domains without rewriting existing commands |
| Renderers | Add formats or terminal views without changing collection logic |
| Services | Compose multiple collectors into higher-level diagnostics |
| Plugins | Load external collectors after the core contracts stabilize |
| Fixtures | Add Linux distributions, kernel versions, and container environments to the test corpus |

## Operating Boundaries

SysKit is read-only. It may inspect resources, correlate data, and report potential problems, but it should not change system configuration, stop services, kill processes, edit kernel parameters, or manage containers in the core command set.

## Further Reading

- [Detailed architecture specification](../specs/architecture.md)
- [Collector architecture](../specs/collectors.md)
- [Rendering architecture](../specs/rendering.md)
- [Plugin architecture](../specs/plugin-architecture.md)
- [Engineering constitution](../specs/constitution.md)
