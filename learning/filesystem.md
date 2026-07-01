# Filesystem — Learning Notes

> Study notes on filesystem internals, VFS, and related Linux interfaces.

---

## Concepts

<!-- Key concepts related to filesystems -->
<!-- - Virtual Filesystem (VFS) layer -->
<!-- - Inodes, dentries, and superblocks -->
<!-- - Filesystem types: ext4, XFS, Btrfs, ZFS, tmpfs, proc, sysfs -->
<!-- - Journaling and write-ahead logging -->
<!-- - File permissions and ACLs -->
<!-- - Hard links vs symbolic links -->
<!-- - Filesystem quotas -->
<!-- - FUSE (Filesystem in Userspace) -->

---

## Linux Internals

<!-- How the Linux kernel manages and exposes filesystem information -->
<!-- - VFS architecture and abstraction layer -->
<!-- - /proc/filesystems — supported filesystem types -->
<!-- - /proc/mounts and /proc/self/mountinfo — current mounts -->
<!-- - statfs/statvfs system calls -->
<!-- - /sys/fs/ hierarchy -->
<!-- - Extended attributes (xattr) -->
<!-- - inotify and fanotify for filesystem events -->

---

## Important Files

<!-- Key files and paths for filesystem data collection -->
<!-- - /proc/filesystems -->
<!-- - /proc/mounts -->
<!-- - /proc/self/mountinfo -->
<!-- - /etc/fstab -->
<!-- - /sys/fs/ -->
<!-- - /sys/fs/ext4/ -->
<!-- - /sys/fs/xfs/ -->

---

## Useful Commands

<!-- Commands useful for understanding and verifying filesystem data -->
<!-- - stat -->
<!-- - df -i (inode usage) -->
<!-- - findmnt -->
<!-- - tune2fs -l -->
<!-- - xfs_info -->
<!-- - dumpe2fs -->

---

## References

<!-- Links to documentation, kernel source, man pages, and articles -->

---

## Personal Notes

<!-- Personal observations, insights, and things learned during research -->
