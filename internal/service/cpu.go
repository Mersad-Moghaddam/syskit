package service

import "github.com/Mersad-Moghaddam/syskit/internal/model"

// CPUCollector is the static CPU collection boundary consumed by CPU.
type CPUCollector interface {
	Collect() (*model.CPUInfo, error)
}

// CPU owns CPU-domain application logic. Utilization sampling is added in the
// following CPU-02 slice; this first slice exposes stable identity data.
type CPU struct{ collector CPUCollector }

func NewCPU(collector CPUCollector) *CPU        { return &CPU{collector: collector} }
func (s *CPU) Collect() (*model.CPUInfo, error) { return s.collector.Collect() }
