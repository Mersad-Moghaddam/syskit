# Disk — Learning Notes

> Study notes on block devices, partitions, and the disk I/O subsystem — the
> layer *below* filesystems. This file leads on **block devices, partitions, and
> `/proc/diskstats`**. Its sibling `filesystem.md` leads on **mounts, inodes, and
> capacity** (mountinfo, `statfs`). Read both before implementing the disk
> collector; together they cover the whole "Disk And Filesystem" checklist.

---

## Concepts

Before touching a single file, get the layers straight. A newcomer to storage
almost always conflates these four things, and the disk collector's whole job is
to keep them distinct.

**Block device.** A kernel abstraction for storage you can address in
fixed-size blocks and seek around in (as opposed to a *character* device like a
serial port, which is a byte stream). `/dev/sda`, `/dev/nvme0n1`, `/dev/vda` are
whole-disk block devices. The kernel exposes each one under `/sys/block/`. A
block device has a **size** and physical/queue properties (rotational or not,
scheduler, sector size) but *no concept of files*.

**Partition.** A contiguous slice of a block device, described by a partition
table (MBR or GPT) written near the start of the disk. `/dev/sda1`,
`/dev/nvme0n1p2` are partitions. A partition is itself a block device — it too
has a size and shows up in sysfs, nested under its parent
(`/sys/block/sda/sda1/`). Partitions carve capacity; they still know nothing
about files.

**Filesystem.** The on-disk format (ext4, xfs, btrfs, vfat) written *into* a
partition (or a whole device, or a logical volume) that turns raw blocks into
files, directories, and inodes. This is where "how much space is used" lives —
and it is **not** a property of the block device. `filesystem.md` owns this.

**Mount.** The act of attaching a filesystem into the directory tree at a mount
point (`/`, `/home`, `/boot`). A filesystem is only reachable through a mount.
Mounts are also owned by `filesystem.md` (`/proc/self/mountinfo`).

So the stack, bottom to top, is:

```
block device (/dev/nvme0n1)
  └─ partition (/dev/nvme0n1p2)        ← a block device too
       └─ filesystem (ext4)            ← usage/inodes live here
            └─ mount point (/)         ← where you can reach it
```

Why this matters for SysKit: the spec's example output puts `DEVICE`, `SIZE`,
`TYPE`, and `MOUNT` on one row. That row is *stitched together* from three
different kernel sources — sysfs for device size, `statfs` on the mount for
usage, mountinfo for the device→mount mapping. If you think of them as one
thing you will fetch usage from the wrong place, or report a device's raw size
as "used space." Keep the sources separate in code and join them deliberately.

Other concepts you'll meet (mostly *future extensions*, but know the words):

- **I/O scheduler** (`none`, `mq-deadline`, `bfq`) — orders requests to the
  device; readable at `/sys/block/<dev>/queue/scheduler`.
- **Device mapper / LVM** — virtual block devices layered over real ones
  (`/dev/mapper/*`, `dm-0`). They appear in `/sys/block/` like any device.
- **RAID (md)** — `md0` etc., also a block device. Out of core scope for now.
- **IOPS vs throughput** — operations/sec vs bytes/sec. Both are *rates*, and
  both come from the cumulative counters discussed below, never from a single
  reading.

---

## Linux Internals

### `/sys/block/` — the block device inventory

`/sys/block/` has one entry per whole block device. Each entry is a directory:

```
/sys/block/sda/
├── size                 # device size in 512-byte SECTORS (not bytes!)
├── stat                 # per-device I/O counters (same fields as diskstats)
├── queue/
│   ├── rotational       # 1 = spinning disk, 0 = SSD/NVMe
│   ├── scheduler        # active I/O scheduler in [brackets]
│   ├── hw_sector_size   # physical sector size, usually 512 or 4096
│   └── logical_block_size
├── device/
│   └── model            # e.g. "Samsung SSD 970"
└── sda1/                # a partition — nested under its parent disk
    ├── size             # partition size, also in 512-byte sectors
    └── start            # starting sector offset
```

Two things to burn into memory:

1. **`size` is in 512-byte sectors, always — regardless of the drive's real
   physical sector size.** This is a fixed kernel ABI unit, not the hardware
   sector size in `queue/hw_sector_size`. The classic newcomer bug is to read
   `size`, see a big number, and multiply by 1024 (assuming KiB) or by
   `hw_sector_size`. Both are wrong. The correct conversion is:

   ```
   bytes = size_field * 512
   ```

   A 500 GB SSD reports `size` ≈ 976773168. Times 512 ≈ 500.1 GB. Times 1024
   would give ~1 TB — an off-by-2x bug that looks plausible enough to ship.

2. **Partitions are nested, not top-level.** To list partitions of `sda` you
   walk `/sys/block/sda/` and pick subdirectories that themselves contain a
   `partition` file (or match the `sdaN` naming). Don't expect `sda1` to appear
   directly under `/sys/block/`.

Distinguish disks from partitions structurally: a directory directly under
`/sys/block/` is a whole device; a directory that has a `partition` file inside
it is a partition. Don't parse names to guess — `nvme0n1` (disk) vs `nvme0n1p1`
(partition) vs `dm-0` vs `mmcblk0` all follow different naming rules.

### `/proc/diskstats` — the I/O counters

This is the heart of the disk collector's I/O half. One line per block device
*and* per partition. Each line has (at least) **14 fields**; modern kernels
append more (flush stats), so parse defensively — read the fields you know by
index and tolerate extras. The classic 14:

| # | Field                    | Meaning                                        |
|---|--------------------------|------------------------------------------------|
| 1 | major                    | device major number                            |
| 2 | minor                    | device minor number                            |
| 3 | device name              | e.g. `sda`, `nvme0n1p2`                         |
| 4 | reads completed          | successfully completed reads                   |
| 5 | reads merged             | adjacent reads merged before hitting device    |
| 6 | sectors read             | **512-byte sectors** read                      |
| 7 | ms spent reading         | total ms all reads waited                      |
| 8 | writes completed         | successfully completed writes                  |
| 9 | writes merged            | adjacent writes merged                         |
| 10| sectors written          | **512-byte sectors** written                   |
| 11| ms spent writing         | total ms all writes waited                     |
| 12| I/Os currently in flight | in-progress request count (the only *gauge*)   |
| 13| ms spent doing I/O       | wall-clock ms the device had I/O in flight     |
| 14| weighted ms doing I/O    | field 13 weighted by queue depth               |

Read/write **bytes** come from the sector fields, and here too **sector = 512
bytes** by convention in diskstats — same trap as `/sys/block/size`, same fix:

```
bytes_read = sectors_read * 512
```

This 512 is a fixed diskstats convention. It does *not* change if the drive
advertises 4096-byte physical sectors. Do not reach for `hw_sector_size` here.

#### Why these counters are cumulative (checklist item)

**Every count in `/proc/diskstats` is a monotonically increasing total
accumulated since boot.** Field 4 is not "reads happening now"; it is "every
read this device has completed since the machine started." (The lone exception
is field 12, in-flight I/Os, which is an instantaneous gauge.)

This is exactly the same design as CPU jiffies in `/proc/stat` and network
byte counters in `/proc/net/dev`, and it exists for the same reason: the kernel
cannot know what time window *you* care about, so it exposes a running total and
lets userspace pick the interval. A running total is cheap, lock-light, and
loses no information.

The consequence for SysKit is the rule the spec states as acceptance criteria:

> **I/O rates are derived only from two samples.**

To report "5 MB/s write throughput" or "300 IOPS" you must:

1. Read `/proc/diskstats` at time *t₀*, remember the counters.
2. Wait a known interval, read again at *t₁*.
3. Rate = (counter₁ − counter₀) / (t₁ − t₀).

A single reading of `sectors_written × 512` is **total bytes since boot**, not
throughput. Reporting that raw counter as a "rate" is simply wrong — on a host
up for 200 days it will read as an absurd, meaningless number. This mirrors the
CPU collector exactly, which is why the architecture spec says the *service
layer* owns the two-snapshot delta and the collector just returns a snapshot.
The disk collector's job is to return honest cumulative counters plus the
timestamp; watch mode / the service does the subtraction.

#### Edge cases the collector must tolerate

- **Counter rollover.** These are unsigned kernel counters of bounded width.
  On a long-lived, busy host a counter can wrap past its maximum back toward
  zero. If you compute `new − old` naively across a wrap you get a huge negative
  (or, with unsigned math, a huge positive) bogus rate. The delta code must
  detect `new < old` for a monotonic field and treat it as a rollover or reset:
  clamp to zero, or skip that sample rather than emit a garbage spike. Rollover
  is rare on 64-bit fields but real on longer-lived counters and older kernels;
  handle it rather than assume monotonicity always holds.

- **Device removed mid-sample.** A removable disk (USB, SD card) or a
  hot-unplugged NVMe can vanish *between* your two reads. Then a device present
  at *t₀* is absent at *t₁* (or a sysfs read fails with ENOENT). Never assume the
  device set is identical across snapshots. Join the two samples by device name
  and simply drop devices that appear in only one — a missing device is not an
  error, it's a race the spec explicitly calls out ("Removable devices disappear
  during collection"). The inverse (a device appearing at *t₁* only) has no *t₀*
  baseline, so it gets no rate until the next interval.

- **Partitions and whole disks both appear.** diskstats lists `sda` *and*
  `sda1`, `sda2`. Summing all lines double-counts, because a partition's I/O is
  already included in its parent disk's totals. Decide deliberately whether you
  report per-partition, per-disk, or both — don't accidentally aggregate.

- **Zero-activity / pseudo devices.** loop devices, ram disks, and idle disks
  show all-zero or static counters. That's valid data, not a failure.

### `/proc/partitions`

An older, simpler view: `major minor #blocks name`, one line per device and
partition. Here `#blocks` is in **1024-byte blocks** (note: *different* unit
from `/sys/block/size`'s 512-byte sectors — another reason to be unit-paranoid).
It predates sysfs and carries less detail. Prefer `/sys/block/` for the block
device inventory (richer: rotational, model, scheduler); `/proc/partitions` is a
useful cross-check and fixture, not the primary source.

---

## Important Files

- `/proc/diskstats` — per-device and per-partition cumulative I/O counters (14+
  fields). **Primary source for the I/O half of the collector.** Sectors are
  512 bytes.
- `/sys/block/` — one directory per whole block device; partitions nested
  inside. Primary source for the device inventory.
- `/sys/block/<dev>/size` — device size in **512-byte sectors** (`× 512` for
  bytes).
- `/sys/block/<dev>/queue/rotational` — `1` spinning, `0` SSD/NVMe.
- `/sys/block/<dev>/queue/scheduler` — active I/O scheduler (in `[brackets]`).
- `/sys/block/<dev>/queue/hw_sector_size` — physical sector size (do **not** use
  this to convert `size` or diskstats sectors — those are always 512).
- `/sys/block/<dev>/stat` — same counter fields as one diskstats line, per
  device.
- `/sys/block/<dev>/device/model` — human-readable model string.
- `/sys/block/<dev>/<part>/{size,start,partition}` — per-partition size/offset;
  presence of `partition` marks a directory as a partition.
- `/proc/partitions` — older `major minor #blocks name` table; `#blocks` is in
  **1024-byte** units. Secondary / cross-check.

---

## Useful Commands

Use these to *verify* what your collector parses — never to source data (the
architecture spec forbids shelling out to `df`/`lsblk`/`iostat` for core data).

- `lsblk` — tree of block devices, partitions, sizes, and mount points. The
  best visual of the device→partition→mount hierarchy this file is about.
  Try `lsblk -o NAME,SIZE,TYPE,ROTA,MOUNTPOINT` and `lsblk -b` (bytes) to
  sanity-check your `size × 512` math against a known-good number.
- `cat /proc/diskstats` — see the raw counters your collector will parse.
  Run it twice a second apart and watch fields grow to *feel* the cumulative
  nature: `cat /proc/diskstats; sleep 1; cat /proc/diskstats`.
- `iostat -dx 1` — per-device extended I/O stats refreshed every second. This is
  the reference implementation of "two samples → rate"; its first block is
  since-boot averages, every block after is a true interval rate. Compare its
  numbers to your own delta computation.
- `cat /sys/block/sda/size` then multiply by 512 — confirm it matches `lsblk -b`
  for that device. This is the single fastest way to catch the ×1024 bug.
- `cat /sys/block/sda/queue/rotational` — confirm SSD vs HDD classification.

---

## References

- Linux kernel: `Documentation/admin-guide/iostats.rst` — authoritative field
  list for `/proc/diskstats` and `/sys/block/<dev>/stat`.
  https://www.kernel.org/doc/html/latest/admin-guide/iostats.html
- Linux kernel: `Documentation/ABI/testing/sysfs-block` — sysfs block device
  attributes (`size`, `queue/*`, partition layout).
  https://www.kernel.org/doc/html/latest/block/index.html
- `man 5 proc` — `/proc/diskstats` and `/proc/partitions` descriptions.
  https://man7.org/linux/man-pages/man5/proc.5.html
- `man 8 lsblk` — https://man7.org/linux/man-pages/man8/lsblk.8.html
- `man 1 iostat` — https://man7.org/linux/man-pages/man1/iostat.1.html
- Kernel block layer overview —
  https://www.kernel.org/doc/html/latest/block/index.html

---

## Personal Notes
