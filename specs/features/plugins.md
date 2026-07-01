# Plugins Feature Specification

## Purpose

Allow advanced users and community maintainers to extend SysKit with custom collectors and command views after the core tool stabilizes.

## User Story

As a power user, I want to add organization-specific system checks or collectors without forking SysKit.

## Motivation

SysKit's core should remain focused, but Linux environments vary widely. Plugins provide a controlled path for domain-specific extensions while preserving core stability.

## Requirements

- Load plugins only from explicit or documented plugin directories.
- Require a manifest that declares plugin identity, version, API compatibility, permissions, and provided collectors.
- Prefer out-of-process execution for isolation.
- Integrate plugin data into normal renderers.
- Provide clear errors for incompatible plugins.

## Linux Concepts

- Executable permissions.
- XDG data directories.
- Process isolation.
- Local IPC.
- Trust boundaries for user-installed code.

## Expected CLI

```sh
syskit plugins list
syskit plugins inspect example
syskit --plugin-dir ./plugins custom-check
```

## Expected Output

```text
NAME           VERSION  API  STATUS      PATH
example-check  0.1.0    v1   compatible  /home/user/.local/share/syskit/plugins/example-check
```

Structured output should expose manifest fields and compatibility status.

## Edge Cases

- Plugin binary is missing or not executable.
- Manifest API version is incompatible.
- Plugin exits with invalid JSON.
- Plugin hangs or exceeds timeout.
- Plugin directory is world-writable.

## Acceptance Criteria

- Plugin loading is opt-in and visible.
- Incompatible plugins do not crash SysKit.
- Plugin output is rendered by core renderers.
- Plugin paths and permissions are shown in verbose diagnostics.
- Security model is documented before release.

## Learning Objectives

- Understand plugin trust and isolation.
- Learn versioned protocols.
- Study process execution and IPC tradeoffs.

## Estimated Complexity

Very High.

## Dependencies

- Plugin architecture.
- Rendering architecture.
- Configuration strategy.
- Security policy.

## Future Extensions

- Signed plugin metadata.
- Plugin registry.
- Sandboxed execution.
- Remote plugin sources with explicit trust prompts.
