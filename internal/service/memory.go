package service

import "github.com/Mersad-Moghaddam/syskit/internal/model"

type MemoryCollector interface {
	Collect() (*model.MemoryInfo, error)
}
type Memory struct{ collector MemoryCollector }

func NewMemory(c MemoryCollector) *Memory             { return &Memory{c} }
func (s *Memory) Collect() (*model.MemoryInfo, error) { return s.collector.Collect() }
