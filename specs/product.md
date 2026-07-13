# SysKit Product Overview

> A concise description of what SysKit is, who it serves, and where it is going.

---

## Vision

To become the go-to command-line toolkit for Linux system inspection, monitoring, and diagnostics — a tool that backend engineers and Linux users reach for every day because it is fast, reliable, and thoughtfully designed.

## Mission

Build an open-source, Linux-first CLI toolkit that provides a modern, unified interface for system information and diagnostics by interacting directly with native kernel interfaces, while also serving as a high-quality reference project for Go engineering and Linux internals.

## Problem Statement

Linux provides an extraordinary amount of system information through its virtual filesystems and kernel interfaces. However, accessing this information today typically requires:

- **Remembering dozens of disparate commands** — `top`, `htop`, `free`, `df`, `lsof`, `ss`, `ip`, `vmstat`, `iostat`, `sar`, each with its own flags and output format.
- **Parsing inconsistent output** — Different tools use different column layouts, units, and formatting conventions.
- **Combining tools with shell pipelines** — Getting a complete picture often requires chaining commands with `grep`, `awk`, and `sort`.
- **Lacking structured output** — Most traditional tools produce human-readable output that is difficult to consume programmatically.
- **No unified experience** — There is no single tool that provides a consistent, modern interface across system subsystems.

SysKit addresses these problems by providing a single, consistent, high-performance CLI that reads directly from native Linux interfaces and presents data in a uniform, structured format.

## Goals

1. **Provide a unified CLI** for inspecting CPU, memory, disk, network, processes, ports, and filesystem information.
2. **Read from native Linux interfaces** (`/proc`, `/sys`, Netlink) instead of wrapping shell commands.
3. **Deliver structured output** in multiple formats (table, JSON, YAML) for both human and programmatic consumption.
4. **Offer real-time monitoring** through an interactive terminal dashboard.
5. **Maintain excellent performance** with fast startup, low memory footprint, and minimal system impact.
6. **Provide an extensible architecture** that supports plugins and custom collectors.
7. **Serve as an educational resource** for Go programming, Linux internals, and CLI design.

## Non-Goals

- **Cross-platform support** — SysKit targets Linux exclusively. macOS, Windows, and BSD are out of scope.
- **Replacing specialized tools** — SysKit does not aim to replace deep-dive tools like `perf`, `strace`, or `bpftrace`. It covers common inspection and monitoring tasks.
- **System administration** — SysKit is a read-only inspection tool. It does not modify system configuration, manage services, or perform administrative actions.
- **Cloud provider integration** — SysKit operates at the OS level. Cloud-specific metadata and APIs are out of scope for the core tool.
- **GUI application** — SysKit is a terminal application. There are no plans for a graphical interface.

## Target Users

### Primary

- **Backend Engineers** — Developers who work with Linux servers daily and need quick access to system state during development, debugging, and incident response.
- **DevOps / SRE Engineers** — Operations professionals who manage Linux infrastructure and need a fast, reliable inspection tool.

### Secondary

- **Linux Enthusiasts** — Users who are curious about their system and want a modern tool for exploration.
- **Students and Learners** — People studying Linux internals, Go programming, or systems engineering who can use SysKit as both a tool and a learning resource.

## Core Features

### System Information
Host details, kernel version, OS release, uptime, load averages, and boot time.

### CPU Monitoring
Core count, architecture, per-core utilization, frequency scaling, CPU time breakdown, and cache information.

### Memory Analysis
Physical and swap memory usage, buffer/cache breakdown, memory pressure indicators, and per-process memory consumption.

### Disk Inspection
Partition layout, filesystem usage, mount points, I/O statistics, and disk health indicators.

### Process Management
Process listing with filtering, process tree visualization, and resource usage per process. SysKit remains read-only and does not send signals.

### Network Monitoring
Interface statistics, active connections, routing tables, DNS configuration, and bandwidth utilization.

### Port Inspection
Listening ports, socket states, associated processes, and protocol filtering.

### Interactive Dashboard
Real-time terminal UI with customizable layouts, live-updating metrics, and keyboard navigation.

### Structured Output
All commands support table, JSON, and YAML output formats via a consistent `--format` flag.

## Long-Term Vision

As SysKit matures, it will expand to include:

- **Container Inspection** — Docker, Podman, and OCI container runtime visibility.
- **Plugin System** — User-defined collectors and custom data sources.
- **Remote Monitoring** — Inspect remote systems over SSH.
- **Watch Mode** — Continuous monitoring with configurable refresh intervals.
- **Health Checks** — Automated system health assessment with configurable thresholds.
- **Historical Data** — Local storage of metrics for trend analysis.

The long-term goal is a tool that is as indispensable to a Linux engineer's workflow as `git` is to a developer's — always available, always reliable, always fast.

---

*This document is maintained alongside the project and updated as the product direction evolves.*
