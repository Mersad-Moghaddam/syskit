# Plugin Architecture

> Planned extension model for custom collectors and community integrations.

Plugins are a future milestone. The core project should first stabilize built-in collectors, output schemas, and service contracts. This document defines the direction so early architecture does not block extensibility.

## Goals

- Allow users to add custom collectors without modifying SysKit core.
- Keep plugins isolated from core process failures where practical.
- Preserve stable output contracts.
- Avoid unsafe dependency on Go compiler or runtime ABI compatibility.
- Make plugin discovery explicit and auditable.

## Preferred Model

SysKit should prefer an out-of-process plugin protocol over Go's in-process `plugin` package. External plugin processes can communicate with SysKit using a documented JSON protocol over stdin/stdout or a local socket.

This model trades a little overhead for better isolation, clearer versioning, and fewer ABI surprises.

## Plugin Manifest

Each plugin should provide a manifest with:

- Plugin name.
- Version.
- Supported SysKit plugin API version.
- Commands or collectors provided.
- Required permissions.
- Output schemas.
- Author and license.

SysKit should refuse to load plugins with incompatible API versions unless the user explicitly opts into experimental behavior.

## Security Boundaries

Plugins are user-installed executable code. SysKit should:

- Never auto-install plugins.
- Never load plugins from world-writable directories.
- Show plugin path and permissions in diagnostic output.
- Make plugin execution opt-in.
- Document the trust model clearly.

## Discovery

Planned discovery locations:

1. Explicit `--plugin-dir` flag.
2. `$SYSKIT_PLUGIN_DIR`.
3. `$XDG_DATA_HOME/syskit/plugins`.
4. `~/.local/share/syskit/plugins`.

Core commands must work without plugins.

## Output Integration

Plugin output should enter SysKit as structured data, not pre-rendered terminal text. SysKit core remains responsible for table, JSON, YAML, and TUI rendering so the user experience stays consistent.

## Protocol v1

An executable declared by a compatible manifest receives one JSON request on
stdin: `{"api_version":"v1","action":"collect"}`. It must write exactly one
JSON value to stdout and use stderr for diagnostics. SysKit enforces a bounded
timeout, rejects executable paths outside the plugin directory, and renders the
returned value itself. Execution occurs only through an explicit `plugins run`.

## Versioning

The plugin protocol should be versioned independently from SysKit's CLI version. A SysKit release may support multiple plugin protocol versions during migration windows.

## Acceptance Criteria

- Core architecture can register new collectors without command rewrites.
- Plugin execution is opt-in and visible.
- Plugin output can be rendered by normal SysKit renderers.
- Plugin protocol compatibility is testable.
- Security limitations are documented before release.
