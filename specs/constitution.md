# SysKit Engineering Constitution

> The foundational engineering principles that guide every decision in the SysKit project.

---

## Purpose

This document defines the core principles that govern how SysKit is designed, built, and maintained. Every contributor, every pull request, and every architectural decision should be measured against these principles.

These are not aspirational guidelines — they are constraints. When two approaches are available and one aligns better with these principles, that is the one we choose.

---

## Principles

### 1. Linux First

SysKit is built exclusively for Linux. We do not introduce cross-platform abstraction layers, compatibility shims, or conditional compilation for other operating systems.

This constraint is deliberate. By targeting a single platform, we can leverage Linux-specific interfaces directly, avoid the complexity of platform abstraction, and produce a tool that is deeply integrated with the operating system it serves.

If a feature requires platform-specific behavior, it uses Linux APIs directly. There is no `runtime.GOOS` branching.

### 2. Native APIs First

SysKit reads system data from native Linux interfaces — `/proc`, `/sys`, Netlink, cgroups, and other kernel-provided mechanisms — rather than parsing the output of shell commands.

Shelling out to external utilities introduces dependencies on specific tool versions, requires parsing human-readable output that may change between distributions, and adds process overhead. Native interfaces provide structured, reliable, and efficient access to system data.

When a native interface is not available or practical for a specific data point, the exception is documented and justified.

### 3. Performance Matters

SysKit must be fast. Users expect CLI tools to respond instantly, and monitoring tools must minimize their own resource footprint to avoid distorting the measurements they report.

This means:

- Minimize memory allocations in hot paths
- Avoid unnecessary I/O and system calls
- Use buffered reading for file-based interfaces
- Benchmark critical paths and track performance regressions
- Prefer simple data structures over complex abstractions

Performance is not an afterthought — it is a design constraint from the beginning.

### 4. Keep It Modular

Every subsystem in SysKit is independent and self-contained. The CPU collector knows nothing about the network collector. The CLI layer knows nothing about how data is gathered.

Modularity enables:

- Independent development and testing of each subsystem
- Easy addition of new collectors without modifying existing ones
- Clear ownership and responsibility boundaries
- The ability to compose features without creating tight coupling

Each module exposes a clean interface and hides its implementation details.

### 5. Test Everything

Every component in SysKit has tests. Unit tests verify individual functions. Integration tests verify that collectors produce correct data on real Linux systems. Benchmarks track performance over time.

Testing is not optional and is not deferred to "later." A feature without tests is not complete.

Test coverage includes:

- Normal operation and expected outputs
- Edge cases and boundary conditions
- Error handling and failure modes
- Performance characteristics via benchmarks

### 6. Documentation First

SysKit follows a Specification-Driven Development workflow. Features begin as specifications, not as code. Architecture is designed before it is implemented. Decisions are documented before they are committed.

This principle applies at every level:

- Project-level decisions are captured in specs
- Public APIs are documented before implementation
- Complex algorithms include explanations of their approach
- The learning journal captures research and understanding

Documentation is a first-class deliverable, not a retroactive chore.

### 7. Clean Go

SysKit code is idiomatic Go. We follow the conventions established by the Go community, the standard library, and the official style guides.

This means:

- No unnecessary frameworks or abstractions
- Explicit error handling — no panics in library code
- Clear naming that communicates intent
- Small functions with single responsibilities
- Standard project layout and package organization
- `go vet`, `go fmt`, and static analysis pass without warnings

Go's simplicity is a feature. We do not fight it.

### 8. Minimal Dependencies

Every external dependency is a liability — it adds build complexity, potential security vulnerabilities, and maintenance burden. SysKit uses the Go standard library wherever possible.

Dependencies are introduced only when:

- The standard library does not provide the required functionality
- Building the equivalent from scratch would be unreasonable
- The dependency is well-maintained, widely used, and stable

When a dependency is added, the decision is documented with a clear justification.

### 9. Consistent CLI Experience

Every SysKit command follows the same patterns. Users should be able to predict how a new command behaves based on their experience with existing commands.

Consistency includes:

- Uniform flag naming and behavior across all commands
- Predictable output formatting (table, JSON, YAML)
- Clear, actionable error messages
- Consistent use of color, alignment, and layout
- A help system that is always accurate and complete

The CLI is the user interface. It deserves the same care as a graphical application.

### 10. Learn Before Build

SysKit is as much a learning project as it is a tool. Before implementing a feature, we invest time in understanding the underlying Linux subsystem, the relevant kernel interfaces, and the trade-offs involved.

This principle is reflected in the `learning/` directory, where research notes, kernel documentation references, and personal understanding are captured before implementation begins.

Building without understanding produces brittle code. Understanding before building produces code that is correct, efficient, and maintainable.

---

## Applying These Principles

When making a decision — whether it is an architectural choice, a code review comment, or a feature design — ask:

1. Does this follow the Linux-first approach?
2. Are we using native interfaces?
3. Is this the most performant approach we can reasonably achieve?
4. Does this maintain modularity?
5. Is this tested?
6. Is this documented?
7. Is this idiomatic Go?
8. Are we introducing unnecessary dependencies?
9. Is this consistent with the rest of the CLI?
10. Do we understand what we are building?

If the answer to any of these is "no," reconsider the approach.

---

*This constitution is a living document. It evolves as the project matures, but its core intent — building well-understood, well-tested, high-quality software — does not change.*
