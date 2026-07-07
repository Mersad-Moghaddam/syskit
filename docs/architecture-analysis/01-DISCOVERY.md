# Phase 1 — Discovery and Full Inventory

*Point-in-time analysis that informed ARCHITECTURE.md; retained here for historical reference.*

> Evidence-based inventory of the SysKit repository as it exists on 2026-07-07. SysKit is in its **design and specification phase**: no production Go code exists yet, and a CI check actively forbids it (`.github/workflows/ci.yml`). Everything below is drawn from the repository's documentation, specs, ADRs, and configuration — not from running code.

## 1. What SysKit Is

SysKit is a planned **Linux-first command-line toolkit for system inspection, resource monitoring, and diagnostics**, to be built in Go. It reads directly from native kernel interfaces (`/proc`, `/sys`, Netlink, cgroups) rather than wrapping shell utilities. It is both a practical tool and a reference/learning project. (`README.md`, `specs/product.md`)

The repository practices **Specification-Driven Development (SDD)**: specs and architecture are written and accepted before any implementation. (`specs/constitution.md` principle 6)

## 2. Complete Directory Map

```text
syskit/
├── .agents/                 # Agent-tooling metadata (contribution assistants)
├── .codex/                  # Agent-tooling metadata
├── .claude/                 # Local agent settings (settings.local.json)
├── .github/                 # Collaboration + CI
│   ├── workflows/ci.yml     # "Repository Checks": enforces planning-phase boundary
│   ├── ISSUE_TEMPLATE/      # bug, feature, docs, design-proposal, config
│   ├── pull_request_template.md
│   └── labels.yml
├── decisions/               # Architecture Decision Records (ADR 001–007 + template)
├── docs/                    # User- and maintainer-facing documentation
├── learning/                # Linux-internals study notes (cpu, memory, disk, etc.)
├── plan/                    # Scrum/Agile delivery plan (backlog, sprints, epics)
│   ├── epics/               # EPIC-00 … EPIC-07
│   ├── sprints/             # sprint-00 … sprint-11
│   └── templates/           # user-story, sprint-plan, retrospective, etc.
├── scripts/                 # Reserved for future automation (only .gitkeep today)
├── specs/                   # Product, architecture, and cross-cutting specifications
│   └── features/            # One spec per planned command/feature
├── standards/               # Engineering process & quality standards
├── AGENTS.md                # Guidance for AI/agent contributors
├── CHANGELOG.md             # Keep-a-Changelog; currently only [Unreleased]
├── CODE_OF_CONDUCT.md
├── GOVERNANCE.md
├── SECURITY.md
├── SUPPORT.md
├── LICENSE                  # MIT
├── README.md
├── .editorconfig
└── .gitignore
```

### Responsibility of each top-level area

| Area | Responsibility |
|---|---|
| `specs/` | Canonical description of **what** SysKit is and how it must behave. Contains the constitution, product overview, architecture, roadmap, and cross-cutting specs (collectors, rendering, plugins, CLI conventions, config, error handling, logging, testing) plus per-feature specs under `features/`. |
| `decisions/` | Architecture Decision Records — the long-term "why" behind binding choices (language, platform, data source, architecture, CLI/TUI frameworks, plugin model). |
| `standards/` | Enforceable engineering process: definition of ready/done, code review, versioning, naming, commit and branch conventions, dependency policy, coding conventions. |
| `docs/` | Human-facing guides: getting started, architecture overview, contributing, onboarding, project structure, release process, documentation standards, implementation-readiness checklist. |
| `learning/` | Study notes on the Linux subsystems behind each feature (Learn Before Build, constitution principle 10). |
| `plan/` | Execution plan (Scrum): product backlog (source of truth for stories/points), release plan, epics, time-boxed sprints, estimation/velocity, risk register, metrics. |
| `.github/` | Contributor workflow: issue/PR templates, labels, and the planning-phase CI guard. |
| `scripts/` | Placeholder for future build/dev automation (e.g. planned `scripts/capture-fixtures.sh`). |
| Root docs | Repository-level governance and metadata (license, changelog, conduct, security, support, governance). |

## 3. Major Modules / Services

**No runtime modules exist yet.** The planned logical modules — ratified as a binding six-layer architecture in ADR 004 — are:

| Planned layer/module | Responsibility (spec) |
|---|---|
| CLI layer | Argument/flag parsing, config load, format selection, terminal + TUI rendering, error presentation, exit codes, logging. |
| Command layer | Thin Cobra command definitions; validate flags, call services. |
| Service layer | Business logic: aggregate collectors, filter/sort, compute derived metrics (rates, deltas, %). |
| Collector layer | Per-domain, independent data gatherers (cpu, memory, disk, process, network, ports, filesystem, …) returning typed structs. |
| Platform abstraction layer | Only layer touching the OS; wraps `/proc`, `/sys`, Netlink, cgroup v1/v2 behind a `SysFS`-style interface. |
| Linux kernel interfaces | The kernel data sources themselves (not code SysKit writes). |

Planned commands (from `specs/roadmap.md` and feature specs): `system`, `cpu`, `memory`, `disk`, `filesystem`, `process` (+ `tree`), `network`, `ports`, `dashboard`, `watch`, `top`, `docker` (+ `inspect`), plus a plugin subsystem.

## 4. Current Technology Stack

| Concern | Choice | Evidence |
|---|---|---|
| Language | Go 1.22+ | ADR 001, `README.md` |
| Platform | Linux only (no `runtime.GOOS` branching) | ADR 002, constitution 1 |
| Data sources | `/proc`, `/sys`, Netlink (`AF_NETLINK` via `golang.org/x/sys/unix`), cgroups v1/v2 | ADR 003 |
| CLI framework | Cobra (`github.com/spf13/cobra`) + pflag | ADR 005 |
| TUI (v0.3+) | Bubble Tea + Lip Gloss (`github.com/charmbracelet/*`) | ADR 006 |
| Config format | TOML, XDG-located, `SYSKIT_*` env, zero-config default | `specs/configuration.md` |
| Logging | stdlib `log/slog`, text handler → stderr | `specs/logging-strategy.md` |
| Testing | stdlib `testing` + testify; fixtures, golden files, `-race`, integration build tag | `specs/testing-strategy.md` |
| Approved deps | cobra, bubbletea, lipgloss, testify, `golang.org/x/sys/unix` (per policy) | `standards/dependency-policy.md` |
| Database / broker / infra | **None** — no persistent store, queue, or server. Read-only local CLI. | product Non-Goals; architecture |
| License | MIT | `LICENSE` |

There is **no `go.mod`** (confirmed absent), no `main.go`, and no `.go` files — the CI boundary check fails the build if any appear during the planning phase.

## 5. Entry Points

| Kind | Present today | Planned |
|---|---|---|
| API (HTTP/RPC) | None | None (out of scope) |
| CLI | None | `syskit <command>` via Cobra; single static binary (ADR 001/005) |
| Workers / daemons | None | None (no background service; `watch`/`top`/`dashboard` are foreground interactive) |
| Cron jobs | None | None |
| CI entry | `.github/workflows/ci.yml` on push/PR to `main` | Will gain Go build/test/lint stages at the implementation-transition PR |

The only executable entry point that exists right now is the **CI workflow**, whose job is to *prevent* code from appearing prematurely and to validate that required planning artifacts are present and free of unresolved task markers.

## 6. Notable Observations for Later Phases

- Two exit-code tables disagree: `specs/cli-conventions.md` defines codes 0–4; `specs/error-handling.md` defines 0–5 (adds `ExitPartial = 5` and renumbers permission/unsupported). Flagged for Phase 3.
- Configuration precedence (`specs/configuration.md`) lists flags > env > file > defaults, with per-command `[section]` overrides inserted between env and global.
- The `plan/` tree (Scrum artifacts) is process, not architecture, but establishes the delivery sequence v0.1 → v1.0.
