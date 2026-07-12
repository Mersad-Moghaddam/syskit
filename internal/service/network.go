package service

import (
	"fmt"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type NetworkCollector interface {
	Collect() (*model.NetworkInfo, error)
}
type Network struct{ collector NetworkCollector }

func NewNetwork(c NetworkCollector) *Network            { return &Network{c} }
func (s *Network) Collect() (*model.NetworkInfo, error) { return s.collector.Collect() }
func (s *Network) Sample(interval time.Duration) (*model.NetworkInfo, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("network sample interval must be positive")
	}
	before, err := s.Collect()
	if err != nil {
		return nil, err
	}
	time.Sleep(interval)
	after, err := s.Collect()
	if err != nil {
		return nil, err
	}
	old := map[string]model.NetworkInterface{}
	for _, n := range before.Interfaces {
		old[n.Name] = n
	}
	elapsed := after.CollectedAt.Sub(before.CollectedAt).Seconds()
	if elapsed <= 0 {
		elapsed = interval.Seconds()
	}
	for i := range after.Interfaces {
		if p, ok := old[after.Interfaces[i].Name]; ok && after.Interfaces[i].RXBytes >= p.RXBytes && after.Interfaces[i].TXBytes >= p.TXBytes {
			rx := float64(after.Interfaces[i].RXBytes-p.RXBytes) / elapsed
			tx := float64(after.Interfaces[i].TXBytes-p.TXBytes) / elapsed
			after.Interfaces[i].RXBytesPerSecond = &rx
			after.Interfaces[i].TXBytesPerSecond = &tx
		}
	}
	return after, nil
}
