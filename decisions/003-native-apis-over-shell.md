# 003. Read native kernel interfaces instead of shelling out

**Status:** Accepted, 2026-07-01

---

## Context

A system-inspection tool must get its data from somewhere. There are two broad
strategies:

1. **Shell out and parse.** Invoke existing utilities (`top`, `free`, `df`,
   `ss`, `ip`, `lsof`, `vmstat`) as subprocesses and parse their textual output.
2. **Read native interfaces.** Read the kernel's own data sources directly —
   the `/proc` and `/sys` pseudo-filesystems, Netlink sockets, and cgroup
   hierarchies — and parse their structured, kernel-defined formats.

The shell-out approach is tempting because it appears to reuse decades of mature
tooling. In practice it is brittle and slow:

- **Version- and distribution-dependent output.** The columns, units, and
  formatting of `ss`, `ip`, or `df` differ across tool versions and distros.
  Parsing human-readable output means chasing those differences indefinitely.
- **Fragile parsing.** Output intended for humans is not a stable contract.
  Locale settings, terminal width, and flag defaults can all change what we must
  parse.
- **Process overhead.** Each invocation forks a subprocess, execs a binary, and
  serialises data to text only for us to deserialise it again — directly at odds
  with constitution principle 3 (*Performance Matters*), especially in the
  real-time refresh loops of the [v0.3 dashboard](../specs/roadmap.md).
- **External dependencies.** The tool would silently depend on whichever
  utilities happen to be installed, undermining the single-static-binary goal
  from [ADR 001](./001-use-go.md).

The kernel interfaces themselves are, by contrast, a comparatively stable and
structured contract. `/proc/stat`, `/proc/meminfo`, `/proc/[pid]/stat`, and the
`RTM_*` Netlink messages have well-documented, slowly-evolving formats.

The constitution makes this a binding rule in principle 2 (*Native APIs First*):
"SysKit reads system data from native Linux interfaces ... rather than parsing the
output of shell commands." This ADR records the reasoning behind that principle.

---

## Decision

We will **read native Linux kernel interfaces directly** rather than executing
and parsing external commands. Specifically:

- Process, CPU, memory, and general statistics come from **`/proc`**.
- Device, topology, and hardware information comes from **`/sys`**.
- Network interfaces, addresses, routes, and socket state come from **Netlink**
  (`AF_NETLINK`), not from parsing `ip`/`ss`. The address adapter currently
  uses the standard library; see [ADR 011](011-stdlib-netlink-addresses.md).
- Container and resource-control data comes from **cgroups** (v1 under
  `/sys/fs/cgroup`, and the v2 unified hierarchy).

SysKit will not fork subprocesses to gather inspection data. Where a native
interface is genuinely unavailable or impractical for a specific data point, the
exception must be documented and justified inline and, if significant, recorded
as its own ADR — the constitution requires the exception to be explicit.

---

## Consequences

### Positive

- Data is structured and defined by the kernel, giving a far more stable contract
  than human-readable tool output.
- No dependency on external utilities being installed or on their versions —
  reinforcing the single-static-binary distribution model.
- Lower latency and allocation overhead: reading a pseudo-file with a buffered
  reader is dramatically cheaper than forking a process and parsing its stdout.
- Access to richer detail than most tools surface (per-CPU `/proc/stat` fields,
  full Netlink routing attributes, cgroup v2 controllers).
- Reading kernel interfaces directly is itself an effective way to learn Linux
  internals (*Learn Before Build*).

### Negative

- We must implement and maintain parsers for each interface format, including
  quirks and kernel-version variations — work that shelling out would have
  delegated.
- Some interfaces are genuinely fiddly: Netlink message framing, cgroup v1-vs-v2
  divergence, and `/proc/[pid]` races when processes exit mid-read all require
  care.
- We take on responsibility for handling missing files, permission errors, and
  malformed data gracefully, rather than inheriting a tool's error handling.

### Neutral

- Netlink access is kept behind the Platform Abstraction Layer. The current
  address-dump adapter uses the standard library; a future dependency requires
  review under the *Minimal Dependencies* policy.
- Kernel-version variation is contained within the Platform Abstraction Layer
  (see [ADR 004](./004-layered-architecture.md)), keeping the quirks in one place.

---

## Alternatives Considered

- **Exec external tools and parse their output.** The traditional approach:
  wrap `free`, `df`, `ss`, `ip`, `lsof`, etc. Rejected for brittleness
  (unstable, version- and locale-dependent text), performance (per-call fork/exec
  overhead), and hidden runtime dependencies. This is precisely the pain the
  product set out to eliminate — the [product overview](../specs/product.md) lists
  "parsing inconsistent output" and "combining tools with shell pipelines" as the
  problems SysKit exists to solve.
- **A hybrid: native where easy, shell out where hard.** Read `/proc` directly
  but fall back to `ip`/`ss` for networking. Rejected because it reintroduces the
  external-dependency and parsing-fragility problems for exactly the subsystems
  where they bite hardest, and it muddies an otherwise clean architectural rule.
  Networking will use Netlink.
- **Link against a C library (e.g. libprocps).** Would offload some parsing to a
  maintained library. Rejected because it requires cgo, breaks the pure-Go static
  binary, and reintroduces an external dependency and its versioning concerns.

---

## References

- [Constitution](../specs/constitution.md) — principle 2 (Native APIs First),
  principle 3 (Performance Matters), principle 8 (Minimal Dependencies)
- [Product Overview](../specs/product.md) — problems SysKit addresses
- [ADR 001](./001-use-go.md) — single static binary goal
- [ADR 004](./004-layered-architecture.md) — Platform Abstraction Layer isolates format quirks
- [proc(5)](https://man7.org/linux/man-pages/man5/proc.5.html),
  [netlink(7)](https://man7.org/linux/man-pages/man7/netlink.7.html),
  [cgroups(7)](https://man7.org/linux/man-pages/man7/cgroups.7.html)
