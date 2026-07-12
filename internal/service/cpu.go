package service

import (
	"fmt"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

// CPUCollector is the static CPU collection boundary consumed by CPU.
type CPUCollector interface {
	Collect() (*model.CPUInfo, error)
}

// CPU owns CPU-domain application logic. Utilization sampling is added in the
// following CPU-02 slice; this first slice exposes stable identity data.
type CPU struct {
	collector CPUCollector
	sleep     func(time.Duration)
}

func NewCPU(collector CPUCollector) *CPU { return &CPU{collector: collector, sleep: time.Sleep} }
func NewCPUWithSleep(collector CPUCollector, sleep func(time.Duration)) *CPU {
	return &CPU{collector: collector, sleep: sleep}
}
func (s *CPU) Collect() (*model.CPUInfo, error) { return s.collector.Collect() }

// Sample collects two snapshots separated by interval and derives utilization.
func (s *CPU) Sample(interval time.Duration) (*model.CPUInfo, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("CPU sample interval must be positive")
	}
	before, err := s.Collect()
	if err != nil {
		return nil, err
	}
	s.sleep(interval)
	after, err := s.Collect()
	if err != nil {
		return nil, err
	}
	previous := make(map[string]model.CPUTime, len(before.Times))
	for _, t := range before.Times {
		previous[t.CPUID] = t
	}
	for i := range after.Times {
		old, ok := previous[after.Times[i].CPUID]
		if !ok || after.Times[i].Total <= old.Total {
			continue
		}
		total := after.Times[i].Total - old.Total
		idle := (after.Times[i].Idle + after.Times[i].IOWait) - (old.Idle + old.IOWait)
		if idle > total {
			continue
		}
		value := float64(total-idle) * 100 / float64(total)
		after.Times[i].Utilization = &value
	}
	return after, nil
}
