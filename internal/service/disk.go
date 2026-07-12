package service

import (
	"fmt"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type DiskCollector interface {
	Collect() (*model.DiskInfo, error)
}
type Disk struct{ collector DiskCollector }

func NewDisk(c DiskCollector) *Disk               { return &Disk{c} }
func (s *Disk) Collect() (*model.DiskInfo, error) { return s.collector.Collect() }
func (s *Disk) Sample(interval time.Duration) (*model.DiskInfo, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("disk sample interval must be positive")
	}
	first, err := s.Collect()
	if err != nil {
		return nil, err
	}
	time.Sleep(interval)
	second, err := s.Collect()
	if err != nil {
		return nil, err
	}
	old := map[string]model.DiskDevice{}
	for _, d := range first.Devices {
		old[d.Name] = d
	}
	elapsed := second.CollectedAt.Sub(first.CollectedAt).Seconds()
	if elapsed <= 0 {
		elapsed = interval.Seconds()
	}
	for i := range second.Devices {
		if prior, ok := old[second.Devices[i].Name]; ok && second.Devices[i].ReadBytes >= prior.ReadBytes {
			r := float64(second.Devices[i].ReadBytes-prior.ReadBytes) / elapsed
			w := float64(second.Devices[i].WrittenBytes-prior.WrittenBytes) / elapsed
			second.Devices[i].ReadBytesPerSecond = &r
			second.Devices[i].WrittenBytesPerSecond = &w
		}
	}
	return second, nil
}
