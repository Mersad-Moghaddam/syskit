# Project Structure

> Repository layout for the planning phase and the expected transition into implementation.

SysKit is intentionally organized before production code exists. The current layout separates product planning, architecture, learning material, engineering standards, and collaboration metadata.

## Current Layout

```text
syskit/
├── .github/              # Issue templates, PR template, repository checks
├── decisions/            # Architecture Decision Records
├── docs/                 # User-facing and maintainer-facing documentation
├── learning/             # Linux internals study material
├── scripts/              # Reserved for future development automation
├── specs/                # Product, architecture, and feature specifications
│   └── features/         # Individual feature specifications
├── standards/            # Engineering process and quality standards
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── GOVERNANCE.md
├── LICENSE
├── README.md
└── SECURITY.md
```

## Pre-Implementation Boundary

Until the implementation readiness checklist is complete, the repository must not contain:

- `cmd/`
- `internal/`
- `pkg/`
- `main.go`
- Production `.go` source files
- Cobra scaffolding
- Implemented collectors, services, or formatters

The GitHub workflow enforces this boundary so the repository remains faithful to the design-first process.

## Future Implementation Layout

The exact implementation layout must be confirmed before code begins, but the planned direction is:

```text
syskit/
├── cmd/                  # CLI entry points after implementation begins
├── internal/             # Application internals
│   ├── cli/              # Cobra command wiring
│   ├── collector/        # Built-in domain collectors
│   ├── platform/         # Linux procfs, sysfs, Netlink adapters
│   ├── render/           # Table, JSON, YAML, TUI rendering
│   └── service/          # Aggregation and domain logic
├── testdata/             # Shared integration fixtures where appropriate
└── docs/, specs/, ...
```

This future structure should be created in the first implementation pull request, not during planning.

## Ownership Rules

| Area | Owner mindset | Review focus |
|---|---|---|
| `docs/` | User clarity | Is the guidance understandable and accurate? |
| `specs/` | Product and architecture correctness | Is the intended behavior complete and testable? |
| `learning/` | Educational value | Does it teach the Linux concepts behind the feature? |
| `standards/` | Engineering consistency | Does it set enforceable expectations? |
| `decisions/` | Long-term design memory | Does the ADR capture context and consequences? |
| `.github/` | Collaboration workflow | Does it help contributors provide useful input? |

## Naming Conventions

- Markdown files use lowercase kebab-case names.
- Feature specs live in `specs/features/<feature>.md`.
- ADRs use `NNN-short-title.md`.
- User-facing docs prefer concise names such as `getting-started.md`.
- Standards documents name the practice they govern, such as `code-review.md`.

## Moving From Planning To Code

Before production code is added, complete [implementation readiness](implementation-readiness.md), then update:

- Repository checks in `.github/workflows/ci.yml`.
- README project status.
- Getting started installation/build instructions.
- Contribution guidelines for code changes.
- Changelog with the first implementation milestone.
