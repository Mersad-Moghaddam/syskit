// Package port parses Linux socket tables without shelling out to ss/netstat.
package port

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

type Collector struct{ fs platform.SysFS }

var _ collector.Collector[*model.PortInfo] = (*Collector)(nil)

func NewCollector(fs platform.SysFS) *Collector { return &Collector{fs} }
func (c *Collector) Collect() (*model.PortInfo, error) {
	info := &model.PortInfo{}
	for _, source := range []struct{ path, protocol string }{{"proc/net/tcp", "tcp"}, {"proc/net/udp", "udp"}} {
		data, err := c.fs.ReadFile(source.path)
		if err != nil {
			continue
		}
		sockets, err := ParseSocketTable(data, source.protocol)
		if err != nil {
			return nil, fmt.Errorf("parsing /%s: %w", source.path, err)
		}
		info.Sockets = append(info.Sockets, sockets...)
	}
	if len(info.Sockets) == 0 {
		return nil, fmt.Errorf("socket tables: %w", collector.ErrFieldMissing)
	}
	return info, nil
}
func ParseSocketTable(data []byte, protocol string) ([]model.Socket, error) {
	var sockets []model.Socket
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "local_address") {
			continue
		}
		f := strings.Fields(line)
		if len(f) < 10 {
			if strings.TrimSpace(line) == "" {
				continue
			}
			return nil, fmt.Errorf("socket row: %w", collector.ErrParse)
		}
		local, lp, err := decodeEndpoint(f[1])
		if err != nil {
			return nil, err
		}
		remote, rp, err := decodeEndpoint(f[2])
		if err != nil {
			return nil, err
		}
		inode, err := strconv.ParseUint(f[9], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("socket inode: %w", collector.ErrParse)
		}
		sockets = append(sockets, model.Socket{Protocol: protocol, LocalAddress: local, LocalPort: lp, RemoteAddress: remote, RemotePort: rp, State: tcpState(f[3], protocol), Inode: inode})
	}
	return sockets, nil
}
func decodeEndpoint(value string) (string, uint16, error) {
	address, port, ok := strings.Cut(value, ":")
	if !ok {
		return "", 0, fmt.Errorf("socket endpoint %q: %w", value, collector.ErrParse)
	}
	p, err := strconv.ParseUint(port, 16, 16)
	if err != nil {
		return "", 0, fmt.Errorf("socket port: %w", collector.ErrParse)
	}
	if len(address) != 8 {
		return address, uint16(p), nil
	}
	raw, err := hex.DecodeString(address)
	if err != nil {
		return "", 0, err
	}
	for i := 0; i < 2; i++ {
		raw[i], raw[3-i] = raw[3-i], raw[i]
	}
	return net.IP(raw).String(), uint16(p), nil
}
func tcpState(raw, protocol string) string {
	if protocol == "udp" {
		return "UNCONN"
	}
	states := map[string]string{"01": "ESTABLISHED", "0A": "LISTEN", "06": "TIME_WAIT", "0C": "NEW_SYN_RECV", "07": "CLOSE", "08": "CLOSE_WAIT"}
	if state, ok := states[raw]; ok {
		return state
	}
	return raw
}
