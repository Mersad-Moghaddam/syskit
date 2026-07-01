# Commit Conventions

> How commit messages are structured so history is readable and changelogs and versions are derivable.

---

## Purpose

SysKit uses **Conventional Commits 1.0.0**. A disciplined commit history is documentation: it explains *why* each change was made, drives the CHANGELOG, and determines the next version per `versioning.md`.

This standard applies to every commit that lands on `main`. Because SysKit squash-merges (see `branch-strategy.md`), the **squash commit title and body** are the artifact that must conform — intermediate work-in-progress commits on a feature branch are exempt.

---

## Format

```text
type(scope): subject

[optional body]

[optional footer(s)]
```

Rules:

- **type** and **subject** are mandatory; **scope** is strongly recommended.
- Subject is imperative mood, lowercase, no trailing period: "add cpu frequency parser", not "Added..." or "adds...".
- Keep the subject line at or under **72 characters**.
- Separate subject, body, and footer with a blank line.
- Wrap the body at ~72 columns. Explain *what* and *why*, not *how*.

---

## Types

| Type | Use for | Version impact |
|---|---|---|
| `feat` | A new feature or command | MINOR |
| `fix` | A bug fix | PATCH |
| `perf` | A change that improves performance | PATCH |
| `refactor` | Code change that neither fixes a bug nor adds a feature | none |
| `docs` | Documentation only | none |
| `test` | Adding or correcting tests | none |
| `build` | Build system, `go.mod`, tooling | none |
| `ci` | CI configuration and scripts | none |
| `chore` | Maintenance with no src/test change | none |
| `style` | Formatting only (gofmt, whitespace) | none |

A commit that introduces a breaking change is MAJOR regardless of type (see below).

---

## Scopes

The scope names the affected domain or layer. Prefer a domain (collector name) over a layer where both apply.

| Category | Scopes |
|---|---|
| Collectors | `cpu`, `memory`, `disk`, `network`, `process`, `system`, `ports` |
| Layers | `cli`, `command`, `service`, `collector`, `platform` |
| Cross-cutting | `output`, `tui`, `config`, `deps`, `ci`, `docs` |

Omit the scope only when a change is genuinely global (e.g., `chore: bump Go to 1.23`).

---

## Body and Footer Rules

- Use the **body** to give context: the problem, the approach, and any trade-offs. Reference `../specs/` or `../decisions/` where a decision is recorded.
- Reference issues in the footer: `Refs: #42`, `Closes: #42`.
- Co-authorship goes in the footer: `Co-authored-by: Name <email>`.

### Breaking changes

A breaking change is signaled in **both** places, one of which is mandatory:

1. A `!` after the type/scope: `feat(cli)!: rename --json to --format`
2. A `BREAKING CHANGE:` footer describing the break and the migration.

```text
feat(cli)!: replace --json flag with --format

The --json boolean is removed in favor of --format {table,json,yaml},
unifying output selection across all commands.

BREAKING CHANGE: --json no longer exists. Use --format json instead.
```

Any commit with `!` or a `BREAKING CHANGE:` footer triggers a MAJOR bump (or a MINOR bump while pre-1.0 — see `versioning.md`).

---

## Examples

### Good

```text
feat(cpu): add per-core frequency reporting

Reads scaling_cur_freq from /sys/devices/system/cpu to report
live per-core frequency alongside the model's base frequency.

Refs: #17
```

```text
fix(memory): correct available-memory calculation on 6.x kernels

MemAvailable was double-counting reclaimable slab. Use the kernel
field directly when present instead of recomputing it.

Closes: #58
```

```text
perf(process): reuse the /proc/[pid]/stat read buffer
```

### Bad

| Commit | Why it fails |
|---|---|
| `Fixed the bug` | No type, past tense, vague. |
| `feat: stuff` | Meaningless subject, no scope. |
| `update` | No type, no subject. |
| `feat(cpu): Added new command.` | Capitalized, past tense, trailing period. |
| `wip` | Not conventional; never squash-merge this. |

---

## Mapping to Changelog and SemVer

- The CHANGELOG groups entries by type: `feat` → **Added**, `fix` → **Fixed**, `perf` → **Performance**, breaking changes → **Changed / Breaking**.
- `docs`, `test`, `chore`, `ci`, `build`, `style`, and `refactor` are omitted from the user-facing CHANGELOG.
- Version selection follows the highest-impact commit since the last tag: any breaking → MAJOR; any `feat` → MINOR; otherwise `fix`/`perf` → PATCH. See `versioning.md` for the pre-1.0 caveat.

---

*A commit message is written once and read many times. Spend the extra minute.*
