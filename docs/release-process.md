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

## Artifact Builds

Run `scripts/build-release.sh vX.Y.Z` on a clean tagged checkout to create
Linux amd64 and arm64 archives with embedded versions and a `SHA256SUMS` file.
Each archive includes the static binary, MIT license, and `syskit(1)` manual.
Tar ownership,
timestamps, ordering, and gzip metadata are normalized using `SOURCE_DATE_EPOCH`.
The tag-triggered release workflow publishes the same artifacts on GitHub.

For Debian-family systems, run `scripts/build-deb.sh vX.Y.Z [amd64|arm64]`.
The resulting package installs the static binary at `/usr/bin/syskit`, the MIT
license under `/usr/share/doc/syskit`, and the compressed manual under
`/usr/share/man/man1`. Package creation does not install or modify the local
system.

For RPM-family systems, install `rpmbuild` and run
`scripts/build-rpm.sh vX.Y.Z [amd64|arm64]`. The result installs the binary under
`/usr/bin` and the license under the distribution's licensedir. Generate AUR
submission metadata after the release archives with
`scripts/build-aur.sh vX.Y.Z [release-dir] [output-dir]`; its archive contains a
`PKGBUILD` and `.SRCINFO` for `syskit-bin`, pinned to both archive checksums.

After every artifact is present, run `scripts/write-checksums.sh [artifact-dir]`
so `SHA256SUMS` covers archives, Debian packages, RPM packages, and AUR metadata.
The release workflow performs these steps in order for both supported architectures.

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
