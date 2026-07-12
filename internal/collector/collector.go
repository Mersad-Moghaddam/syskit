package collector

import "github.com/Mersad-Moghaddam/syskit/internal/platform"

// Collector is the shared contract every domain collector (cpu, memory, disk,
// process, network, ports, fs) implements. It generalizes the informal
// signature in ARCHITECTURE.md §4 — `Collect() (*model.T, error)` — into a
// single, type-safe statement.
//
// It is generic (Go 1.22) rather than a plain interface because each domain
// returns a *different* typed snapshot. A non-generic `Collect() (any, error)`
// would erase that type at the source and force every caller to type-assert;
// a per-domain hand-written interface would restate the same one-method shape
// in every package. `Collector[T]` states the contract exactly once while
// preserving each domain's concrete snapshot type end to end.
//
// A domain constructor follows the convention `func NewCollector(platform.SysFS)`
// returning the concrete collector struct ("accept interfaces, return structs")
// and asserts `var _ collector.Collector[Snapshot] = (*Collector)(nil)` so the
// contract is checked at compile time. See internal/collector/example for a
// reference implementation.
//
// Collect returns a point-in-time snapshot. Rate-based metrics that need two
// snapshots and a time delta are the service layer's concern (specs/collectors.md
// "Snapshot Model"); collectors stay stateless and reusable.
type Collector[T any] interface {
	Collect() (T, error)
}

// CollectFunc is the type-erased form of a domain collector, used only at the
// registration seam. It builds a collector against fs, runs it, and returns the
// snapshot as any so the Registry can hold heterogeneous domains in one map and
// a command can run "any registered collector by name" without knowing T at
// compile time. Callers that need the concrete snapshot type use the domain's
// own typed Collector[T] directly instead of the registry.
type CollectFunc func(fs platform.SysFS) (any, error)

// Adapt bridges a typed domain constructor to the non-generic registration
// seam. It captures newCollector, and on each call injects fs, runs Collect,
// and erases the typed snapshot to any at exactly one boundary. Type safety is
// preserved everywhere inside the collector; erasure happens only where the
// heterogeneous Registry requires it.
func Adapt[T any](newCollector func(fs platform.SysFS) Collector[T]) CollectFunc {
	return func(fs platform.SysFS) (any, error) {
		return newCollector(fs).Collect()
	}
}
