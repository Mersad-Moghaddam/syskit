// Package process reads process snapshots from procfs without invoking ps.
package process

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

type Collector struct{ fs platform.SysFS }

var _ collector.Collector[*model.ProcessList] = (*Collector)(nil)

func NewCollector(fs platform.SysFS) *Collector { return &Collector{fs} }
func (c *Collector) Collect() (*model.ProcessList, error) {
	entries, err := c.fs.ReadDir("proc")
	if err != nil {
		return nil, fmt.Errorf("listing /proc: %w", err)
	}
	users := c.users()
	list := &model.ProcessList{}
	if data, err := c.fs.ReadFile("proc/stat"); err == nil {
		list.CPUTimeTotal = ParseCPUTotal(data)
	}
	if data, err := c.fs.ReadFile("proc/meminfo"); err == nil {
		list.TotalMemoryBytes = ParseMemoryTotal(data)
	}
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid < 1 {
			continue
		}
		p, err := c.collectPID(pid)
		if err != nil {
			if errors.Is(err, platform.ErrNotFound) || errors.Is(err, platform.ErrPermission) {
				list.Partial = list.Partial || errors.Is(err, platform.ErrPermission)
				continue
			}
			return nil, err
		}
		p.User = users[p.UID]
		list.Processes = append(list.Processes, *p)
	}
	return list, nil
}

// ParseCPUTotal returns aggregate clock ticks from the first cpu row.
func ParseCPUTotal(data []byte) uint64 {
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != "cpu" {
			continue
		}
		var total uint64
		for _, field := range fields[1:] {
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return 0
			}
			total += value
		}
		return total
	}
	return 0
}

// ParseMemoryTotal returns MemTotal in bytes; procfs reports KiB.
func ParseMemoryTotal(data []byte) uint64 {
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok || key != "MemTotal" {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			return 0
		}
		amount, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return 0
		}
		return amount * 1024
	}
	return 0
}
func (c *Collector) collectPID(pid int) (*model.Process, error) {
	base := fmt.Sprintf("proc/%d/", pid)
	stat, err := c.fs.ReadFile(base + "stat")
	if err != nil {
		return nil, err
	}
	p, err := ParseStat(stat)
	if err != nil {
		return nil, err
	}
	status, err := c.fs.ReadFile(base + "status")
	if err == nil {
		applyStatus(p, status)
	}
	cmdline, err := c.fs.ReadFile(base + "cmdline")
	if err == nil && len(cmdline) > 0 {
		p.Command = strings.ReplaceAll(strings.TrimRight(string(cmdline), "\x00"), "\x00", " ")
	}
	if cgroup, err := c.fs.ReadFile(base + "cgroup"); err == nil {
		for _, membership := range platform.ParseCgroupMembership(cgroup) {
			if id := platform.ContainerIDFromCgroupPath(membership.Path); id != "" {
				p.ContainerID = id
				break
			}
		}
	}
	return p, nil
}

// ParseStat uses the final ')' to safely handle command names containing spaces or parentheses.
func ParseStat(data []byte) (*model.Process, error) {
	raw := strings.TrimSpace(string(data))
	left := strings.IndexByte(raw, '(')
	right := strings.LastIndexByte(raw, ')')
	if left < 0 || right < left {
		return nil, fmt.Errorf("process stat command: %w", collector.ErrParse)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(raw[:left]))
	if err != nil {
		return nil, fmt.Errorf("process PID: %w", collector.ErrParse)
	}
	rest := strings.Fields(raw[right+1:])
	if len(rest) < 20 {
		return nil, fmt.Errorf("process stat fields: %w", collector.ErrFieldMissing)
	}
	ppid, err := strconv.Atoi(rest[1])
	if err != nil {
		return nil, fmt.Errorf("process PPID: %w", collector.ErrParse)
	}
	utime, _ := strconv.ParseUint(rest[11], 10, 64)
	stime, _ := strconv.ParseUint(rest[12], 10, 64)
	threads, _ := strconv.ParseUint(rest[17], 10, 64)
	startTime, _ := strconv.ParseUint(rest[19], 10, 64)
	return &model.Process{PID: pid, PPID: ppid, State: rest[0], Command: raw[left+1 : right], CPUTime: utime + stime, StartTimeTicks: startTime, Threads: threads}, nil
}

func (c *Collector) users() map[uint64]string {
	data, err := c.fs.ReadFile("etc/passwd")
	if err != nil {
		return nil
	}
	return ParsePasswd(data)
}

// ParsePasswd returns the login name for each syntactically valid passwd row.
// It intentionally accepts entries with an empty password field and ignores
// malformed rows so an optional identity lookup cannot break process listing.
func ParsePasswd(data []byte) map[uint64]string {
	users := map[uint64]string{}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Split(line, ":")
		if len(fields) < 3 || fields[0] == "" {
			continue
		}
		uid, err := strconv.ParseUint(fields[2], 10, 64)
		if err == nil {
			users[uid] = fields[0]
		}
	}
	return users
}
func applyStatus(p *model.Process, data []byte) {
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}
		switch key {
		case "Uid":
			p.UID, _ = strconv.ParseUint(fields[0], 10, 64)
		case "VmRSS":
			v, _ := strconv.ParseUint(fields[0], 10, 64)
			p.ResidentBytes = v * 1024
		case "Threads":
			p.Threads, _ = strconv.ParseUint(fields[0], 10, 64)
		}
	}
}
