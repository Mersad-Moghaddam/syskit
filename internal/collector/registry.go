package collector

import (
	"errors"
	"fmt"
	"sort"
)

// Registration describes one domain collector as seen by the discovery seam: a
// stable domain name, a one-line human summary, and a type-erased CollectFunc
// that builds and runs the collector against an injected platform.SysFS.
//
// Registrations are inert data. Building a Registration never touches the OS
// and never references another collector, which is what keeps registered
// collectors independent: a command discovers them through the Registry by name
// and never through a compile-time import of one collector by another
// (specs/collectors.md "Keep collectors independent").
type Registration struct {
	// Name is the domain identifier used for lookup, e.g. "cpu", "memory".
	Name string
	// Summary is a short human description for `--help`/listing output.
	Summary string
	// Collect runs the collector against an injected SysFS and returns its
	// snapshot in erased form.
	Collect CollectFunc
}

// Registration/lookup sentinels, testable with errors.Is.
var (
	// ErrAlreadyRegistered is returned by Register when a name is already taken.
	ErrAlreadyRegistered = errors.New("collector already registered")
	// ErrInvalidRegistration is returned by Register for a registration with an
	// empty Name or a nil Collect function.
	ErrInvalidRegistration = errors.New("invalid collector registration")
)

// Registry maps a domain name to its Registration. It is an explicitly
// constructed value passed by reference, never a package-level global mutated
// by init(); this satisfies the "no package-level mutable state" rule in
// specs/collectors.md "Concurrency" and standards/coding-conventions.md
// "Avoiding Globals". The intended lifecycle is: construct once with
// NewRegistry, Register every domain during program wiring on a single
// goroutine, then treat the Registry as read-only. Concurrent reads (Lookup,
// Names, All) after wiring are safe; concurrent Register calls are not, by
// design, because registration is a one-time startup step.
type Registry struct {
	byName map[string]Registration
}

// NewRegistry returns an empty Registry ready for Register.
func NewRegistry() *Registry {
	return &Registry{byName: make(map[string]Registration)}
}

// Register adds reg to the registry. It returns ErrInvalidRegistration if the
// name is empty or Collect is nil, and ErrAlreadyRegistered if the name is
// already taken. Duplicate registration is rejected rather than silently
// overwriting so that a wiring mistake (two domains claiming one name) surfaces
// immediately instead of shadowing a collector.
func (r *Registry) Register(reg Registration) error {
	if reg.Name == "" || reg.Collect == nil {
		return fmt.Errorf("%q: %w", reg.Name, ErrInvalidRegistration)
	}
	if _, exists := r.byName[reg.Name]; exists {
		return fmt.Errorf("%q: %w", reg.Name, ErrAlreadyRegistered)
	}
	r.byName[reg.Name] = reg
	return nil
}

// Lookup returns the registration for name and whether it was found.
func (r *Registry) Lookup(name string) (Registration, bool) {
	reg, ok := r.byName[name]
	return reg, ok
}

// Names returns the registered domain names in lexical order, so listing and
// help output are deterministic.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.byName))
	for name := range r.byName {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// All returns every registration ordered by Name, for deterministic iteration.
func (r *Registry) All() []Registration {
	regs := make([]Registration, 0, len(r.byName))
	for _, name := range r.Names() {
		regs = append(regs, r.byName[name])
	}
	return regs
}
