// Package system collects host identity and load information from procfs and
// the standard Linux os-release file.
package system

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// Collector reads a single system snapshot through the injected platform seam.
type Collector struct {
	fs  platform.SysFS
	now func() time.Time
}

var _ collector.Collector[*model.SystemInfo] = (*Collector)(nil)

// NewCollector constructs a collector using the current clock for boot-time
// derivation. Tests that need a fixed clock use NewCollectorWithClock.
func NewCollector(fs platform.SysFS) *Collector {
	return NewCollectorWithClock(fs, time.Now)
}

// NewCollectorWithClock is the injectable-clock constructor used by tests.
func NewCollectorWithClock(fs platform.SysFS, now func() time.Time) *Collector {
	return &Collector{fs: fs, now: now}
}

// Collect reads required procfs files and optional os-release metadata.
func (c *Collector) Collect() (*model.SystemInfo, error) {
	uptimeData, err := c.fs.ReadFile("proc/uptime")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/uptime: %w", err)
	}
	uptime, err := ParseUptime(uptimeData)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/uptime: %w", err)
	}

	loadData, err := c.fs.ReadFile("proc/loadavg")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/loadavg: %w", err)
	}
	load1, load5, load15, err := ParseLoadAverage(loadData)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/loadavg: %w", err)
	}

	hostname, err := readTrimmed(c.fs, "proc/sys/kernel/hostname")
	if err != nil {
		return nil, err
	}
	release, err := readTrimmed(c.fs, "proc/sys/kernel/osrelease")
	if err != nil {
		return nil, err
	}
	version, err := readTrimmed(c.fs, "proc/sys/kernel/version")
	if err != nil {
		return nil, err
	}

	info := &model.SystemInfo{
		Hostname: hostname, KernelRelease: release, KernelVersion: version,
		Architecture: runtime.GOARCH, UptimeSeconds: uptime,
		BootTime:     c.now().Add(-time.Duration(uptime * float64(time.Second))).UTC(),
		LoadAverage1: load1, LoadAverage5: load5, LoadAverage15: load15,
	}
	if data, err := c.fs.ReadFile("etc/os-release"); err == nil {
		info.OSName, info.OSVersion = ParseOSRelease(data)
	} else if !errors.Is(err, platform.ErrNotFound) {
		return nil, fmt.Errorf("reading /etc/os-release: %w", err)
	}
	return info, nil
}

func readTrimmed(fs platform.SysFS, path string) (string, error) {
	data, err := fs.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading /%s: %w", path, err)
	}
	value := strings.TrimSpace(string(data))
	if value == "" {
		return "", fmt.Errorf("parsing /%s: %w", path, collector.ErrFieldMissing)
	}
	return value, nil
}

// ParseUptime extracts the first, seconds-since-boot field from /proc/uptime.
func ParseUptime(data []byte) (float64, error) {
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("uptime field: %w", collector.ErrFieldMissing)
	}
	v, err := strconv.ParseFloat(fields[0], 64)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("uptime %q: %w", fields[0], collector.ErrParse)
	}
	return v, nil
}

// ParseLoadAverage extracts the 1-, 5-, and 15-minute load averages.
func ParseLoadAverage(data []byte) (float64, float64, float64, error) {
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return 0, 0, 0, fmt.Errorf("load average fields: %w", collector.ErrFieldMissing)
	}
	values := [3]float64{}
	for i := range values {
		v, err := strconv.ParseFloat(fields[i], 64)
		if err != nil || v < 0 {
			return 0, 0, 0, fmt.Errorf("load average %q: %w", fields[i], collector.ErrParse)
		}
		values[i] = v
	}
	return values[0], values[1], values[2], nil
}

// ParseOSRelease extracts NAME and VERSION_ID from freedesktop os-release
// syntax. Malformed optional lines are ignored so absent distribution metadata
// never prevents useful kernel information from being shown.
func ParseOSRelease(data []byte) (name, version string) {
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		value = strings.Trim(value, "\"")
		switch key {
		case "NAME":
			name = value
		case "VERSION_ID":
			version = value
		}
	}
	return name, version
}
