package service

import "github.com/Mersad-Moghaddam/syskit/internal/model"

// SystemCollector is the narrow dependency required by System. It makes the
// service testable without coupling it to a concrete collector package.
type SystemCollector interface {
	Collect() (*model.SystemInfo, error)
}

// System coordinates the host-summary collection. It deliberately contains no
// rendering or platform I/O; its seam leaves room for future aggregation.
type System struct{ collector SystemCollector }

// NewSystem constructs the system-summary service.
func NewSystem(collector SystemCollector) *System { return &System{collector: collector} }

// Collect returns one host snapshot.
func (s *System) Collect() (*model.SystemInfo, error) { return s.collector.Collect() }
