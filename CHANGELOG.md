# Changelog

> All notable changes to SysKit are recorded in this file.

---

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

SysKit is in its early design and specification phase. No versioned releases have been
published yet. Until the first release, changes are tracked under the `Unreleased`
section below. Once implementation begins and versions are tagged, each release will be
recorded here with its date and categorized changes.

---

## [Unreleased]

### Added

- Initial project scaffolding and repository structure (`docs/`, `specs/`, `learning/`, `standards/`, `decisions/`, `checklists/`, `scripts/`, `.github/`).
- Engineering Constitution (`specs/constitution.md`) defining the project's core engineering principles.
- Product overview, roadmap, and system architecture specifications under `specs/`.
- Foundational project documentation, including the README, getting-started guide, and contributing guide under `docs/`.
- Community and governance documents: Code of Conduct, Security Policy, and project Governance.
- MIT License.
- **Implementation transition (EPIC-00):** Go module `github.com/Mersad-Moghaddam/syskit` (Go 1.22+) and the approved package layout (`cmd/syskit`, `internal/{cli,collector,platform,render,service,model}`, `testdata/`), replacing the planning-phase repository boundary with a Go CI pipeline.
- **CLI bootstrap (FND-03):** Cobra root command `syskit` with a persistent `--format` flag (`table`/`json`/`yaml`, default `table`, validated), a `version` subcommand, and CLI-boundary exit-code mapping (success `0`, usage error `2`).
- **Platform abstraction (FND-04):** the `SysFS` interface with `RealFS()` (rooted at `/`, reads pseudo-files to EOF) and fixture-backed `TestFS(fs.FS)`, plus platform sentinel errors (`ErrNotFound`, `ErrPermission`, `ErrUnsupported`) — the injectable seam that makes every collector testable against fixtures.
- **CI pipeline (FND-10):** Go stages — gofmt, goimports, `go vet`, build, `go test -race`, integration (`-tags=integration`), coverage, benchmarks, and `govulncheck`.
- **Configuration format decision (ADR 010):** `github.com/BurntSushi/toml` (MIT, zero transitive deps) approved for parsing the optional TOML config file, confined to the CLI layer.

### Changed

- Project status moved from *Design & Specification Phase* to *Implementation (v0.1 Foundation)* across the README and contributing guide; implementation-readiness checklist signed off.

---

<!--
Release entries will follow this structure once versions are tagged:

## [X.Y.Z] - YYYY-MM-DD

### Added
- New features.

### Changed
- Changes in existing functionality.

### Deprecated
- Soon-to-be removed features.

### Removed
- Now removed features.

### Fixed
- Bug fixes.

### Security
- Vulnerability fixes.
-->

*This changelog is maintained by hand until release automation is in place. Every user-visible change should be recorded here before it ships.*
