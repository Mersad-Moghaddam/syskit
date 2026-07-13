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
| Production code | v1.0.0 stable release |
| Installable releases | Archives, deb/rpm packages, AUR metadata, and checksums |

## Development Environment

Build and test on Linux with:

- Go 1.26.3 or newer.
- A Linux kernel with procfs and sysfs mounted in the standard locations.
- Standard development tools such as `git`, `make` or shell, and a POSIX-compatible terminal.
- Normal user permissions for inspection commands. Commands that read restricted process details may report partial data unless run with elevated permissions.

The project is Linux-first by design. macOS, Windows, and BSD support are out of scope for the core tool.

## Install a Tagged Release

Each tagged release publishes static Linux amd64 and arm64 archives, Debian and
RPM packages, AUR build metadata, and `SHA256SUMS`. Verify a downloaded artifact
before installing it:

```sh
sha256sum -c SHA256SUMS --ignore-missing
```

For a portable archive, extract the file matching your architecture and install
the versioned binary as `syskit`:

```sh
tar -xzf syskit_VERSION_linux_amd64.tar.gz
sudo install -m 0755 syskit_VERSION_linux_amd64 /usr/local/bin/syskit
```

On Debian-family systems, install the matching `.deb` with
`sudo apt install ./syskit_VERSION_amd64.deb`. On RPM-based systems, use
`sudo dnf install ./syskit-VERSION-1.x86_64.rpm`. The AUR archive contains the
reviewable `PKGBUILD` and `.SRCINFO` for the `syskit-bin` package; extract it and
run `makepkg -si` on Arch Linux. Replace `amd64`/`x86_64` with
`arm64`/`aarch64` where appropriate.

## Build and Test from Source

```sh
go build ./...
go test -race ./...
go run ./cmd/syskit --help
```

The CLI exposes the v0.1 commands plus process, network, port, live-monitoring,
and cgroup-derived container commands. For example:

```sh
go run ./cmd/syskit
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

The bare command opens SysKit's interactive control center when run in a real
terminal. The responsive wordmark animates once on entry (any input skips it),
and each option's icon and accent continue into its loading and output view.
Use the keyboard or mouse to browse domain submenus, Enter to select, and Escape
or Left to return. Result screens support vertical and horizontal scrolling;
the menu reopens at the previous location after completion. Redirected bare
output continues to show standard Cobra help.

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
restricted scans with `owner_mapping_partial: true`. Check `--help` for the
current supported flags.

## How to Explore This Repository

Start with the documents in this order:

1. [README](../README.md) for the project vision and feature map.
2. [Product overview](../specs/product.md) for goals, users, and non-goals.
3. [Architecture overview](architecture.md) for the system shape.
4. [Feature specifications](../specs/features/) for expected behavior.
5. [Learning Center](../learning/README.md) for the complete Linux, Go, and
   SysKit engineering course; use its [roadmap](../learning/roadmap.md) to choose
   a starting stage.
6. [Implementation readiness record](implementation-readiness.md) for the design-to-code transition.

## Command Roadmap

All commands in the [command reference](command-reference.md) are implemented.

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
