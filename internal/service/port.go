package service

import (
	"fmt"
	"strconv"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type PortCollector interface {
	Collect() (*model.PortInfo, error)
}
type PortOptions struct {
	Listening bool
	Protocol  string
	Port      int
}
type Port struct{ collector PortCollector }

func NewPort(c PortCollector) *Port { return &Port{c} }
func (s *Port) List(o PortOptions) (*model.PortInfo, error) {
	if o.Port < 0 || o.Port > 65535 {
		return nil, fmt.Errorf("port must be between 0 and 65535")
	}
	info, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	out := &model.PortInfo{}
	for _, socket := range info.Sockets {
		if o.Listening && socket.State != "LISTEN" {
			continue
		}
		if o.Protocol != "" && socket.Protocol != o.Protocol {
			continue
		}
		if o.Port > 0 && socket.LocalPort != uint16(o.Port) {
			continue
		}
		out.Sockets = append(out.Sockets, socket)
	}
	return out, nil
}
func ParsePort(raw string) (int, error) {
	if raw == "" {
		return 0, nil
	}
	return strconv.Atoi(raw)
}
