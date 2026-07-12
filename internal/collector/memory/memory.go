// Package memory parses Linux memory accounting and optional PSI data.
package memory

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

var _ collector.Collector[*model.MemoryInfo] = (*Collector)(nil)

func NewCollector(fs platform.SysFS) *Collector { return &Collector{fs: fs} }

func (c *Collector) Collect() (*model.MemoryInfo, error) {
	data, err := c.fs.ReadFile("proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/meminfo: %w", err)
	}
	info, err := ParseMemInfo(data)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/meminfo: %w", err)
	}
	psi, err := c.fs.ReadFile("proc/pressure/memory")
	if err == nil {
		pressure, parseErr := ParsePSI(psi)
		if parseErr != nil {
			return nil, fmt.Errorf("parsing /proc/pressure/memory: %w", parseErr)
		}
		info.Pressure = pressure
	} else if !errors.Is(err, platform.ErrNotFound) {
		return nil, fmt.Errorf("reading /proc/pressure/memory: %w", err)
	}
	return info, nil
}

func ParseMemInfo(data []byte) (*model.MemoryInfo, error) {
	fields := map[string]uint64{}
	for _, line := range strings.Split(string(data), "\n") {
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		parts := strings.Fields(value)
		if len(parts) == 0 {
			continue
		}
		n, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s %q: %w", name, parts[0], collector.ErrParse)
		}
		if len(parts) > 1 && parts[1] == "kB" {
			n *= 1024
		}
		fields[name] = n
	}
	total, ok := fields["MemTotal"]
	if !ok {
		return nil, fmt.Errorf("MemTotal: %w", collector.ErrFieldMissing)
	}
	info := &model.MemoryInfo{TotalBytes: total, FreeBytes: fields["MemFree"], BuffersBytes: fields["Buffers"], CacheBytes: fields["Cached"] + fields["SReclaimable"], SwapTotalBytes: fields["SwapTotal"], SwapFreeBytes: fields["SwapFree"]}
	if available, ok := fields["MemAvailable"]; ok {
		info.AvailableBytes = &available
		used := total - available
		info.UsedBytes = &used
	}
	if info.SwapFreeBytes <= info.SwapTotalBytes {
		info.SwapUsedBytes = info.SwapTotalBytes - info.SwapFreeBytes
	}
	return info, nil
}

func ParsePSI(data []byte) (*model.MemoryPSI, error) {
	result := &model.MemoryPSI{}
	found := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		for _, part := range parts[1:] {
			key, value, ok := strings.Cut(part, "=")
			if !ok {
				continue
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("PSI %s: %w", key, collector.ErrParse)
			}
			switch parts[0] + "." + key {
			case "some.avg10":
				result.SomeAvg10, found["some"] = v, true
			case "some.avg60":
				result.SomeAvg60 = v
			case "full.avg10":
				result.FullAvg10, found["full"] = v, true
			case "full.avg60":
				result.FullAvg60 = v
			}
		}
	}
	if !found["some"] && !found["full"] {
		return nil, fmt.Errorf("PSI fields: %w", collector.ErrFieldMissing)
	}
	return result, nil
}
