# Sprint 05 — Process Tree & Networking

**Dates:** TBD → +2 weeks · **Milestone:** v0.2 (Processes & Networking) · **Committed points:** 29 ⚠️

## Sprint goal

Ship `syskit process tree` and bring networking online: Netlink platform integration plus the `syskit network` command.

## Capacity

- **⚠️ Overflow-risk sprint (29 pts > nominal ~23).** Scheduled tight on purpose so the risk is visible. NET-01 (Netlink, 13) is high-uncertainty.

## Committed backlog

| ID | Story | Pts | Status |
|---|---|---|---|
| PRC-02 | `syskit process tree`. | 8 | Committed |
| NET-01 | Netlink platform integration. | 13 | Committed |
| NET-02 | `syskit network` interfaces/connections/routing. | 8 | Committed |

## Task breakdowns

**PRC-02** — spec `../../specs/features/process.md`
- [ ] service: build parent→child tree from PID/PPID; handle disappearing processes mid-build.
- [ ] render: tree/indented output + JSON nesting (+ golden).
- [ ] tests: orphan/re-parent cases; docs + CHANGELOG.

**NET-01** — spec `../../specs/features/network.md`, ADR `../../decisions/003-native-apis-over-shell.md`
- [ ] platform: open Netlink socket; send/receive; parse messages (no shelling to `ss`/`ip`).
- [ ] fixtures: captured Netlink payloads for parse unit tests.
- [ ] integration: real Netlink query on Linux CI; benchmark parse.

**NET-02** — spec `../../specs/features/network.md`
- [ ] collector/service: interface stats, active connections, routing table via NET-01.
- [ ] command + render (table + JSON, + golden); reuse FLT-01 filters; docs + CHANGELOG.

## Definition of Ready / Done

Standard gates. NET-01 entered *Ready* only after **SPK-NET** (Sprint 4 refinement) confirmed the approach.

## Risks this sprint

- **R-01 (Netlink) — HIGH.** Mitigation: if NET-01 slips, **defer NET-02 to Sprint 6** and pull PRT-01 forward; the v0.2 tag moves with the work. Do not compromise Netlink test quality to hit the date.

## Dependencies

Blocked by Sprint 4 (FLT-01 for NET-02). NET-01 → NET-02 (command needs the platform capability).

## Notes

This is the plan's first deliberately over-committed sprint. If burndown is flat by day 6, invoke the NET-02 deferral immediately — don't wait for day 10.
