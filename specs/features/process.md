# Process Feature Specification

## Purpose

List processes, show process trees, expose resource usage, and support filtering by PID, user, name, command, and state.

## User Story

As an engineer debugging a host, I want to see what is running and how resources are distributed so I can identify suspicious or expensive processes.

## Motivation

Process inspection requires careful parsing of `/proc/[pid]` files and must handle processes disappearing while they are being read. SysKit should provide a reliable, scriptable process view without wrapping `ps`.

## Requirements

- List PID, PPID, user, state, command, CPU, memory, start time, and thread count.
- Provide tree view grouped by parent-child relationships.
- Support sorting, filtering, and limiting.
- Map process resource metrics from procfs fields.
- Report partial data gracefully when permissions restrict access.

## Linux Concepts

- `/proc/[pid]/stat`
- `/proc/[pid]/status`
- `/proc/[pid]/cmdline`
- `/proc/[pid]/fd`
- Process states
- Parent and child relationships
- UID to user lookup

## Expected CLI

```sh
syskit process
syskit process --sort cpu --limit 20
syskit process --pid 1234
syskit process --user mersad
syskit process tree
syskit process --format json
```

## Expected Output

```text
PID     USER    STATE  CPU%  MEM%  THREADS  COMMAND
1284    mersad  S      4.2   1.1   18       code
2201    root    S      0.8   0.3   4        sshd
```

JSON output should expose numeric PID, PPID, UID, CPU counters, memory bytes, state code, and command fields.

## Edge Cases

- Process exits between directory listing and file read.
- Command name contains spaces or parentheses.
- `/proc/[pid]/cmdline` is empty for kernel threads.
- Permissions hide environment or file descriptor details.
- PID namespaces change what the user can see.

## Acceptance Criteria

- Disappearing processes do not fail the entire command.
- Parser correctly handles `/proc/[pid]/stat` command names with spaces.
- Filtering and sorting use documented field names.
- Tree output handles orphaned or reparented processes.
- No external `ps` command is executed.

## Learning Objectives

- Understand procfs process representation.
- Learn process states and lifecycle races.
- Learn how process memory and CPU accounting differ from system totals.

## Estimated Complexity

High.

## Dependencies

- Collector architecture.
- CPU and memory services for derived metrics.
- Process learning material.

## Future Extensions

- Signal sending as a separate opt-in command if project scope expands.
- Process environment inspection.
- Open file and socket views.
- Container-aware process grouping.
