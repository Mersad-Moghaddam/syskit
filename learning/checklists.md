# Learning Checklists

> Practical study checks for implementers before writing each feature.

## System And CPU

- [ ] Explain the difference between `/proc/stat` and `/proc/cpuinfo`.
- [ ] Explain why CPU utilization needs two samples.
- [ ] Identify per-core and aggregate CPU counters.
- [ ] Explain jiffies and why raw counters are not percentages.

## Memory

- [ ] Explain `MemTotal`, `MemFree`, `MemAvailable`, buffers, and cache.
- [ ] Explain swap usage and why it does not always mean active pressure.
- [ ] Understand PSI memory pressure signals.
- [ ] Identify which memory fields may be missing on older kernels.

## Disk And Filesystem

- [ ] Explain block devices, partitions, mounts, and filesystems.
- [ ] Read `/proc/self/mountinfo` and identify mount options.
- [ ] Explain inode exhaustion.
- [ ] Explain why disk I/O counters are cumulative.

## Processes

- [ ] Read `/proc/[pid]/stat` safely despite command names with spaces.
- [ ] Explain process states and parent-child relationships.
- [ ] Handle processes that disappear during collection.
- [ ] Understand permission limits for other users' processes.

## Networking And Ports

- [ ] Explain interface counters and common rollover concerns.
- [ ] Interpret TCP socket states.
- [ ] Map socket inodes to process file descriptors.
- [ ] Understand the difference between routing, interfaces, DNS, and sockets.

## Dashboard And Watch Mode

- [ ] Explain sampling intervals and jitter.
- [ ] Avoid mixing stdout data with stderr diagnostics.
- [ ] Understand terminal resize events.
- [ ] Keep collection and rendering separated.

## Containers And Plugins

- [ ] Explain cgroup membership and resource accounting.
- [ ] Identify cgroup v1 and v2 layouts.
- [ ] Describe plugin trust boundaries.
- [ ] Explain why output schemas need versioning.
