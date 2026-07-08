# Process — Learning Notes

> Study notes on process management, Linux process interfaces, and related
> internals. Written for the implementer of SysKit's `process` collector. Read
> this before you write a line of the collector — the goal is not to memorize
> field names but to parse `/proc/[pid]` the way it actually behaves under a
> hostile process name and a moving target.

---

## Concepts

A process collector looks trivial ("list `/proc`, read a few files") and is one
of the most bug-prone collectors in the project. The traps are not in the data —
they are in the *shape* of the data (a command name that fights your parser) and
in the fact that the thing you are measuring **exits while you measure it**. Get
the mental model right before touching procfs.

**Process vs. thread vs. task.** In Linux these are the same kernel object wearing
different labels. The kernel schedules **tasks** (`task_struct`), each with a
**TID** (thread ID). A **process** is a group of tasks that share an address space
and other resources; the process's **PID** is the TID of its first/leader task
(the *thread group leader*, `tgid`). A **thread** is a non-leader task in that
group. This matters concretely for procfs:

- `/proc/[pid]` where `pid == tgid` is the **process** view — this is what
  `syskit process` lists. The top-level `/proc` directory contains only these
  thread-group leaders, so a plain directory listing already gives you one entry
  per process, not one per thread.
- `/proc/[pid]/task/[tid]` is the **per-thread** view. The thread *count* SysKit
  reports is simply the number of entries under `task/` (also field 20,
  `num_threads`, in `stat`). For the first cut, list processes, not threads —
  but know that "18 threads" comes from here, not from counting PIDs.

**Process states.** A process is always in exactly one state, exposed as a single
letter (field 3 of `stat`, and `State:` in `status`). SysKit reports the raw code
and should not collapse them — an operator triaging a hung box needs to tell `D`
from `S`.

- **`R` — running (or runnable).** On a CPU right now, or on a run queue waiting
  for one. "Runnable," not necessarily "currently executing."
- **`S` — interruptible sleep.** The common case: waiting for an event (I/O, a
  timer, a socket) and *will* wake on a signal. The overwhelming majority of
  processes on an idle box are `S`. Not a problem.
- **`D` — uninterruptible sleep.** Sleeping in a kernel call that must not be
  interrupted — almost always waiting on I/O (a slow or stuck disk, an NFS mount
  that went away). Here is the part newcomers miss: a `D` process **cannot be
  killed**, not even with `SIGKILL`, because it is not accepting signals until the
  kernel operation completes. `kill -9` does nothing; the process is stuck until
  the I/O returns or the box reboots. That is exactly why `D` matters for a
  diagnostic tool: a pile of `D` processes is the signature of failing storage or
  a dead network filesystem, and it is the single most useful state to surface
  during an incident. Do not treat `D` as "just sleeping."
- **`Z` — zombie.** The process has *exited* — its code is done, its memory is
  gone — but its parent has not yet **reaped** it with `wait()`, so the kernel
  keeps a stub entry alive to hold the exit status. A zombie consumes no CPU or
  memory; it holds a slot in the process table. This is why **PPID matters**: a
  zombie only clears when its *parent* reaps it, so a growing pile of zombies
  points at a buggy parent that forks children and never calls `wait()`. To
  diagnose it you must report the PPID (field 4) so the operator can find the
  guilty parent. You cannot kill a zombie (it is already dead); you fix the
  parent. Note the collector *will* still see zombies in `/proc` — they have a
  directory — but many of their fields read as empty/zero.
- **`T` — stopped.** Suspended, typically by a job-control signal (`SIGSTOP`,
  `SIGTSTP`, or `Ctrl-Z`) or by a debugger. It will not run until sent `SIGCONT`.
  Distinct from sleeping: a stopped process is deliberately frozen, not waiting on
  an event. (Modern kernels also use `t` for "tracing stop.")
- **`I` — idle kernel thread.** A kernel worker thread that is idle. Introduced so
  that idle kernel threads (which used to show as `D`) do not inflate load
  average. You will see many `I` entries with names in brackets like `[kworker/…]`.
  They are normal; do not mistake `I` for a problem or confuse it with `R`.

(You may also encounter `X`/`x` (dead) transiently and `S`/`D` with a trailing
`+`/`l`/`s` in `ps` output — those extra flags are `ps` decorations, not part of
the `stat` state field. Parse only the single leading letter from procfs.)

**Parent-child relationships and the tree.** Every process except `init`/`systemd`
(PID 1) has a parent, named by its **PPID** (field 4 of `stat`, `PPid:` in
`status`). The collector's job is *not* to build the tree — per the collector
architecture it returns a flat, point-in-time snapshot of processes, each carrying
its own PID and PPID. The **service layer** assembles the parent→child tree for
`syskit process tree`. Two realities the tree builder must handle, so the
collector must faithfully preserve PPID even when it looks odd:

- **Reparenting / orphans.** When a parent dies before its child, the child is
  re-parented (to PID 1, or to a nearer "subreaper"). So a child's PPID can point
  at a process that is not its original creator, and PID 1 accumulates a broad set
  of children. The tree builder must tolerate a PPID whose process it has already
  seen *or* has not — never assume the parent appears before the child in your
  snapshot.
- **PID 1 is the root.** The tree's root is PID 1; anything whose PPID is 0 is a
  kernel thread parented under the kernel (`[kthreadd]`, PID 2, and its children).

**The snapshot is a lie by the time you print it — and that is fine.** The
collector returns a point-in-time snapshot. Between listing `/proc` and reading a
process, the world changes: PIDs appear, exit, and (eventually) get reused.
SysKit does not promise a consistent instantaneous cross-section of the kernel;
it promises a best-effort snapshot with **fields you could not read marked
unavailable, never faked as zero**. A CPU counter you were denied is not `0`; a
command you could not read is not `""`. This "unavailable is not zero" rule is the
same discipline the memory and network notes hammer on, and it is the difference
between an honest tool and one that invents data.

**Disappearing processes are the normal case, not an error.** A PID directory can
vanish at any instant — the process exits and the kernel tears its `/proc` entry
down while you are mid-read. You will list PID 4321, then `open("/proc/4321/stat")`
returns `ENOENT`, or a read returns `ESRCH`. This is **expected and routine** on a
busy host, not a failure. The rule: treat `ENOENT`/`ESRCH` on any `/proc/[pid]/*`
path as "this process exited — skip it and move on." One vanished process must
never abort the whole collection; you drop that entry and continue. This is
called out in the collector error-classification rules as "race with disappearing
resources, especially processes," and it is the single most important robustness
property of this collector.

**Permission limits — what you can and cannot read about other users.** Most
per-PID files are **world-readable**: any user can read `stat`, `status`,
`cmdline`, and `comm` for *any* process, which is why `ps aux` works unprivileged.
But some data is **restricted** to the owning user (or a process with
`CAP_SYS_PTRACE` / root):

- `/proc/[pid]/fd/` and `/proc/[pid]/fdinfo/` — open file descriptors. Readable
  only for your own processes without `CAP_SYS_PTRACE`.
- `/proc/[pid]/environ` — the process's environment (can contain secrets).
  Owner/ptrace only.
- Certain `status` fields and `wchan`/`stack` details for other users' processes.

The `ptrace_scope` sysctl (`/proc/sys/kernel/yama/ptrace_scope`) can tighten this
further. The collector rule mirrors the disappearing-process rule: hitting
`EACCES`/`EPERM` on a restricted field of another user's process yields
**unavailable / partial data for that field**, never a fatal error and never a
fabricated value. `syskit process` run as a normal user should still list every
process (names, states, PPIDs are world-readable) and simply mark the
privileged bits it could not see as unavailable.

---

## Linux Internals

### The `/proc/[pid]/stat` safe-parse rule — read this twice

`/proc/[pid]/stat` is a single line of space-separated fields, and it is the
richest cheap source of per-process data. It also contains the one parsing trap
that breaks more process tools than anything else, so this is the headline of the
whole document.

**Field 2 (`comm`) is the executable name wrapped in parentheses — and it can
contain spaces *and* parentheses itself.** The comm field is up to 15 characters
copied from the program, and a program can name itself almost anything. A process
can legitimately be called `foo bar`, or `)( )(`, or — the adversarial case that
matters — a name crafted to look like the rest of the line. So a real `stat` line
can look like:

```text
1234 (evil) 0 R) 1 1234 1234 0 -1 4194304 …
```

Here the actual command is literally `evil) 0 R` — a name chosen to impersonate
fields. **Naive whitespace splitting is broken.** If you `strings.Fields(line)`
and index by position, the spaces and the `)` inside `comm` shift every field
after it, so field 3 is no longer the state, field 4 is no longer the PPID, and
every number you read for that process is wrong. On a benign host you will never
notice; the day a process has a space in its name, your PPIDs and states silently
corrupt — and a security tool that can be blinded by naming a process cleverly is
worse than no tool.

**The robust parse — anchor on the LAST `')'`, not the first.** The trick the
kernel's own format guarantees: field 1 (PID) is a plain integer before the first
`(`, and every field from 3 onward is a well-behaved space-separated token with no
parentheses. Only `comm` (field 2) contains arbitrary bytes, and it is the *only*
thing between the first `(` and the last `)`. So:

1. Read the whole line.
2. Take everything **before the first `(`** → that is field 1, the **PID**.
3. Take everything **between the first `(` and the LAST `)`** → that is field 2,
   **`comm`** (the command name, verbatim, spaces and inner `)` included).
4. Take everything **after the last `)`** → a clean, space-delimited list of the
   remaining positional fields, and **the first token of that remainder is field
   3, the state.**

Worked example on a hostile line for a process whose name is literally
`evil) 0 R` (chosen to imitate the "comm) state ppid" shape):
`1234 (evil) 0 R) S 1 1234 1234 0 -1 4194304 …`

- Before first `(` → `1234` → **PID = 1234**. ✓
- Between first `(` and **last** `)` → `evil) 0 R` → **comm = "evil) 0 R"**. ✓
  (A first-`)`-based parser would wrongly stop at `evil` and treat `0 R) …` as
  positional fields — exactly the corruption we are avoiding.)
- After last `)` → `S 1 1234 1234 0 -1 4194304 …` → split on spaces; the **first
  token is the genuine state → `S`**, then ppid `1`, and so on. A first-`)` parser
  would instead have read the fake `0`/`R` embedded in the name as state and ppid —
  every downstream field shifted by the injected tokens. Anchor on the **last** `)`
  and every field lines up; anchor on the first and everything shifts.

Everything after the last `)` is then positional and safe to index. The fields
SysKit cares about, counting the state as field 3 (1-indexed, matching `proc(5)`):

- **Field 3 — `state`** — the single state letter (`R`, `S`, `D`, `Z`, `T`, `I`).
- **Field 4 — `ppid`** — parent PID (drives the tree and zombie diagnosis).
- **Field 14 — `utime`**, **Field 15 — `stime`** — CPU time in **user** and
  **kernel** mode, in **clock ticks** (jiffies; `sysconf(_SC_CLK_TCK)`, typically
  100/sec). These are *cumulative counters*, not percentages — the same species as
  the CPU and network counters in the sibling notes. The collector returns the raw
  ticks; the **service layer** does the two-sample delta to derive CPU%. Do not
  compute a percentage inside the collector.
- **Field 22 — `starttime`** — when the process started, in clock ticks **since
  boot**. To get a wall-clock start time you combine it with the system boot time
  (`btime` in `/proc/stat`) and `_SC_CLK_TCK` — again service-layer territory.
- **Field 24 — `rss`** — resident set size in **pages** (multiply by page size,
  `sysconf(_SC_PAGESIZE)`, usually 4096, to get bytes). Note this is pages, not
  kB — different unit convention from `status`'s `VmRSS`.

(Field numbers follow `proc(5)`; there are ~52 fields. Read the man page rather
than trusting a hard-coded index you half-remember.)

### `/proc/[pid]/status` — the human-readable alternative

Same information, laid out as `Key:\tValue` lines — far friendlier to parse and
**immune to the comm-parentheses trap** because keys are fixed and values are on
their own lines. Useful fields: `Name:` (the comm, still truncated to 15 chars but
without the parenthesis-embedding hazard), `State:` (`S (sleeping)` — letter plus
word), `PPid:`, `Uid:` (real/effective/saved/fs — the **real UID is the first
number**, which you map to a username), `VmRSS:` (resident memory in **kB**),
`Threads:`, `VmSize:`. Trade-off: `status` is easier and safer to parse but does
**not** give you the CPU tick counters (`utime`/`stime`) or `starttime` — those
live only in `stat`. Practical approach: read `stat` for the numeric
accounting/CPU fields, and `status` for name/UID/memory where its clarity and
safety win. Parse `status` by key, never by line number — kernels add and reorder
lines across versions (same discipline as `/proc/meminfo`).

### `/proc/[pid]/cmdline` — NUL-separated, another gotcha

This is the full command line **as the process was invoked** (argv), and it is the
right thing to display as `COMMAND` — richer than the 15-char `comm`. The trap:
**arguments are separated by NUL bytes (`\0`), not spaces**, and the string
usually ends with a trailing `\0`. If you read it as text and split on spaces you
get one giant blob and mangle any argument that itself contains a space (a file
path with a space collapses into its neighbors). **Split on `\0`**, drop the empty
trailing element, and join with a single space for display. Two edge cases:

- **Kernel threads have an empty `cmdline`** — zero bytes. That is how you tell a
  kernel thread from a userspace process without a heuristic: empty `cmdline` →
  kernel thread → fall back to `comm` (which for kernel threads reads like
  `kworker/0:1`), and by convention display it in brackets (`[kworker/0:1]`), the
  way `ps` does.
- A zombie's `cmdline` is also typically empty (its memory is gone) — fall back to
  `comm` there too.

### `/proc/[pid]/comm` — the short name, writable and truthful

The bare command name (the same 15-char value as `stat`'s comm and `status`'s
`Name`), one line, no parentheses, no NULs — the cleanest source for the short
name. Note a process can change its own comm at runtime (via `prctl(PR_SET_NAME)`
or writing this file), so `comm` and `cmdline[0]` can legitimately disagree; that
is not corruption, just a process that renamed its thread.

### Putting the sources together (collector shape)

The point-in-time snapshot per process is: list numeric entries in `/proc` (skip
non-numeric names like `self`, `net`, `meminfo`); for each PID, read `stat`
(numeric/CPU fields via the last-`)` parse), optionally `status` (UID → user,
VmRSS), and `cmdline` (display command, NUL-split, with the kernel-thread
fallback). Wrap **every** open/read so that `ENOENT`/`ESRCH` skips the process and
`EACCES`/`EPERM` marks that field unavailable. Return raw counters and let the
service layer own CPU%, MEM%, wall-clock start time, and the tree.

---

## Important Files

- **`/proc/[pid]/stat`** — single-line, space-separated per-process accounting.
  Fields: 1 `pid`, 2 `comm` (**parenthesized, parse via last `)`**), 3 `state`,
  4 `ppid`, 14 `utime`, 15 `stime`, 20 `num_threads`, 22 `starttime`, 24 `rss`
  (pages). CPU times in clock ticks; RSS in pages. Raw counters only — no rates in
  the collector. The trap file; write the hostile-comm test first.
- **`/proc/[pid]/status`** — human-readable `Key:\tValue`. `Name`, `State`,
  `PPid`, `Uid` (real UID = first field), `VmRSS` (kB), `Threads`. Safe from the
  comm trap; parse by key, not position. No CPU tick counters here — use `stat`
  for those.
- **`/proc/[pid]/cmdline`** — full argv, **NUL-separated**, trailing `\0`. Split
  on `\0`, not spaces. **Empty for kernel threads and zombies** → fall back to
  `comm`, display in brackets. The display `COMMAND` source.
- **`/proc/[pid]/comm`** — short (≤15 char) command name, one clean line. Can
  differ from `cmdline[0]` if the process renamed itself; not an error.
- **`/proc/[pid]/fd/`** — symlinks to the process's open files; socket fds point at
  `socket:[inode]` (the join key the `ports` collector uses). **Restricted**:
  readable only for your own processes without `CAP_SYS_PTRACE`/root → mark
  unavailable for others, do not error. `ls -l` here to see targets.
- **`/proc/[pid]/cgroup`** — which cgroup(s) the process belongs to (the container
  grouping key). v2 shows a single `0::/…` line. Basis for future container-aware
  grouping (a listed future extension, not first-cut scope).
- **`/proc/[pid]/limits`** — the process's resource limits (`Max open files`,
  `Max stack size`, etc.), soft and hard, human-readable. Context for "why did
  this process hit a ceiling"; not core to the first list view.
- **`/proc/loadavg`** — system-wide, not per-PID: the 1/5/15-minute load averages,
  then `runnable/total` process counts, then the last-created PID. A cheap
  sanity check on how many processes should exist and how many are `R`/`D`.

---

## Useful Commands

These are for **verification only** — run them by hand to check your parser
against ground truth. **SysKit must not shell out to any of them** (the process
spec is explicit: "No external `ps` command is executed," per ADR-003). The
collector reads `/proc` directly; these tools just let you diff your output
against a trusted reference.

- `ps aux` — the classic per-process table (user, PID, %CPU, %MEM, STAT, START,
  COMMAND). The reference your `syskit process` columns should broadly match.
  Watch how `ps` decorates the state (`Ss`, `R+`, `D<`) — those extra letters are
  `ps` flags, not the `stat` state field.
- `ps -eLf` — one row **per thread** (`-L`), with `NLWP` (thread count) and `LWP`
  (TID). Cross-check your `num_threads` / `task/` counting.
- `ps -o pid,ppid,state,comm -p <pid>` — read one process's raw fields to compare
  directly against your `stat`/`status` parse, especially for a process with a
  space in its name.
- `pstree -p` — the parent→child tree with PIDs. Ground truth for
  `syskit process tree`, including reparented children hanging off PID 1.
- `top` / `htop` — live process view; useful to watch states change (spot a `D`
  process during simulated I/O stalls) and to confirm CPU%/MEM% derivation.
- `ls -l /proc/<pid>/fd` — see the fd symlinks (and `socket:[inode]` targets);
  run it against another user's PID unprivileged to observe the `EACCES` you must
  handle as unavailable.
- `cat /proc/1/status` — inspect PID 1 (`systemd`) directly: `Name`, `State`,
  `PPid` (0), `Uid`, `VmRSS`, `Threads`. A stable, always-present target for a
  fixture.
- `cat /proc/self/stat` — see a real `stat` line for the reader itself; a safe way
  to eyeball the field layout.

To *manufacture* the hostile case for a test fixture: start a process that renames
its comm to include a space and a `)`, then `cat` its `stat` and confirm your
last-`)` parser still recovers the correct state and PPID.

---

## References

- `proc(5)` — the authoritative `/proc/[pid]/stat`, `status`, `cmdline`, `comm`
  field reference (note the explicit warning about the parenthesized comm field):
  https://man7.org/linux/man-pages/man5/proc.5.html
- `proc_pid_stat(5)` — the dedicated `/proc/pid/stat` field-by-field page:
  https://man7.org/linux/man-pages/man5/proc_pid_stat.5.html
- Linux kernel — `Documentation/filesystems/proc.rst` (procfs process files and
  the `stat` field table): https://www.kernel.org/doc/html/latest/filesystems/proc.html
- `ps(1)` — process state codes (`R`/`S`/`D`/`Z`/`T`/`I`) and their meanings:
  https://man7.org/linux/man-pages/man1/ps.1.html
- `credentials(7)` / `ptrace(2)` — real vs. effective UID and the
  `CAP_SYS_PTRACE`/`ptrace_scope` limits on reading other users' process data:
  https://man7.org/linux/man-pages/man7/credentials.7.html
- ADR-003 — Read native kernel interfaces instead of shelling out:
  `../decisions/003-native-apis-over-shell.md`

---

## Personal Notes
