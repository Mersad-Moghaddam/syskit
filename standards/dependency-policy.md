# Dependency Policy

> How external dependencies are evaluated, approved, and maintained.

---

## Purpose

SysKit treats every external dependency as a liability. This standard operationalizes the **Minimal Dependencies** principle from `../specs/constitution.md`: prefer the standard library, add a dependency only when it clearly earns its place, and record the decision.

The Go standard library, together with direct reads of `/proc`, `/sys`, Netlink, and cgroups, covers the overwhelming majority of SysKit's needs. A dependency is the exception, not the default.

---

## Standard Library First

Before proposing any dependency, confirm the standard library cannot do the job at acceptable cost. In particular:

- Data collection reads kernel interfaces directly — **never** shell out to external binaries and **never** add a library that wraps them.
- Parsing, formatting (JSON/YAML via stdlib where feasible), and I/O use stdlib.
- If building the equivalent in-repo is small and well-understood, build it.

---

## Evaluation Criteria

Every proposed dependency is judged against all of the following:

| Criterion | Question |
|---|---|
| **Necessity** | Does the stdlib genuinely not cover this? Is the cost of building it ourselves unreasonable? |
| **License** | Is it compatible with SysKit's **MIT** license? (MIT, BSD, Apache-2.0 are compatible; GPL/AGPL are not.) |
| **Maintenance** | Is it actively maintained, with recent releases and responsive maintainers? |
| **Popularity** | Is it widely used and battle-tested, or niche and unproven? |
| **Transitive deps** | How large is its dependency tree? Each transitive dep is also our liability. |
| **Security** | Any known advisories? Does `govulncheck` come back clean? |
| **Footprint** | Does it meaningfully increase binary size or build time? |

If a candidate fails any single criterion, it is rejected or the concern is resolved before adoption.

---

## Currently Approved Dependencies

These are the only external dependencies sanctioned for SysKit. Each maps to a capability the stdlib does not provide.

| Dependency | Purpose | License | Justification |
|---|---|---|---|
| `github.com/spf13/cobra` | CLI command framework | Apache-2.0 | The de facto standard for Go CLIs; provides consistent commands, flags, and help required by the **Consistent CLI Experience** principle. |
| `github.com/charmbracelet/bubbletea` | Interactive TUI runtime | MIT | Powers the dashboard/interactive mode; building an equivalent event loop and renderer is unreasonable. |
| `github.com/charmbracelet/lipgloss` | Terminal styling and layout | MIT | Declarative styling for TUI and table output; keeps rendering consistent and maintainable. |
| `github.com/stretchr/testify` | Test assertions and mocks | MIT | Readable assertions (`require`, `assert`) for the table-driven tests mandated by **Test Everything**; test-only, not in the shipped binary. |
| `github.com/goccy/go-yaml` | YAML output encoding | MIT | stdlib has no YAML encoder; YAML is a committed output format (roadmap v0.2 / FR-10). Actively maintained, pure-Go, minimal transitive deps; confined to the render layer. Recorded in [ADR 009](../decisions/009-yaml-encoding.md). |
| `github.com/BurntSushi/toml` | Configuration file parsing | MIT | stdlib has no TOML decoder; TOML is the committed config-file format (`../specs/configuration.md`, FND-09). Actively maintained, pure-Go, **zero** transitive deps; confined to the CLI layer (input only, no output-schema obligations). Recorded in [ADR 010](../decisions/010-toml-config.md). |

No other direct dependencies are permitted without going through the process below.

---

## Proposing a New Dependency

A new dependency is an architectural decision. It requires an **ADR** in `../decisions/`.

1. Open an ADR that states the need, the stdlib alternatives considered, and the candidate evaluated against every criterion above.
2. Include the license, the transitive dependency count, and a clean `govulncheck` result.
3. Get maintainer approval on the ADR (see `code-review.md`).
4. Only then add it, in a `build(deps): add <dep>` commit that references the ADR.

Adding a dependency without an accepted ADR is grounds to block a PR.

---

## Maintenance

- **`go.mod` tidiness:** run `go mod tidy` before every PR that touches dependencies. `go.mod` and `go.sum` must contain no unused or missing entries.
- **Pinning:** dependencies are pinned by `go.sum`. Never edit `go.sum` by hand.
- **Update cadence:** review dependency updates on a regular cadence (at least quarterly) and immediately when a security advisory lands. Update in a dedicated `build(deps):` PR, never bundled with feature work.
- **Security scanning:** `govulncheck ./...` runs in CI and must be clean. A new advisory against a pinned dependency is treated as a `fix` with priority.
- **Removal:** if a dependency's justification no longer holds (stdlib gains the capability, or usage is removed), drop it and run `go mod tidy`.

---

*Every import is a promise to maintain someone else's code. Make few promises.*
