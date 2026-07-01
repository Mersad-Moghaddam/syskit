# 001. Use Go as the implementation language

**Status:** Accepted, 2026-07-01

---

## Context

SysKit is a Linux system-inspection CLI that reads directly from kernel
interfaces (`/proc`, `/sys`, Netlink, cgroups) and presents the data through
both one-shot commands and, eventually, an interactive dashboard. The choice of
implementation language is the most consequential and least reversible decision
in the project, so it is recorded first.

The forces that bear on this decision:

- **Fast startup.** A system-inspection tool is invoked interactively and often
  scripted in loops. Interpreter or VM warm-up time is unacceptable; users expect
  the command to return "instantly" (constitution principle 3, *Performance
  Matters*).
- **Single static binary.** Distribution should be a single file with no runtime
  or dependency to install. This matters for the packaging goals in the
  [roadmap v1.0](../specs/roadmap.md) (deb, rpm, AUR, plain binary).
- **Strong concurrency.** The dashboard and `watch`/`top` modes
  ([roadmap v0.3](../specs/roadmap.md)) refresh multiple independent collectors
  on a timer; independent per-domain collectors map naturally onto lightweight
  concurrency.
- **A capable standard library for file and socket I/O.** The bulk of the work is
  reading pseudo-files and speaking Netlink over `AF_NETLINK` sockets. A strong
  stdlib for buffered file reading and raw sockets reduces the dependency
  surface (constitution principle 8, *Minimal Dependencies*).
- **Learning goal.** SysKit is also a reference project for high-quality
  engineering (constitution principle 10, *Learn Before Build*). The language
  should be productive enough that most learning effort goes into Linux internals
  rather than fighting the toolchain.

---

## Decision

We will implement SysKit in **Go, targeting version 1.22 or later**.

Go 1.22 is the baseline because it stabilises the revised `for` loop variable
semantics (eliminating a long-standing class of closure bugs) and includes
routing enhancements and math/rand improvements we are comfortable depending on.
We will track the toolchain via the `go` directive in `go.mod` and will not use
language features newer than the declared minimum.

---

## Consequences

### Positive

- Compiles to a single statically linked binary with no external runtime,
  satisfying the distribution goal directly.
- Sub-millisecond process startup keeps interactive and scripted use snappy.
- Goroutines and channels give us cheap, readable concurrency for parallel
  collectors and the real-time refresh pipeline.
- The standard library covers buffered file reading, raw sockets (needed for
  Netlink via `golang.org/x/sys/unix`), JSON/text encoding, and testing,
  keeping the dependency count low.
- Idiomatic Go is simple and widely known, lowering the contribution barrier and
  supporting the constitution's *Clean Go* principle.
- Cross-compilation is trivial (`GOARCH=arm64 go build`), useful for shipping
  multi-architecture Linux binaries.

### Negative

- The garbage collector introduces non-deterministic pauses. For a read-mostly
  inspection tool the impact is negligible, but it rules Go out of hard
  real-time use cases we are not targeting anyway.
- Go offers less low-level control than C or Rust; a few kernel interfaces are
  accessed through `golang.org/x/sys/unix` rather than pure stdlib.
- The runtime and GC add a fixed floor to binary size and memory footprint
  compared to C.
- Generics and error handling are more verbose than some alternatives, though
  this is consistent with Go's explicit-over-implicit philosophy.

### Neutral

- We commit to the Go release cadence and its backward-compatibility promise.
- `gofmt` removes formatting debates entirely — a constraint, not a preference.

---

## Alternatives Considered

- **Rust.** Excellent performance, no GC, single static binary, and a strong
  safety story. Rejected because its steeper learning curve and slower iteration
  would shift effort away from the primary learning goal (Linux internals), and
  the borrow checker's benefits are less pronounced for a read-only tool whose
  concurrency is coarse-grained. Rust remains the strongest runner-up.
- **C.** Maximum control and the most direct mapping to kernel interfaces, with
  effectively zero runtime overhead. Rejected for its memory-safety burden, weak
  standard library, manual dependency and build management, and the amount of
  boilerplate required for what should be ergonomic CLI code.
- **Python.** Fastest to prototype and has rich libraries. Rejected because
  interpreter startup and per-call overhead conflict with *Performance Matters*,
  it requires a runtime to be installed (no single static binary), and it hides
  rather than teaches the systems-level detail we want to learn.
- **Zig.** Compelling C-interop, no hidden control flow, and small binaries.
  Rejected on maturity grounds: the language and ecosystem are pre-1.0, the
  tooling and library support are thinner than Go's, and stability matters for a
  long-lived reference project.

---

## References

- [Constitution](../specs/constitution.md) — principles 3 (Performance Matters),
  7 (Clean Go), 8 (Minimal Dependencies), 10 (Learn Before Build)
- [Product Overview](../specs/product.md) — reference-project and performance goals
- [Roadmap](../specs/roadmap.md) — v0.3 real-time monitoring, v1.0 packaging
- [Effective Go](https://go.dev/doc/effective_go)
- [Go 1.22 release notes](https://go.dev/doc/go1.22)
