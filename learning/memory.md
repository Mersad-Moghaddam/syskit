# Memory — Learning Notes

> Study notes on Linux memory accounting and the kernel interfaces SysKit's
> memory collector will read. Read this before you write a line of the collector.
> The goal is not to memorize field names — it is to report memory the way an
> operator actually needs to interpret it.

---

## Concepts

Linux memory accounting is the single most misread area of system inspection,
and SysKit exists partly to stop that misreading. Get these mental models right
before touching `/proc`.

**Physical vs. reclaimable vs. free.** On a healthy Linux box, "free" memory is
almost always small — the kernel deliberately spends idle RAM on the page cache,
because unused RAM is wasted RAM. That cached data can be dropped the instant a
program needs the pages. So there are three quantities, not two:

- **Free** — pages that are wholly unused right now.
- **Reclaimable** — pages holding cache/buffers/reclaimable slab that the kernel
  can hand back without going to disk.
- **Available** — the kernel's estimate of what a new workload could get its
  hands on = roughly free + most of reclaimable, minus what must be kept.

**The classic newcomer bug — do not report MemFree when the user wants
"available."** This is the mistake SysKit is built to avoid, so it is worth
stating bluntly:

> `MemAvailable` is the kernel's estimate of how much memory a new process could
> allocate **without swapping**. It already accounts for reclaimable page cache
> and reclaimable slab. `MemFree` is only the wholly-unused pages — it ignores
> the many gigabytes of cache the kernel would happily evict on demand.

If the collector reports `MemFree` under an "available" or "used" label, then on
a normal server showing 200 MiB free and 12 GiB of reclaimable cache, the user
concludes they are out of RAM and starts killing processes or adding swap — when
in reality ~12 GiB is instantly allocatable. "Used" computed as
`MemTotal - MemFree` produces the same false alarm: it counts reclaimable cache
as used. Compute practical "used" as `MemTotal - MemAvailable`, and surface
buffers/cache separately so the user sees where the rest went. This is the whole
reason `free -h` grew its "available" column and why `free` output confused
people for a decade before it did.

**Buffers vs. cache.** Both are page-cache-backed and both are largely
reclaimable, but they hold different things:

- **Buffers** — block-device / filesystem metadata (raw block I/O, superblocks,
  directory data). Historically the "buffer cache." Usually small.
- **Cached** — the page cache: contents of files that have been read or written.
  Usually the large number. Note `Cached` in `/proc/meminfo` includes tmpfs and
  shared memory, which are *not* trivially reclaimable, which is one reason the
  kernel computes `MemAvailable` for you rather than expecting you to derive it.

Report both, but do not lecture the user by trying to sum them into "used" —
that math is exactly what `MemAvailable` already does more correctly.

**Swap usage is not memory pressure.** Seeing swap in use does not mean the
system is thrashing. Under low pressure the kernel proactively swaps out pages
it judges cold (tunable via `vm.swappiness`) to free RAM for cache. Long-idle
daemon pages can sit in swap for weeks on a machine with gigabytes free. What
signals *pressure* is the *rate* of swapping (paging in/out under contention) and
— far better — PSI (below). So: report swap used as a fact, but never equate a
nonzero swap figure with "under pressure."

**Pages, and why units matter.** The kernel accounts memory in pages (typically
4 KiB). `/proc/meminfo` pre-converts most fields to kB for you, so the collector
mostly multiplies by 1024. The trap is that `/proc/meminfo`'s "kB" means
**kibibytes (1024 bytes)** despite the lowercase label — the kernel is lying
about the unit name, not the value. Multiply by 1024, store bytes, done.

---

## Linux Internals

**`/proc/meminfo` is the primary source.** One line per field:
`Name:` then whitespace then a number then (for most fields) `kB`. A handful of
fields (e.g. `HugePages_Total`) are bare counts with no unit. Parse defensively:
split on the colon, trim, and treat a trailing `kB` as "multiply by 1024." Never
assume a fixed set of lines or a fixed order — kernels add and reorder fields
across versions. Read the fields you need by name; ignore the rest.

**`MemAvailable` is a computed estimate, and it is version-gated.** The kernel
computes it from free pages, the low watermark, and the reclaimable portions of
the page cache and slab. Crucially:

> `MemAvailable` was added in **Linux 3.14** (commit `34e431b0ae39`). On any
> older kernel the line is simply **absent** from `/proc/meminfo`.

This is the collector's most important edge case. When `MemAvailable` is missing,
SysKit must represent it as **UNAVAILABLE** — a distinct state — and must **not**:

- fall back to `0` (implies zero allocatable memory — a lie), or
- silently substitute `MemFree` (the exact bug this whole document warns about).

The domain type should be able to say "this field was not present" so rendering
can print `n/a` and structured output can emit `null`, not a fabricated number.
The spec's acceptance criterion "Missing PSI reports unavailable pressure data"
applies in spirit to `MemAvailable` too: absence is information, not zero.

**Which fields may be missing.** Plan the type so each of these can be absent:

- `MemAvailable` — kernels < 3.14.
- `SwapTotal` / `SwapFree` — present but `0` when no swap is configured (report
  zero, not an error, per the spec).
- `SReclaimable` / `SUnreclaim` — split slab accounting; older kernels may show
  only a combined `Slab`.
- PSI (`/proc/pressure/*`) — see below; whole interface may not exist.

**`/proc/vmstat`** holds the cumulative counters behind pressure: `pgpgin`,
`pgpgout` (block I/O pages), `pswpin`, `pswpout` (swap pages), `pgfault`,
`pgmajfault`. These are monotonic counters — meaningful only as a rate between
two snapshots, which per the collector model is the service layer's job, not the
collector's. You likely will not need vmstat for the first cut, but it is where
"is it actively swapping *right now*" really lives if PSI is unavailable.

**Slab.** `Slab = SReclaimable + SUnreclaim`. `SReclaimable` (dentry/inode
caches, etc.) folds into `MemAvailable`; `SUnreclaim` is kernel memory that
cannot be handed back. Report `Slab` and, when present, the split.

**cgroup memory (context matters).** Inside a container, `/proc/meminfo` usually
still shows the *host's* memory, not the container's limit — a well-known
footgun. The true limit lives in the cgroup files (below). Full cgroup-limit
reporting is a listed future extension, not first-cut scope, but be aware the
numbers from `/proc/meminfo` can be misleading in a container and note it rather
than presenting host memory as if it were the container's.

---

## PSI — Pressure Stall Information

PSI is the modern, correct answer to "is memory actually a problem?" Instead of
inferring pressure from free memory (wrong) or swap usage (also wrong), PSI
measures the **time tasks were stalled waiting on memory**.

**File:** `/proc/pressure/memory`. Format is two lines:

```text
some avg10=0.00 avg60=0.00 avg300=0.00 total=0
full avg10=0.00 avg60=0.00 avg300=0.00 total=0
```

- **`some`** — the share of time *at least one* task was stalled waiting for
  memory (paging, reclaim). Early warning: some work is being delayed.
- **`full`** — the share of time *all* runnable tasks were stalled
  simultaneously — i.e. nobody was making progress. This is the serious signal;
  sustained nonzero `full` means real trouble.
- **`avg10/avg60/avg300`** — percentages of stalled time over the last 10, 60,
  and 300 seconds (already averaged by the kernel — no two-sample math needed).
- **`total`** — cumulative stall time in **microseconds** since boot.

**Why PSI beats free memory as a pressure signal.** Free memory tells you a
level; it cannot distinguish "low free because of harmless cache" from "low free
and thrashing." PSI measures the actual symptom — stalled work — so a rising
`some.avg10` catches pressure building before the OOM killer fires, and
`full.avg60` distinguishes a busy-but-fine box from one grinding on reclaim.
This maps directly to the spec's "PRESSURE: low/…" summary column: derive that
label from PSI, not from a free-memory ratio.

**Version and config gate — PSI may be absent.**

> PSI landed in **Linux 4.20** and requires **`CONFIG_PSI=y`**. Some distros
> ship it disabled by default and require the boot parameter **`psi=1`** to
> enable it.

So `/proc/pressure/memory` can be missing entirely (old kernel or PSI compiled
out) or present-but-disabled. The collector must treat this as
**pressure = unavailable**, exactly as the spec requires ("Missing PSI reports
unavailable pressure data"), never as `0.00` (which would read as "no
pressure" — a materially different and misleading claim).

---

## Important Files

| Path | What it gives you | Notes for the collector |
|------|-------------------|-------------------------|
| `/proc/meminfo` | Totals, free, available, buffers, cache, swap, slab | Primary source. Values in kibibytes despite the `kB` label; multiply by 1024. Parse by name, not position. |
| `/proc/pressure/memory` | PSI `some`/`full` stall averages + totals | Absent if kernel < 4.20 or `CONFIG_PSI`/`psi=1` off → report pressure unavailable. `total` is microseconds. |
| `/proc/swaps` | Per-swap-area table (device, type, size, used, priority) | Sizes here are in **kB (1024-byte)** rows; useful to enumerate swap devices. If only the header line exists, no swap is configured → report zero, not error. |
| `/proc/vmstat` | Cumulative paging/swap counters (`pswpin`, `pswpout`, `pgmajfault`, …) | Rates only; a two-snapshot concern owned by the service layer. Optional for first cut. |
| `/sys/fs/cgroup/memory.current` | cgroup v2: current memory usage in bytes | Already in **bytes**. Relevant inside containers where `/proc/meminfo` shows the host. |
| `/sys/fs/cgroup/memory.max` | cgroup v2: memory hard limit (or `max`) | Literal string `max` means unlimited — parse it, do not choke. |
| `/sys/fs/cgroup/memory.stat` | cgroup v2: detailed breakdown (anon, file, slab, …) | For future per-cgroup accounting. |
| `/sys/fs/cgroup/memory.pressure` | cgroup v2: PSI scoped to this cgroup | Same format as `/proc/pressure/memory`; per-container pressure. |

(cgroup v1 uses `/sys/fs/cgroup/memory/memory.*` with different names, e.g.
`memory.usage_in_bytes` / `memory.limit_in_bytes`. Detecting v1 vs. v2 is a
listed future extension — know the layout differs, do not build it yet.)

---

## Useful Commands

These are for *you*, to verify the collector against ground truth. Per the
collector architecture SysKit must **not** shell out to these — they are study
and test-verification tools only.

- `cat /proc/meminfo` — see the raw fields and units the collector parses.
  Check whether `MemAvailable` is present on your test host/kernel.
- `free -h` — human summary. Compare its **available** column against your
  `MemTotal - MemAvailable` "used" math; they should agree. Watch how `free`
  keeps free small and available large on a cache-warm box — that is the whole
  lesson in one command.
- `free -w` — widens output to show buffers and cache as separate columns
  (matches the fields SysKit reports separately).
- `vmstat 1` — live paging: the `si`/`so` (swap in/out) and `bi`/`bo` columns
  show whether the box is *actively* swapping vs. merely holding pages in swap.
  This is the difference between "swap used" and "under pressure."
- `cat /proc/pressure/memory` — inspect PSI directly; on a machine without it
  the file is absent, which is exactly the unavailable case to handle.
- `cat /proc/swaps` — enumerate swap areas; run with swap off to see the
  header-only "no swap" shape.
- `slabtop` — live slab breakdown; context for `Slab` / `SReclaimable`.

---

## References

- Linux kernel — `/proc/meminfo` field reference (`Documentation/filesystems/proc.rst`):
  https://www.kernel.org/doc/html/latest/filesystems/proc.html
- `proc(5)` man page (meminfo, vmstat, swaps):
  https://man7.org/linux/man-pages/man5/proc.5.html
- Linux kernel — PSI documentation (`Documentation/accounting/psi.rst`):
  https://www.kernel.org/doc/html/latest/accounting/psi.html
- `MemAvailable` — original patch and rationale (kernel commit `34e431b0ae39`, Linux 3.14):
  https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=34e431b0ae398fc54ea69ff85ec700722c9da773
- `free(1)` man page (available vs. free columns):
  https://man7.org/linux/man-pages/man1/free.1.html
- cgroup v2 memory controller (`memory.current`, `memory.max`, `memory.pressure`):
  https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html#memory
- Facebook/Meta PSI overview (why stall time beats free-memory heuristics):
  https://facebookmicrosites.github.io/psi/docs/overview

---

## Personal Notes
