// Package example is a reference collector, not a real domain. It exists only
// to prove that the collector contract (collector.Collector[T]), the domain
// constructor convention (NewCollector taking platform.SysFS by injection), the
// error-classification sentinels (collector.ErrParse / collector.ErrFieldMissing
// and the optional-missing-as-partial pattern), and the registration seam
// (collector.Registry via collector.Adapt) compile and are testable together.
//
// Real domain collectors (cpu, memory, disk, ...) are delivered by EPIC-01 and
// follow this exact shape. Nothing here should be imported by production code.
package example

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// loadAvgPath is the single kernel interface this reference collector reads,
// expressed as a slash-relative SysFS name (never an absolute /proc path).
const loadAvgPath = "proc/loadavg"

// Reading is the typed snapshot the example collector returns. Load figures are
// dimensionless kernel-reported averages over 1, 5, and 15 minutes.
type Reading struct {
	// One, Five, Fifteen are the required 1/5/15-minute load averages.
	One     float64
	Five    float64
	Fifteen float64

	// Running/Total come from the OPTIONAL "runnable/total" token (e.g. "1/234").
	// When that token is absent or unrecognized these stay zero and EntitiesKnown
	// is false: optional missing data is represented as unavailable, never as an
	// error (specs/collectors.md "Error Classification").
	Running       int
	Total         int
	EntitiesKnown bool
}

// Collector reads and parses proc/loadavg through an injected SysFS. It holds
// only its immutable dependency and no mutable state, so it is safe to call
// Collect concurrently.
type Collector struct {
	fs platform.SysFS
}

// Compile-time proof that *Collector satisfies the shared generic contract.
var _ collector.Collector[Reading] = (*Collector)(nil)

// NewCollector returns a Collector that reads through fs. Injection is the only
// way to supply a data source: the struct has no OS-touching code path, so a
// production wiring passes platform.RealFS() and tests pass platform.TestFS(...)
// with the same collector unchanged.
func NewCollector(fs platform.SysFS) *Collector {
	return &Collector{fs: fs}
}

// Collect reads proc/loadavg via the injected SysFS and parses it into a
// Reading. Platform errors (not found, permission, unsupported) pass through
// wrapped for context; malformed data yields collector.ErrParse; a missing
// required load field yields collector.ErrFieldMissing.
func (c *Collector) Collect() (Reading, error) {
	data, err := c.fs.ReadFile(loadAvgPath)
	if err != nil {
		return Reading{}, fmt.Errorf("reading %s: %w", loadAvgPath, err)
	}
	return parseLoadAvg(data)
}

// parseLoadAvg parses the space-separated /proc/loadavg format:
//
//	0.42 0.35 0.30 1/234 5678
//
// The three leading floats are required; the "runnable/total" token and the
// last-PID token are optional.
func parseLoadAvg(data []byte) (Reading, error) {
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return Reading{}, fmt.Errorf("parsing %s: need 3 load fields, got %d: %w",
			loadAvgPath, len(fields), collector.ErrFieldMissing)
	}

	loads := [3]float64{}
	for i := 0; i < 3; i++ {
		v, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return Reading{}, fmt.Errorf("parsing %s load field %d %q: %w",
				loadAvgPath, i+1, fields[i], collector.ErrParse)
		}
		loads[i] = v
	}

	r := Reading{One: loads[0], Five: loads[1], Fifteen: loads[2]}

	// Optional runnable/total token: best-effort. Absent or unrecognized leaves
	// the entity counts unavailable (EntitiesKnown false) without erroring.
	if len(fields) >= 4 {
		if running, total, ok := parseEntities(fields[3]); ok {
			r.Running, r.Total, r.EntitiesKnown = running, total, true
		}
	}
	return r, nil
}

// parseEntities parses a "runnable/total" token such as "1/234". It reports ok
// false for any shape it does not recognize, so a malformed optional token is
// treated as unavailable rather than a parse error.
func parseEntities(token string) (running, total int, ok bool) {
	before, after, found := strings.Cut(token, "/")
	if !found {
		return 0, 0, false
	}
	running, err := strconv.Atoi(before)
	if err != nil {
		return 0, 0, false
	}
	total, err = strconv.Atoi(after)
	if err != nil {
		return 0, 0, false
	}
	return running, total, true
}

// Register adds the example collector to r under the name "example". It wraps
// NewCollector into the non-generic seam with collector.Adapt, demonstrating
// how a real domain registers itself without any other collector importing it.
func Register(r *collector.Registry) error {
	return r.Register(collector.Registration{
		Name:    "example",
		Summary: "reference collector: parses load averages from proc/loadavg",
		Collect: collector.Adapt(func(fs platform.SysFS) collector.Collector[Reading] {
			return NewCollector(fs)
		}),
	})
}
