# Developer Onboarding

> A practical path for contributors who want to understand SysKit before writing code.

SysKit is designed to be implemented deliberately. This onboarding path helps a new contributor build context in the same order the project expects implementation to happen.

## Step 1: Understand The Product

Read:

- [README](../README.md)
- [Product overview](../specs/product.md)
- [Roadmap](../specs/roadmap.md)

Be able to answer:

- Who is SysKit for?
- What problems does it solve?
- What is explicitly out of scope?

## Step 2: Understand The Architecture

Read:

- [Architecture overview](architecture.md)
- [Architecture specification](../specs/architecture.md)
- [Collector architecture](../specs/collectors.md)
- [Rendering architecture](../specs/rendering.md)

Be able to explain:

- Why CLI code must not read `/proc` directly.
- Why services compute derived metrics.
- Why renderers must not collect data.

## Step 3: Understand The CLI Contract

Read:

- [CLI conventions](../specs/cli-conventions.md)
- [Configuration](../specs/configuration.md)
- [Error handling](../specs/error-handling.md)
- [Logging strategy](../specs/logging-strategy.md)

Be able to explain:

- How flags override configuration.
- Where errors should be printed.
- How partial data should be represented.

## Step 4: Pick A Feature

Read the feature spec in `specs/features/` and the related learning note in `learning/`.

Before writing code, confirm:

- Linux data sources are known.
- Fixtures are planned.
- Output fields are defined.
- Acceptance criteria are testable.

## Step 5: Follow Project Standards

Read:

- [Definition of Ready](../standards/definition-of-ready.md)
- [Definition of Done](../standards/definition-of-done.md)
- [Code review](../standards/code-review.md)
- [Coding conventions](../standards/coding-conventions.md)
- [Naming conventions](../standards/naming-conventions.md)
- [Dependency policy](../standards/dependency-policy.md)

## First Good Contributions

Useful first contributions include:

- Correcting Linux explanations.
- Extending fixture and error-path coverage for an existing collector.
- Improving command examples or man-page wording.
- Adding benchmark coverage for a measured hot path.
- Clarifying feature acceptance criteria or edge cases.

Production changes must follow an accepted specification and the
[Definition of Done](../standards/definition-of-done.md).
