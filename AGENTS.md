# AGENTS.md

Guidance for AI coding agents and automated contributors working in this repository.

## Project Context

SysKit is a Linux-first command-line toolkit for system inspection, monitoring, and diagnostics, planned in Go 1.22+. The project follows Specification-Driven Development: design, documentation, and architecture come before implementation.

The repository is currently in the design and specification phase. Treat the documents under `specs/`, `docs/`, `standards/`, `decisions/`, and `learning/` as the source of truth.

## Current Phase Rules

- Do not add production Go code unless the implementation readiness documents explicitly permit it.
- Do not create `cmd/`, `internal/`, `pkg/`, or a Go module as incidental scaffolding.
- Prefer improving specifications, architecture notes, ADRs, standards, onboarding material, and learning notes.
- Keep related documents consistent when changing requirements, architecture, or process.
- If a change contradicts an existing ADR, add or update an ADR instead of silently changing direction.

## Read First

Before making non-trivial changes, read the relevant subset of these files:

- `README.md` for project vision, status, structure, and feature scope.
- `specs/constitution.md` for non-negotiable engineering principles.
- `docs/contributing.md` for planning-phase contribution rules.
- `specs/architecture.md` and `docs/architecture.md` for the planned layered design.
- `specs/cli-conventions.md` for command and output behavior.
- `specs/testing-strategy.md` for test expectations.
- `standards/definition-of-done.md` and `standards/code-review.md` for completion criteria.
- `standards/commit-conventions.md` before committing.

For feature work, also read the matching file in `specs/features/` and any related notes in `learning/`.

## Architectural Constraints

SysKit is intentionally Linux-only. Do not introduce cross-platform abstractions, compatibility shims, or `runtime.GOOS` branching for non-Linux behavior.

Prefer native Linux interfaces over shelling out:

- `/proc`
- `/sys`
- Netlink
- cgroups
- other kernel-provided APIs

Shell commands may be used for development tasks, but SysKit's planned implementation should not parse external command output unless a documented exception explains why.

The planned dependency direction is strict:

```text
CLI -> Command -> Service -> Collector -> Platform -> Linux kernel interfaces
```

Lower layers must not import or know about higher layers. Collectors must not depend on Cobra, terminal colors, table renderers, or CLI flags. Renderers must not collect system data.

## Documentation Standards

- Write clear, precise English.
- Prefer concrete acceptance criteria over vague intent.
- Keep terminology aligned with `docs/glossary.md`.
- Update cross-references when moving or renaming documents.
- Preserve the project's design-first posture: explain why a behavior exists, not only what it does.
- Avoid duplicating large sections across documents; link to the canonical source instead.

## Future Go Implementation Standards

When implementation begins, follow these rules:

- Use idiomatic Go and the standard library by default.
- Use Cobra for CLI commands, Bubble Tea for the TUI, and Lip Gloss for terminal styling as documented.
- Keep packages small, domain-focused, and named after what they provide.
- Avoid `util`, `common`, and `helpers` grab-bag packages.
- Return errors instead of panicking in library code.
- Wrap errors with `%w` and add useful lowercase context without trailing punctuation.
- Inject filesystem and platform dependencies so collectors can be tested with fixtures.
- Avoid mutable package-level state and observable `init()` side effects.
- Keep dependencies minimal and justify additions through the dependency policy and ADR process.

## Testing Expectations

Documentation-only changes should be checked for consistency, broken links where practical, and alignment with the constitution.

When Go code exists, feature completion requires:

- Unit tests for parsing, transformation, filtering, formatting, and error paths.
- Fixture-backed collector tests instead of depending on the developer's live `/proc` or `/sys`.
- Linux integration tests for real kernel interfaces, guarded with appropriate build tags.
- Golden-file tests for user-facing output contracts.
- Benchmarks for hot paths and allocation-sensitive code.
- Passing `gofmt`, `goimports`, `go vet`, vulnerability checks, and race tests where applicable.

Do not claim a feature is complete without the tests required by `specs/testing-strategy.md` and `standards/definition-of-done.md`.

## Working Safely

- Keep changes focused and small.
- Do not rewrite unrelated prose or reorganize directories unless the task requires it.
- Do not remove constraints, non-goals, or architectural boundaries without documenting the decision.
- Preserve user changes in the working tree. If unrelated files are dirty, leave them alone.
- Prefer `rg` for searching and inspect existing patterns before editing.
- Use structured formats and parsers when available instead of brittle string manipulation.

## Commit Guidance

Use Conventional Commits as defined in `standards/commit-conventions.md`.

Examples:

- `docs(project): add agent contribution guidance`
- `docs(cpu): clarify collector fixture requirements`
- `feat(cli): add cpu command skeleton`

For this repository's current phase, most changes should use the `docs` type.

Commit subjects must be imperative, lowercase, under 72 characters, and without a trailing period.

## Agent Completion Checklist

Before handing work back:

- Confirm the change matches the current project phase.
- Confirm affected docs do not contradict `specs/constitution.md` or ADRs.
- Check relevant links and file paths.
- Run the smallest meaningful validation available.
- Summarize what changed and any checks that were not possible.
