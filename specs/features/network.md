# Network Feature Specification

## Purpose

Inspect network interfaces, addresses, counters, routes, DNS configuration, and basic bandwidth rates.

## User Story

As an operator diagnosing connectivity, I want one consistent view of network state so I can quickly see interfaces, traffic, routes, and resolver configuration.

## Motivation

Linux networking data is split across procfs, sysfs, Netlink, and configuration files. SysKit should prioritize Netlink for authoritative network state and use procfs/sysfs where they provide useful counters.

## Requirements

- List network interfaces with state, MTU, MAC address, addresses, RX/TX bytes, packets, errors, and drops.
- Report routes and default gateway.
- Report DNS resolver configuration.
- Support bandwidth rates when two samples are available.
- Support table, JSON, YAML, and watch mode.

## Linux Concepts

- Netlink route messages.
- `/sys/class/net`
- `/proc/net/dev`
- `/proc/net/route`
- `/etc/resolv.conf`
- Interface flags and operational state.

## Expected CLI

```sh
syskit network
syskit network interfaces
syskit network routes
syskit network dns
syskit network --watch
syskit network --format json
```

## Expected Output

```text
IFACE   STATE  MTU   RX        TX        ERRORS  ADDRESSES
eth0    up     1500  8.2 GiB   1.4 GiB   0       192.168.1.20/24
lo      up     65536 120 MiB   120 MiB   0       127.0.0.1/8
```

Structured output should separate interfaces, addresses, counters, routes, and DNS data.

## Edge Cases

- Interfaces appear or disappear during collection.
- Network namespaces hide host interfaces.
- `/etc/resolv.conf` is a symlink managed by a resolver service.
- Counter values reset after interface restart.
- Wireless-specific data is unavailable or outside the initial scope.

## Acceptance Criteria

- Interface counters are numeric and unit-stable in JSON.
- Bandwidth rates require two samples.
- Routes are represented separately from interfaces.
- DNS parsing handles comments and multiple nameservers.
- No external `ip`, `ifconfig`, `ss`, or `netstat` command is executed.

## Learning Objectives

- Understand interfaces, addresses, routes, and DNS as separate concepts.
- Learn Netlink's role in Linux networking.
- Study cumulative counters and rate calculation.

## Estimated Complexity

High.

## Dependencies

- Collector architecture.
- Netlink adapter design.
- Network learning material.

## Future Extensions

- Wireless interface details.
- Network namespace selection.
- Per-process network correlation.
- Packet loss and latency diagnostics.
