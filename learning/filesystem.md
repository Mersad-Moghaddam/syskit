# Filesystem: Mounts, Capacity, And Inodes

> Study notes on filesystem internals, VFS, and related Linux interfaces.
> Focus for SysKit: how to enumerate mounts robustly, read mount options, and
> surface both byte usage AND inode usage so we catch failures that a naive
> "disk full?" check would miss.

---

| Attribute | Value |
|---|---|
| Level | Domain |
| Prerequisites | [Disk](disk.md), [kernel interfaces](kernel-interfaces.md) |
| Time | 2–3 hours |
| Product contract | [Filesystem feature](../specs/features/filesystem.md) |

## Learning Objectives

- Explain VFS, filesystem type, mount, bind mount, and mount namespace.
- Parse mountinfo including separators and octal escapes.
- Calculate byte and inode usage using the correct statfs fields.
- Distinguish virtual, remote, overlay, and block-backed filesystems.
- Test zero totals, vanished mounts, long options, and observer-relative views.

```mermaid
flowchart LR
    M[/proc/self/mountinfo] --> E[Mount entries]
    E --> S[statfs per mount]
    S --> B[Byte capacity]
    S --> I[Inode capacity]
    E --> J[Joined filesystem view]
    B --> J
    I --> J
```

## Concepts

**What the collector actually produces.** For each mounted filesystem SysKit
reports: mount point, source (backing device or fs name), filesystem type,
mount options (as an array, not a comma-joined string), and usage — both bytes
and inodes with a used-percent for each. Everything else in this note exists to
help you produce those fields correctly.

**Real vs virtual/pseudo filesystems.** Not everything in the mount table is a
disk-backed filesystem. Broadly:

- *Real / disk-backed*: `ext4`, `xfs`, `btrfs`, `f2fs`, `vfat`. Backed by a
  block device; byte and inode usage are meaningful.
- *Virtual / pseudo*: `proc`, `sysfs`, `cgroup`/`cgroup2`, `devtmpfs`,
  `tmpfs`, `overlay`, `squashfs`, `fuse.*`, `debugfs`, `tracefs`, `mqueue`,
  `bpf`. These have no real block device. Some (`tmpfs`) do report usage;
  others report zero or garbage.

Why we care: a typical host has 30–40 mounts, most of them pseudo. If we dump
them all, the signal (the ext4 root filling up) drowns in noise. The spec says
pseudo filesystems are **hidden by default** (`--show-pseudo` reveals them) or
**clearly marked**. Decide the pseudo set explicitly — don't guess from the fs
name at render time.

**Bind mounts and overlays duplicate sources.** A bind mount makes the same
underlying data appear at two mount points; an overlay (containers) stacks
lower+upper dirs. The same `major:minor` device can appear many times. That is
correct, not a bug — but if you sum usage across all mounts to get "total disk
used" you will double-count. Deduplicate by `major:minor` when aggregating.

**Mount namespaces.** `/proc/self/mountinfo` shows the mounts *of the process
that reads it*. Inside a container, that is the container's view, not the
host's. SysKit reads its own `self` view intentionally — it reports the
filesystems visible to SysKit. Just know that "why don't I see the host's
mounts?" inside a container is expected behavior, not a collector failure.

---

## Linux Internals

### Reading the mount table: mountinfo vs mounts vs mtab

There are three ways to ask "what is mounted?" They are not equal.

- **`/etc/mtab`** — historically a real file the `mount` command wrote to. On
  modern systems it is a symlink to `/proc/self/mounts`. Never write to it,
  never trust it as authoritative. Legacy.
- **`/proc/mounts`** (a.k.a. `/proc/self/mounts`) — kernel-generated,
  fstab-style: `source mountpoint fstype options dump pass`. Better than mtab
  because the kernel maintains it, but it is *ambiguous*: it does not give you
  device numbers, mount IDs, or propagation info, and the option string mixes
  per-superblock and per-mount options together.
- **`/proc/self/mountinfo`** — the robust source, and what SysKit should
  parse. It is unambiguous, exposes `major:minor` device numbers (so you can
  join to block-device / diskstats data), separates per-mount options from
  per-superblock options, and carries propagation/namespace info.

**Why mountinfo is preferred, concretely:** it lets you distinguish two mounts
of the same device, correlate a mount to its block device via `major:minor`,
and see mount options (`ro`, `relatime`) separately from superblock options.
`/proc/mounts` cannot do the device correlation at all. Prefer mountinfo;
fall back to `/proc/mounts` only if mountinfo is somehow unavailable.

### mountinfo line format

Each line has fixed leading fields, then a variable number of optional fields,
then a ` - ` separator, then three trailing fields:

```
36 35 98:0 /mnt1 /mnt2 rw,noatime shared:1 - ext4 /dev/root rw,errors=remount-ro
│  │  │    │     │     │          │        │ │    │         │
1  2  3    4     5     6          7        │ 9    10        11
                                          8 = separator " - "
```

1. **mount ID** — unique id for this mount.
2. **parent ID** — mount id of the parent (builds the mount tree).
3. **major:minor** — device number. Join key to block devices / diskstats.
4. **root** — the pathname *within* the filesystem that is mounted (e.g. `/`
   for a normal mount, a subdir for a bind mount).
5. **mount point** — where it is mounted in *our* namespace.
6. **mount options** — *per-mount* options: `ro`/`rw`, `relatime`, `nosuid`,
   `nodev`, `noexec`, etc.
7. **optional fields** — zero or more tags like `shared:1`, `master:2`,
   `propagate_from:`, `unbindable`. This is the propagation info. **Its length
   is variable** — this is the field that trips people up.
8. **separator** — a single ` - `. This is how you find the end of the
   optional fields: scan forward until you hit the token `-`.
9. **filesystem type** — `ext4`, `tmpfs`, `overlay`, ...
10. **mount source** — backing device (`/dev/nvme0n1p2`) or fs name (`tmpfs`).
11. **super options** — *per-superblock* options, e.g. `rw,errors=remount-ro`.

**Parsing rule:** split on whitespace, take fields 1–6 positionally, then walk
forward collecting optional fields until you reach the literal `-` token; the
three fields after `-` are type, source, super options. Do **not** assume a
fixed column count — the optional-fields section makes that wrong.

**Mount options for the report:** fields 6 and 11 together describe the mount.
Field 6 (`ro,relatime,nosuid`) is what users usually mean by "mount options."
Keep them as an array. The checklist asks you to *identify mount options* from
mountinfo — that is field 6, sometimes enriched by field 11.

### The newcomer bug: spaces and octal escaping

Mount points and paths can contain spaces, tabs, and newlines. The kernel does
**not** quote them — it octal-escapes them in these fields (root, mount point,
and the source). The mapping:

- space → `\040`
- tab → `\011`
- newline → `\012`
- backslash → `\134`

So a mount at `/mnt/my backup` appears in the file as `/mnt/my\040backup`.

If you naively split the line on spaces and treat field 5 as the mount point,
a path with a space silently shifts every subsequent field by one — your
"filesystem type" becomes half a path, your options vanish, and the parse is
quietly corrupt. **You must octal-decode fields 4, 5, and 10 after splitting.**
Split on whitespace first (the escapes guarantee no real whitespace remains
inside a field), then decode `\ooo` sequences back to bytes. This is the same
class of bug as the process collector's "command names with spaces" — the
kernel gives you an escaped, split-safe token; honor the escaping.

### Usage: statfs / statvfs

Mountinfo tells you *what* is mounted; it does not tell you *how full*. For
that you call `statfs(2)` (or POSIX `statvfs(3)`) on the mount point. The
fields that matter:

Bytes:
- **`f_bsize`** / **`f_frsize`** — block size. Multiply block counts by this to
  get bytes. (`f_frsize` is the fragment/fundamental size; use it for the math.)
- **`f_blocks`** — total data blocks in the filesystem.
- **`f_bfree`** — free blocks, *including* root-reserved blocks.
- **`f_bavail`** — free blocks available to **unprivileged** users.

Inodes:
- **`f_files`** — total inodes.
- **`f_ffree`** — free inodes.

**The `bfree` vs `bavail` trap.** ext-family filesystems reserve ~5% of blocks
for root (so a full disk doesn't lock out root and fragment badly). `f_bfree`
counts that reserve as free; `f_bavail` does not. A normal user cannot use the
reserved blocks. If you report `f_bfree` as "free space" you **overstate**
what users can actually write — and your used-percent won't match `df`. This
is exactly how `df` computes it:

- Used = `f_blocks - f_bfree`
- Available (to users) = `f_bavail`
- Used% = `used / (used + f_bavail)` rounded up.

Report **`f_bavail`** as the free/available figure. Use `f_bfree` only if you
explicitly want the raw kernel-free number, and label it as such.

**`statfs` failure modes.** Some filesystems don't report inodes: `f_files`
and `f_ffree` come back as `0`. Per the spec, "inode unavailable" must be
**distinct from zero inodes** — model it as absent/optional, not the number 0,
or you'll render a real filesystem as "100% inodes used" (0 free of 0). A
`statfs` call can also block or fail on a stale/hung network mount (NFS, FUSE);
be prepared for the syscall to error or hang, and classify that as
missing/partial data rather than crashing the whole collection.

### Inode exhaustion (headline point)

An inode is the on-disk record for one file/dir/symlink — it holds metadata
and block pointers, but **not** the file name. Most disk-backed filesystems
allocate a **fixed number of inodes at mkfs time**, proportional to size (ext4
default is roughly one inode per 16 KB of space). That pool does not grow.

Consequence: you can exhaust inodes while the disk is nearly empty by bytes.
A directory full of millions of tiny files (mail spools, session/cache files,
`node_modules`, PID/lock files) consumes one inode each but few bytes. When
`f_ffree` hits 0, **every `create`/`open(O_CREAT)`/`mkdir` fails with ENOSPC —
"No space left on device"** — even though `df -h` shows the disk at, say, 3%.
This is a genuinely confusing production incident: the error screams "disk
full," the byte gauge says otherwise, and the real culprit is invisible unless
you look at inodes.

**This is why SysKit reports inode usage as a first-class column, not just
bytes.** A byte-only tool would show a healthy filesystem right up until writes
mysteriously fail. Compute inode used-percent from `f_files`/`f_ffree` the same
way as bytes, and surface it alongside. Being able to explain this scenario is
a checklist item — the one-liner: *a filesystem can be 0% full by bytes yet
reject all writes because its fixed inode pool is exhausted by many tiny files.*

(Note: XFS allocates inodes dynamically and Btrfs shares metadata space, so
their reported inode totals are estimates and can even change — still worth
reporting, just don't treat the total as a hard mkfs-time constant everywhere.)

---

## Important Files

- **`/proc/self/mountinfo`** — the authoritative, unambiguous mount table for
  the reading process's namespace. **Primary source for SysKit.** Fields:
  mount id, parent id, `major:minor`, root, mount point, per-mount options,
  optional propagation fields, ` - `, fs type, source, super options. Paths
  are octal-escaped.
- **`/proc/mounts`** — kernel mount table, fstab-style, `source mountpoint
  fstype options dump pass`. Fallback source. No device numbers; ambiguous.
- **`/proc/self/mounts`** — same as `/proc/mounts` for the current process;
  what `/etc/mtab` points to today.
- **`/etc/mtab`** — legacy; now a symlink to `/proc/self/mounts`. Do not treat
  as authoritative.
- **`/proc/filesystems`** — filesystem types the kernel currently supports; a
  leading `nodev` marks types with no backing device (pseudo). Useful as a
  reference for classifying pseudo filesystems.
- **`/etc/fstab`** — *desired* mounts at boot, not the live state. Not a
  collector source; useful only for comparing intent vs reality.
- **`/sys/fs/`** — per-filesystem-type sysfs trees (`/sys/fs/ext4/<dev>/`,
  `/sys/fs/xfs/<dev>/`) with tunables and stats. Out of scope for the core
  report but where you'd go for fs-specific detail later.

Note there is no file that hands you usage numbers — usage comes from the
`statfs` syscall on each mount point, not from reading a file.

---

## Useful Commands

Use these to hand-verify what the collector produces. SysKit itself must **not**
shell out to them (collectors read kernel interfaces directly — see
`specs/collectors.md`); they are for your own cross-checking.

- **`df -h`** — human-readable **byte** usage per mount. Note it reports
  available space using `f_bavail` (the user-visible free), so your numbers
  should match it, not `f_bfree`.
- **`df -i`** — **inode** usage per mount (total / used / free / IUse%). This
  is your ground truth for the inode columns and for demonstrating exhaustion.
- **`findmnt`** — pretty tree of mounts straight from mountinfo; great for
  seeing parent/child, source, fstype, and options. `findmnt --real` hides
  pseudo filesystems (compare to your `--show-pseudo` default). `findmnt -o
  TARGET,SOURCE,FSTYPE,OPTIONS` mirrors the collector's columns.
- **`mount`** (no args) — prints the current mount table with options; quick
  sanity check, but it's the ambiguous view — prefer `findmnt`.
- **`cat /proc/self/mountinfo`** — see the raw format you're parsing,
  including ` - ` separators and `\040` escapes if any path has a space.
- **`stat -f <path>`** — shows `statfs` results (block size, total/free blocks,
  total/free inodes) for the filesystem containing `<path>` — the exact fields
  the collector uses for usage.

To reproduce inode exhaustion safely in a throwaway VM: make a small
filesystem, then create files until `df -i` hits 100% and watch `touch` fail
with "No space left on device" while `df -h` still shows free bytes.

---

## References

- man 5 proc — `/proc/self/mountinfo` and `/proc/mounts` format:
  https://man7.org/linux/man-pages/man5/proc.5.html
- Kernel doc, `filesystems/proc.rst` (mountinfo fields & propagation):
  https://www.kernel.org/doc/html/latest/filesystems/proc.html
- Shared subtrees / mount propagation (`shared:`, `master:` optional fields):
  https://www.kernel.org/doc/html/latest/filesystems/sharedsubtree.html
- man 2 statfs — `f_blocks`, `f_bfree`, `f_bavail`, `f_files`, `f_ffree`:
  https://man7.org/linux/man-pages/man2/statfs.2.html
- man 3 statvfs — POSIX filesystem statistics:
  https://man7.org/linux/man-pages/man3/statvfs.3.html
- man 8 findmnt — https://man7.org/linux/man-pages/man8/findmnt.8.html
- man 1 df — https://man7.org/linux/man-pages/man1/df.1.html
- man 5 fstab — https://man7.org/linux/man-pages/man5/fstab.5.html
- mount namespaces, man 7 mount_namespaces:
  https://man7.org/linux/man-pages/man7/mount_namespaces.7.html

---

## Practical Lab

Select `/` and one virtual or overlay mount. Parse their mountinfo rows by hand,
decode escaped fields, and compare statfs byte/inode results with `syskit
filesystem`, `df`, and `findmnt`. Explain why the two mounts require different
operational interpretation.

## Failure-Mode Matrix

| Case | Correct behavior |
|---|---|
| Mount disappears before statfs | Skip or partial; continue other mounts |
| total blocks or inodes is zero | Avoid division by zero; percentage unavailable |
| spaces in mount point | Decode mountinfo octal escapes |
| bind or overlay mount | Preserve source/type/options; avoid false physical mapping |
| permission denied below mount | Do not infer that mount statistics are absent |
| namespace differs | State observer-relative scope, not global truth |

## Checkpoint

Demonstrate why bytes-full, inodes-full, device-full, and slow I/O are four
different diagnoses and identify the sources needed for each.
