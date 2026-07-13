# CLI Conventions

> Command naming, flags, output behavior, exit codes, and terminal interaction rules.

SysKit should feel like one coherent tool, not a bundle of unrelated commands. The same flag names, output formats, sorting behavior, and error conventions should work across every subsystem.

## Command Naming

- Use singular nouns for primary domains: `cpu`, `memory`, `disk`, `process`, `network`, `ports`, `filesystem`.
- Use subcommands for narrower views: `network interfaces`, `process tree`, `disk io`.
- Prefer clear words over abbreviations.
- Avoid commands that imply mutation unless the command is explicitly out of scope.

## Global Flags

| Flag | Values | Purpose |
|---|---|---|
| `--format` | `table`, `json`, `yaml` | Select output format |
| `--config` | path | Use a specific config file |
| `--color` | `auto`, `always`, `never` | Control color output |
| `--no-header` | boolean | Suppress table headers |
| `--verbose` | boolean | Show diagnostic detail |
| `--debug` | boolean | Show debug-level diagnostic detail |
| `--quiet` | boolean | Suppress non-essential messages |

`--interval` is a consistent command-local flag on sampling and live views.
Continuous generic refresh uses the `watch <command>` subcommand; there is no
global `--watch` boolean.

## Filtering And Sorting

Commands that list resources should use consistent flags:

| Flag | Meaning |
|---|---|
| `--sort <field>` | Sort by a named field |
| `--reverse` | Reverse sort order |
| `--limit <n>` | Limit result count |
| `--filter <expr>` | Apply a simple filter expression where supported |

Invalid fields should return a usage error that lists valid fields for that command.

`--filter` uses `field=value` equality predicates. Repeated filters are combined
with AND. `--sort` accepts one documented field; `--reverse` reverses that
field's natural order, and `--limit 0` means no limit.

## Output Formats

Table output is optimized for humans. JSON and YAML are stable automation contracts.

| Format | Stability expectation |
|---|---|
| `table` | Columns may evolve before `v1.0`, but should remain predictable within a release |
| `json` | Field names and types are part of the output contract |
| `yaml` | Mirrors JSON structure for human-editable pipelines |

Structured output should not include terminal color, progress text, or warnings mixed into stdout. Diagnostics belong on stderr.

## Color

Color defaults to `auto`:

- Enabled when stdout is an interactive terminal.
- Disabled when stdout is redirected or piped.
- Disabled when `NO_COLOR` is set.
- Forced by `--color always`.
- Disabled by `--color never`.

Color must never be the only carrier of meaning.

## Exit Codes

| Code | Name | Meaning |
|---|---|---|
| `0` | Success | Command completed successfully |
| `1` | General error | An unspecified runtime error occurred |
| `2` | Usage error | Invalid flags, arguments, or command usage |
| `3` | Permission | Insufficient privilege to read a kernel interface |
| `4` | Unsupported | Required kernel interface is missing or unsupported |
| `5` | Partial failure | Some data collected; one or more collectors failed |

This table is canonical and defined in full in [error-handling.md](error-handling.md); the two files must always match. Commands may return partial data with exit code `5` only when the output clearly identifies what could not be read.

## Error Presentation

Errors should be concise by default:

```text
syskit: cannot read /proc/1234/status: permission denied
```

Verbose mode may include source paths, wrapped error causes, and troubleshooting hints. Machine-readable output should represent errors as structured fields when a command supports partial results.

## Help Text

Help text should include:

- One-sentence command purpose.
- Common examples.
- Supported output fields.
- Notes about required permissions or kernel support.

Help text should not mention implementation details unless they help users understand limitations.
