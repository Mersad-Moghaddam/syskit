# System Feature Specification

## Purpose

Provide a concise overview of the host: identity, kernel, OS release, uptime, boot time, architecture, load averages, and virtualization hints where available.

## User Story

As a Linux user, I want one command that summarizes the machine I am inspecting so I can quickly understand the host before digging into specific resources.

## Motivation

Incident response and debugging often begin with basic orientation: what host is this, what kernel is running, how long has it been up, and is the system under load? SysKit should make that first glance consistent and scriptable.

## Requirements

- Report hostname, kernel release, kernel version, architecture, OS name, OS version, uptime, boot time, and load averages.
- Read from Linux-native sources where possible.
- Support table, JSON, and YAML output.
- Show partial data when optional files are unavailable.
- Avoid shelling out to `uname`, `uptime`, or `hostname`.

## Linux Concepts

- `/proc/uptime`
- `/proc/loadavg`
- `/etc/os-release`
- `/proc/sys/kernel/hostname`
- `/proc/sys/kernel/osrelease`
- `/proc/sys/kernel/version`
- Kernel architecture from runtime or system interfaces

## Expected CLI

```sh
syskit system
syskit system --format json
syskit system --verbose
```

## Expected Output

Default table output should fit in a normal terminal:

```text
HOST        OS              KERNEL          ARCH    UPTIME      LOAD
workstation Ubuntu 24.04    6.8.0-31-generic x86_64 3d 04h 12m  0.42 0.35 0.30
```

JSON output should expose raw fields such as `hostname`, `os_name`, `os_version`, `kernel_release`, `kernel_version`, `architecture`, `uptime_seconds`, `boot_time`, and `load_average`.

## Edge Cases

- `/etc/os-release` is missing or malformed.
- Host is inside a container with limited kernel metadata.
- Load average file contains unexpected fields.
- Uptime is available but boot time cannot be derived reliably.
- Hostname differs between `/proc/sys/kernel/hostname` and resolver configuration.

## Acceptance Criteria

- The command succeeds on a standard Linux host without elevated privileges.
- Missing optional OS metadata does not fail the command.
- JSON output contains numeric uptime in seconds.
- Table output remains readable under 100 columns.
- No external Linux command is executed to gather data.

## Learning Objectives

- Understand the difference between kernel identity and distribution identity.
- Learn how uptime and load average are exposed through procfs.
- Learn why load average is not the same as CPU utilization.

## Estimated Complexity

Low.

## Dependencies

- CLI conventions.
- Platform file-reading abstraction.
- Table and structured renderers.

## Future Extensions

- Virtualization detection.
- Hardware vendor and model.
- Boot ID and machine ID reporting.
- Optional timezone and locale summary.
