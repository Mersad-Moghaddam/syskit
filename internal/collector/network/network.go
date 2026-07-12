package network

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

type Collector struct{ fs platform.SysFS }

var _ collector.Collector[*model.NetworkInfo] = (*Collector)(nil)

func NewCollector(fs platform.SysFS) *Collector { return &Collector{fs} }
func (c *Collector) Collect() (*model.NetworkInfo, error) {
	data, err := c.fs.ReadFile("proc/net/dev")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/net/dev: %w", err)
	}
	interfaces, err := ParseDev(data)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/net/dev: %w", err)
	}
	return &model.NetworkInfo{Interfaces: interfaces}, nil
}
func ParseDev(data []byte) ([]model.NetworkInterface, error) {
	var interfaces []model.NetworkInterface
	for _, line := range strings.Split(string(data), "\n") {
		name, rest, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(rest)
		if len(fields) < 16 {
			return nil, fmt.Errorf("network counters: %w", collector.ErrParse)
		}
		values := make([]uint64, 16)
		for i := range values {
			v, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("network counter %q: %w", fields[i], collector.ErrParse)
			}
			values[i] = v
		}
		interfaces = append(interfaces, model.NetworkInterface{Name: strings.TrimSpace(name), RXBytes: values[0], RXPackets: values[1], RXErrors: values[2], RXDrops: values[3], TXBytes: values[8], TXPackets: values[9], TXErrors: values[10], TXDrops: values[11]})
	}
	if len(interfaces) == 0 {
		return nil, fmt.Errorf("network interfaces: %w", collector.ErrFieldMissing)
	}
	return interfaces, nil
}
