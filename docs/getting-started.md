# Getting Started

> How to approach SysKit during the design phase and what users can expect once implementation begins.

SysKit is currently a planning repository. It contains architecture, specifications, standards, learning notes, and project governance, but no production Go source code yet. This is intentional: the repository is being prepared so implementation can start from a clear, reviewed design.

## Current State

| Area | Status |
|---|---|
| Product vision | Defined |
| Architecture | Defined |
| Feature specifications | Defined for planned core features |
| Engineering standards | Defined |
| Production code | Not started |
| Installable releases | Not available yet |

## Expected Development Environment

When implementation begins, contributors should expect to work on Linux with:

- Go 1.22 or newer.
- A Linux kernel with procfs and sysfs mounted in the standard locations.
- Standard development tools such as `git`, `make` or shell, and a POSIX-compatible terminal.
- Normal user permissions for inspection commands. Commands that read restricted process details may report partial data unless run with elevated permissions.

The project is Linux-first by design. macOS, Windows, and BSD support are out of scope for the core tool.

## How to Use This Repository Now

Start with the documents in this order:

1. [README](../README.md) for the project vision and feature map.
2. [Product overview](../specs/product.md) for goals, users, and non-goals.
3. [Architecture overview](architecture.md) for the system shape.
4. [Feature specifications](../specs/features/) for expected behavior.
5. [Learning roadmap](../learning/roadmap.md) for Linux concepts to study before implementation.
6. [Implementation readiness checklist](implementation-readiness.md) before creating production code.

## Expected First Commands

These commands are planned contracts, not currently executable commands:

```sh
syskit system
syskit cpu --format json
syskit memory --watch
syskit process --sort cpu
syskit network interfaces
syskit ports --listening
syskit dashboard
```

Each command must be implemented only after its specification and acceptance criteria are reviewed.

## Configuration Model

SysKit is designed to work with zero configuration. Future configuration will use this precedence:

1. Command-line flags.
2. `SYSKIT_*` environment variables.
3. XDG config file at `$XDG_CONFIG_HOME/syskit/config.toml` or `~/.config/syskit/config.toml`.
4. Built-in defaults.

See [configuration](../specs/configuration.md) for the detailed policy.

## Next Steps

- [Architecture overview](architecture.md)
- [CLI conventions](../specs/cli-conventions.md)
- [Contributing](contributing.md)
- [Roadmap](roadmap.md)
