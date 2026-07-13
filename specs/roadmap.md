# SysKit Roadmap

> Development milestones and planned feature progression.

---

## Overview

SysKit follows a milestone-based development approach. Each milestone represents a coherent set of features that are designed, implemented, tested, and documented together before moving to the next.

Milestones are sequential — each builds on the foundation established by the previous ones.

---

## v0.1 — Foundation

*Core infrastructure and basic system inspection.*

**Goals:**
- Establish project structure, build system, and CI pipeline
- Implement the CLI framework and command routing
- Build the collector abstraction layer
- Deliver the first set of system inspection commands

**Features:**
- [x] `syskit system` — Host information, kernel version, OS release, uptime, load averages
- [x] `syskit cpu` — Core count, architecture, model, frequency, cache info
- [x] `syskit memory` — Physical/swap usage, buffers, caches, available memory
- [x] `syskit disk` — Partition layout, filesystem usage, mount points

**Technical:**
- [x] CLI framework with Cobra
- [x] Collector interface and platform abstraction layer
- [x] Table and JSON output formatters
- [x] Unit test framework and CI integration
- [x] Error handling patterns and conventions

---

## v0.2 — Processes & Networking

*Process inspection and network visibility.*

**Goals:**
- Implement process-related commands
- Add network and port inspection
- Introduce filtering and sorting across commands

**Features:**
- [x] `syskit process` — Process listing, filtering by name/PID/user, resource usage
- [x] `syskit process tree` — Process tree visualization
- [x] `syskit network` — Interface statistics, addresses, routes, and DNS configuration
- [x] `syskit ports` — Listening ports, socket states, associated processes

**Technical:**
- [x] Netlink integration for network data
- [x] `/proc/[pid]` parsing for process data
- [x] Filtering and sorting framework
- [x] YAML output formatter

---

## v0.3 — Real-Time Monitoring

*Interactive dashboard and live monitoring capabilities.*

**Goals:**
- Build the interactive terminal UI
- Add real-time data refresh and live monitoring
- Implement top-like process monitoring

**Features:**
- [ ] `syskit dashboard` — Interactive terminal dashboard with real-time metrics
- [ ] `syskit watch <command>` — Continuous monitoring with configurable refresh interval
- [ ] `syskit top` — Interactive process monitor with sorting and filtering

**Technical:**
- [ ] Bubble Tea integration for terminal UI
- [ ] Lip Gloss styling system
- [ ] Real-time data refresh pipeline
- [ ] Keyboard navigation and interaction model
- [ ] Layout system for dashboard widgets

---

## v0.4 — Containers

*Container runtime inspection and visibility.*

**Goals:**
- Add Docker and container runtime support
- Provide container-aware process and resource views
- Support cgroup-based resource monitoring

**Features:**
- [ ] `syskit docker` — Container listing, resource usage, status
- [ ] `syskit docker inspect <id>` — Detailed container inspection
- [ ] Container-aware process views
- [ ] Cgroup resource monitoring

**Technical:**
- [ ] Docker API client integration
- [ ] Cgroup v1/v2 parsing
- [ ] Container-to-process mapping

---

## v0.5 — Extensibility

*Plugin system and custom collectors.*

**Goals:**
- Design and implement the plugin architecture
- Support user-defined collectors
- Enable community extensions

**Features:**
- [ ] Plugin discovery and loading
- [ ] Plugin API and SDK
- [ ] Custom collector registration
- [ ] Plugin configuration system

**Technical:**
- [ ] Plugin interface definition
- [ ] Dynamic loading mechanism
- [ ] Plugin isolation and security model
- [ ] Plugin documentation and examples

---

## v1.0 — Stable Release

*Production-ready release with complete documentation and stability guarantees.*

**Goals:**
- Stabilize all public APIs and CLI interfaces
- Complete comprehensive documentation
- Achieve thorough test coverage across all subsystems
- Performance optimization and benchmarking
- Community-ready release process

**Deliverables:**
- [ ] Stable CLI interface with semantic versioning
- [ ] Complete user documentation
- [ ] Installation packages (deb, rpm, AUR, binary releases)
- [ ] Performance benchmarks and optimization
- [ ] Contributing guide and community guidelines
- [ ] Release automation and changelog generation

---

## Future Considerations

The following features are under consideration for post-1.0 development:

- **Remote Monitoring** — Inspect remote systems over SSH
- **Historical Data** — Local metric storage and trend analysis
- **Health Checks** — Automated system health assessment with configurable thresholds
- **Alerting** — Threshold-based notifications for monitored metrics
- **Filesystem Deep Dive** — Inode analysis, large file discovery, directory size breakdown
- **Hardware Information** — PCI devices, USB devices, DMI/SMBIOS data
- **Kubernetes Integration** — Pod and node inspection for Kubernetes clusters

---

*This roadmap is a living document. Priorities and scope may adjust as the project evolves and the community grows.*
