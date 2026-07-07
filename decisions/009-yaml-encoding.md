# 009. Use goccy/go-yaml for YAML output encoding

**Status:** Accepted, 2026-07-07

> **Verification (2026-07-07):** Maintenance/license claims re-checked against
> live sources. `goccy/go-yaml` is actively maintained (latest release
> 2026-01-08) and MIT — confirmed. `gopkg.in/yaml.v3` (go-yaml/yaml) was
> archived/marked unmaintained by its author on 2025-04-01 — confirmed. One
> correction applied: the maintained continuation of yaml.v3 is now the YAML
> org's drop-in fork **`go.yaml.in/yaml/v3`** (Apache-2.0/MIT), so the fallback
> below points there rather than at the frozen `gopkg.in/yaml.v3` import path.
> The core decision (adopt `goccy/go-yaml`) is unchanged.

---

## Context

SysKit commits to YAML as a first-class output format alongside table and JSON
(`specs/features/output-formats.md`, `specs/rendering.md`, roadmap **v0.2 /
FR-10**). The rendering spec fixes the contract: YAML **mirrors the JSON
structure**, uses `snake_case` field names, carries explicit units, and never
emits lossy human formatting in numeric fields.

The Go standard library has **no YAML encoder** — `encoding/json` exists, but
there is no `encoding/yaml`. YAML's specification (anchors, flow vs. block
styles, multi-document streams, quoting rules) is large enough that a
hand-rolled in-repo encoder would be unreasonable to build and maintain, which
is exactly the exception the **Minimal Dependencies** policy
(`standards/dependency-policy.md`) permits: "the standard library does not
provide the required functionality" and "building the equivalent from scratch
would be unreasonable."

`specs/features/output-formats.md` and the baseline `ARCHITECTURE.md` §8 both
flagged the YAML library choice as an unresolved open question ("YAML dependency
policy must be reviewed before implementation"). This ADR resolves it. Only the
*decision* is needed now; YAML output ships in v0.2, not v0.1.

## Decision

We will add **`github.com/goccy/go-yaml`** (MIT) as an approved direct
dependency and use it for YAML encoding in the render layer.

- It is confined to the **render layer** (`internal/render`), mirroring how
  Cobra is confined to the CLI/command layers and Bubble Tea to the CLI layer.
  Collectors, services, and the platform layer never import it.
- The YAML renderer marshals the **same domain models** the JSON renderer uses.
  The renderer is configured so YAML field names mirror the JSON contract
  (`snake_case`, same shapes) — either by honoring the existing `json:` tags or
  by a JSON-first conversion — so the model layer maintains **one** canonical
  tag set, not a divergent YAML schema. The precise marshalling mechanism is an
  implementation detail for the v0.2 render work; the contract (YAML mirrors
  JSON, per `specs/rendering.md`) is binding regardless.

### Adoption gate (per dependency-policy.md)

This ADR records the decision; the actual `go get` happens in the v0.2 render
PR. At that point the adopter must re-confirm, in the `build(deps):` commit:

- `govulncheck ./...` is clean for the pinned version.
- The library is still actively maintained at adoption time.
- The transitive-dependency count is still minimal.

If `goccy/go-yaml`'s health has materially regressed by v0.2, the documented
fallbacks — in preference order — are `sigs.k8s.io/yaml` (JSON-first wrapper,
guarantees JSON mirroring by construction) then `go.yaml.in/yaml/v3` (the YAML
org's actively-maintained, drop-in continuation of the archived
`gopkg.in/yaml.v3`). Switching to a fallback requires only superseding this ADR,
not re-litigating the format.

## Evaluation Against the 7-Criteria Test

Held to the same bar as cobra, bubbletea, lipgloss, and testify.

| Criterion | `goccy/go-yaml` (chosen) | Assessment |
|---|---|---|
| **Necessity** | stdlib has no YAML encoder; a spec-mandated format; in-repo build unreasonable | **Pass** |
| **License** | **MIT** | **Pass** — MIT-compatible, no GPL/AGPL |
| **Maintenance** | Actively maintained, pure-Go project (latest release 2026-01-08); chosen over `gopkg.in/yaml.v3` specifically because the long-standing `go-yaml/yaml` project was archived/marked unmaintained on 2025-04-01, weakening its Maintenance score | **Pass** (differentiator) |
| **Popularity** | Widely adopted, increasingly the community default for actively-maintained Go YAML | **Pass** |
| **Transitive deps** | Minimal — pure Go, no heavy transitive tree | **Pass** |
| **Security** | No known advisories at decision time; re-verified via `govulncheck` at the adoption gate above | **Pass** (gated) |
| **Footprint** | Modest; render-path only, not in the hot collection path; test/output surface | **Pass** |

## Consequences

### Positive
- SysKit can honor its YAML output contract with a maintained, MIT-licensed,
  pure-Go library, preserving the single static binary (no cgo).
- Confined to the render layer, so it is replaceable and the lower layers stay
  dependency-free and independently testable.
- Reusing the JSON domain models keeps one tag set and guarantees YAML mirrors
  JSON, satisfying `specs/rendering.md` by construction.

### Negative
- One more direct dependency to track, update, and scan — accepted under the
  documented Minimal Dependencies exception.
- YAML's format flexibility (quoting, block vs. flow) needs golden-file tests to
  keep output deterministic (already required by `specs/testing-strategy.md`).

### Neutral
- The choice is scoped to v0.2+; v0.1 ships with table + JSON only and takes on
  no YAML dependency.
- We adopt the library's default encoding conventions except where the render
  contract dictates otherwise.

## Alternatives Considered

- **`gopkg.in/yaml.v3` (go-yaml/yaml).** The historical de-facto standard,
  MIT/Apache-2.0. Rejected as the primary choice because the upstream project
  was archived and marked unmaintained by its author on 2025-04-01, which fails
  the spirit of the **Maintenance** criterion for a long-lived reference
  project. Its maintained continuation — the YAML org's drop-in fork
  `go.yaml.in/yaml/v3` — is retained as a documented fallback in its place.
- **`sigs.k8s.io/yaml`.** A JSON-first wrapper that marshals via `encoding/json`
  then converts to YAML, guaranteeing "YAML mirrors JSON" by construction and
  letting structs carry only `json:` tags. Rejected as primary because it
  transitively depends on the older `yaml.v2`, enlarging the transitive-dep
  surface; kept as the first fallback because its JSON-mirroring guarantee is
  architecturally attractive.
- **Hand-rolled in-repo YAML encoder.** Zero dependencies. Rejected: YAML is far
  too large to implement and maintain correctly for a formatting concern — the
  textbook case the dependency policy says to delegate.
- **Drop YAML support entirely.** Would remove the need for any dependency.
  Rejected because YAML is a committed product feature (product overview,
  roadmap v0.2/FR-10, output-formats spec).

## References

- [Dependency Policy](../standards/dependency-policy.md) — 7-criteria test, Minimal Dependencies operationalized
- [Constitution](../specs/constitution.md) — principle 8 (Minimal Dependencies), principle 9 (Consistent CLI Experience)
- [Output Formats](../specs/features/output-formats.md) — YAML as a first-class format
- [Rendering Architecture](../specs/rendering.md) — "YAML mirrors JSON structure", snake_case, explicit units
- [Roadmap](../specs/roadmap.md) — v0.2 YAML output formatter (FR-10)
- [ARCHITECTURE.md](../ARCHITECTURE.md) — §8 Open Question resolved by this ADR
