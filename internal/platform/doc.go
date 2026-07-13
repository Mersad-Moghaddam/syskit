// Package platform is the only layer permitted to touch the operating system.
//
// It exposes the SysFS seam (ReadFile/Open/ReadDir/ReadLink) with a real implementation
// rooted at "/" and a fixture-backed implementation rooted at test data, so
// every collector reads through an injectable interface rather than the host
// filesystem directly. It also owns the Netlink client and the cgroup v1/v2
// reader.
//
// Dependency rule (ADR-004): platform is the lowest SysKit layer; it must not
// import cli, command, service, collector, render, or model. SysKit is
// Linux-only (ADR-002): no runtime.GOOS branching lives here.
package platform
