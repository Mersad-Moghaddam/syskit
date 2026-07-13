# Getting Started

> How to build and explore the current SysKit development preview.

SysKit is a Linux-first Go project. The current development preview is available
to build from source with system, process, network, port, live-monitoring, and
cgroup-derived container inspection commands.

## Current State

| Area | Status |
|---|---|
| Product vision | Defined |
| Architecture | Defined |
| Feature specifications | Defined for planned core features |
| Engineering standards | Defined |
| Production code | v0.5.0 released; v1 stabilization in development |
| Installable releases | Source build available; packages planned |

## Development Environment

Build and test on Linux with:

- Go 1.26.3 or newer.
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

The CLI exposes the v0.1 commands plus process, network, port, live-monitoring,
and cgroup-derived container commands. For example:

```sh
go run ./cmd/syskit system
go run ./cmd/syskit system --format json
go run ./cmd/syskit cpu
go run ./cmd/syskit memory --format json
go run ./cmd/syskit filesystem --show-pseudo
go run ./cmd/syskit disk --io --interval 1s
go run ./cmd/syskit process --user root --limit 20
go run ./cmd/syskit process --sort cpu --interval 1s --limit 20
go run ./cmd/syskit process tree
go run ./cmd/syskit network
go run ./cmd/syskit network interfaces --format json
go run ./cmd/syskit network routes
go run ./cmd/syskit network dns --format yaml
go run ./cmd/syskit ports --listening --pid 1234
go run ./cmd/syskit ports --address 127.0.0.1 --state listen
go run ./cmd/syskit dashboard --interval 2s
go run ./cmd/syskit dashboard --panel processes
go run ./cmd/syskit watch network --interval 2s
go run ./cmd/syskit top --sort memory --limit 20
go run ./cmd/syskit process --containers
go run ./cmd/syskit containers
go run ./cmd/syskit containers inspect <container-id>
go run ./cmd/syskit plugins list --plugin-dir ./plugins
go run ./cmd/syskit plugins inspect example --plugin-dir ./plugins
go run ./cmd/syskit plugins run example --plugin-dir ./plugins --timeout 5s
go run ./cmd/syskit diagnostics --severity warning
```

`dashboard` requires an interactive terminal; use the one-shot commands with
`--format json` or `--format yaml` when redirecting output.
It switches to a compact resize notice below 48×12 cells instead of overlapping
dashboard panels. Its overview derives CPU utilization and aggregate network
RX/TX throughput after the first refresh, and includes memory, swap, disk, and
the top memory process.
`watch` also requires an interactive terminal and refreshes a one-shot command
as a table until Ctrl-C.
`top` is an interactive process monitor; use `c`, `m`, `n`, or `p` to change
the CPU, memory, name, or PID sort respectively, and `j`/`k` to scroll rows.
`containers` is a best-effort, cgroup-derived listing: it reports recognized
runtime-style IDs, process counts, and available cgroup memory/CPU/I/O counters
without requiring runtime socket access.
`containers inspect` expands that mapping into its associated processes; it
does not claim runtime names or status, and unavailable cgroup controllers are
omitted rather than reported as zero.
`plugins list` discovers manifests from `--plugin-dir`, `SYSKIT_PLUGIN_DIR`,
or the XDG data path. Discovery never executes plugin code.

When permissions hide a process, structured process output sets `partial: true`
while retaining every process it could read.

`ports` reads TCP, UDP, IPv6, and Unix socket tables directly from procfs. It
best-effort maps socket inodes to process IDs and commands; inaccessible or
short-lived processes simply remain unmapped. JSON and YAML mark permission-
restricted scans with `owner_mapping_partial: true`. The remaining commands below include both shipped and planned contracts; check
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

See the [command reference](command-reference.md) for every implemented command,
interactive keybinding, and safety boundary.

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
