package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

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
	c.enrichInterfaces(interfaces)
	info := &model.NetworkInfo{Interfaces: interfaces, CollectedAt: time.Now().UTC()}
	if routes, err := c.fs.ReadFile("proc/net/route"); err == nil {
		info.Routes, _ = ParseRoutes(routes)
	}
	if resolv, err := c.fs.ReadFile("etc/resolv.conf"); err == nil {
		info.Nameservers = ParseResolvConf(resolv)
	}
	return info, nil
}

// enrichInterfaces supplements procfs counters with stable per-interface
// attributes exposed by sysfs. Missing attributes are optional because an
// interface may disappear between the procfs snapshot and these reads.
func (c *Collector) enrichInterfaces(interfaces []model.NetworkInterface) {
	for i := range interfaces {
		base := "sys/class/net/" + interfaces[i].Name + "/"
		if data, err := c.fs.ReadFile(base + "operstate"); err == nil {
			interfaces[i].State = strings.TrimSpace(string(data))
		}
		if data, err := c.fs.ReadFile(base + "address"); err == nil {
			interfaces[i].MACAddress = strings.TrimSpace(string(data))
		}
		if data, err := c.fs.ReadFile(base + "mtu"); err == nil {
			if value, parseErr := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 32); parseErr == nil {
				mtu := uint32(value)
				interfaces[i].MTU = &mtu
			}
		}
	}
}
func ParseRoutes(data []byte) ([]model.Route, error) {
	var routes []model.Route
	for _, line := range strings.Split(string(data), "\n") {
		f := strings.Fields(line)
		if len(f) < 3 || f[0] == "Iface" {
			continue
		}
		destination, err := decodeIPv4LE(f[1])
		if err != nil {
			return nil, err
		}
		gateway, err := decodeIPv4LE(f[2])
		if err != nil {
			return nil, err
		}
		routes = append(routes, model.Route{Interface: f[0], Destination: destination, Gateway: gateway, Default: f[1] == "00000000"})
	}
	return routes, nil
}
func decodeIPv4LE(raw string) (string, error) {
	v, err := strconv.ParseUint(raw, 16, 32)
	if err != nil {
		return "", fmt.Errorf("route address %q: %w", raw, collector.ErrParse)
	}
	return net.IPv4(byte(v), byte(v>>8), byte(v>>16), byte(v>>24)).String(), nil
}
func ParseResolvConf(data []byte) []string {
	var result []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		f := strings.Fields(line)
		if len(f) == 2 && f[0] == "nameserver" {
			result = append(result, f[1])
		}
	}
	return result
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
