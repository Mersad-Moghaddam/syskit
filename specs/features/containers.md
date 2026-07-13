# Containers Feature Specification

## Purpose

Inspect containers, container-associated processes, and cgroup resource usage in Docker, Podman, and OCI-style environments.

## User Story

As an operator on a Linux host, I want to connect process and resource usage to containers so I can diagnose containerized workloads without leaving SysKit.

## Motivation

Containers are Linux processes organized through namespaces, cgroups, and runtime metadata. SysKit should expose container context while keeping the core read-only and Linux-native.

## Requirements

- Detect cgroup v1 and cgroup v2 layout.
- Map processes to container identifiers where possible.
- Report container CPU, memory, and I/O usage when cgroup files expose it.
- Optionally integrate with Docker or Podman metadata when runtime sockets are available.
- Degrade gracefully when runtime access is unavailable.

## Linux Concepts

- cgroup v1 and v2.
- PID namespaces.
- Mount namespaces.
- Container runtime metadata.
- OCI container identifiers.

## Expected CLI

```sh
syskit containers
syskit containers inspect <id>
syskit process --containers
syskit containers --format json
```

## Expected Output

```text
CONTAINER  RUNTIME    PIDS  MEMORY     CPU NS       READ       WRITE
abc123...  containerd 12    440401920  3240000000   1048576    2097152
```

Structured output should separate cgroup-derived metrics from runtime metadata.
Runtime names and status remain unavailable unless a future optional metadata
adapter can prove them; cgroup reporting never depends on such an adapter.

## Edge Cases

- Host uses cgroup v2 unified hierarchy.
- Runtime socket is unavailable or permission denied.
- Container ID appears in cgroup path but metadata is unavailable.
- Processes exit during mapping.
- Rootless containers use user-scoped paths.

## Acceptance Criteria

- cgroup version detection is explicit.
- Missing runtime metadata does not prevent cgroup-based reporting.
- Container process mapping is best-effort and marks uncertainty.
- The feature performs no container management actions.
- Tests include cgroup v1 and v2 fixture layouts.

## Learning Objectives

- Understand cgroups as the foundation of container resource accounting.
- Learn how container runtimes map metadata onto Linux primitives.
- Study namespace limitations for host inspection.

## Estimated Complexity

Very High.

## Dependencies

- Process feature.
- CPU, memory, and disk collectors.
- Cgroup platform adapter.
- Container learning material.

## Future Extensions

- Kubernetes pod and namespace mapping.
- Container log metadata.
- Runtime-specific detail panels.
- Historical container resource trends.
