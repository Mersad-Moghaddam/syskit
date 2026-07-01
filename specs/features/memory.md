# Memory Feature Specification

## Purpose

Report physical memory, swap, buffers, cache, available memory, and memory pressure signals in a way that is accurate and easy to interpret.

## User Story

As a Linux operator, I want to understand memory usage without misreading cache as wasted memory so I can identify real pressure and avoid false alarms.

## Motivation

Linux memory accounting is frequently misunderstood. SysKit should expose raw memory fields while also presenting practical summaries that distinguish used memory, reclaimable cache, swap, and pressure.

## Requirements

- Read memory totals and usage from `/proc/meminfo`.
- Report total, used, available, free, buffers, cache, swap total, swap used, and swap free.
- Report PSI memory pressure when `/proc/pressure/memory` is available.
- Support table, JSON, YAML, and watch mode.
- Preserve numeric values in bytes for structured output.

## Linux Concepts

- `/proc/meminfo`
- `MemAvailable` versus `MemFree`
- Page cache and buffers
- Swap accounting
- Pressure Stall Information

## Expected CLI

```sh
syskit memory
syskit memory --format json
syskit memory --watch --interval 2s
syskit memory --pressure
```

## Expected Output

```text
TOTAL     USED      AVAILABLE  FREE      BUFFERS  CACHE     SWAP USED  PRESSURE
31.2 GiB  12.4 GiB  16.8 GiB   2.1 GiB   512 MiB  13.9 GiB  0 B        low
```

JSON output should include byte-valued fields and optional PSI windows such as `some.avg10`, `some.avg60`, `full.avg10`, and `full.avg60`.

## Edge Cases

- Older kernels do not expose `MemAvailable`.
- PSI is disabled or unavailable.
- Swap is not configured.
- Values are present in kB and must be converted.
- Containers may expose host memory or cgroup-limited memory depending on context.

## Acceptance Criteria

- Memory byte conversion is exact and tested.
- Missing swap reports zero total and zero used, not an error.
- Missing PSI reports unavailable pressure data.
- Structured output never uses human-formatted strings for numeric values.
- Documentation explains how cache affects "used" memory.

## Learning Objectives

- Learn Linux memory accounting fields.
- Understand reclaimable memory.
- Understand PSI as a pressure signal rather than usage.

## Estimated Complexity

Medium.

## Dependencies

- Platform file-reading abstraction.
- Memory learning material.
- Rendering architecture.

## Future Extensions

- Per-process memory correlation.
- Cgroup memory limits.
- Huge page reporting.
- NUMA memory statistics.
