# Containers And Plugins — Learning Notes

> Study notes on cgroups, container identity, and SysKit's out-of-process plugin model.
> Written for the implementer of SysKit's `containers` collector and the Stage 6 (v0.4–v0.5) plugin system.

---

## Concepts

A "container" is not a kernel object. There is no `struct container` in Linux.
What actually exists is a set of processes that a runtime (Docker, Podman,
containerd, CRI-O) has wrapped in **namespaces** (isolation — what the process
can *see*) and **cgroups** (accounting and limits — what the process can *use*),
plus some runtime metadata on the side. SysKit's job is to expose that context
read-only; it performs no container management actions. Keep three concerns
separate — newcomers blur them:

- **Namespaces** answer *"what does this process see?"* — its own PID space, its
  own mount tree, its own network stack. Relevant to why a container "feels"
  like a separate machine, but namespaces do **not** track resource usage.
- **cgroups** (control groups) answer *"what is this group of processes allowed
  to use, and how much has it used?"* — CPU time, memory bytes, I/O. **This is
  where all container resource accounting lives.** When you want CPU% or memory
  for a container, you are reading cgroup files, not summing the process tree.
- **Runtime metadata** answers *"what does a human call this?"* — the friendly
  name `web-api`, the image, the status `running`. This lives in a runtime
  socket (`/var/run/docker.sock`) or state files, is optional, and may be
  permission-denied. It must never be a prerequisite for cgroup-based reporting.

**cgroups are the real unit of accounting, not the process tree.** A junior's
instinct is to find a container's "main" PID, walk its children via
`/proc/[pid]/stat` PPIDs, and sum their CPU/RSS. That is wrong for two reasons:
processes fork, exit, and re-parent constantly (you race the tree and
double-count or miss), and the kernel already maintains the authoritative total
for the whole group in the cgroup's own files (`memory.current`, `cpu.stat`).
The cgroup is a stable membership boundary the kernel accounts against directly.
Read the cgroup; don't reconstruct it from the process tree.

**Container identity is derived, and the derivation is a heuristic.** The clean
signal that ties a process to a container is its cgroup path, because runtimes
encode the container ID into that path. But the encoding differs per runtime and
per version, so treat the mapping as best-effort and mark uncertainty in the
output rather than asserting a container ID you inferred by pattern-matching.

**Plugins are a separate axis.** Stage 6 also adds custom collectors via
plugins. The load-bearing decision (ADR 007) is that plugins run **out of
process** — as separate executables SysKit talks to over a versioned protocol —
not as in-process Go `plugin` objects. That choice is a trust and isolation
boundary, covered below.

---

## Linux Internals

### cgroup membership — reading it from `/proc`

Every process belongs to a cgroup, and the kernel tells you which one at
`/proc/[pid]/cgroup` (and `/proc/self/cgroup` for the current process). That
file *is* the membership record — you do not guess it, you read it. Its format
is the single most important thing to get right because it also tells you which
cgroup version is in effect:

- **cgroup v2 (unified):** exactly **one** line, always beginning `0::`:

  ```
  0::/system.slice/docker-3f9a...c1.scope
  ```

  The `0` is the hierarchy ID (always 0 for v2), the middle field is empty, and
  the path after `::` is the cgroup relative to the v2 root mount.

- **cgroup v1:** **multiple** lines, one per controller hierarchy, of the form
  `N:controller:/path`:

  ```
  11:memory:/docker/3f9a...c1
  10:cpu,cpuacct:/docker/3f9a...c1
  4:pids:/docker/3f9a...c1
  ```

  Each controller (memory, cpu, pids, …) is a *separate* hierarchy with its own
  path, which is exactly why v1 is more painful — the same process has one path
  per controller and you must read the right one for each metric.

**Detecting which version is active.** Do not hardcode an assumption. Two robust
signals: (1) if `/sys/fs/cgroup/cgroup.controllers` exists, the unified v2
hierarchy is mounted; (2) `/proc/self/cgroup` containing a single `0::` line
means v2, whereas multiple `N:controller:/path` lines means v1. A **hybrid**
setup also exists (v2 mounted but some controllers still on v1 hierarchies), so
the honest check is "which hierarchy actually carries the controller I need,"
not "is the machine v1 or v2."

> **Newcomer bug: assuming cgroup v2 everywhere.** Modern distros default to v2,
> so it is tempting to only parse the `0::` line and read `memory.current`. But
> plenty of hosts still run pure v1 or hybrid (older LTS distros, some managed
> Kubernetes nodes, custom kernels). Code that only handles v2 will silently
> mis-read or crash on those. The acceptance criterion is explicit: **detect the
> version, handle both, or clearly report the layout as unsupported** — never
> assume.

### cgroup resource accounting — where the numbers live

Once you have the cgroup path, resource usage is read from files inside that
cgroup directory under the mount root. The filenames differ sharply between
versions, which is the other reason version detection must come first:

- **v2 (unified, all under `/sys/fs/cgroup/<path>/`):**
  - `cgroup.controllers` — which controllers are enabled here.
  - `memory.current` — current memory usage in bytes (single number).
  - `cpu.stat` — cumulative CPU accounting; `usage_usec` is total CPU time in
    microseconds. Like every counter in SysKit, it is a monotonic total — a
    single read is a total, a *rate* needs two samples and the service layer's
    delta math, never the collector's.
  - `io.stat` — per-device I/O counters (`rbytes`, `wbytes`, …).
- **v1 (separate hierarchies, one mount per controller):**
  - `memory/memory.usage_in_bytes` — current memory usage (the v1 spelling of
    `memory.current`).
  - `cpu/cpuacct.usage` (a.k.a. under `cpuacct/`) — cumulative CPU time in
    nanoseconds.

The takeaway: **same concept, different files.** Build a small platform adapter
that, given a version and a controller, returns the right path and unit
(microseconds vs nanoseconds, note the difference), so the rest of the collector
speaks one vocabulary.

### Container-to-process mapping — heuristic by design

Runtimes encode the container ID into the cgroup path, so the mapping is: read
`/proc/[pid]/cgroup` → extract the path → recognize the runtime's pattern →
pull the container ID. Real-world path shapes:

- **Docker (systemd cgroup driver):** `.../system.slice/docker-<id>.scope`
- **Docker (cgroupfs driver):** `.../docker/<id>`
- **Kubernetes pods:** `.../kubepods/.../pod<uid>/<container-id>` (and
  `kubepods.slice/...` variants under systemd).
- **Podman / rootless:** user-scoped paths such as
  `.../user.slice/user-1000.slice/user@1000.service/.../<id>`

**Say it plainly in the code and the docs: deriving container identity from the
cgroup path is a heuristic that varies across runtimes and versions.** There is
no kernel API that says "PID 4021 is in container web-api." So: match known
patterns, extract the ID, and when runtime metadata (the friendly name from the
Docker/Podman socket) is unavailable, still report the container by its cgroup
ID and mark the identity as uncertain rather than dropping it. Expect processes
to exit mid-scan — handle `ENOENT`/`ESRCH` on `/proc/[pid]/*` as "process gone,"
not fatal, exactly as the process collector does.

### Plugin trust boundaries — why out-of-process (ADR 007)

SysKit runs plugins as **separate executables** communicating over a versioned
protocol (JSON over stdin/stdout or a local socket), **not** as Go in-process
`plugin` packages. This is a deliberate trust boundary. Why it matters:

- **Fault isolation / crash containment.** A plugin is user-installed code of
  unknown quality. In-process, a segfault, panic, or memory corruption in the
  plugin takes down SysKit itself. Out-of-process, the plugin lives in its own
  address space — **a misbehaving plugin cannot corrupt SysKit's memory** or
  crash the core; SysKit just sees a dead child process and reports it.
- **Language-agnostic.** Anything that can read stdin and write structured
  output on stdout can be a plugin — Python, Rust, a shell script — because the
  contract is a wire protocol, not the Go ABI.
- **No Go ABI / toolchain coupling.** Go's `plugin` package demands the plugin
  be built with a matching Go version, matching dependency versions, and a
  matching platform, or it fails to load. The out-of-process model sidesteps all
  of that.
- **Explicit trust.** Plugin execution is opt-in and visible: SysKit shows the
  plugin's path and permissions in diagnostic output, never auto-installs, and
  **never loads plugins from world-writable directories** (a world-writable
  plugin dir means any local user could drop in code that SysKit would run — an
  obvious privilege-escalation vector). Discovery is explicit: `--plugin-dir`,
  `$SYSKIT_PLUGIN_DIR`, `$XDG_DATA_HOME/syskit/plugins`, then
  `~/.local/share/syskit/plugins`. Core commands must work with zero plugins.

The cost (accepted in ADR 007) is protocol design work, per-call overhead, and
the need for timeout/lifecycle management and careful validation of plugin
output — a worthwhile trade for isolation and a clear trust model.

### Output schema versioning — why an explicit version is mandatory

Plugins and the SysKit core that consumes their output evolve on **independent
schedules**: a plugin author ships v2 of their collector next month; a user's
SysKit is still last quarter's release; a downstream script parses the JSON. If
the output schema carries no explicit version, a field rename or type change
**silently breaks downstream parsers** — they read stale field names, get
`null`, and produce wrong results with no error. An explicit schema/protocol
version turns that silent breakage into a loud, handleable event: SysKit can
refuse to load a plugin whose API version it doesn't support (unless the user
opts into experimental behavior), consumers can branch on the version, and
migration windows can support multiple versions at once. The plugin **protocol**
version is also versioned independently from SysKit's CLI version for the same
reason. Rule of thumb: **any data crossing a trust or process boundary carries
its own version** — unversioned schemas are a future silent-failure bug.

Note the division of labor: plugin output enters SysKit as **structured data,
not pre-rendered terminal text**, so the core keeps owning table/JSON/YAML/TUI
rendering and the user experience stays consistent.

---

## Important Files

- `/proc/[pid]/cgroup` — the cgroup membership of a process. **v2:** one `0::/path`
  line. **v1:** multiple `N:controller:/path` lines. Primary source for both
  version detection and container-ID derivation.
- `/proc/self/cgroup` — same, for the current process; the easy thing to `cat`
  while learning the format.
- `/sys/fs/cgroup/` — the cgroup filesystem root. Under **v2** this is the single
  unified hierarchy; under **v1** it holds one subdirectory per controller mount.
- `/sys/fs/cgroup/cgroup.controllers` — **existence signals v2**; lists the
  controllers available in the unified hierarchy.
- `/sys/fs/cgroup/<path>/memory.current` — **v2** current memory usage in bytes.
- `/sys/fs/cgroup/<path>/cpu.stat` — **v2** cumulative CPU accounting
  (`usage_usec` = total CPU microseconds). Two-sample rate source.
- `/sys/fs/cgroup/<path>/io.stat` — **v2** per-device I/O counters.
- `/sys/fs/cgroup/memory/memory.usage_in_bytes` — **v1** current memory usage in
  bytes.
- `/sys/fs/cgroup/cpu/cpuacct.usage` — **v1** cumulative CPU time in nanoseconds
  (also exposed under a `cpuacct/` hierarchy).

---

## Useful Commands

These are for *learning and verifying* your collector against ground truth —
SysKit itself reads the files natively and does not shell out (ADR-003). Run
them by hand and diff against your parser.

- `cat /proc/self/cgroup` — see your shell's cgroup line(s); confirm at a glance
  whether the host is v2 (`0::/...`) or v1 (multiple `N:controller:/...`).
- `cat /sys/fs/cgroup/cgroup.controllers` — if this file exists and lists
  controllers, the unified v2 hierarchy is active. A "no such file" is itself a
  signal (pure v1).
- `systemd-cgls` — the cgroup tree rendered as a hierarchy with the processes in
  each group; ground truth for "which PIDs are in which cgroup."
- `cat /sys/fs/cgroup/<path>/memory.current` — read a container's live memory
  (v2) and cross-check your collector's number.
- `docker ps` — **external reference only** (needs the Docker runtime/socket).
  Lists container IDs and friendly names so you can confirm the ID you extracted
  from a cgroup path maps to the container you think it does.

---

## References

- Linux kernel `Documentation/admin-guide/cgroup-v2.rst` — the unified hierarchy, controllers, and interface files: https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html
- `cgroups(7)` — cgroup v1 and v2 concepts, hierarchies, and `/proc/[pid]/cgroup` format: https://man7.org/linux/man-pages/man7/cgroups.7.html
- `namespaces(7)` — the isolation half of what makes a container: https://man7.org/linux/man-pages/man7/namespaces.7.html
- `proc(5)` — `/proc/[pid]/cgroup` and related files: https://man7.org/linux/man-pages/man5/proc.5.html
- OCI Runtime Specification — how runtimes map metadata onto Linux primitives: https://github.com/opencontainers/runtime-spec/blob/main/spec.md
- ADR-007 — Prefer out-of-process plugins: `../decisions/007-out-of-process-plugins.md`
- Plugin architecture spec — trust model, discovery, versioning: `../specs/plugin-architecture.md`
- ADR-003 — Read native kernel interfaces instead of shelling out: `../decisions/003-native-apis-over-shell.md`

---

## Personal Notes
