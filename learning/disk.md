# Disk — Learning Notes

> Study notes on disk management, I/O subsystems, and related Linux internals.

---

## Concepts

<!-- Key concepts related to disk and storage -->
<!-- - Block devices vs character devices -->
<!-- - Partitions, partition tables (MBR, GPT) -->
<!-- - Filesystems and mount points -->
<!-- - I/O schedulers (mq-deadline, bfq, none) -->
<!-- - Disk I/O: reads, writes, IOPS, throughput -->
<!-- - RAID levels and device mapper -->
<!-- - LVM (Logical Volume Manager) -->
<!-- - NVMe vs SATA vs SCSI -->

---

## Linux Internals

<!-- How the Linux kernel manages and exposes disk information -->
<!-- - /proc/diskstats structure and fields -->
<!-- - /proc/partitions -->
<!-- - /proc/mounts and /proc/self/mountinfo -->
<!-- - /sys/block/ hierarchy -->
<!-- - /sys/block/[dev]/queue/ for I/O scheduler info -->
<!-- - /sys/block/[dev]/stat for device statistics -->
<!-- - Device mapper and /dev/mapper/ -->
<!-- - udev and device discovery -->

---

## Important Files

<!-- Key files and paths for disk data collection -->
<!-- - /proc/diskstats -->
<!-- - /proc/partitions -->
<!-- - /proc/mounts -->
<!-- - /proc/self/mountinfo -->
<!-- - /sys/block/ -->
<!-- - /sys/block/[dev]/stat -->
<!-- - /sys/block/[dev]/queue/scheduler -->
<!-- - /sys/block/[dev]/device/model -->

---

## Useful Commands

<!-- Commands useful for understanding and verifying disk data -->
<!-- - lsblk -->
<!-- - fdisk -l -->
<!-- - df -h -->
<!-- - mount -->
<!-- - iostat -->
<!-- - blkid -->
<!-- - smartctl -->

---

## References

<!-- Links to documentation, kernel source, man pages, and articles -->

---

## Personal Notes

<!-- Personal observations, insights, and things learned during research -->
