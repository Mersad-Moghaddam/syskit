package service

import "github.com/Mersad-Moghaddam/syskit/internal/model"

type DiskCollector interface {
	Collect() (*model.DiskInfo, error)
}
type Disk struct{ collector DiskCollector }

func NewDisk(c DiskCollector) *Disk               { return &Disk{c} }
func (s *Disk) Collect() (*model.DiskInfo, error) { return s.collector.Collect() }
