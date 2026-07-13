# Changelog

> All notable changes to SysKit are recorded in this file.

---

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

SysKit is in active development. Changes are tracked under `Unreleased` until
their milestone is tagged, then recorded in a dated release entry below.

---

## [Unreleased]

### Added

- **Command reference:** consolidated documentation now covers every implemented
  command, interactive keybinding, and plugin/container safety boundary.

- **Release artifacts:** tagged releases build Linux amd64/arm64 archives with
  embedded versions and SHA-256 checksums.

- **Diagnostics:** `syskit diagnostics` reports explainable, read-only memory
  pressure, swap-use, and filesystem-capacity findings.

- **Plugin discovery:** `syskit plugins list` reads declared plugin manifests
  from explicit or documented directories without executing plugin code.
- **Plugin inspection:** `syskit plugins inspect <name>` displays manifest
  compatibility and declared permissions without executing the plugin.
- **Plugin protocol:** `syskit plugins run <name>` explicitly executes a
  compatible, path-confined plugin with bounded JSON stdin/stdout and timeout.
- **Plugin SDK example:** a buildable protocol-v1 collector and authoring guide
  document manifests, schemas, permissions, lifecycle, and trust boundaries.
- **Manifest validation:** discovery rejects manifests missing executable,
  collector, permission, schema, author, or license declarations.

## [0.4.0] - 2026-07-13

### Added

- **Cgroup foundation:** explicit cgroup v1/v2 detection, normalized process
  memberships, and optional memory, CPU, read, and write counters.
- **Container-aware processes:** `syskit process --containers` limits results
  to processes with a recognized runtime-style container ID.
- **Container commands:** `syskit containers` groups processes and reports
  cgroup counters; `containers inspect <id>` expands one mapping into processes.
- **Partial data:** structured container output marks permission-restricted
  mappings with `partial: true` and omits unavailable controllers.

## [0.3.0] - 2026-07-13

### Added

- **Live monitoring:** interactive dashboard, generic watch mode, and top-like
  process monitoring with bounded refresh intervals and clean terminal exit.
- **Live summaries:** dashboard CPU, memory, swap, disk, network throughput,
  and top-process summaries; `top` supports keyboard sorting and scrolling.
- **Interactive top:** `syskit top` refreshes a filterable process view with
  CPU, memory, name, and PID sort controls plus `j`/`k` scrolling.
- **Watch mode:** `syskit watch <command> --interval` continuously refreshes
  the same in-process table command until interrupted.
- **Dashboard foundation:** `syskit dashboard` uses the existing services,
  handles terminal resizing and collection errors, refuses non-TTY output, and
  prevents overlapping refresh collection.

## [0.2.0] - 2026-07-13

### Added

- **Process identities:** `syskit process` resolves UID values to names from
  `/etc/passwd`, supports `--user <name>`, and includes raw start-time ticks.
- **Process resource usage:** `syskit process` reports memory percentage and
  can derive aggregate CPU percentage with `--interval`.
- **Process partial data:** structured process output marks permission-restricted
  procfs snapshots as partial while preserving readable rows.
- **Network interface metadata:** `syskit network` and `syskit network
  interfaces` now report sysfs operational state, MTU, and MAC address with
  procfs traffic counters.
- **Network views:** `syskit network interfaces`, `syskit network routes`, and
  `syskit network dns` expose the collected interface, route, and resolver
  data as focused table, JSON, or YAML views.
- **Network addresses:** `syskit network interfaces` includes IPv4 and IPv6
  CIDR addresses collected through a native `RTM_GETADDR` Netlink dump.
- **Port ownership:** `syskit ports` now reads TCP, UDP, IPv6, and Unix socket
  tables and best-effort associates socket inodes with owning process IDs and
  commands. Use `--pid`, `--address`, or `--state` to limit results; structured
  output reports permission-restricted ownership scans explicitly.

## [0.1.0] - 2026-07-12

### Added

- **YAML output (OUT-03):** `--format yaml` now mirrors the JSON output schema
  for all non-interactive commands through the ADR-009 approved encoder.

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
