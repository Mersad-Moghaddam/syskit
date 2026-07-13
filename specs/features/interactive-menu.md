# Interactive Menu Feature Specification

## Purpose

Provide a modern, discoverable terminal control center when a user runs
`syskit` without a subcommand in an interactive terminal.

## User Story

As a Linux user who does not remember every command and flag, I want to browse
SysKit's capabilities by domain, run one, and return to the menu so I can
inspect a host without repeatedly consulting help text.

## Motivation

The command-oriented interface remains the stable automation contract, but it
is less approachable for interactive exploration. A hierarchical menu can
compose the existing commands without duplicating collectors or services and
make advanced variants such as per-core CPU sampling and focused diagnostics
visible.

## Requirements

- Start automatically for a bare `syskit` invocation when both input and
  output are interactive terminals.
- Preserve the existing help output for redirected or piped bare invocations.
- Group every implemented user-facing capability into domain submenus.
- Support arrow keys, Vim-style navigation, Enter, Escape, Backspace, and mouse
  selection.
- Show a breadcrumb and description for the current selection.
- Prompt inside the menu for values required by container and plugin commands.
- Run existing Cobra commands rather than collecting or rendering data in the
  menu layer.
- Return to the menu after a selected command exits; allow quitting from any
  menu level without corrupting terminal state.
- Handle small terminals by keeping the selected row visible.

## Expected CLI

```sh
syskit
```

Explicit commands such as `syskit cpu` and `syskit --help` retain their current
non-menu behavior.

## Expected Interaction

```text
SYSKIT // CONTROL CENTER
Home / CPU

> CPU overview       topology, model, frequency, and utilization
  Per-core view      utilization for every logical CPU
```

Selecting a leaf runs the corresponding existing command. Static output waits
for confirmation before returning; interactive views return after their own
quit key is used.

## Edge Cases

- stdin or stdout is redirected.
- The terminal is resized below the number of menu entries.
- A selected collector returns an error or partial data.
- A required container ID or plugin name is empty.
- The user exits an interactive child view with Ctrl-C.

## Acceptance Criteria

- Bare interactive execution opens the menu while bare non-interactive
  execution still prints help.
- Every stable command family in `docs/command-reference.md` is reachable.
- CPU opens a submenu with overview and per-core choices.
- Escape, Left, or Backspace returns one level; it quits only at the root.
- Static and live selections return to the menu after completion.
- Keyboard, mouse, input, resize, and quit transitions have unit tests.
- Existing CLI contract and command tests remain green under the race detector.

## Dependencies

- Cobra command tree.
- Bubble Tea and Lip Gloss adopted by ADR 006.
- Existing command, service, collector, and render layers.

## Non-Goals

- Replacing explicit commands or their stable flags.
- Reimplementing collection or rendering in the menu.
- Persisting menu state between invocations.
- Adding non-Linux behavior.
