# Diagnostics Feature Specification

## Purpose

Combine system metrics into read-only health checks that identify likely resource bottlenecks and configuration risks.

## User Story

As an engineer responding to a problem, I want SysKit to highlight suspicious system conditions so I can decide where to investigate first.

## Motivation

Raw metrics are useful, but many users need a first-pass interpretation. Diagnostics should provide careful, explainable hints without pretending to be an automated repair tool.

## Requirements

- Evaluate CPU, memory, disk, filesystem, process, network, and port signals.
- Produce severity levels such as info, warning, and critical.
- Explain why each finding was produced and which data sources were used.
- Never modify system state.
- Support table, JSON, and YAML output.

## Linux Concepts

- Load average and CPU utilization.
- Memory pressure and swap.
- Filesystem usage and inode exhaustion.
- Disk I/O saturation signals.
- Process resource concentration.
- Listening port exposure.

## Expected CLI

```sh
syskit diagnostics
syskit diagnostics --category memory
syskit diagnostics --severity warning
syskit diagnostics --format json
```

## Expected Output

```text
SEVERITY  CATEGORY    FINDING                         EVIDENCE
warning   filesystem  Root filesystem above 85% used   / is 88% used
info      memory      No memory pressure detected      PSI unavailable, swap unused
```

Structured output should include `id`, `severity`, `category`, `summary`, `evidence`, `sources`, and `recommendation`.

## Edge Cases

- Missing optional data prevents a check from running.
- Containers expose limited or misleading host-level metrics.
- Thresholds vary by workload.
- A single symptom may appear in multiple categories.
- Diagnostics must avoid false precision.

## Acceptance Criteria

- Each finding includes evidence.
- Missing data produces an unavailable check, not a fabricated result.
- Severity thresholds are documented.
- JSON output is suitable for automation.
- The command performs no writes and executes no repair action.

## Initial Thresholds

- CPU load: warning above one-minute load per logical CPU; critical above twice
  logical CPU count.
- Memory PSI: warning when full-stall `avg10` is at least 10%.
- Swap: warning at 80% used.
- Filesystem capacity: warning at 85% used and critical at 95%.
- Process concentration: warning when one process holds at least 50% of visible
  physical memory.
- Network: warning when cumulative interface errors or drops are non-zero.
- Ports: informational when listening sockets bind wildcard addresses.
- Disk saturation: informational unavailable finding until device busy-time
  utilization is collected; throughput alone is not treated as saturation.

## Learning Objectives

- Learn the difference between metric collection and diagnosis.
- Understand how thresholds can mislead.
- Study resource bottleneck patterns across subsystems.

## Estimated Complexity

High.

## Dependencies

- System, CPU, memory, disk, filesystem, process, network, and ports services.
- Error handling strategy.
- Rendering architecture.

## Future Extensions

- User-configurable thresholds.
- Diagnostic profiles.
- Exportable reports.
- Historical comparison when local metric storage exists.
