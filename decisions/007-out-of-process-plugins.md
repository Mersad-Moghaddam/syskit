# ADR 007: Prefer Out-of-Process Plugins

## Status

Accepted

## Date

2026-07-01

## Context

SysKit plans to support plugins after the core CLI, collectors, and output contracts stabilize. Plugins should allow external collectors and custom checks without requiring users to fork the main project.

Go offers an in-process `plugin` package, but it has practical constraints: ABI compatibility, toolchain matching, platform limitations, and process-safety concerns. SysKit also needs a clear trust boundary because plugins are executable code installed by users.

## Decision

SysKit will prefer an out-of-process plugin model. Plugins should run as separate executables and communicate with SysKit through a versioned protocol, likely using JSON over stdin/stdout or a local socket.

The plugin protocol will be defined before the plugin milestone begins. Core SysKit commands must not require plugins.

## Consequences

Positive:

- Better fault isolation.
- Clearer version compatibility.
- Easier language-agnostic plugin development.
- More explicit trust model.
- Avoids Go plugin ABI constraints.

Negative:

- More protocol design work.
- Higher overhead than in-process calls.
- Requires timeout and process lifecycle management.
- Requires careful validation of plugin output.

## Alternatives Considered

### Go `plugin` package

Rejected as the default because it tightly couples plugin builds to Go versions, platforms, and ABI behavior.

### Static linking only

Rejected because it requires rebuilding SysKit for every extension and does not support user-installed plugins.

### No plugin support

Rejected for the long-term roadmap because extensibility is a core project goal, even though plugins are intentionally delayed until after the core stabilizes.
