# Versioning

> How SysKit versions its releases and what compatibility users can rely on.

---

## Purpose

SysKit follows **Semantic Versioning 2.0.0**. This standard defines what MAJOR, MINOR, and PATCH mean for a *CLI tool*, what constitutes the public contract, and how deprecations are handled.

Versions are derived from commit history (`commit-conventions.md`) and applied as tags on `main` (`branch-strategy.md`).

---

## SemVer Format

Versions are `MAJOR.MINOR.PATCH`, optionally with a pre-release suffix:

```text
v0.2.0
v1.0.0
v1.1.0-rc.1
```

| Segment | Increment when |
|---|---|
| **MAJOR** | An incompatible change to the public contract. |
| **MINOR** | Backward-compatible new functionality (a new command, flag, or field). |
| **PATCH** | Backward-compatible bug or performance fix. |

---

## The Public Contract

For a library, the public API is its exported symbols. For SysKit — a CLI tool — the public contract is what users and scripts depend on:

| Contract surface | Examples |
|---|---|
| **CLI commands and flags** | `syskit cpu`, `--format`, `--no-color` |
| **Output formats** | The JSON/YAML schema and table columns for each command |
| **Exit codes** | `0` success, non-zero documented failure codes |

A change is **breaking** if it could break a script that consumes SysKit: removing or renaming a command or flag, removing or renaming a field in machine-readable output, changing an exit code's meaning, or changing default behavior in an incompatible way.

Adding a new command, a new flag, or a new field to output is backward-compatible (MINOR). Fixing incorrect output values or crashes is PATCH.

> Internal Go package structure is **not** part of the public contract. SysKit is consumed as a binary, not imported as a library. Refactors that do not change CLI behavior, output, or exit codes are non-versioned.

---

## Pre-1.0 Caveat

While SysKit is on a `0.x` line (see `../specs/roadmap.md`), the API is explicitly unstable:

- A **MINOR** bump (`0.1.0` → `0.2.0`) **may** include breaking changes.
- A **PATCH** bump (`0.1.0` → `0.1.1`) remains backward-compatible.
- Breaking changes are still announced in the CHANGELOG and, where practical, given a deprecation window.

The first stable contract is frozen at `v1.0.0`. From then on, breaking changes require a MAJOR bump.
The candidate inventory and its CI enforcement are documented in the
[v1 compatibility contract](../docs/compatibility.md); the canonical manifests
live under `contracts/`.

---

## Pre-release Tags

Release candidates and previews use dot-separated pre-release identifiers, which sort *before* the final release:

```text
v1.0.0-alpha.1  <  v1.0.0-beta.1  <  v1.0.0-rc.1  <  v1.0.0
```

Pre-releases carry no compatibility guarantee and are for validation only.

---

## Compatibility Guarantees

At and after `v1.0.0`, within a MAJOR line:

- Existing commands, flags, output fields, and exit codes keep working.
- New fields may be **added** to JSON/YAML output; consumers must ignore unknown fields.
- Table output layout may change (it is human-facing); scripts must consume `--format json`, not scrape tables.

---

## Deprecation Policy

Nothing in the public contract is removed abruptly. Removal follows a three-stage path across releases:

1. **Deprecate** — Announce in the CHANGELOG and docs, provide the replacement, and mark the item deprecated. No behavior change yet.
2. **Warn** — Emit a deprecation notice to **stderr** when the deprecated command/flag is used, so machine-readable **stdout** stays clean. Keep the behavior working.
3. **Remove** — Delete the item. This is a breaking change (MAJOR post-1.0; permitted on a MINOR pre-1.0) and is recorded with a `BREAKING CHANGE:` footer.

A deprecated feature remains in the **Warn** stage for at least one MINOR release before removal, giving users a migration window.

```text
release N     release N+1        release N+2
deprecate  →  warn on use    →   remove (MAJOR)
```

---

## Output-Format Versioning

Machine-readable output is part of the contract and evolves under the same rules:

- Adding a field is MINOR; consumers must tolerate unknown fields.
- Renaming, removing, or changing the type of a field is breaking.
- If output schemas ever need independent evolution, a top-level `"schemaVersion"` field will be introduced and that decision recorded in `../decisions/`.

---

*A version number is a promise. SemVer is how we keep it precise.*
