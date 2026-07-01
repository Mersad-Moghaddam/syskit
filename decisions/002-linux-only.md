# 002. Target Linux exclusively

**Status:** Accepted, 2026-07-01

---

## Context

SysKit's entire value proposition rests on reading data directly from native
kernel interfaces: `/proc`, `/sys`, Netlink sockets, and cgroups (see
[ADR 003](./003-native-apis-over-shell.md)). These interfaces are
**Linux-specific**. macOS exposes system data through `sysctl`, Mach APIs, and
IOKit; the BSDs have their own `sysctl` MIBs and `kvm` interfaces; Windows uses
the registry, WMI, and performance counters. There is no meaningful overlap in
the underlying data sources.

Supporting multiple operating systems would therefore require a platform
abstraction that hides these fundamentally different mechanisms behind a common
interface — the approach taken by libraries such as `gopsutil`. That abstraction
carries a real cost:

- Every collector must be implemented and tested once per platform.
- The abstraction must reduce data to a lowest common denominator, discarding the
  Linux-specific richness (cgroup hierarchies, Netlink routing detail, per-CPU
  `/proc/stat` fields) that makes the tool worth building.
- Depth of understanding — a core learning goal — is diluted across three or four
  operating systems instead of concentrated on one.

The constitution settles the direction in principle 1 (*Linux First*): "SysKit is
built exclusively for Linux ... There is no `runtime.GOOS` branching." This ADR
records and justifies that constraint at the decision level.

---

## Decision

We will **target Linux exclusively**. Concretely:

- No `runtime.GOOS` branching and no `//go:build` OS-conditional source files for
  non-Linux platforms.
- No cross-platform compatibility shims or lowest-common-denominator
  abstractions.
- The Platform Abstraction Layer (see [ADR 004](./004-layered-architecture.md))
  abstracts over *Linux kernel-version and cgroup-version variation*, not over
  operating systems.
- Code may freely assume a Linux filesystem layout and Linux system-call
  semantics.

Running SysKit on a non-Linux OS is unsupported and expected to fail cleanly, not
to degrade gracefully.

---

## Consequences

### Positive

- Collectors can use the full richness of Linux interfaces without flattening
  them to a portable subset.
- No per-platform implementation, testing, or CI matrix — engineering effort
  compounds on a single target.
- The codebase stays simpler: one filesystem layout, one set of syscall
  semantics, no conditional compilation to reason about.
- Concentrated depth directly serves the *Learn Before Build* goal.

### Negative

- macOS, BSD, and Windows users cannot run SysKit. This excludes a portion of
  potential users, including developers on macOS laptops.
- Reversing this decision later would be expensive: adding cross-platform support
  would require introducing the very abstraction layer we are choosing to omit.
- Contributors developing on non-Linux machines need a Linux VM, container, or
  remote host to run and test the tool.

### Neutral

- SysKit competes in the Linux tooling space specifically, not as a universal
  cross-platform monitor — a positioning choice, not merely a technical one.
- Go's easy cross-compilation still applies *within* Linux (multiple
  architectures), just not across operating systems.

---

## Alternatives Considered

- **Cross-platform via an abstraction library (e.g. `gopsutil`).** Would let
  SysKit run on macOS, Windows, and the BSDs by delegating platform specifics to
  a maintained library. Rejected because it forces the tool into `gopsutil`'s
  data model, hides the Linux internals we specifically want to expose and learn,
  adds a substantial dependency (conflicting with *Minimal Dependencies*), and
  turns the project into a thin wrapper rather than a direct reader of kernel
  interfaces.
- **Linux-first now, portable later.** Build for Linux but keep the door open by
  routing all access through an OS-neutral interface from day one. Rejected
  because a "portable-ready" interface is an abstraction tax paid up front for a
  capability that may never be exercised, and it would compromise the directness
  that distinguishes SysKit. If cross-platform support ever becomes a goal, it
  will be a deliberate, separately recorded decision — not a latent assumption.
- **Linux plus BSD (Unix-like only).** The BSDs share Unix conventions but not
  `/proc` semantics or Netlink. Rejected for the same reason as full
  cross-platform support, at a smaller scale.

---

## References

- [Constitution](../specs/constitution.md) — principle 1 (Linux First),
  principle 8 (Minimal Dependencies)
- [ADR 003](./003-native-apis-over-shell.md) — Read native kernel interfaces
- [ADR 004](./004-layered-architecture.md) — Platform Abstraction Layer scope
- [Architecture](../specs/architecture.md) — Linux Kernel Interfaces layer
- Linux kernel documentation: [proc(5)](https://man7.org/linux/man-pages/man5/proc.5.html),
  [sysfs](https://www.kernel.org/doc/html/latest/filesystems/sysfs.html)
