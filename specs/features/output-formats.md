# Output Formats Feature Specification

## Purpose

Provide consistent table, JSON, YAML, and plain text output behavior across SysKit commands.

## User Story

As a user, I want SysKit output to be predictable whether I am reading it in a terminal or piping it into automation.

## Motivation

Traditional Linux tools often mix human text, warnings, and data in ways that are hard to parse. SysKit should make output format a first-class contract.

## Requirements

- Support `--format table`, `--format json`, and `--format yaml` across all non-interactive commands.
- Keep diagnostics and warnings out of stdout for structured output.
- Use stable field names and numeric types.
- Support `--no-header` for table output.
- Respect color and TTY conventions.

## Linux Concepts

- stdout and stderr separation.
- TTY detection.
- Pipe-friendly command design.
- Exit status conventions.

## Expected CLI

```sh
syskit cpu --format table
syskit cpu --format json
syskit cpu --format yaml
syskit process --no-header
syskit disk --color never
```

## Expected Output

Table output:

```text
NAME      VALUE
uptime    3d 04h 12m
load_1m   0.42
```

JSON output:

```json
{
  "uptime_seconds": 274320,
  "load_average": {
    "one_minute": 0.42,
    "five_minutes": 0.35,
    "fifteen_minutes": 0.30
  }
}
```

## Edge Cases

- Broken pipe when output is piped to `head`.
- Terminal width is too narrow.
- Structured output includes partial data warnings.
- Color requested for non-TTY output.
- YAML dependency policy must be reviewed before implementation.

## Acceptance Criteria

- Each command declares supported formats.
- JSON output is valid and deterministic for the same input.
- YAML mirrors JSON structure.
- Table output is readable at common terminal widths.
- Warnings never corrupt JSON or YAML stdout.

## Learning Objectives

- Learn CLI output contracts.
- Understand why stdout and stderr separation matters.
- Learn golden-file testing for user-facing output.

## Estimated Complexity

Medium.

## Dependencies

- Rendering architecture.
- CLI conventions.
- Dependency policy for YAML support.

## Future Extensions

- NDJSON for streaming watch output.
- Prometheus exposition format.
- Field selection flags.
- Custom templates.
