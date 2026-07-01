# SysKit Architecture

> The planned system architecture for SysKit, describing each layer, its responsibilities, and how data flows through the system.

---

## Overview

SysKit uses a layered architecture that separates concerns cleanly between user interaction, business logic, data collection, and platform-specific access. Each layer communicates only with its immediate neighbors, and dependencies flow strictly downward.

```text
┌─────────────────────────────────┐
│             CLI Layer           │   ← User interaction, flags, output formatting
├─────────────────────────────────┤
│          Command Layer          │   ← Command definitions, input validation
├─────────────────────────────────┤
│          Service Layer          │   ← Business logic, data aggregation
├─────────────────────────────────┤
│         Collector Layer         │   ← Data collection from system interfaces
├─────────────────────────────────┤
│    Platform Abstraction Layer   │   ← OS-specific interface adapters
├─────────────────────────────────┤
│    Linux Kernel Interfaces      │   ← /proc, /sys, Netlink, cgroups
└─────────────────────────────────┘
```

---

## Layers

### CLI Layer

**Responsibility:** Handle all user-facing interaction — argument parsing, flag processing, help text, output formatting, and terminal rendering.

The CLI layer is the entry point for every user interaction. It receives raw input from the user, delegates to the appropriate command, and formats the result for display.

**Key concerns:**
- Parse command-line arguments and flags using Cobra
- Validate user input before passing it to commands
- Select the appropriate output formatter (table, JSON, YAML)
- Render output to the terminal, including color and alignment
- Handle the interactive terminal UI (Bubble Tea) for dashboard mode
- Display errors in a consistent, user-friendly format

The CLI layer knows how to present data but has no knowledge of how that data is collected.

### Command Layer

**Responsibility:** Define the available commands, map user intent to service calls, and coordinate the execution of each command.

Each command represents a specific user action — inspecting CPU information, listing processes, showing disk usage. Commands translate parsed CLI input into service method calls and return structured results to the CLI layer for formatting.

**Key concerns:**
- Define command structure, subcommands, and flags
- Validate command-specific input and flag combinations
- Call the appropriate service methods
- Return structured data to the CLI layer
- Handle command-specific error cases

Commands are thin — they coordinate but do not contain business logic or data collection code.

### Service Layer

**Responsibility:** Implement business logic, aggregate data from multiple collectors, and prepare structured results for the command layer.

The service layer is where data from different collectors is combined, filtered, sorted, and transformed into the shapes that commands need. A service might combine CPU utilization data with process data to show per-process CPU usage.

**Key concerns:**
- Aggregate data from one or more collectors
- Apply filtering, sorting, and transformation logic
- Compute derived metrics (e.g., percentages, rates, deltas)
- Define the data structures that commands consume
- Handle cross-cutting concerns like caching and rate limiting

Services depend on collectors but are independent of the CLI and output format.

### Collector Layer

**Responsibility:** Gather raw system data from platform abstractions and return it in normalized, typed data structures.

Each collector is responsible for a single domain — CPU, memory, disk, network, processes. Collectors read raw data from platform abstractions and parse it into well-defined Go structs.

**Key concerns:**
- Read raw data from platform abstraction interfaces
- Parse and normalize data into typed structures
- Handle parsing errors and missing data gracefully
- Provide consistent data shapes regardless of kernel version variations
- Expose a clean, testable interface

Collectors are independent of each other. The CPU collector does not depend on the memory collector. This independence enables parallel development and isolated testing.

### Platform Abstraction Layer

**Responsibility:** Provide a stable interface over OS-specific data sources, isolating the rest of the application from the details of how data is read from the kernel.

This layer wraps the specifics of reading from `/proc`, `/sys`, Netlink sockets, and other Linux interfaces behind clean Go interfaces. If a kernel interface changes format between versions, the change is contained within this layer.

**Key concerns:**
- Abstract file reading from `/proc` and `/sys`
- Manage Netlink socket connections and message parsing
- Handle cgroup v1 and v2 differences
- Provide consistent interfaces despite kernel version variations
- Enable testing through interface-based design (mock/stub in tests)

The platform abstraction layer is the only layer that directly interacts with the operating system.

### Linux Kernel Interfaces

**Responsibility:** The underlying data sources provided by the Linux kernel.

These are not code that SysKit writes — they are the system interfaces that SysKit reads from:

- **`/proc`** — Process information, CPU statistics, memory information, network statistics, and more. A virtual filesystem maintained by the kernel.
- **`/sys`** — Device and driver information, hardware topology, power management, and kernel parameters.
- **Netlink** — A socket-based interface for communication between kernel and userspace. Used for network configuration, routing, and device management.
- **Cgroups** — Control group interfaces for resource management and container inspection (v1 at `/sys/fs/cgroup`, v2 at unified hierarchy).

---

## Data Flow

A typical request flows through the system as follows:

1. **User** runs `syskit cpu --format json`
2. **CLI Layer** parses the command (`cpu`) and flags (`--format json`), invokes the CPU command
3. **Command Layer** validates input, calls the CPU service
4. **Service Layer** requests data from the CPU collector, computes derived metrics
5. **Collector Layer** calls platform abstractions to read raw CPU data
6. **Platform Abstraction** reads from `/proc/cpuinfo`, `/proc/stat`, and `/sys/devices/system/cpu/`
7. **Data returns** upward through each layer: raw bytes → parsed structs → enriched models → formatted output

Each layer transforms the data as it passes through — from raw bytes to user-facing output — with each transformation adding structure, meaning, and presentation.

---

## Design Decisions

### Why Layers?

A layered architecture enforces separation of concerns. The CLI can be redesigned without touching data collection. Collectors can be rewritten without affecting business logic. New output formats can be added without modifying any service code.

### Why Not a Flat Architecture?

A flat architecture (commands directly reading from `/proc`) is simpler initially but becomes unmaintainable as the project grows. Layering adds a small amount of indirection in exchange for significant gains in testability, modularity, and long-term maintainability.

### Why Interfaces Between Layers?

Go interfaces between layers enable testing with mocks and stubs, allow alternative implementations (e.g., a mock platform layer for CI environments), and make dependencies explicit and inspectable.

### Why Independent Collectors?

Independent collectors can be developed, tested, and maintained by different contributors without coordination. They can be loaded on demand, and new collectors can be added without modifying existing ones. This is essential for the plugin system planned in v0.5.

---

*This architecture document describes the planned design. Implementation details will be refined as development progresses, but the layered structure and separation of concerns are foundational decisions that will not change.*
