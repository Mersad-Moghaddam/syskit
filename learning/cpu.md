# CPU: Topology, Time, And Utilization

> Study notes on CPU architecture, Linux CPU interfaces, and related internals.
> Written to prepare you to implement SysKit's `cpu` collector. Read the CPU
> feature spec (`specs/features/cpu.md`) and the collector model
> (`specs/collectors.md`) alongside this.

---

| Attribute | Value |
|---|---|
| Level | Domain |
| Prerequisites | [Kernel interfaces](kernel-interfaces.md) |
| Time | 2–3 hours plus sampling lab |
| Product contract | [CPU feature](../specs/features/cpu.md) |

## Learning Objectives

- Distinguish identity/topology gauges from cumulative activity counters.
- Derive aggregate and per-core utilization from timestamped snapshots.
- Interpret jiffies, load average, frequency scaling, steal, and iowait.
- Handle CPU hotplug, missing cpufreq, resets, and guest accounting honestly.
- Design fixture and service tests for raw and derived CPU values.

```mermaid
flowchart LR
    A[/proc/cpuinfo and sysfs topology] --> I[Static identity]
    B[/proc/stat at t1] --> D[Counter deltas]
    C[/proc/stat at t2] --> D
    D --> U[Utilization and state percentages]
    I --> V[CPU view]
    U --> V
```

## Concepts

The `cpu` collector has to report two very different kinds of data, and the
single biggest source of bugs is treating them the same way.

**1. Static identity vs. sampled activity.** Some CPU facts are *static*: the
model name, how many logical and physical cores exist, the socket count, the
instruction-set flags. You read them once and you're done. Other facts are
*activity over time*: how busy the CPU was, at what frequency it ran. These only
have meaning as a *rate*, which means you cannot get them from a single read.
SysKit's spec makes this split explicit — "keep raw counters separate from
derived percentages" — and the collector/service split enforces it: the
collector returns raw snapshots, the service layer computes rates from two of
them (`specs/collectors.md`, "Snapshot Model").

**2. Logical vs. physical cores (and why you'll confuse them).** A modern CPU
with SMT / Hyper-Threading exposes more *logical* CPUs than it has *physical*
cores. On the box these notes were verified against, one socket has 4 physical
cores but the OS also could present 8 logical CPUs if SMT were on; here
`siblings` = 4 and `cpu cores` = 4, so SMT is effectively 1:1. The rule:
- **Logical cores** = number of `processor:` blocks in `/proc/cpuinfo` = `nproc`
  = number of `cpuN` lines in `/proc/stat`. This is what the scheduler runs on.
- **Physical cores** = `cpu cores` field. **Sockets** = count of distinct
  `physical id` values. `siblings` = logical CPUs per socket.
SysKit must report logical cores always (easy, always present) and physical
cores "where detectable" — because in many VMs `physical id` / `cpu cores` are
missing or faked. Do not compute physical cores by dividing logical by 2; SMT
is not always 2x and not always on.

**3. CPU time is accounted in "jiffies".** The kernel does not store "CPU was
47% busy." It stores cumulative *tick counts* since boot, split into states:
user, nice, system, idle, iowait, irq, softirq, steal, guest, guest_nice. Each
tick is one *jiffy*. See Linux Internals below — this is the concept that trips
up every newcomer.

**4. CPU states worth knowing:**
- **user / nice** — time in userspace (nice = user time for re-niced/low-prio
  processes).
- **system** — time in the kernel on behalf of processes.
- **idle** — doing nothing.
- **iowait** — idle *while* a task was blocked on I/O. It is a soft hint, not a
  reliable metric; don't build alarms on it.
- **irq / softirq** — servicing hardware / software interrupts.
- **steal** — time the hypervisor gave to *other* VMs while this vCPU wanted to
  run. High steal on a VM means a noisy neighbour, not a busy guest. This is why
  steal matters for SysKit's VM story.
- **guest / guest_nice** — time spent running guest VMs (for hosts running KVM).
  Note: guest time is *already counted inside* user/nice — see the gotcha below.

**5. Frequency scaling (cpufreq).** Modern CPUs change clock speed at runtime
under a *governor* (e.g. `powersave`, `performance`, `schedutil`). The kernel
exposes current/min/max frequency through sysfs — **but only if a cpufreq driver
is loaded.** In many VMs and containers there is no cpufreq directory at all.
SysKit must represent missing frequency as **UNAVAILABLE, never 0** (spec
Acceptance Criteria). Reporting 0 MHz is a lie that looks like a broken CPU.

**6. Load average.** `/proc/loadavg` gives 1/5/15-minute run-queue averages.
It is *not* a percentage and *not* directly comparable to core count without
context: a load of 4.0 on a 4-core box is ~saturated; on a 64-core box it's
nearly idle. It counts runnable *and* uninterruptible (D-state, usually disk)
tasks, so it is a coarse pressure signal, not pure CPU utilization.

---

## Linux Internals

### Jiffies, USER_HZ, and why raw counters are not percentages

The columns in `/proc/stat` are **cumulative counts of clock ticks since boot**,
not percentages and not seconds. The tick rate is `USER_HZ`, discoverable at
runtime:

```
$ getconf CLK_TCK
100
```

So `100` means 100 jiffies = 1 second of that state *per logical CPU*. A raw
value like `user = 1423178` means ~14231 CPU-seconds of userspace time
accumulated since boot. **This is the classic newcomer mistake:** taking a raw
counter, or the ratio `user/total` from a single read, and printing it as
"CPU %". At best you get a meaningless lifetime average; at worst garbage. The
number only becomes utilization when you diff two reads.

### Why utilization needs TWO samples (the core lesson)

Because the counters are monotonic totals, "how busy is the CPU *right now*" is a
**derivative**: how much did each counter grow between time T1 and T2.

The algorithm the SysKit service layer implements:
1. Read `/proc/stat` at T1 → save every column of the `cpu` line (and each
   `cpuN` line).
2. Wait the sample interval (SysKit default watch interval, e.g. 1s).
3. Read `/proc/stat` again at T2.
4. For each field compute `delta = value_T2 - value_T1`.
5. `total_delta = sum of all deltas`. `busy_delta = total_delta - (idle_delta +
   iowait_delta)`.
6. `utilization = busy_delta / total_delta * 100`.

**What breaks if you skip the second sample:** you have no `total_delta`, so you
either divide by the lifetime total (dilutes any spike into near-zero) or invent
a percentage from a single counter (nonsense). This is exactly why the spec says
"Utilization calculations require two timestamped samples" and why the collector
returns snapshots while the *service* owns the subtraction. Keep the timestamp
with each snapshot — you need the real elapsed time, not the intended interval,
if the process was delayed.

### `/proc/stat` layout

The first `cpu` line is the **aggregate** across all logical CPUs. Each
subsequent `cpuN` line is **one logical CPU**. Same 10 columns, same order:

```
cpu  1423178 1892 323404 7421365 34906 0 9836 0 0 0
cpu0 354308  580  80847  1857960  7246 0 3398 0 0 0
cpu1 364031  396  81089  1852685  9629 0  673 0 0 0
```

Columns, left to right: `user nice system idle iowait irq softirq steal guest
guest_nice`. Older kernels have fewer columns (no steal/guest); parse
defensively and treat absent trailing columns as 0. Note the aggregate `cpu`
line is the *sum* of the per-core lines — don't double-count by adding both.

**guest gotcha:** `guest` is already included in `user`, and `guest_nice` in
`nice`. If you sum all ten columns naively you slightly overcount total. The
common convention (and what `top` does) is to treat the ten columns as the total
anyway, or to subtract guest back out; be consistent and document your choice.
For per-core *utilization* it rarely matters, but know it exists.

### `/proc/cpuinfo` structure

One block per **logical** CPU, blocks separated by a blank line, each line
`key<tab>: value`. Real fields you will parse (verified):

```
processor	: 0
model name	: Intel(R) Core(TM) i5-6500 CPU @ 3.20GHz
cpu MHz		: 2699.991
siblings	: 4
cpu cores	: 4
physical id	: 0
core id		: 0
flags		: fpu vme de pse tsc msr ... sse4_2 avx2 ...
```

- `cpu MHz` here is an **instantaneous** reading of the current frequency of
  *that* logical CPU at parse time — not a max, not stable. Don't present it as
  the CPU's rated speed; use cpufreq for min/max, or the model-name string's
  advertised base clock.
- `flags` is the feature list (SSE, AVX, virtualization, mitigations). SysKit
  reports a "flags summary," so decide which flags are interesting (e.g. `avx2`,
  `vmx`/`svm`, `aes`) rather than dumping all of them.

### `/sys/devices/system/cpu/` topology and cpufreq

sysfs is the structured, one-value-per-file interface preferred by the collector
model. Key paths:

- `/sys/devices/system/cpu/online` → e.g. `0-3`, the currently online logical
  CPUs. Relevant to **hotplug**: core count can change between your two samples.
  Match `cpuN` lines by ID across samples; do not assume core *i* in sample 1 is
  core *i* in sample 2 by array position.
- `/sys/devices/system/cpu/cpuN/topology/core_id`,
  `physical_package_id`, `thread_siblings_list` → authoritative topology, better
  than parsing `/proc/cpuinfo` for physical mapping.
- `/sys/devices/system/cpu/cpuN/cpufreq/` → frequency data **if present**.

Frequencies in cpufreq are in **kHz**, not MHz or Hz (another easy mistake):

```
$ cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq
2700009            # 2,700,009 kHz ≈ 2.70 GHz
$ cat /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq
3600000            # 3.60 GHz
$ cat /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_min_freq
800000             # 0.80 GHz
```

If the whole `cpufreq/` directory is missing (common in cloud VMs/containers),
frequency is UNAVAILABLE for that CPU — surface it as such, per the collector
error model's "missing optional data" class.

### Collector implications (tie it back to SysKit)

- The collector must go through the **platform adapter / injected filesystem**,
  not read `/proc` and `/sys` directly, so parsers run against fixtures (small
  VM, multi-core host, container, malformed data) — see `specs/collectors.md`.
- Never shell out to `lscpu`/`mpstat` from the collector. Those tools are for
  *you* to verify your parser, not for SysKit to call.
- Return raw counter structs from the collector; let the service subtract two
  snapshots. Don't smuggle percentages out of the collector.

---

## Important Files

| Path | What it gives you | Notes for the collector |
|------|-------------------|-------------------------|
| `/proc/stat` | Cumulative CPU time counters | `cpu` = aggregate, `cpuN` = per-core. 10 jiffy columns: `user nice system idle iowait irq softirq steal guest guest_nice`. Sample twice for utilization. |
| `/proc/cpuinfo` | Static per-logical-CPU identity | Fields: `processor`, `model name`, `cpu MHz`, `siblings`, `cpu cores`, `physical id`, `core id`, `flags`. Blank line separates blocks. |
| `/proc/loadavg` | Run-queue load averages | Format: `1min 5min 15min running/total lastpid` → e.g. `0.81 1.05 1.18 6/1335 188339`. Not a percentage. |
| `/sys/devices/system/cpu/online` | Online logical CPUs | e.g. `0-3`; watch for hotplug between samples. |
| `/sys/devices/system/cpu/present` | Populated logical CPUs | e.g. `0-3`. |
| `/sys/devices/system/cpu/cpuN/cpufreq/scaling_cur_freq` | Current freq (kHz) | Absent ⇒ frequency UNAVAILABLE, never 0. |
| `/sys/devices/system/cpu/cpuN/cpufreq/cpuinfo_max_freq` | Rated max freq (kHz) | Divide by 1000 for MHz. |
| `/sys/devices/system/cpu/cpuN/cpufreq/cpuinfo_min_freq` | Rated min freq (kHz) | |
| `/sys/devices/system/cpu/cpuN/cpufreq/scaling_governor` | Active governor | e.g. `powersave`, `performance`, `schedutil`. |
| `/sys/devices/system/cpu/cpuN/topology/core_id` | Physical core ID | Authoritative topology for logical→physical mapping. |
| `/sys/devices/system/cpu/cpuN/topology/physical_package_id` | Socket ID | Count distinct values for socket count. |

Runtime constant, not a file but needed: `getconf CLK_TCK` → `USER_HZ` (jiffy
rate, usually 100). In Go you would obtain this via `sysconf`/`_SC_CLK_TCK`
rather than shelling out, but the value is the same.

---

## Useful Commands

These are for **you** to understand and cross-check your parser output. SysKit
itself must not run them.

| Command | Use |
|---------|-----|
| `lscpu` | Human-readable topology: sockets, cores per socket, threads per core, model, flags, min/max MHz. Best sanity check for your topology logic. |
| `nproc` | Number of online logical CPUs. Should equal your `cpuN`-line count and your `processor:`-block count. |
| `cat /proc/stat` | See the raw counters your parser consumes. |
| `cat /proc/cpuinfo` | Inspect real field names/formatting for the parser. |
| `cat /proc/loadavg` | Confirm the 5-field format. |
| `getconf CLK_TCK` | Confirm `USER_HZ` for your jiffy math. |
| `mpstat -P ALL 1` | Per-core utilization once per second (from `sysstat`). This is essentially the two-sample algorithm you are implementing — compare your numbers to it. |
| `mpstat 1` | Aggregate utilization once per second; validates your `cpu`-line math. |
| `watch -n1 'grep cpu0 /proc/stat'` | Watch a single counter climb — makes "cumulative, not percentage" obvious. |
| `turbostat` | (root; Intel) live per-core frequency and C-states; cross-check `cpufreq` readings. |
| `cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq` | Verify kHz units and that the file exists at all. |

Handy check for the two-sample idea: read `/proc/stat`, wait a second, read
again, and diff the `cpu` line by hand — the growth in each column *is* the
per-second activity, and `busy/total` of the deltas is the utilization.

---

## References

- proc(5) man page (procfs, incl. `/proc/stat`, `/proc/cpuinfo`, `/proc/loadavg`):
  https://man7.org/linux/man-pages/man5/proc.5.html
- Kernel `Documentation/filesystems/proc.rst`:
  https://www.kernel.org/doc/html/latest/filesystems/proc.html
- Kernel CPU frequency scaling (cpufreq) docs:
  https://www.kernel.org/doc/html/latest/admin-guide/pm/cpufreq.html
- Kernel CPU topology sysfs (`Documentation/admin-guide/cputopology.rst`):
  https://www.kernel.org/doc/html/latest/admin-guide/cputopology.html
- CPU hotplug sysfs interface:
  https://www.kernel.org/doc/html/latest/core-api/cpu_hotplug.html
- sysconf(3) (for `_SC_CLK_TCK` / `USER_HZ`):
  https://man7.org/linux/man-pages/man3/sysconf.3.html
- lscpu(1): https://man7.org/linux/man-pages/man1/lscpu.1.html
- mpstat(1): https://man7.org/linux/man-pages/man1/mpstat.1.html

---

## Practical Lab

Capture the aggregate `cpu` line twice with timestamps one second apart.
Calculate every field delta, total delta, idle delta, and busy percentage. Then
compare with `syskit cpu` and `mpstat 1` as independently timed observations.
Record why a difference is expected and which layer owns parsing versus
subtraction.

## Failure-Mode Matrix

| Case | Correct behavior |
|---|---|
| cpufreq directory absent | Frequency unavailable, never zero |
| CPU appears or disappears | Match by CPU ID; derived value may be unavailable |
| Second counter is smaller | Treat as reset/replacement; avoid underflow |
| `total_delta` is zero | Utilization unavailable for the interval |
| Older line has fewer trailing fields | Accept the documented minimum safely |
| `/proc/stat` malformed | Contextual parse error, not a partial percentage |

## Checkpoint

Explain and demonstrate the two-snapshot algorithm, including real elapsed time,
counter reset, hotplug, and the distinction between load average and CPU
utilization. Record evidence in [checklists.md](checklists.md).
