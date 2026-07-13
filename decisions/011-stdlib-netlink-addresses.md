# 011. Use the standard-library Netlink adapter for interface addresses

**Status:** Accepted, 2026-07-13

## Context

ADR 003 selected native Netlink for network addresses and anticipated
`golang.org/x/sys/unix` because older standard-library APIs did not provide a
convenient raw interface. The current Go baseline exposes the small subset
needed for an address dump: `syscall.NetlinkRIB`, `ParseNetlinkMessage`, and
`ParseNetlinkRouteAttr`.

Adding a dependency solely for this bounded `RTM_GETADDR` operation would not
meet the standard-library-first policy.

## Decision

The platform layer implements an injectable `AddressSource` using those
standard-library Netlink APIs. It decodes IPv4 and IPv6 address attributes and
returns CIDR strings associated with interface names. Collectors receive the
adapter as a dependency; fixture tests use a stub instead of a live socket.

## Consequences

- Network address collection stays Linux-native, read-only, and free of a new
  module dependency.
- The implementation is deliberately limited to address dumps. A future route
  or link-message adapter may adopt `x/sys/unix` if the standard library no
  longer provides an adequate, maintainable surface; that change requires a
  dependency-policy review and a superseding ADR.

## References

- [ADR 003](003-native-apis-over-shell.md)
- [Dependency Policy](../standards/dependency-policy.md)
- [Network specification](../specs/features/network.md)
