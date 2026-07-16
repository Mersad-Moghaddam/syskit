# SysKit

> **Native Linux intelligence for the terminal.**
>
> Inspect live system state with one fast, read-only CLI.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26.3+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Linux](https://img.shields.io/badge/Platform-Linux-FCC624?logo=linux&logoColor=black)](https://kernel.org)
[![Status](https://img.shields.io/badge/Status-v1.0.0%20stable-brightgreen)](docs/releases/v1.0.0.md)

<p align="center">
  <a href="#install">Install</a> ·
  <a href="#what-you-get">Features</a> ·
  <a href="docs/command-reference.md">Commands</a> ·
  <a href="docs/getting-started.md">Documentation</a>
</p>

```text
███████╗██╗   ██╗███████╗██╗  ██╗██╗████████╗
██╔════╝╚██╗ ██╔╝██╔════╝██║ ██╔╝██║╚══██╔══╝
███████╗ ╚████╔╝ ███████╗█████╔╝ ██║   ██║
╚════██║  ╚██╔╝  ╚════██║██╔═██╗ ██║   ██║
███████║   ██║   ███████║██║  ██╗██║   ██║
╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝   ╚═╝

● SYSKIT // CONTROL CENTER  native Linux intelligence  •  read-only  •  zero shell-outs
```

SysKit reads `/proc`, `/sys`, Netlink, and cgroups directly—never the
human-formatted output of system utilities. It is built for Linux engineers who
need fast answers without changing the host they are investigating.

## Install

Install the latest stable release on Linux with one command:

```sh
curl -fsSL https://raw.githubusercontent.com/Mersad-Moghaddam/syskit/main/scripts/install.sh | sh
```

The installer detects amd64/arm64, verifies the release SHA-256 checksum, and
installs the binary and manual page under `/usr/local` (using `sudo` only when
needed).

```sh
syskit version
syskit system
```

Need a package for Debian/Ubuntu, Fedora/RHEL, Arch, a pinned version, a
user-local install, or source build? See the [complete installation guide](docs/getting-started.md#install-a-tagged-release).

## What you get

| Explore | Monitor | Automate |
|---|---|---|
| System, CPU, memory, disk, filesystems | Interactive dashboard, `top`, and `watch` | Stable table, JSON, and YAML output |
| Processes, ports, network, and DNS | Live resource and process views | Shell completion and TOML configuration |
| Cgroup-derived containers and diagnostics | Keyboard and mouse control center | Explicit, bounded external plugins |

```sh
# Human-friendly inspection
syskit process --sort cpu --limit 20

# Script-friendly output
syskit memory --format json

# Live terminal dashboard
syskit dashboard --interval 2s
```

## See it in action

<p align="center">
  <img src="assets/readme/syskit-terminal-preview.svg" alt="SysKit control center and system inspection in a terminal" width="960">
</p>

Run bare `syskit` from an interactive terminal to open the control center.
Use `syskit <command> --help` to discover flags, or go straight to the
[command reference](docs/command-reference.md).

## Designed for Linux, built with care

- **Linux-native:** direct kernel interfaces; no cross-platform abstraction layer.
- **Read-only by default:** inspect and diagnose without mutating host state.
- **Fast and predictable:** one static Go binary, low overhead, stable v1 CLI contracts.
- **Terminal-first:** beautiful interactive views when a TTY is available; clean structured output when it is not.

SysKit supports Linux `amd64`/`x86_64` and `arm64`/`aarch64`. Some process and
socket ownership data may be permission-restricted; commands report partial
results rather than inventing missing data.

## Learn and contribute

- New user? Start with [Getting started](docs/getting-started.md).
- Need every flag and command? Read the [command reference](docs/command-reference.md).
- Curious how it works? Explore the [architecture](docs/architecture.md),
  [feature specifications](specs/features/), and [Learning Center](learning/README.md).
- Want to contribute? Read [Contributing](docs/contributing.md) and the
  [engineering constitution](specs/constitution.md).

SysKit is released under the [MIT License](LICENSE).
