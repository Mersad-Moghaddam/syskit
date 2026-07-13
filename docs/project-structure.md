# Project Structure

> Repository layout for the implemented SysKit command-line application.

SysKit separates production code, public contracts, tests, product specifications,
engineering standards, and contributor documentation.

## Current Layout

```text
syskit/
├── .github/              # Issue templates, PR template, repository checks
├── cmd/syskit/           # CLI entry point
├── contracts/            # Machine-readable v1 compatibility manifests
├── decisions/            # Architecture Decision Records
├── docs/                 # User-facing and maintainer-facing documentation
├── internal/             # CLI, services, collectors, models, platform, rendering
├── learning/             # Practical Linux, Go, and SysKit engineering course
├── plan/                 # Delivery plan, backlog, epics, and sprint records
├── scripts/              # Fixture, release, and packaging automation
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

## Implementation Boundary

The readiness checklist is complete and production code is allowed in the approved
layout. New code must preserve these boundaries:

- `cmd/syskit` remains a minimal entry point.
- CLI → command → service → collector → platform dependency direction is strict.
- Collectors use injected platform interfaces and do not shell out.
- No public `pkg/` API is provided; SysKit's public API is its CLI contract.
- Shared behavior belongs in domain-focused packages, not grab-bag helpers.

## Internal Layout

```text
syskit/
├── cmd/syskit/           # CLI entry point
├── internal/             # Application internals
│   ├── cli/              # Cobra wiring, config, logging, exit mapping
│   ├── collector/        # Built-in domain collectors
│   ├── model/            # Typed structured-output models
│   ├── platform/         # Linux procfs, sysfs, Netlink, cgroup adapters
│   ├── plugin/           # Out-of-process plugin discovery and protocol
│   ├── render/           # Table, JSON, YAML, TUI rendering
│   └── service/          # Aggregation and domain logic
├── testdata/             # Shared fixture data
└── contracts/            # v1 compatibility manifests enforced by tests
```

## Ownership Rules

| Area | Owner mindset | Review focus |
|---|---|---|
| `docs/` | User clarity | Is the guidance understandable and accurate? |
| `specs/` | Product and architecture correctness | Is the intended behavior complete and testable? |
| `learning/` | Educational value | Does it teach the Linux/Go concepts, investigation method, and evidence needed behind the feature? |
| `standards/` | Engineering consistency | Does it set enforceable expectations? |
| `decisions/` | Long-term design memory | Does the ADR capture context and consequences? |
| `contracts/` | Public compatibility | Does a CLI or schema change obey SemVer? |
| `.github/` | Collaboration workflow | Does it help contributors provide useful input? |

## Naming Conventions

- Markdown files use lowercase kebab-case names.
- Feature specs live in `specs/features/<feature>.md`.
- ADRs use `NNN-short-title.md`.
- User-facing docs prefer concise names such as `getting-started.md`.
- Standards documents name the practice they govern, such as `code-review.md`.

## Adding Production Code

Read [implementation readiness](implementation-readiness.md) for the historical
transition decision, then follow the current [contributing guide](contributing.md),
[testing strategy](../specs/testing-strategy.md), and
[Definition of Done](../standards/definition-of-done.md).
