# Disk Feature Specification

## Purpose

Inspect block devices, partitions, mount usage, and disk I/O counters.

## User Story

As a developer or operator, I want to see storage capacity and I/O activity so I can detect full filesystems, busy disks, and unusual device layouts.

## Motivation

Storage data spans mount tables, filesystem statistics, block device sysfs entries, and cumulative disk counters. SysKit should present capacity and activity without forcing users to combine `df`, `lsblk`, and `iostat`.

## Requirements

- List mounted filesystems with size, used, available, use percent, mount point, and filesystem type.
- List block devices and partitions where available.
- Report disk I/O counters from `/proc/diskstats`.
- Support filtering by mount point, filesystem type, and device.
- Support table, JSON, YAML, and watch mode for I/O rates.

## Linux Concepts

- `/proc/self/mountinfo`
- `/proc/mounts`
- `statfs`
- `/proc/diskstats`
- `/sys/block`
- Block devices and partitions

## Expected CLI

```sh
syskit disk
syskit disk usage --mount /
syskit disk io --watch
syskit disk --format json
```

## Expected Output

```text
DEVICE    TYPE  SIZE      USED      AVAIL     USE%  MOUNT
/dev/nvme0n1p2 ext4  468.0 GiB  210.4 GiB  233.8 GiB  47%   /
```

I/O output should show read/write operations, bytes, and derived rates when sampling is active.

## Edge Cases

- Bind mounts and overlay filesystems duplicate usage.
- Network filesystems may not expose normal block-device data.
- Removable devices disappear during collection.
- Diskstats counters are cumulative and require deltas for rates.
- Permission or namespace restrictions hide devices.

## Acceptance Criteria

- Capacity data comes from filesystem statistics, not parsed `df` output.
- I/O rates are derived only from two samples.
- Mountinfo parsing handles escaped spaces in mount paths.
- Duplicate or pseudo filesystems are identified clearly.
- Structured output exposes raw bytes and counters.

## Learning Objectives

- Understand mount tables and block devices.
- Learn the difference between filesystem usage and device I/O.
- Study cumulative disk counters and rate calculation.

## Estimated Complexity

High.

## Dependencies

- Filesystem feature.
- Collector architecture.
- Disk learning material.

## Future Extensions

- SMART health integration through optional adapters.
- Filesystem growth trends.
- Per-process I/O correlation.
- RAID and LVM visibility.
