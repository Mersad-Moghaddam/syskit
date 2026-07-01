# Ports Feature Specification

## Purpose

Show listening ports, active sockets, socket states, protocols, local and remote addresses, and associated processes where discoverable.

## User Story

As a developer or operator, I want to know which processes are listening on which ports so I can debug service conflicts, exposure, and connectivity issues.

## Motivation

Port inspection is commonly done with `ss`, `netstat`, or `lsof`. SysKit should provide the common read-only view directly from Linux socket tables and process file descriptors.

## Requirements

- List TCP, UDP, TCP6, UDP6, and Unix sockets where supported.
- Show protocol, local address, local port, remote address, remote port, state, socket inode, PID, and process command.
- Support filters for listening sockets, protocol, port, address, PID, and state.
- Map socket inodes to processes by scanning `/proc/[pid]/fd`.
- Report partial process mapping when permissions restrict fd access.

## Linux Concepts

- `/proc/net/tcp`
- `/proc/net/tcp6`
- `/proc/net/udp`
- `/proc/net/udp6`
- `/proc/net/unix`
- Socket states
- Socket inode to process fd mapping

## Expected CLI

```sh
syskit ports
syskit ports --listening
syskit ports --protocol tcp --port 8080
syskit ports --pid 1234
syskit ports --format json
```

## Expected Output

```text
PROTO  LOCAL            REMOTE           STATE   PID    COMMAND
tcp    0.0.0.0:22       0.0.0.0:*        LISTEN  898    sshd
tcp    127.0.0.1:5432   0.0.0.0:*        LISTEN  1420   postgres
```

JSON output should include raw socket state, decoded state, addresses, ports, inode, PID when known, and process command when known.

## Edge Cases

- Permission denied while scanning another user's file descriptors.
- Process exits after socket table is read.
- IPv6 addresses require correct decoding.
- UDP sockets have no TCP-style connection state.
- Multiple processes may share a socket after fork.

## Acceptance Criteria

- Listening TCP sockets are correctly identified.
- UDP sockets are represented without misleading TCP states.
- Socket inode mapping handles partial permission failures.
- IPv4 and IPv6 addresses are decoded accurately.
- No external `ss`, `netstat`, or `lsof` command is executed.

## Learning Objectives

- Understand procfs socket table formats.
- Learn TCP state names and UDP differences.
- Learn how file descriptor symlinks expose socket inodes.

## Estimated Complexity

High.

## Dependencies

- Process feature.
- Network feature.
- Collector architecture.

## Future Extensions

- Netlink socket collection.
- Namespace-aware socket inspection.
- Connection grouping by process.
- Port exposure diagnostics.
