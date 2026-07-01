# Naming Conventions

> Naming rules for commands, packages, files, fields, tests, and documentation.

Consistent names make SysKit easier to learn and easier to review. Names should be boring, searchable, and close to the Linux concept they represent.

## CLI Names

- Use lowercase command names.
- Use singular domain names: `cpu`, `memory`, `disk`, `process`, `network`, `filesystem`.
- Use plural only when the Linux concept is naturally plural, such as `ports`.
- Prefer full words over abbreviations.
- Use subcommands for narrower views: `process tree`, `network routes`, `disk io`.

## Flag Names

- Use kebab-case for multi-word flags, such as `--no-header`.
- Reuse global flags from [CLI conventions](../specs/cli-conventions.md).
- Prefer positive names unless the negative form is a common CLI convention.
- Boolean flags should read clearly when present.

## Go Names

These rules apply once implementation begins:

- Package names are short, lowercase, and singular.
- Interfaces describe behavior when useful, such as `Reader` or `Collector`.
- Avoid package names that repeat their parent directory.
- Avoid vague names such as `common`, `utils`, and `helpers`.
- Exported names need documentation when they become public package API.

## Structured Output Fields

- Use snake_case.
- Include units in field names when the value is numeric and unit-specific, such as `memory_bytes`.
- Avoid human-formatted strings in structured output.
- Use arrays for repeated values and objects for grouped fields.

## Documentation Files

- Use lowercase kebab-case.
- Keep names short and descriptive.
- Feature specs live under `specs/features/`.
- ADRs use `NNN-short-title.md`.

## Test Names

Once code exists:

- Test names should describe behavior, not implementation.
- Table-driven test cases should have scenario names.
- Fixture directories should identify the host shape or kernel behavior they represent.

## Review Checklist

- Does the name match an existing project pattern?
- Can a new contributor guess what it means?
- Is it searchable?
- Does it avoid leaking implementation details into user-facing contracts?
