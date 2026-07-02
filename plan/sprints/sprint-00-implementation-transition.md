# Sprint 00 — Implementation Transition

**Dates:** TBD → +2 weeks · **Milestone:** v0.1 (Foundation) · **Committed points:** 21

## Sprint goal

Flip SysKit from planning-phase to implementation-phase: sign off readiness, create the Go module and approved layout, bootstrap the Cobra CLI, and land the `SysFS` platform abstraction that makes every future collector testable — on a green Go CI.

## Capacity

- Developers available: core team (warm-up sprint — tooling and unknowns dominate).
- Commitment intentionally **below** nominal velocity; this is the team's first Go delivery.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| FND-01 | Implementation-readiness sign-off. | 3 | Committed |
| FND-02 | Transition PR: Go module + approved layout + README/CI/contributing updates. | 5 | Committed |
| FND-03 | Cobra CLI bootstrap (root, `--format`, `--help`, `version`). | 5 | Committed |
| FND-04 | `SysFS` interface + `RealFS` + fixture-backed `TestFS`. | 8 | Committed |
| FND-10a | Go CI pipeline — first half: fmt, vet, `go build`, `go test`. | (part of FND-10) | Committed |

_Authoritative IDs/points in ../product-backlog.md. FND-10 spans Sprints 0–1._

## Task breakdowns

**FND-01 — readiness sign-off**
- [ ] Walk `../../docs/implementation-readiness.md` "Required Before First Code"; check every box or record the gap.
- [ ] PO signs off; record decision.

**FND-02 — transition PR**
- [ ] `go mod init`; choose module path; Go 1.22+.
- [ ] Create `cmd/`, `internal/{cli,collector,platform,render,service}`, `testdata/` per `../../docs/project-structure.md`.
- [ ] Replace the planning-boundary check in `../../.github/workflows/ci.yml` with Go stages.
- [ ] Update README status badge/section; update `../../docs/contributing.md` for code-phase; CHANGELOG entry.

**FND-03 — CLI bootstrap**
- [ ] Root command with Cobra (`../../decisions/005-cobra-for-cli.md`).
- [ ] Persistent `--format` flag (default table), `--help`, `version` subcommand.
- [ ] Wire to render layer placeholder; unit test flag parsing.

**FND-04 — platform abstraction**
- [ ] Define `SysFS` (`ReadFile`/`Open`/`ReadDir`) per `../../specs/testing-strategy.md`.
- [ ] `RealFS()` rooted at `/`; `TestFS(fs.FS)` rooted at fixtures.
- [ ] Unit tests: `TestFS` resolves `proc/stat` → `testdata/...`.

## Definition of Ready (entry gate)

All committed stories passed `../../standards/definition-of-ready.md`. FND-01 itself gates the rest.

## Definition of Done (exit gate)

`../../standards/definition-of-done.md`. Additionally: `main` compiles, `go test ./...` passes, CI green with Go stages, planning-boundary check removed.

## Risks this sprint

- **R-11 (estimation immaturity)** — treat this sprint's velocity as warm-up, not a baseline.
- **R-03 (fixture drift)** — establish the `testdata/` layout and capture conventions now (feeds FND-11 next sprint).

## Dependencies

FND-01 unblocks everything. FND-04 unblocks all collectors. This sprint has no external blockers.

## Notes

Do not add any user-facing feature this sprint — foundation only, per `../../docs/implementation-readiness.md` "v0.1 Implementation Entry Criteria."
