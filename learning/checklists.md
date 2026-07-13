# Learning Checklists

> Evidence-based completion gates for the SysKit course. Check a box only when
> you can link to a lab report, command transcript, diagram, fixture, or test.

## How To Record Progress

| Evidence | Example |
|---|---|
| Explanation | Short note that states semantics, units, and limitations |
| Calculation | Two timestamped samples and shown arithmetic |
| Diagram | Source-to-model or correlation flow with join keys |
| Fixture/test | Sanitized case with provenance and expected behavior |
| Validation | Commands and results another contributor can reproduce |

```mermaid
flowchart LR
    K[Knowledge] --> E[Reproducible evidence]
    E --> R[Peer or self review]
    R -->|gap found| K
    R -->|criteria met| C[Competency complete]
```

## Foundation Gate

- [ ] I can distinguish kernel, userspace, process, thread, namespace, and cgroup.
- [ ] I can classify counters, gauges, rates, ratios, estimates, and labels.
- [ ] I can convert sectors, pages, ticks, kHz, KiB, and time units safely.
- [ ] I can explain why zero, unavailable, denied, transient, and malformed differ.
- [ ] I can select procfs, sysfs, Netlink, or cgroups for a stated observation.
- [ ] I can trace SysKit's CLI → Command → Service → Collector → Platform path.
- [ ] I can explain why verification commands are not production data sources.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## CPU Gate

- [ ] I can distinguish topology/identity from activity.
- [ ] I have manually derived utilization from two `/proc/stat` samples.
- [ ] I can explain jiffies, guest accounting, iowait, steal, and load average.
- [ ] I handle zero delta, reset, hotplug, and absent cpufreq in my reasoning.
- [ ] I can identify collector parse tests and service delta tests.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Memory Gate

- [ ] I can explain free, available, cache, reclaimable slab, and used memory.
- [ ] I can distinguish swap occupancy, swap activity, and actual pressure.
- [ ] I can interpret PSI `some`, `full`, averages, and total.
- [ ] I handle missing `MemAvailable`, missing PSI, and no-swap correctly.
- [ ] I can distinguish host `/proc/meminfo` from cgroup-relative memory.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Disk And Filesystem Gate

- [ ] I can map device → partition → filesystem → mount without collapsing them.
- [ ] I know why sysfs/diskstats sectors are converted using 512 bytes.
- [ ] I can derive IOPS and throughput from matched diskstats samples.
- [ ] I can parse mountinfo structure and escaped paths.
- [ ] I can diagnose byte exhaustion separately from inode exhaustion.
- [ ] I can explain virtual, overlay, remote, and namespace-relative mounts.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Process Gate

- [ ] I can distinguish PID, TID, TGID, PPID, process, and thread.
- [ ] I can safely parse a `stat` line with hostile parentheses/spaces in `comm`.
- [ ] I understand process states, especially `D` and `Z`.
- [ ] I handle NUL-separated command lines and empty kernel-thread command lines.
- [ ] I treat PID disappearance and permission failures as partial observation.
- [ ] I can explain why PID plus start time is safer for matching samples.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Network And Ports Gate

- [ ] I can separate links, addresses, routes, sockets, ports, and DNS.
- [ ] I can explain routing Netlink dumps and socket diagnostics conceptually.
- [ ] I can derive interface rates with reset/recreation protection.
- [ ] I can map a socket to a process and name every race/permission boundary.
- [ ] I understand wildcard binds, IPv4/IPv6, TCP states, and network namespaces.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Live Monitoring Gate

- [ ] I use real elapsed time and can explain refresh jitter.
- [ ] I can draw loading, ready, failed, refresh, resize, and quit transitions.
- [ ] I understand backpressure when collection exceeds the interval.
- [ ] I keep collection/service logic separate from TUI rendering.
- [ ] I can show test evidence for resize, cancellation, error, and no-color paths.
- [ ] I preserve non-TTY behavior and stdout/stderr separation.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Containers And Plugins Gate

- [ ] I can separate namespaces, cgroups, and runtime metadata.
- [ ] I can identify cgroup v1/v2 from mounts and membership.
- [ ] I can interpret usage, numeric limits, `max`, and unavailable controllers.
- [ ] I understand hierarchy, delegation, and process movement during collection.
- [ ] I can explain out-of-process plugin discovery and trust boundaries.
- [ ] I can explain timeout, invalid output, and unsupported protocol versions.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Engineering Gate

- [ ] I can map each acceptance criterion to concrete evidence.
- [ ] I can define a field contract with source, type, unit, optionality, and scope.
- [ ] I can place parsing, derivation, validation, rendering, and exit policy correctly.
- [ ] I can design unit, fixture, integration, golden, race, fuzz, and benchmark coverage.
- [ ] I can distinguish portable integration invariants from host-specific constants.
- [ ] I review structured output changes as compatibility changes.
- [ ] I can complete the [feature review worksheet](engineering.md#10-feature-review-worksheet).
- [ ] I have run applicable Definition of Done checks and recorded results.

Evidence link/notes:

> Evidence: _add a reproducible artifact here._

## Capstone Gate

- [ ] I completed at least one incident investigation from [labs.md](labs.md).
- [ ] I completed at least one parser/test or engineering lab.
- [ ] My report states a falsifiable hypothesis and considers alternatives.
- [ ] Every observation records source, unit, metric type, timestamp, and scope.
- [ ] I used at least three signals and separated primary from verification sources.
- [ ] I proposed or added a fixture/test for the learned edge case.
- [ ] Every assessment dimension is at least **Competent**.

Final evidence link/notes:

> Evidence: _add a reproducible artifact here._
