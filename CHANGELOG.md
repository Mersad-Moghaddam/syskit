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

## [0.1.0] - 2026-07-12

### Added

- **Disk I/O and filters (DSK-01):** `syskit disk --io` derives per-device
  read/write rates from two `/proc/diskstats` snapshots; capacity output now
  filters by mount point, filesystem type, or source device.
- **Filesystem command (FS-01):** `syskit filesystem` shows mount sources,
  types, options, and inode availability; pseudo filesystems are hidden by
  default and can be included with `--show-pseudo`.
- **Memory command (MEM-01):** `syskit memory` reports byte-normalized memory,
  cache, swap, and optional PSI pressure data from procfs. Missing
  `MemAvailable` and PSI remain unavailable rather than being fabricated.
- **CPU utilization (CPU-02):** `syskit cpu` samples `/proc/stat` twice to
  report aggregate utilization, with `--per-core` and configurable `--interval`
  for logical-CPU rows. Raw counters remain available in JSON output.
- **CPU static command (CPU-01):** `syskit cpu` reports logical/physical
  topology, sockets, model, architecture, flags, and optional cpufreq values
  from native Linux interfaces. Timed utilization and `--per-core` land in the
  follow-up CPU-02 slice.
- **System command (SYS-01):** `syskit system` reports host, distribution,
  kernel, architecture, uptime, boot time, and load averages in table or JSON
  output from native Linux interfaces, with fixture, golden, integration, and
  benchmark coverage.
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
- **Collector contract (FND-05):** a generic `Collector[T]` snapshot interface plus a mutable-state-free `Registry` for name-based collector discovery, with domain sentinels (`ErrParse`, `ErrFieldMissing`) and a documented "optional-missing is unavailable, not an error" rule. Collectors take `SysFS` by injection and never touch the OS directly.
- **Render layer (FND-06):** a `Renderer` interface with deterministic, golden-testable table and JSON formatters (snake_case, explicit units, no color in structured output); YAML kept as a distinct deferred seam for v0.2.
- **Error handling & exit codes (FND-07):** CLI-boundary `present()` mapping sentinels to the canonical exit codes (0 success, 1 general, 2 usage, 3 permission, 4 unsupported, 5 partial) with a `PartialError` type for partial-collection failures.
- **Logging (FND-08):** structured `log/slog` diagnostics on stderr, silent by default, raised by `--verbose`/`--debug` and silenced by `--quiet` (precedence quiet > debug > verbose); lower layers never log.
- **Configuration loading (FND-09):** optional TOML config with precedence flags > env (`SYSKIT_*`) > per-command `[section]` > global > defaults, XDG discovery, `--config`, and the documented env-outranks-per-command-section rule; missing file is silent, malformed file is an error.
- **Test harness (FND-11):** a `golden` helper (`Assert`/`Read`, `-update` regeneration), the cross-package `testdata/` layout, and a read-only `scripts/capture-fixtures.sh` that records fixture provenance.

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
