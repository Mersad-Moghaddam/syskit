# Contributing to SysKit

> How to propose changes, review designs, and prepare implementation work.

SysKit is in the design and specification phase. Contributions are welcome when they improve the planning foundation: clarifying requirements, strengthening architecture, adding Linux references, tightening acceptance criteria, or improving project process. Production Go code should wait until the implementation readiness checklist is satisfied.

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

For future implementation work, use Linux and Go 1.22 or newer. No Go module or production package is expected during the planning phase.

## Workflow

1. Open or select an issue that describes the change.
2. Create a topic branch using the [branch strategy](../standards/branch-strategy.md).
3. Keep the change focused on one documentation, planning, or design concern.
4. Update related documents together so the repository does not contradict itself.
5. Run the repository checks locally where practical.
6. Open a pull request using the template.

## Planning-Phase Rules

- Do not add production Go source files.
- Do not create `cmd/`, `internal/`, or `pkg/` yet.
- Do not scaffold Cobra commands before the feature specs are approved.
- Do not add external dependencies without updating the dependency policy and ADRs.
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
