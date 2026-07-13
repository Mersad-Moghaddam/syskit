# Getting Started

> How to build and explore SysKit during its v0.1 implementation phase.

SysKit is a Linux-first Go project. v0.1.0 is available to build from source
with system, CPU, memory, disk, and filesystem inspection commands.

## Current State

| Area | Status |
|---|---|
| Product vision | Defined |
| Architecture | Defined |
| Feature specifications | Defined for planned core features |
| Engineering standards | Defined |
| Production code | v0.1.0 released |
| Installable releases | Source build available; packages planned |

## Development Environment

Build and test on Linux with:

- Go 1.22 or newer.
- A Linux kernel with procfs and sysfs mounted in the standard locations.
- Standard development tools such as `git`, `make` or shell, and a POSIX-compatible terminal.
- Normal user permissions for inspection commands. Commands that read restricted process details may report partial data unless run with elevated permissions.

The project is Linux-first by design. macOS, Windows, and BSD support are out of scope for the core tool.

## Build and Test

```sh
go build ./...
go test -race ./...
go run ./cmd/syskit --help
```

The CLI currently exposes the v0.1 commands plus in-progress v0.2 process,
network, and ports commands. For example:

```sh
go run ./cmd/syskit system
go run ./cmd/syskit system --format json
go run ./cmd/syskit cpu
go run ./cmd/syskit memory --format json
go run ./cmd/syskit filesystem --show-pseudo
go run ./cmd/syskit disk --io --interval 1s
go run ./cmd/syskit process --limit 20
go run ./cmd/syskit process tree
go run ./cmd/syskit network
go run ./cmd/syskit ports --listening --pid 1234
```

`ports` reads TCP, UDP, IPv6, and Unix socket tables directly from procfs. It
best-effort maps socket inodes to process IDs and commands; inaccessible or
short-lived processes simply remain unmapped. The remaining commands below include both shipped and planned contracts; check
`--help` for the current supported flags.

## How to Explore This Repository

Start with the documents in this order:

1. [README](../README.md) for the project vision and feature map.
2. [Product overview](../specs/product.md) for goals, users, and non-goals.
3. [Architecture overview](architecture.md) for the system shape.
4. [Feature specifications](../specs/features/) for expected behavior.
5. [Learning roadmap](../learning/roadmap.md) for Linux concepts to study before implementation.
6. [Implementation readiness checklist](implementation-readiness.md) before creating production code.

## Command Roadmap

These commands show the intended product shape; `system`, `cpu`, `memory`,
`disk`, `filesystem`, `process`, `network`, and `ports` are executable today.

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

SysKit works with zero configuration. The implemented configuration loader uses
this precedence:

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
