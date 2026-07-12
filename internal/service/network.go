package service

import (
	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type NetworkCollector interface {
	Collect() (*model.NetworkInfo, error)
}
type Network struct{ collector NetworkCollector }

func NewNetwork(c NetworkCollector) *Network            { return &Network{c} }
func (s *Network) Collect() (*model.NetworkInfo, error) { return s.collector.Collect() }
