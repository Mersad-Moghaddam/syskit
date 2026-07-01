# CPU Feature Specification

## Purpose

Inspect CPU topology, model information, frequency data, cache hints, and utilization derived from Linux CPU counters.

## User Story

As an engineer diagnosing performance, I want to see CPU capacity and current utilization so I can tell whether the host is CPU-bound or misconfigured.

## Motivation

CPU data is spread across `/proc/stat`, `/proc/cpuinfo`, and sysfs topology files. SysKit should normalize the important facts and explain the difference between static CPU identity and sampled CPU activity.

## Requirements

- Report logical cores, physical cores where detectable, sockets, model name, architecture, and CPU flags summary.
- Report aggregate and per-core utilization when two samples are available.
- Report current, minimum, and maximum frequencies where sysfs exposes them.
- Support table, JSON, YAML, and watch mode.
- Keep raw counters separate from derived percentages.

## Linux Concepts

- `/proc/stat` CPU time counters.
- `/proc/cpuinfo` processor metadata.
- `/sys/devices/system/cpu/` topology.
- CPU frequency scaling through cpufreq sysfs files.
- Jiffies and cumulative counters.

## Expected CLI

```sh
syskit cpu
syskit cpu --per-core
syskit cpu --watch --interval 1s
syskit cpu --format json
```

## Expected Output

```text
CPU    MODEL                         CORES  THREADS  UTIL   USER  SYSTEM  IDLE
all    AMD Ryzen 7 7840U             8      16       18.4%  9.1%  4.2%    81.6%
cpu0   AMD Ryzen 7 7840U             -      -        21.0%  11.3% 5.0%    79.0%
```

Structured output should include `cpu_id`, `user`, `system`, `idle`, `iowait`, `steal`, `guest`, `total`, and derived utilization fields.

## Edge Cases

- CPU hotplug changes core count between samples.
- Virtual machines hide topology details.
- cpufreq files are unavailable.
- Counters wrap or appear inconsistent.
- Very short sample windows produce noisy utilization.

## Acceptance Criteria

- Static CPU data works with one sample.
- Utilization calculations require two timestamped samples.
- Per-core rows preserve CPU IDs from `/proc/stat`.
- Missing frequency data is represented as unavailable, not zero.
- Parser tests cover standard, virtualized, and malformed fixture data.

## Learning Objectives

- Understand cumulative CPU counters.
- Learn why utilization is a rate, not a value directly stored by the kernel.
- Study CPU topology and frequency scaling interfaces.

## Estimated Complexity

Medium.

## Dependencies

- Collector architecture.
- Rendering architecture.
- Watch mode conventions.
- CPU learning material.

## Future Extensions

- NUMA topology.
- CPU pressure stall information.
- Thermal throttling hints.
- Scheduler statistics.
