# Rendering Architecture

> How SysKit should present structured results across table, JSON, YAML, and interactive terminal views.

Rendering is a separate concern from collection and services. A renderer receives structured data and converts it into a user-facing representation without reading Linux interfaces or applying business rules.

## Rendering Goals

- Keep structured output stable and scriptable.
- Make table output readable in common terminal widths.
- Keep diagnostics on stderr.
- Respect color and TTY conventions.
- Share domain models between static commands and live views where possible.

## Renderer Types

| Renderer | Audience | Contract |
|---|---|---|
| Table | Humans in terminals | Readable, aligned, width-aware |
| JSON | Automation and scripts | Stable field names and types |
| YAML | Humans and configuration workflows | Mirrors JSON structure |
| TUI | Interactive monitoring | Responsive, keyboard-driven, live-updating |

## Table Rendering

Table renderers should:

- Use consistent units and column names.
- Align numbers to the right and text to the left.
- Truncate only when necessary and visibly.
- Respect `--no-header`.
- Disable color when output is not a TTY.

Tables may choose a smaller default column set than JSON, but hidden values must be documented.

## Structured Rendering

JSON and YAML output should:

- Emit only valid machine-readable data on stdout.
- Use snake_case field names.
- Use explicit units in field names or nested metadata.
- Represent timestamps in RFC 3339 format when needed.
- Avoid lossy human formatting such as `"1.2 GB"` in numeric fields.

Warnings and partial-data notices belong in structured fields or stderr, depending on command behavior.

## TUI Rendering

The future dashboard should use the same service layer as non-interactive commands. The TUI renderer owns:

- Layout.
- Keyboard navigation.
- Focus state.
- Refresh timing.
- Color and styling.
- Resize handling.

It must not read from `/proc`, `/sys`, Netlink, or cgroups directly.

## Output Stability

Before `v1.0`, output contracts may evolve with changelog entries. At `v1.0`, JSON and YAML fields become compatibility commitments. Removing or changing field types after `v1.0` requires a major version unless the field was clearly experimental.

## Acceptance Criteria

- Renderers are deterministic for the same input.
- Renderers are testable with golden files.
- Structured output contains no terminal control sequences.
- Terminal color follows `NO_COLOR` and `--color`.
- Rendering code does not collect system data.
