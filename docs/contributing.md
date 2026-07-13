# Contributing to SysKit

> How to propose changes, review designs, and prepare implementation work.

SysKit is in its implementation phase, with v0.1 released and later milestones
under active development. Contributions are welcome across code and planning:
implementing features against accepted specs, clarifying requirements,
strengthening architecture, adding Linux references, tightening acceptance
criteria, or improving project process. Production Go code must build on the
approved layout and satisfy the Definition of Done.

## Before You Start

Read these documents first:

- [Engineering constitution](../specs/constitution.md)
- [Architecture specification](../specs/architecture.md)
- [CLI conventions](../specs/cli-conventions.md)
- [Testing strategy](../specs/testing-strategy.md)
- [Definition of Done](../standards/definition-of-done.md)
- [Code of Conduct](../CODE_OF_CONDUCT.md)

## Development Setup

For documentation and planning work, install:

- Git.
- A Markdown-capable editor.
- A terminal on Linux or any environment that can edit Markdown files.

For implementation work, use Linux and Go 1.26.3 or newer. Build with
`go build ./...`, and run the test suite with `go test -race ./...`. The module
path is `github.com/Mersad-Moghaddam/syskit`.

## Workflow

1. Open or select an issue that describes the change.
2. Create a topic branch using the [branch strategy](../standards/branch-strategy.md).
3. Keep the change focused on one documentation, planning, or design concern.
4. Update related documents together so the repository does not contradict itself.
5. Run the repository checks locally where practical.
6. Open a pull request using the template.

## Implementation-Phase Rules

- Production Go source lives under `cmd/` and `internal/` on the approved layout; no `pkg/` grab-bag and no `util`/`common`/`helpers` packages.
- Respect the strict downward dependency direction (CLI → Command → Service → Collector → Platform → kernel); lower layers never import higher ones.
- Collectors read only through the injected `platform.SysFS` seam — never the OS filesystem directly — and never shell out to external commands.
- SysKit is Linux-only: no `runtime.GOOS` branching or OS build tags.
- Do not add external dependencies without updating the dependency policy and adding an ADR.
- Do not remove architecture constraints unless a new ADR explains the decision.

## Pull Request Expectations

A good pull request explains:

- What changed.
- Why the change is needed.
- Which specs, standards, or docs were updated.
- Which tradeoffs were considered.
- How reviewers can validate the result.

Design changes should link to the affected feature spec or ADR. Documentation-only changes should still be reviewed for accuracy, consistency, and broken links.

## Review Standards

Reviewers should focus on:

- Correctness of Linux concepts.
- Alignment with project scope and non-goals.
- Consistency with the architecture and CLI conventions.
- Testability and acceptance criteria.
- Whether the change makes future implementation clearer.

See [code review standards](../standards/code-review.md) for the full review model.

## Reporting Issues

Use the issue templates for bugs, feature proposals, design proposals, and documentation improvements. Since SysKit has no released implementation yet, bug reports should usually be about documentation errors, conflicting requirements, broken workflows, or repository process.

## Conduct

All participation follows the [Code of Conduct](../CODE_OF_CONDUCT.md). Technical disagreement is expected; disrespect is not.
