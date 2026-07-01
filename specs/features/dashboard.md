# Dashboard Feature Specification

## Purpose

Provide an interactive terminal dashboard that shows live system metrics across CPU, memory, disk, process, network, and diagnostics views.

## User Story

As a Linux user monitoring a host, I want a single interactive view of important resource metrics so I can observe changes without running multiple commands.

## Motivation

Static commands are excellent for snapshots, but live troubleshooting needs refresh, navigation, sorting, and visual grouping. The dashboard should reuse services and render live data without bypassing the architecture.

## Requirements

- Render live CPU, memory, disk, process, and network summaries.
- Support keyboard navigation between panels.
- Support configurable refresh interval.
- Handle terminal resize events.
- Keep collection separate from rendering.
- Exit cleanly without corrupting terminal state.

## Linux Concepts

- Sampling intervals.
- Rate calculations from cumulative counters.
- Terminal dimensions and TTY behavior.
- Long-running process lifecycle.

## Expected CLI

```sh
syskit dashboard
syskit dashboard --interval 2s
syskit dashboard --panel processes
```

## Expected Output

The dashboard is interactive, but the first screen should include:

- Host and uptime summary.
- CPU utilization.
- Memory and swap usage.
- Disk capacity warnings.
- Network throughput.
- Top processes by CPU or memory.

## Edge Cases

- Terminal is too small for all panels.
- SSH session disconnects.
- Data collection takes longer than the refresh interval.
- Metrics become unavailable during the session.
- User pipes dashboard output to a non-TTY.

## Acceptance Criteria

- Dashboard refuses or degrades gracefully when stdout is not a TTY.
- Refresh interval is bounded to avoid excessive system overhead.
- Resize does not panic or overlap panels.
- Collection errors are displayed without killing the session unless fatal.
- Services used by dashboard are the same services used by static commands.

## Learning Objectives

- Learn event loops and terminal rendering constraints.
- Understand sampling, refresh cadence, and jitter.
- Study separation between data and view state.

## Estimated Complexity

Very High.

## Dependencies

- Rendering architecture.
- CPU, memory, disk, process, and network services.
- Configuration strategy.
- Logging strategy.

## Future Extensions

- Custom layouts.
- Saved dashboard profiles.
- Drill-down panels.
- Remote dashboard over SSH.
