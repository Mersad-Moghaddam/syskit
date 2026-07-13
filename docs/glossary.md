# Glossary

> Terms used throughout SysKit documentation and specifications.

| Term | Meaning |
|---|---|
| Collector | Component that reads one Linux data domain and returns normalized structured data. |
| Service | Component that combines collector data, applies domain rules, and prepares results for commands. |
| Renderer | Component that turns structured results into table, JSON, YAML, or interactive terminal output. |
| procfs | Virtual filesystem mounted at `/proc` that exposes kernel and process information. |
| sysfs | Virtual filesystem mounted at `/sys` that exposes devices, drivers, topology, and kernel objects. |
| Netlink | Linux socket interface used for communication between userspace and the kernel. |
| cgroup | Linux control group used to account and limit process resource usage. |
| PSI | Pressure Stall Information, a Linux interface for CPU, memory, and I/O pressure signals. |
| Golden file | Expected output fixture used to detect user-visible output changes in tests. |
| TUI | Terminal user interface used by dashboard and live monitoring views. |
| XDG | Linux desktop directory convention used for configuration and state file locations. |
| Read-only inspection | SysKit's core policy that commands inspect system state without modifying it. |
| Output contract | Stable structure and semantics of command output, especially JSON and YAML. |
| ADR | Architecture Decision Record that captures an important decision and its consequences. |
