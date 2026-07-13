# Command Reference

> Implemented SysKit commands and their stable user-facing purpose.

All one-shot commands support `--format table`, `--format json`, and
`--format yaml`. Use structured output for automation; table output is intended
for people. SysKit is Linux-only and reads native kernel interfaces without
executing system utilities.

## Host and resources

| Command | Purpose |
|---|---|
| `system` | Host, operating-system, kernel, uptime, and load summary. |
| `cpu` | CPU topology, identity, and sampled utilization. |
| `memory` | Memory, swap, cache, and optional PSI pressure. |
| `disk` | Mounted storage and optional sampled device I/O. |
| `filesystem` | Mount, capacity, inode, type, and option details. |

## Processes and networking

| Command | Purpose |
|---|---|
| `process` | Process listing with `--filter`, `--sort`, `--limit`, and `--containers`. |
| `process tree` | Parent/child process tree. |
| `network` | Interface counters and metadata. |
| `network interfaces` | Focused interface and address view. |
| `network routes` | IPv4 routing-table view. |
| `network dns` | Resolver nameserver view. |
| `ports` | TCP, UDP, IPv6, and Unix socket views; supports protocol, port, PID, address, and state filters. |

Permission-restricted procfs scans are best-effort. Structured process and port
output indicates partial data rather than fabricating values.

## Live views

| Command | Purpose |
|---|---|
| `dashboard` | Interactive host summary; `Tab` switches panels and `q` exits. |
| `watch <command>` | Refreshes a table command until Ctrl-C. |
| `top` | Interactive process view; `c/m/n/p` change sort and `j/k` scroll. |

These commands require an interactive terminal. They refuse redirected output
instead of emitting terminal control sequences into a file or pipe.

## Containers

| Command | Purpose |
|---|---|
| `containers` | Cgroup-derived IDs, process counts, and available resource counters. |
| `containers inspect <id>` | Processes associated with one recognized cgroup container ID. |

Container runtime names and status are not inferred. Missing cgroup controller
files remain unavailable rather than being reported as zero. Structured output
sets `partial: true` when permissions prevent a complete process mapping.

## Plugins and diagnostics

| Command | Purpose |
|---|---|
| `plugins list` | Discover manifests from documented or explicit plugin directories. |
| `plugins inspect <name>` | View plugin API compatibility and requested permissions. |
| `plugins run <name>` | Explicitly execute one compatible plugin through the bounded JSON protocol. |
| `diagnostics` | Explainable CPU, memory, disk, filesystem, process, network, and port findings; `--category` selects one domain and `--severity` accepts `info`, `warning`, or `critical`. |

Diagnostics never silently converts missing optional signals to healthy values.
It emits informational unavailable findings with evidence and source paths, and
a category filter avoids collecting unrelated domains.

Plugin discovery never executes plugin code. Use `--plugin-dir` to inspect a
specific directory; world-writable directories are rejected.
Executables must stay inside their plugin directory, be regular executable
files, return exactly one JSON value, and complete before `--timeout`.

## Utilities

| Command | Purpose |
|---|---|
| `version` | Print the embedded SysKit version. |
| `completion bash` | Generate Bash completion source. |
| `completion fish` | Generate Fish completion source. |
| `completion powershell` | Generate PowerShell completion source. |
| `completion zsh` | Generate Zsh completion source. |

## Global flags

| Flag | Meaning |
|---|---|
| `--format` | Output format: `table`, `json`, or `yaml`. |
| `--config` | Explicit TOML configuration file. |
| `--color` | Table color: `auto`, `always`, or `never`; `NO_COLOR` disables it. |
| `--no-header` | Suppress headers in table output. |
| `--verbose`, `--debug`, `--quiet` | Control diagnostics written to stderr. |

Run `syskit <command> --help` for the complete current flag contract.
The machine-readable command/flag inventory and structured schemas are frozen in
the [v1 compatibility contract](compatibility.md).
