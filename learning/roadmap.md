# Learning Roadmap

> Concepts to study before and during SysKit implementation.

SysKit is both a tool and a learning project. The roadmap below follows the order in which concepts are most useful for implementation.

## Stage 1: Linux Inspection Basics

- Understand procfs and sysfs as virtual filesystems.
- Learn how counters differ from gauges.
- Study `/proc/uptime`, `/proc/loadavg`, `/proc/stat`, and `/proc/meminfo`.
- Practice identifying units, monotonic counters, and kernel-provided fields.

Recommended docs:

- [CPU](cpu.md)
- [Memory](memory.md)

## Stage 2: Filesystems And Storage

- Learn mount tables, block devices, partitions, and filesystem statistics.
- Compare `/proc/mounts`, `/proc/self/mountinfo`, `/proc/diskstats`, and `/sys/block`.
- Understand inode exhaustion and why free bytes are not the whole story.

Recommended docs:

- [Disk](disk.md)
- [Filesystem](filesystem.md)

## Stage 3: Processes

- Study `/proc/[pid]` lifetime races.
- Learn process states, parent-child relationships, command lines, and file descriptors.
- Understand why process inspection can produce partial data under normal permissions.

Recommended docs:

- [Process](process.md)

## Stage 4: Networking

- Learn interface counters, routing tables, sockets, and DNS resolver configuration.
- Compare procfs networking files with Netlink.
- Understand IPv4, IPv6, TCP states, UDP sockets, and socket inode mapping.

Recommended docs:

- [Network](network.md)

## Stage 5: Live Monitoring

- Learn sampling, deltas, refresh intervals, and terminal rendering constraints.
- Understand how to avoid misleading rates from short sampling windows.
- Study terminal resize behavior and keyboard interaction models.

Recommended specs:

- [Dashboard feature](../specs/features/dashboard.md)
- [Rendering architecture](../specs/rendering.md)

## Stage 6: Containers And Extensions

- Learn cgroup v1 and v2 differences.
- Study container runtime metadata and process-to-container mapping.
- Understand the trust model for executable plugins.

Recommended specs:

- [Containers feature](../specs/features/containers.md)
- [Plugin architecture](../specs/plugin-architecture.md)

## Study Loop

For each feature:

1. Read the feature spec.
2. Read the related learning document.
3. Identify Linux data sources.
4. Capture representative fixtures.
5. Define parser tests before implementation.
6. Implement the smallest useful vertical slice.
