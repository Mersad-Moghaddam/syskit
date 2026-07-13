# SysKit

> A modern, Linux-first command-line toolkit for system inspection, resource monitoring, and diagnostics — built with Go.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26.3+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Linux](https://img.shields.io/badge/Platform-Linux-FCC624?logo=linux&logoColor=black)](https://kernel.org)
[![Status](https://img.shields.io/badge/Status-v0.2.0%20released-brightgreen)]()

---

## Overview

SysKit is an open-source command-line toolkit designed for backend engineers, DevOps professionals, and Linux enthusiasts who need fast, reliable, and consistent access to system information.

Rather than wrapping existing Linux utilities, SysKit interacts directly with native Linux interfaces — `/proc`, `/sys`, Netlink, and other kernel APIs — to collect and present system data. This approach provides better performance, richer detail, and a deeper understanding of the underlying operating system.

SysKit is both a practical daily-use tool and a long-term educational project for mastering Go, Linux internals, CLI development, and systems programming.

## Philosophy

SysKit follows a **Specification-Driven Development (SDD)** workflow. Implementation comes after documentation and architecture. Every feature begins as a specification, is reviewed for correctness and consistency, and only then moves into code.

This approach ensures:

- Deliberate, well-understood design decisions
- A codebase that remains maintainable as it grows
- Documentation that stays in sync with the implementation
- A project that serves as both a tool and a learning resource

## Goals

- **Unified Interface** — Provide a single, consistent CLI for common Linux inspection and monitoring tasks.
- **Native Data Collection** — Read directly from kernel interfaces instead of parsing shell command output.
- **Performance** — Deliver fast startup, low memory footprint, and minimal overhead.
- **Modularity** — Build an architecture that supports independent, composable subsystems.
- **Extensibility** — Support plugins, custom collectors, and multiple output formats.
- **Education** — Serve as a reference project for Go engineering, Linux internals, and CLI design.

## Planned Features

| Category | Features |
|---|---|
| **System** | Host information, kernel version, uptime, load averages |
| **CPU** | Core count, utilization, frequency, per-core statistics |
| **Memory** | Physical/swap usage, buffers, caches, memory pressure |
| **Disk** | Partition layout, usage, I/O statistics, mount points |
| **Process** | Process listing, tree view, resource usage, signals |
| **Network** | Interface statistics, connections, routing, DNS |
| **Ports** | Listening ports, socket states, associated processes |
| **Filesystem** | Inode usage, filesystem types, mount options |
| **Diagnostics** | System health checks, resource bottleneck detection |
| **Dashboard** | Interactive terminal UI with real-time monitoring |
| **Output** | Table, JSON, YAML, and plain-text output formats |
| **Plugins** | User-defined collectors and custom extensions |
| **Containers** | Docker and container runtime inspection |

## Design Principles

- **Linux First** — Built exclusively for Linux. No cross-platform abstraction layers.
- **Native APIs First** — Prefer `/proc`, `/sys`, and Netlink over shelling out to external commands.
- **Performance Matters** — Minimize allocations, avoid unnecessary work, benchmark critical paths.
- **Keep It Modular** — Each subsystem is independent and self-contained.
- **Test Everything** — Unit tests, integration tests, and benchmarks for every component.
- **Documentation First** — Specs before code. Every feature is designed before it is built.
- **Clean Go** — Idiomatic Go. No frameworks, no magic, no unnecessary abstractions.
- **Minimal Dependencies** — Rely on the standard library wherever possible.
- **Consistent CLI Experience** — Predictable flags, uniform output, clear error messages.

## Technology Stack

| Component | Technology |
|---|---|
| Language | [Go](https://go.dev) 1.26.3+ |
| Data Sources | `/proc`, `/sys`, Netlink, kernel APIs |
| CLI Framework | [Cobra](https://github.com/spf13/cobra) |
| Terminal UI | [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| Styling | [Lip Gloss](https://github.com/charmbracelet/lipgloss) |
| Testing | Go standard `testing` package, [testify](https://github.com/stretchr/testify) |

## Project Status

**v0.2.0 released**

SysKit v0.2.0 adds process listing and trees, network interfaces/routes/DNS,
and port inspection with table, JSON, and YAML output. The layered
architecture, fixture-backed collectors, golden output contracts, and Linux
integration coverage are in place for later milestones.

See the [Roadmap](specs/roadmap.md) for planned milestones.

## Project Structure

```text
syskit/
├── cmd/syskit/         # CLI entry point (main)
├── internal/           # Application internals (not importable externally)
│   ├── cli/            # Cobra command wiring, config, logger, exit mapping
│   │   └── command/    # One file per subcommand
│   ├── collector/      # Built-in domain collectors
│   ├── platform/       # Linux procfs, sysfs, Netlink, cgroup adapters (SysFS)
│   ├── render/         # Table, JSON, YAML, TUI rendering
│   ├── service/        # Aggregation and domain logic
│   └── model/          # Shared typed domain structs
├── testdata/           # Shared fixtures
├── .github/            # GitHub templates and CI workflows
├── docs/               # User-facing documentation and maintainer guides
├── specs/              # Specifications and architecture documents
│   ├── constitution.md # Engineering principles
│   ├── product.md      # Product overview
│   ├── roadmap.md      # Development milestones
│   ├── architecture.md # System architecture
│   └── features/       # Individual feature specifications
├── learning/           # Study notes on Linux internals
├── standards/          # Engineering standards and review policies
├── decisions/          # Architecture Decision Records
├── scripts/            # Development and build scripts
├── LICENSE
├── README.md
└── .gitignore
```

The Go module (`github.com/Mersad-Moghaddam/syskit`) targets Go 1.26.3+ and builds
to a single static binary. Dependencies flow strictly downward
(CLI → Command → Service → Collector → Platform → kernel); lower layers never
import higher ones.

## Documentation Map

- [Getting started](docs/getting-started.md)
- [Command reference](docs/command-reference.md)
- [Architecture overview](docs/architecture.md)
- [Developer onboarding](docs/developer-onboarding.md)
- [Product overview](specs/product.md)
- [Feature specifications](specs/features/)
- [Collector architecture](specs/collectors.md)
- [Rendering architecture](specs/rendering.md)
- [Plugin architecture](specs/plugin-architecture.md)
- [Learning roadmap](learning/roadmap.md)
- [Implementation readiness](docs/implementation-readiness.md)

## Vision

SysKit aims to become a reliable daily companion for backend engineers and Linux users — a tool that provides a modern, fast, and thoughtful command-line experience for system inspection, monitoring, and diagnostics.

Beyond the tool itself, SysKit is designed to be a living reference for building well-architected Go applications that interact deeply with the Linux kernel.

## Contributing

SysKit is in its implementation phase. Contributions are
welcome across code, specs, architecture, Linux explanations, documentation, and
repository process. Production code must follow the architecture boundaries and
meet the [Definition of Done](standards/definition-of-done.md). See
[docs/contributing.md](docs/contributing.md).

## License

This project is licensed under the [MIT License](LICENSE).

---

*Built with care for Linux, Go, and the terminal.*
