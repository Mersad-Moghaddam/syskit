# Filesystem Feature Specification

## Purpose

Inspect mounted filesystems, inode usage, filesystem types, mount options, and namespace-aware mount details.

## User Story

As a Linux user, I want to understand how filesystems are mounted and whether inode or option-related issues may affect the system.

## Motivation

A filesystem can fail operationally even when free bytes remain, especially when inodes are exhausted or mount options are wrong. SysKit should make filesystem health visible beyond basic disk capacity.

## Requirements

- Report mount point, source, filesystem type, mount options, total inodes, free inodes, and inode use percent.
- Parse mount information from `/proc/self/mountinfo`.
- Use filesystem statistics for inode counts.
- Support filters for filesystem type, mount point, and pseudo filesystems.
- Support table, JSON, and YAML output.

## Linux Concepts

- `/proc/self/mountinfo`
- Mount namespaces
- Inodes
- Filesystem types
- Mount options
- Pseudo filesystems such as procfs, sysfs, tmpfs, and cgroupfs

## Expected CLI

```sh
syskit filesystem
syskit filesystem --type ext4
syskit filesystem --show-pseudo
syskit filesystem --format json
```

## Expected Output

```text
MOUNT       TYPE  SOURCE            INODES USED  INODES FREE  IUSE%  OPTIONS
/           ext4  /dev/nvme0n1p2    1.9M         28.2M        6%     rw,relatime
/run        tmpfs tmpfs             1.9M         1.9M         1%     rw,nosuid,nodev
```

Structured output should keep mount options as arrays, not comma-joined strings.

## Edge Cases

- Mount points contain spaces or escaped characters.
- Filesystems do not expose inode counts.
- Overlay and bind mounts duplicate sources.
- Containers expose a different mount namespace from the host.
- Very long option lists exceed terminal width.

## Acceptance Criteria

- Mountinfo parsing handles escaped paths.
- Pseudo filesystems are hidden by default or clearly marked.
- Inode unavailable is distinct from zero inodes.
- JSON output represents options as an array.
- Table output remains readable with long mount points.

## Learning Objectives

- Understand Linux mount namespaces.
- Learn why inode usage matters.
- Study mount option semantics and pseudo filesystems.

## Estimated Complexity

Medium.

## Dependencies

- Disk feature.
- Platform file-reading abstraction.
- Filesystem learning material.

## Future Extensions

- Directory-level usage analysis.
- Large file discovery.
- Filesystem-specific warnings.
- Namespace comparison.
