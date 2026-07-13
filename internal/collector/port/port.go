// Package port parses Linux socket tables without shelling out to ss/netstat.
package port

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sort"
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
	for _, source := range []struct{ path, protocol string }{{"proc/net/tcp", "tcp"}, {"proc/net/tcp6", "tcp6"}, {"proc/net/udp", "udp"}, {"proc/net/udp6", "udp6"}} {
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
	if data, err := c.fs.ReadFile("proc/net/unix"); err == nil {
		sockets, parseErr := ParseUnixSocketTable(data)
		if parseErr != nil {
			return nil, fmt.Errorf("parsing /proc/net/unix: %w", parseErr)
		}
		info.Sockets = append(info.Sockets, sockets...)
	}
	if len(info.Sockets) == 0 {
		return nil, fmt.Errorf("socket tables: %w", collector.ErrFieldMissing)
	}
	info.OwnerMappingPartial = c.mapOwners(info)
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
		sockets = append(sockets, model.Socket{Protocol: protocol, LocalAddress: local, LocalPort: lp, RemoteAddress: remote, RemotePort: rp, State: tcpState(f[3], protocol), RawState: f[3], Inode: inode})
	}
	return sockets, nil
}
func ParseUnixSocketTable(data []byte) ([]model.Socket, error) {
	var sockets []model.Socket
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "Num       RefCount") || strings.TrimSpace(line) == "" {
			continue
		}
		f := strings.Fields(line)
		if len(f) < 7 {
			return nil, fmt.Errorf("unix socket row: %w", collector.ErrParse)
		}
		inode, err := strconv.ParseUint(f[6], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unix socket inode: %w", collector.ErrParse)
		}
		path := ""
		if len(f) > 7 {
			path = strings.Join(f[7:], " ")
		}
		sockets = append(sockets, model.Socket{Protocol: "unix", LocalAddress: path, State: unixState(f[5]), RawState: f[5], Inode: inode})
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
	raw, err := hex.DecodeString(address)
	if err != nil {
		return "", 0, fmt.Errorf("socket address: %w", collector.ErrParse)
	}
	switch len(raw) {
	case net.IPv4len:
		reverse4(raw)
	case net.IPv6len:
		for i := 0; i < net.IPv6len; i += net.IPv4len {
			reverse4(raw[i : i+net.IPv4len])
		}
	default:
		return "", 0, fmt.Errorf("socket address length: %w", collector.ErrParse)
	}
	return net.IP(raw).String(), uint16(p), nil
}
func reverse4(value []byte) {
	value[0], value[3], value[1], value[2] = value[3], value[0], value[2], value[1]
}
func tcpState(raw, protocol string) string {
	if protocol == "udp" || protocol == "udp6" {
		return "UNCONN"
	}
	states := map[string]string{
		"01": "ESTABLISHED", "02": "SYN_SENT", "03": "SYN_RECV", "04": "FIN_WAIT1",
		"05": "FIN_WAIT2", "06": "TIME_WAIT", "07": "CLOSE", "08": "CLOSE_WAIT",
		"09": "LAST_ACK", "0A": "LISTEN", "0B": "CLOSING", "0C": "NEW_SYN_RECV",
	}
	if state, ok := states[raw]; ok {
		return state
	}
	return raw
}
func unixState(raw string) string {
	if raw == "01" {
		return "LISTEN"
	}
	return raw
}

func (c *Collector) mapOwners(info *model.PortInfo) bool {
	partial := false
	byInode := make(map[uint64]*model.Socket, len(info.Sockets))
	for i := range info.Sockets {
		byInode[info.Sockets[i].Inode] = &info.Sockets[i]
	}
	entries, err := c.fs.ReadDir("proc")
	if err != nil {
		return errors.Is(err, platform.ErrPermission)
	}
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid < 1 {
			continue
		}
		fds, err := c.fs.ReadDir(fmt.Sprintf("proc/%d/fd", pid))
		if err != nil {
			partial = partial || errors.Is(err, platform.ErrPermission)
			continue
		}
		command := c.command(pid)
		for _, fd := range fds {
			target, err := c.fs.ReadLink(fmt.Sprintf("proc/%d/fd/%s", pid, fd.Name()))
			if err != nil {
				if errors.Is(err, platform.ErrPermission) || errors.Is(err, platform.ErrNotFound) {
					partial = partial || errors.Is(err, platform.ErrPermission)
					continue
				}
				continue
			}
			inode, ok := socketInode(target)
			if !ok {
				continue
			}
			if socket, ok := byInode[inode]; ok {
				socket.Owners = append(socket.Owners, model.SocketOwner{PID: pid, Command: command})
			}
		}
	}
	for i := range info.Sockets {
		sort.Slice(info.Sockets[i].Owners, func(a, b int) bool { return info.Sockets[i].Owners[a].PID < info.Sockets[i].Owners[b].PID })
	}
	return partial
}
func (c *Collector) command(pid int) string {
	data, err := c.fs.ReadFile(fmt.Sprintf("proc/%d/cmdline", pid))
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(strings.TrimRight(string(data), "\x00"), "\x00", " ")
}
func socketInode(target string) (uint64, bool) {
	value, ok := strings.CutPrefix(target, "socket:[")
	if !ok || !strings.HasSuffix(value, "]") {
		return 0, false
	}
	inode, err := strconv.ParseUint(strings.TrimSuffix(value, "]"), 10, 64)
	return inode, err == nil
}
