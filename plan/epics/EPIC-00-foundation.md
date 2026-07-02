# EPIC-00 — Foundation & Delivery Infrastructure

> **Milestone:** v0.1 · **Sprints:** 0–1 · **Points:** 58
> Turn the design-first planning repository into a compiling, tested Go project — without shipping user-facing features prematurely.

---

## Goal

Establish everything the first feature needs to exist: the Go module, the approved repository layout, the CLI entry point, the platform abstraction that makes collectors testable, the render layer, cross-cutting patterns (errors, logging, config), and a Go-aware CI pipeline. When this epic is Done, a contributor can add a new command as a clean vertical slice and have it tested and merged.

## Why it comes first

`../docs/implementation-readiness.md` forbids production Go code until readiness is signed off, and `../docs/project-structure.md` forbids `cmd/`/`internal/`/`pkg/` during planning. This epic executes the **transition pull request** those documents describe, flipping the repo from planning-phase to implementation-phase, then builds the scaffolding every later epic depends on.

## Success criteria

- The implementation-readiness checklist is fully checked and signed off.
- `main` compiles, `go test -race ./...` passes, CI is green with the Go pipeline.
- The platform `SysFS` interface exists with `RealFS` (rooted at `/`) and a fixture-backed `TestFS`.
- Table and JSON formatters render a trivial model behind the `Formatter` interface.
- Error, logging, and configuration conventions are code, not just specs.
- README status, getting-started build steps, and contributing rules reflect the implementation phase.

## Stories

See [`../product-backlog.md`](../product-backlog.md#epic-00--foundation--delivery-infrastructure) for the authoritative list. Summary:

| ID | Story | Pts | Sprint |
|---|---|---|---|
| FND-01 | Implementation-readiness sign-off. | 3 | 0 |
| FND-02 | Transition PR: Go module + layout + status/CI/contributing updates. | 5 | 0 |
| FND-03 | Cobra CLI bootstrap (root, `--format`, `--help`, `version`). | 5 | 0 |
| FND-04 | `SysFS` interface + `RealFS` + `TestFS`. | 8 | 0 |
| FND-05 | Collector interface + registration. | 5 | 1 |
| FND-06 | Render layer: `Formatter` + table + JSON. | 8 | 1 |
| FND-07 | Error-handling patterns (sentinels, `%w`, exit codes). | 3 | 1 |
| FND-08 | Logging scaffolding (structured, off by default). | 3 | 1 |
| FND-09 | Configuration loading (flags > env > file > default). | 5 | 1 |
| FND-10 | Go CI pipeline (fmt/vet/race/integration/coverage/bench/vuln). | 8 | 0–1 |
| FND-11 | Test harness (golden helper, `testdata/`, capture script). | 5 | 1 |

## Key dependencies & references

- `../decisions/001-use-go.md`, `../decisions/004-layered-architecture.md`, `../decisions/005-cobra-for-cli.md`
- `../specs/architecture.md`, `../specs/collectors.md`, `../specs/rendering.md`
- `../specs/error-handling.md`, `../specs/logging-strategy.md`, `../specs/configuration.md`
- `../specs/testing-strategy.md`, `../.github/workflows/ci.yml`

## Architectural guardrails (must hold from the first commit)

- Dependency direction is strict: `CLI → Command → Service → Collector → Platform → kernel`. Lower layers never import higher ones.
- Collectors receive `SysFS` by injection and never touch the OS filesystem directly.
- No `runtime.GOOS` branching; SysKit is Linux-only (`../decisions/002-linux-only.md`).
- No `util`/`common`/`helpers` grab-bag packages (`../standards/naming-conventions.md`).

## Definition of Done for the epic

Every story meets `../standards/definition-of-done.md`, **and** the CI boundary check in `../.github/workflows/ci.yml` is replaced by real Go stages, **and** the README no longer advertises "Design Phase."

## Risks

- **R-03 fixture drift** and **R-11 estimation immaturity** are highest here — this is the team's first Go delivery. Treat Sprint 0–1 velocity as warm-up.
