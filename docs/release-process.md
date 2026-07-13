# Release Process

> How SysKit should move from unreleased planning work to stable open-source releases.

SysKit released its first preview, v0.1.0, on 2026-07-12. This process governs
subsequent releases as implementation milestones are completed.

## Versioning

SysKit follows semantic versioning as described in [versioning standards](../standards/versioning.md).

- `0.x` releases may refine CLI contracts while the project is still stabilizing.
- `1.0` means core CLI commands, output schemas, configuration behavior, and extension contracts are stable.
- Breaking changes after `1.0` require a major version.

## Release Types

| Type | Purpose | Example |
|---|---|---|
| Planning milestone | Documentation and architecture readiness | `planning-2026-07` |
| Preview release | Early binaries for testing | `v0.1.0` |
| Patch release | Bug fixes and documentation corrections | `v0.1.1` |
| Stable release | Compatibility commitment | `v1.0.0` |

## Release Checklist

- [ ] All tests pass on Linux CI.
- [ ] Changelog is updated.
- [ ] User documentation matches released behavior.
- [ ] Feature specs reflect implemented behavior.
- [ ] Version number follows the versioning standard.
- [ ] Known limitations are documented.
- [ ] Security notes are reviewed.
- [ ] Release artifacts are reproducible.

## Changelog Policy

The changelog should group user-visible changes by:

- Added
- Changed
- Deprecated
- Removed
- Fixed
- Security

Internal refactors should appear only when they affect users, contributors, or extension authors.

## Release Notes

Release notes should explain:

- What changed.
- Who should upgrade.
- Whether any CLI or output contract changed.
- Known limitations.
- Links to relevant documentation.

## Post-Release Tasks

- Verify release artifacts can be downloaded and executed.
- Open follow-up issues for deferred work.
- Update roadmap status.
- Announce the release in project channels.
