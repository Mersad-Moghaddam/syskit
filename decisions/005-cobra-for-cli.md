# 005. Use Cobra for the CLI framework

**Status:** Accepted, 2026-07-01

---

## Context

SysKit is, first and foremost, a command-line tool, and the constitution treats
the CLI as the product's user interface (principle 9, *Consistent CLI
Experience*). The tool's surface is a tree of subcommands — `syskit system`,
`syskit cpu`, `syskit memory`, `syskit disk`, later `syskit process tree`,
`syskit network`, `syskit ports`, `syskit dashboard`, `syskit top`
([roadmap](../specs/roadmap.md)) — each with shared and command-specific flags
(`--format json`, filters, sort keys) and consistent, accurate help.

Building this on the standard library's `flag` package alone would mean writing,
by hand, a subcommand dispatcher, nested flag scoping, help/usage generation,
and — if we want it — shell completion. That is a meaningful amount of
undifferentiated plumbing to build and, more importantly, to keep consistent as
the command tree grows. The constitution's *Minimal Dependencies* principle
(principle 8) permits a dependency precisely when "the standard library does not
provide the required functionality" and "building the equivalent from scratch
would be unreasonable" — CLI command routing is a textbook case.

The requirements for the framework:

- Nested subcommands with per-command flags and inherited/persistent flags.
- Automatic, always-accurate help and usage text (supports principle 9).
- Shell completion (bash/zsh/fish) as a low-cost usability win.
- POSIX-style flag parsing with sensible defaults.
- Maturity and wide adoption, so the dependency is stable and well-understood.

---

## Decision

We will use **Cobra** (`github.com/spf13/cobra`) as the CLI framework.

Cobra provides the subcommand tree, persistent and local flag handling
(via `pflag`, giving POSIX/GNU-style flags), auto-generated help and usage, and
generated shell completion. It is the de facto standard for Go CLIs (kubectl,
Hugo, GitHub CLI, Docker CLI components), which makes it a low-risk, well-
maintained choice.

Cobra is confined to the **CLI and Command layers**
([ADR 004](./004-layered-architecture.md)). Services, collectors, and the
platform layer have no knowledge of Cobra; commands translate parsed input into
plain service calls. This keeps the framework replaceable and the lower layers
independently testable.

---

## Consequences

### Positive

- Eliminates hand-written subcommand routing, flag scoping, and help generation —
  effort goes into collectors and services instead of CLI plumbing.
- Auto-generated help and shell completion directly serve the *Consistent CLI
  Experience* principle: help is always accurate because it is derived from the
  command definitions.
- Persistent flags (e.g. a global `--format`) are defined once and inherited,
  keeping flag behaviour uniform across commands.
- Mature and ubiquitous, so contributors likely already know it and the
  dependency is stable and actively maintained.

### Negative

- Introduces a dependency (and its transitive `pflag`), a cost we accept under
  the documented *Minimal Dependencies* exception.
- Cobra is opinionated and comparatively heavyweight; a very small tool could get
  by with less. SysKit's planned command tree justifies it.
- Some Cobra idioms (e.g. `init()`-based command registration) can encourage
  patterns that need discipline to keep tidy and testable.

### Neutral

- We adopt Cobra's conventions for command structure and help formatting.
- Because Cobra is isolated to the top two layers, replacing it later would be
  contained — a deliberate consequence of the layered architecture.

---

## Alternatives Considered

- **Standard library `flag` only.** Zero dependencies and full control. Rejected
  because it has no native subcommand model; we would hand-build dispatch, nested
  flag scoping, help generation, and completion — reimplementing Cobra, poorly,
  and diverting effort from the tool's actual purpose.
- **`urfave/cli`.** A capable, popular alternative with a lighter feel. A
  reasonable choice, but rejected because Cobra's ecosystem, adoption, and
  completion/help maturity are stronger, and its subcommand model maps more
  cleanly onto SysKit's deep command tree.
- **`alecthomas/kong`.** Elegant struct-tag-driven definitions and strong typing.
  Rejected on adoption and familiarity grounds relative to Cobra; the struct-tag
  approach is appealing but less battle-tested at the scale and ubiquity of Cobra
  for a long-lived reference project.
- **Build a minimal in-house command router.** Keeps dependencies at zero.
  Rejected as undifferentiated work that the *Minimal Dependencies* policy
  explicitly says should be delegated to a well-maintained dependency when
  building it ourselves would be unreasonable.

---

## References

- [Constitution](../specs/constitution.md) — principle 8 (Minimal Dependencies,
  the dependency-justification policy), principle 9 (Consistent CLI Experience)
- [Architecture](../specs/architecture.md) — CLI and Command layers
- [ADR 004](./004-layered-architecture.md) — Cobra confined to the top two layers
- [Roadmap](../specs/roadmap.md) — v0.1 "CLI framework with Cobra"
- [Cobra documentation](https://cobra.dev/)
- [spf13/pflag](https://github.com/spf13/pflag) — POSIX/GNU flag parsing
