# Network — Learning Notes

> Study notes on networking, Linux network interfaces, and related internals.
> Written for the implementer of SysKit's `network` and `ports` collectors.

---

## Concepts

Before touching a single file, separate four concerns that newcomers habitually
lump together as "the network." They live in different places, change on
different schedules, and belong to different SysKit subcommands.

- **Interfaces.** A network device — `eth0`, `wlan0`, `lo`, a bridge, a bond, a
  veth pair. It has a MAC address, an MTU, an operational state (`up`/`down`),
  and a set of cumulative traffic counters (bytes/packets/errors/drops). This is
  hardware/link-level identity and volume. → `syskit network interfaces`.
- **Addresses.** IPv4/IPv6 addresses *assigned to* an interface (`192.168.1.20/24`).
  One interface can carry many addresses; an address is not the interface. This
  is layer-3 identity.
- **Routing.** The kernel's decision table for *where a packet goes next* — the
  destination prefix, the gateway, and which interface to send it out of. The
  default route (`0.0.0.0/0`) is your gateway to everything not on a local
  subnet. Routing is a *global* kernel table, not a property of any one
  interface. → `syskit network routes`.
- **DNS.** Name-to-address resolution config — which resolver to ask. It is pure
  userspace configuration (`/etc/resolv.conf`), has nothing to do with the kernel
  routing table, and is often a symlink managed by `systemd-resolved`. →
  `syskit network dns`.
- **Sockets.** A live endpoint of communication owned by a process — a TCP
  connection or a listening port with a local/remote address, a state, and a
  kernel socket inode. Sockets are what `syskit ports` reports. A *listening
  socket* is a service waiting for clients; an *established socket* is an active
  connection.

Keep this mental model sharp: an interface is a door, an address is the door's
number, a route is the map telling you which door to use, DNS is the phone book,
and a socket is an actual conversation happening through a door. Four files,
four subcommands, four concerns.

**Cumulative counters, again.** Interface RX/TX counters are the same species of
number as CPU jiffies and disk I/O sectors: monotonic totals since boot (or
since the interface last came up). A single read tells you a total, never a
rate. Rate = (sample2 − sample1) / (t2 − t1). Per the collector architecture,
the *collector* returns the raw counter snapshot and the *service layer* owns the
two-sample rate math. Do not compute bandwidth inside the collector.

---

## Linux Internals

### Why Netlink, not `ss`/`netstat` — and the two-path story

ADR-003 is binding here: SysKit reads native kernel interfaces and does **not**
shell out to `ss`, `ip`, `netstat`, or `lsof`. Human-readable tool output is not
a stable contract — columns, units, and locale shift across versions and
distros, and every invocation forks a process just to re-serialise data we then
re-parse. For networking specifically, ADR-003 names **Netlink** as the
authoritative source.

There are effectively two native paths for socket data, and you should
understand both because SysKit's journey goes from one to the other:

1. **`/proc/net/{tcp,tcp6,udp,udp6}` — the procfs path.** A plain text table,
   trivial to read and parse with a buffered reader. This is the pragmatic
   starting point and what the field examples below decode. Its limitations are
   real:
   - Addresses and ports are **hex-encoded** (and the address is byte-swapped —
     see below), so naive readers decode them wrong.
   - The whole table is materialised as text on every read. On a host with tens
     of thousands of sockets this is slow and can **truncate**/scale poorly.
   - It exposes limited per-socket detail compared to the diag API.
2. **Netlink `sock_diag` (`NETLINK_INET_DIAG`, `AF_NETLINK`) — the robust
   production path.** You send a `SOCK_DIAG_BY_FAMILY` request and the kernel
   streams back structured `inet_diag_msg` records: addresses as raw bytes (no
   hex parsing), state, inode, UID, and rich attributes — the same source `ss`
   itself uses. It scales to large socket tables and is the intended endpoint for
   SysKit's networking layer (via `golang.org/x/sys/unix`).

The migration is deliberate: **start by parsing `/proc/net/tcp` to learn the
socket model concretely, then move to `sock_diag` for the production
collector.** The domain records the collector returns should be identical either
way — only the adapter underneath changes. Interface enumeration and routing
follow the same principle: `RTM_GETLINK`/`RTM_GETADDR`/`RTM_GETROUTE` over
RTNETLINK are the authoritative sources, with `/proc/net/dev` and
`/sys/class/net/*/statistics/` available as easy counter sources.

### `/proc/net/` layout

- `/proc/net/dev` — one line per interface: `iface: rx_bytes rx_packets rx_errs
  rx_drop ... tx_bytes tx_packets tx_errs tx_drop ...`. The 16 fields after the
  colon are 8 RX then 8 TX.
- `/proc/net/tcp`, `/proc/net/tcp6`, `/proc/net/udp`, `/proc/net/udp6` — the
  socket tables (see decoding below).
- `/proc/net/route` — the IPv4 routing table, also hex-encoded and little-endian.
  A destination of `00000000` with a gateway is the default route.
- `/proc/net/unix` — Unix domain sockets (no IP address, path-named).

### `/sys/class/net/<iface>/`

Per-interface sysfs tree. `statistics/rx_bytes`, `statistics/tx_bytes`,
`statistics/rx_packets`, `statistics/rx_errors`, `statistics/rx_dropped`, and the
`tx_*` equivalents give the same counters as `/proc/net/dev` but one value per
file — cleaner to read for a single interface. `operstate`, `mtu`, `address`
(MAC), and `flags` give interface metadata.

### Counter rollover — the classic bug

The counters are unsigned and finite. On modern 64-bit kernel counters wrap is
essentially never a concern, but **32-bit counters** (older kernels, some
drivers, some virtualised NICs) wrap at 2^32 bytes (~4 GiB). When they wrap,
`sample2 < sample1`, and a naive `sample2 - sample1` on signed arithmetic goes
**hugely negative** — producing an absurd negative or (if unsigned-wrapped) an
enormous positive bandwidth spike. This is *the* classic interface-counter bug.

The collector/service must detect it: if the current reading is less than the
previous one, treat it as a wrap (or an interface reset — see edge cases) rather
than emitting a garbage delta. The defensive rule: **a decrease means "counter
reset or wrapped," so skip that interval rather than report a negative rate.**
Note that an interface bouncing (`down`/`up`) also zeroes the counters, which
looks exactly like a wrap — both are handled by the same "counter went backwards,
don't compute a delta this cycle" guard.

### Socket inode → process mapping

`/proc/net/tcp` gives you a socket's **inode** but not the owning process. The
kernel exposes the link the other way around: every open file descriptor of a
process is a symlink under `/proc/[pid]/fd/`, and a socket fd's symlink target is
the string `socket:[INODE]`. So to answer "which PID owns port 8080?" you:

1. Parse the socket tables to get `(local addr:port, state, inode)`.
2. Walk every `/proc/[pid]/fd/*`, `readlink` each entry, and match targets of the
   form `socket:[<inode>]` to build an `inode → pid` map.
3. Join: the listening socket's inode → the owning PID → `/proc/[pid]/comm` for
   the command name.

This is exactly how `ss -p` and `lsof` associate sockets with processes, done
natively. **Permissions matter:** you can only `readlink` the fds of processes
you own unless you have `CAP_SYS_PTRACE` / run as root. For other users' sockets
the inode simply won't be found in your map, and the collector must report the
socket with **PID unknown** rather than failing — partial data, per the collector
error-classification rules. Also expect races: a process can exit between reading
the socket table and scanning its fds (handle `ENOENT`/`ESRCH` on
`/proc/[pid]/*` as "process gone," not fatal). After a `fork`, several PIDs may
share one inode.

---

## Important Files

- `/proc/net/dev` — per-interface cumulative RX/TX counters (bytes, packets,
  errors, drops). Rate source; watch for rollover.
- `/proc/net/tcp` — IPv4 TCP socket table: hex local/remote addr:port, hex state,
  socket inode, owning UID. The learning path for `ports`.
- `/proc/net/tcp6` — IPv6 TCP sockets (32-hex-char addresses).
- `/proc/net/udp`, `/proc/net/udp6` — UDP sockets. UDP is connectionless: its
  "state" field is not a TCP-style connection state, so do not render UDP entries
  with TCP state names — represent them as UDP/no-connection-state.
- `/proc/net/route` — IPv4 routing table (hex, little-endian). Default route =
  destination `00000000`.
- `/sys/class/net/<iface>/statistics/{rx_bytes,tx_bytes,rx_packets,rx_errors,rx_dropped,tx_bytes,tx_packets,...}`
  — one counter per file; alternative to `/proc/net/dev` for a single interface.
- `/sys/class/net/<iface>/{operstate,mtu,address,flags}` — interface state, MTU,
  MAC, flags.
- `/proc/[pid]/fd/` — a process's open file descriptors; socket fds symlink to
  `socket:[inode]`. The join key for inode → PID mapping.
- `/proc/[pid]/comm` — the process command name for the mapped PID.
- `/etc/resolv.conf` — DNS resolver config (`nameserver`, `search`, `options`).
  Often a symlink to a `systemd-resolved`-managed file; parse comments (`#`/`;`)
  and multiple `nameserver` lines.

### Decoding a `/proc/net/tcp` line

Fields of interest per row: `local_address`, `rem_address`, `st` (state), and
`inode`. Addresses look like `0100007F:1F90`.

- Split on `:` → address `0100007F`, port `1F90`.
- **Port is big-endian hex:** `0x1F90` = **8080**. Straightforward.
- **Address is little-endian hex (byte-swapped):** `0100007F` read as bytes is
  `01 00 00 7F`; reverse to `7F 00 00 01` = `127.0.0.1`.
- Result: `0100007F:1F90` = **127.0.0.1:8080**. This byte-swap is the mistake
  every newcomer makes — write a test that asserts this exact example.
- IPv6 (`/proc/net/tcp6`) uses 32 hex chars, byte-swapped per 4-byte word — more
  fiddly, which is another argument for `sock_diag`, where the kernel hands you
  raw address bytes directly.

### TCP socket states (the `st` hex field)

| Hex | State | Meaning |
|-----|-------|---------|
| 01 | ESTABLISHED | Active, open connection |
| 02 | SYN_SENT | Connecting: SYN sent, awaiting reply |
| 03 | SYN_RECV | Received a SYN, handshake in progress |
| 04 | FIN_WAIT1 | Local close initiated |
| 05 | FIN_WAIT2 | Local close, awaiting remote FIN |
| 06 | TIME_WAIT | Closed, waiting out lingering packets |
| 07 | CLOSE | Socket unused/closed |
| 08 | CLOSE_WAIT | Remote closed; local app must still close |
| 09 | LAST_ACK | Closing, awaiting final ACK |
| 0A | LISTEN | **A listening port** — a service accepting connections |
| 0B | CLOSING | Both sides closing simultaneously |

For `syskit ports`, **`0A` LISTEN is the state that identifies a bound,
listening service** — it is what `--listening` filters to. A pile of `06`
TIME_WAIT entries is normal after churn (they age out); a growing pile of `08`
CLOSE_WAIT usually means an application is leaking connections by not calling
`close()`.

---

## Useful Commands

These are for *learning and verifying* your parser against ground truth — SysKit
itself must not execute them (ADR-003). Run them by hand and diff against your
collector's output.

- `cat /proc/net/tcp` — see the raw hex socket table you are parsing; decode a
  line by hand and confirm your byte-swap.
- `ss -tlnp` — listening (`-l`) TCP (`-t`) sockets, numeric (`-n`), with process
  info (`-p`). The reference output your `syskit ports --listening` should match.
- `ss -tan` — all TCP sockets with states; cross-check your state decoding.
- `ip -s link` — per-interface stats (`-s`) including RX/TX bytes, packets,
  errors, drops. Ground truth for `/proc/net/dev` counters.
- `ip route` — the routing table with the default gateway; verify your
  `/proc/net/route` / RTNETLINK route parsing.
- `ip addr` — interfaces with their assigned addresses; confirms the
  interface-vs-address distinction.
- `cat /proc/net/dev` — raw interface counter table.
- `cat /etc/resolv.conf` — DNS resolver config; check symlink target with
  `readlink -f /etc/resolv.conf`.
- `ls -l /proc/$$/fd` — see the `socket:[inode]` symlinks for your own shell's
  sockets; the mechanism behind inode → PID mapping.

---

## References

- Linux kernel `Documentation/networking/` — https://www.kernel.org/doc/html/latest/networking/index.html
- `proc(5)` — `/proc/net/*` file formats: https://man7.org/linux/man-pages/man5/proc.5.html
- `sock_diag(7)` — the Netlink socket-diagnostics API (`NETLINK_INET_DIAG`): https://man7.org/linux/man-pages/man7/sock_diag.7.html
- `netlink(7)` — Netlink socket protocol: https://man7.org/linux/man-pages/man7/netlink.7.html
- `rtnetlink(7)` — interface/address/route messages (`RTM_GETLINK`, `RTM_GETADDR`, `RTM_GETROUTE`): https://man7.org/linux/man-pages/man7/rtnetlink.7.html
- `resolv.conf(5)` — DNS resolver configuration: https://man7.org/linux/man-pages/man5/resolv.conf.5.html
- `tcp(7)` — TCP states and semantics: https://man7.org/linux/man-pages/man7/tcp.7.html
- ADR-003 — Read native kernel interfaces instead of shelling out: `../decisions/003-native-apis-over-shell.md`

---

## Personal Notes
