# Process — Learning Notes

> Study notes on process management, Linux process interfaces, and related internals.

---

## Concepts

<!-- Key concepts related to process management -->
<!-- - Process vs thread vs task -->
<!-- - Process states (R, S, D, Z, T) -->
<!-- - Process lifecycle: fork, exec, exit, wait -->
<!-- - Process groups and sessions -->
<!-- - Signals and signal handling -->
<!-- - Namespaces and cgroups -->
<!-- - Scheduling policies and priorities -->
<!-- - File descriptors and resource limits -->

---

## Linux Internals

<!-- How the Linux kernel manages and exposes process information -->
<!-- - /proc/[pid]/ directory structure -->
<!-- - /proc/[pid]/stat and /proc/[pid]/status -->
<!-- - /proc/[pid]/cmdline and /proc/[pid]/comm -->
<!-- - /proc/[pid]/fd/ for open file descriptors -->
<!-- - /proc/[pid]/maps and /proc/[pid]/smaps -->
<!-- - /proc/[pid]/cgroup -->
<!-- - Task struct and kernel representation -->
<!-- - Process tree and parent-child relationships -->

---

## Important Files

<!-- Key files and paths for process data collection -->
<!-- - /proc/[pid]/stat -->
<!-- - /proc/[pid]/status -->
<!-- - /proc/[pid]/cmdline -->
<!-- - /proc/[pid]/comm -->
<!-- - /proc/[pid]/io -->
<!-- - /proc/[pid]/fd/ -->
<!-- - /proc/[pid]/maps -->
<!-- - /proc/[pid]/limits -->
<!-- - /proc/[pid]/cgroup -->
<!-- - /proc/[pid]/ns/ -->

---

## Useful Commands

<!-- Commands useful for understanding and verifying process data -->
<!-- - ps aux -->
<!-- - pstree -->
<!-- - top / htop -->
<!-- - strace -->
<!-- - lsof -->
<!-- - /proc/[pid]/ exploration -->

---

## References

<!-- Links to documentation, kernel source, man pages, and articles -->

---

## Personal Notes

<!-- Personal observations, insights, and things learned during research -->
