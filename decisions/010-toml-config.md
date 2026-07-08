# 010. Use BurntSushi/toml for configuration file parsing

**Status:** Accepted, 2026-07-08

> **Verification (2026-07-08):** `github.com/BurntSushi/toml` v1.6.0 adopted.
> License **MIT** — confirmed. Transitive dependencies: **zero** (`go mod graph`
> shows no edges out of the module). Actively maintained, pure-Go. This ADR
> records the decision required by `standards/dependency-policy.md` before the
> dependency is used in FND-09 (configuration loading).

---

## Context

SysKit's configuration spec (`specs/configuration.md`) commits to **TOML** as
the configuration file format: "SysKit uses TOML for its configuration file.
TOML is unambiguous, comment-friendly, and maps cleanly onto typed configuration
structs." The loading model in that spec is written directly against a
`toml.Unmarshal(data, cfg)` call.

The Go standard library has **no TOML decoder** — `encoding/json` exists, but
there is no `encoding/toml`. Story FND-09 (configuration loading, precedence
`flags > env > file > default`) requires reading a TOML file into a typed
`Config` struct. The config file is optional and read once at CLI startup, so
the parser is confined to the CLI layer and never touches the collection hot
path.

Hand-rolling a TOML parser is unreasonable: TOML's grammar (tables, arrays of
tables, inline tables, typed values, datetimes, quoting rules) is large enough
that a correct in-repo implementation would be a maintenance liability — exactly
the exception the **Minimal Dependencies** policy
(`standards/dependency-policy.md`) permits: "the standard library does not
provide the required functionality" and "building the equivalent from scratch
would be unreasonable."

`ARCHITECTURE.md` §4 fixes the boundary: configuration is "loaded once at CLI
startup and threaded down as plain values"; lower layers never read config
files. This ADR resolves which library performs that load.

## Decision

We will add **`github.com/BurntSushi/toml`** (MIT) as an approved direct
dependency and use it for decoding the optional TOML configuration file in the
CLI layer.

- It is confined to the **CLI layer** (`internal/cli`), mirroring how Cobra is
  confined to the CLI/command layers and `goccy/go-yaml` (ADR 009) to the render
  layer. Collectors, services, the platform layer, and the render layer never
  import it.
- It is used only for **input** (parsing `config.toml`). It is not part of any
  output contract, so it carries none of the schema-stability obligations that
  govern the JSON/YAML renderers.
- A missing config file is not an error (defaults apply); only a malformed file
  is surfaced to the user, per `specs/configuration.md`.

### Adoption gate (per dependency-policy.md)

Confirmed at adoption in the `build(deps):` change:

- License is **MIT** — compatible, no GPL/AGPL.
- Transitive-dependency count is **zero** (`go mod graph` shows no outbound
  edges), the smallest possible footprint.
- `govulncheck ./...` is clean for the pinned version (v1.6.0) in CI; the local
  sandbox blocks egress to `vuln.go.dev`, so the authoritative check runs on the
  CI runner.

## Evaluation Against the 7-Criteria Test

Held to the same bar as cobra, bubbletea, lipgloss, testify, and goccy/go-yaml.

| Criterion | `BurntSushi/toml` (chosen) | Assessment |
|---|---|---|
| **Necessity** | stdlib has no TOML decoder; a spec-mandated config format; in-repo build unreasonable | **Pass** |
| **License** | **MIT** | **Pass** — MIT-compatible |
| **Maintenance** | Actively maintained, pure-Go, the long-standing reference TOML implementation for Go | **Pass** |
| **Popularity** | The de facto standard Go TOML library, widely battle-tested | **Pass** |
| **Transitive deps** | **Zero** — no outbound module edges | **Pass** (best possible) |
| **Security** | No known advisories; re-verified via `govulncheck` at the CI adoption gate | **Pass** (gated) |
| **Footprint** | Modest; CLI-startup path only, never in the hot collection path | **Pass** |

## Consequences

### Positive
- SysKit can honor its TOML configuration contract with a maintained,
  MIT-licensed, pure-Go library, preserving the single static binary (no cgo).
- Confined to the CLI layer, so it is replaceable and the lower layers stay
  dependency-free and independently testable.
- Zero transitive dependencies — the smallest possible addition to the
  dependency surface.

### Negative
- One more direct dependency to track, update, and scan — accepted under the
  documented Minimal Dependencies exception.

### Neutral
- Scoped to configuration input only; it takes on no output-schema obligations.
- We adopt the library's default decoding conventions except where the config
  spec's precedence rules dictate otherwise (precedence is applied in SysKit
  code after the file is decoded, not by the library).

## Alternatives Considered

- **`pelletier/go-toml/v2`.** A strong, actively-maintained alternative with a
  fast decoder. Rejected as the primary choice only on footprint parity:
  `BurntSushi/toml` has zero transitive dependencies and is the historical
  reference implementation the config spec's examples read most naturally
  against. `go-toml/v2` remains a viable drop-in fallback if `BurntSushi/toml`'s
  health regresses; switching requires only superseding this ADR.
- **Hand-rolled in-repo TOML parser.** Zero dependencies. Rejected: TOML is far
  too large to implement and maintain correctly for a configuration-input
  concern — the textbook case the dependency policy says to delegate.
- **Switch the config format to JSON (stdlib only).** Would remove the need for
  any dependency. Rejected because `specs/configuration.md` deliberately chose
  TOML over JSON ("easier to read and less error-prone by hand than JSON") as a
  committed design decision; changing the user-facing config format to dodge a
  zero-transitive-dep, MIT library is the wrong trade.
- **Drop configuration-file support entirely.** Would leave only flags and env
  vars. Rejected because a user-persistent config file is a committed feature of
  `specs/configuration.md` (FND-09).

## References

- [Dependency Policy](../standards/dependency-policy.md) — 7-criteria test, Minimal Dependencies operationalized
- [Constitution](../specs/constitution.md) — principle 8 (Minimal Dependencies)
- [Configuration](../specs/configuration.md) — TOML config format, precedence, loading model
- [ADR 009 — YAML encoding](./009-yaml-encoding.md) — the parallel decision for the render layer's YAML output
- [ARCHITECTURE.md](../ARCHITECTURE.md) — §4, configuration loaded once at the CLI layer and threaded down as plain values
