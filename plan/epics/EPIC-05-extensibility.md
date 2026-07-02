# EPIC-05 — Extensibility (Plugins)

> **Milestone:** v0.5 · **Sprints:** 12–13 · **Points:** 45
> An out-of-process plugin system that lets users add custom collectors safely — culminating in **v0.5.0**.

---

## Goal

Design and implement the plugin architecture per `../decisions/007-out-of-process-plugins.md`: a plugin interface and wire protocol, discovery and loading, custom collector registration, an isolation/security model, plugin configuration, and an SDK with an example and authoring docs.

## Why out-of-process

ADR-007 commits SysKit to out-of-process plugins (not Go `plugin` buildmode) for isolation, stability, and security: a misbehaving plugin cannot crash or corrupt the host. This epic implements that decision; it does not revisit it.

## Success criteria

- The plugin protocol is defined and **frozen before** dependent stories build on it (PLG-02+).
- Plugins are discovered, loaded, and can register custom collectors that flow through the normal service/render path.
- The isolation model is reviewed against `../SECURITY.md`: a plugin runs with least privilege and cannot escalate through SysKit.
- A published SDK lets an external author build the example plugin from documentation alone.
- `v0.5.0` is tagged on green `main`.

## Stories

Authoritative list: [`../product-backlog.md`](../product-backlog.md#epic-05--extensibility--plugins-v05).

| ID | Story | Pts | Sprint | Spec |
|---|---|---|---|---|
| PLG-01 | Plugin interface + out-of-process protocol. | 13 | 12 | `../specs/plugin-architecture.md` |
| PLG-02 | Discovery & loading. | 8 | 12 | `../specs/plugin-architecture.md` |
| PLG-03 | Custom collector registration. | 5 | 12 | `../specs/features/plugins.md` |
| PLG-04 | Isolation & security model. | 8 | 13 | `../specs/plugin-architecture.md`, `../SECURITY.md` |
| PLG-05 | Plugin configuration. | 3 | 13 | `../specs/configuration.md` |
| PLG-06 | SDK + example plugin + authoring docs. | 5 | 13 | `../specs/features/plugins.md` |
| REL-v05 | Release v0.5.0. | 3 | 13 | `../docs/release-process.md` |

## Dependencies & risk

- Blocked by EPIC-01 (collector/service/render must be stable — plugins extend them).
- **R-02 (protocol churn)** is dominant. Mitigation: run spike **SPK-PLG** in Sprint 11 refinement to prototype the wire protocol and confirm ADR-007's assumptions before committing PLG-02. PLG-03 is the first story to defer if Sprint 12 overflows.

## Definition of Done for the epic

All stories meet the DoD; the protocol has versioned, tested framing; the security model is reviewed and documented; the example plugin builds against the published SDK in CI; `v0.5.0` tagged.
