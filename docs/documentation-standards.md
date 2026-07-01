# Documentation Standards

> How SysKit documentation should be written, reviewed, and maintained.

SysKit documentation is part of the product. It should help users understand what the tool does, help contributors implement it correctly, and preserve architectural decisions over time.

## Documentation Types

| Type | Location | Purpose |
|---|---|---|
| User docs | `docs/` | Explain how users and contributors work with SysKit |
| Product specs | `specs/` | Define expected behavior and constraints |
| Feature specs | `specs/features/` | Define one feature in enough detail to implement |
| Learning notes | `learning/` | Teach Linux concepts needed for implementation |
| Standards | `standards/` | Define engineering and collaboration rules |
| ADRs | `decisions/` | Record important architectural decisions |

## Writing Principles

- Prefer precise, direct language over vague promises.
- Write in the present tense for stable project policy.
- Mark future behavior as planned when no implementation exists yet.
- Avoid placeholders in committed documentation.
- Keep examples realistic and tied to SysKit's Linux-first scope.
- Use tables when comparing contracts, formats, or responsibilities.
- Link related documents so readers can move from overview to detail.

## Feature Specification Requirements

Every feature specification must include:

- Purpose
- User story
- Motivation
- Requirements
- Linux concepts
- Expected CLI
- Expected output
- Edge cases
- Acceptance criteria
- Learning objectives
- Estimated complexity
- Dependencies
- Future extensions

Acceptance criteria must be testable. A reviewer should be able to tell whether an implementation satisfies the spec without interpreting intent.

## Output Examples

Output examples should:

- Prefer representative data over synthetic filler.
- Include table output for human workflows.
- Include JSON when a command is expected to support automation.
- Avoid promising fields that have not been accepted as part of the output contract.
- Clearly state when examples are illustrative.

## Link Hygiene

When adding or renaming documents:

- Update the README documentation map if the document is a major entry point.
- Update nearby "Further reading" lists.
- Check relative links from both the source and target document.
- Avoid links to files that do not exist yet.

## Review Checklist

- The document has a clear audience.
- The document has no unresolved placeholders.
- The content matches existing architecture and standards.
- Linux terminology is accurate.
- Examples are plausible on real Linux systems.
- Future work is clearly labeled as planned.
- Related docs are linked.
