// Package cpu collects static processor topology and optional frequency data
// from procfs and sysfs.
package cpu

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// Collector reads one static CPU snapshot through the injected filesystem.
type Collector struct{ fs platform.SysFS }

var _ collector.Collector[*model.CPUInfo] = (*Collector)(nil)

// NewCollector returns a CPU collector using fs for all Linux reads.
func NewCollector(fs platform.SysFS) *Collector { return &Collector{fs: fs} }

// Collect parses /proc/cpuinfo and enriches it with optional cpufreq values.
func (c *Collector) Collect() (*model.CPUInfo, error) {
	data, err := c.fs.ReadFile("proc/cpuinfo")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/cpuinfo: %w", err)
	}
	info, ids, err := ParseCPUInfo(data)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/cpuinfo: %w", err)
	}
	info.Architecture = runtime.GOARCH
	info.Caches = c.collectCaches()
	info.Frequencies = c.collectFrequencies(ids)
	stat, err := c.fs.ReadFile("proc/stat")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/stat: %w", err)
	}
	times, err := ParseCPUStat(stat)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/stat: %w", err)
	}
	info.Times, info.CollectedAt = times, time.Now().UTC()
	return info, nil
}

// ParseCPUStat extracts aggregate and per-core counters from /proc/stat.
// Missing trailing fields on older kernels are normalized to zero.
func ParseCPUStat(data []byte) ([]model.CPUTime, error) {
	var times []model.CPUTime
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || (fields[0] != "cpu" && !strings.HasPrefix(fields[0], "cpu")) {
			continue
		}
		id := fields[0]
		if id == "cpu" {
			id = "all"
		} else {
			if _, err := strconv.Atoi(strings.TrimPrefix(id, "cpu")); err != nil {
				continue
			}
			id = strings.TrimPrefix(id, "cpu")
		}
		if len(fields) < 2 {
			return nil, fmt.Errorf("%s counters: %w", fields[0], collector.ErrFieldMissing)
		}
		values := [10]uint64{}
		for i := 1; i < len(fields) && i <= len(values); i++ {
			v, err := strconv.ParseUint(fields[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%s counter %q: %w", fields[0], fields[i], collector.ErrParse)
			}
			values[i-1] = v
		}
		t := model.CPUTime{CPUID: id, User: values[0], Nice: values[1], System: values[2], Idle: values[3], IOWait: values[4], IRQ: values[5], SoftIRQ: values[6], Steal: values[7], Guest: values[8], GuestNice: values[9]}
		// guest values are already included in user/nice and therefore excluded.
		t.Total = t.User + t.Nice + t.System + t.Idle + t.IOWait + t.IRQ + t.SoftIRQ + t.Steal
		times = append(times, t)
	}
	if len(times) == 0 {
		return nil, fmt.Errorf("CPU counters: %w", collector.ErrFieldMissing)
	}
	return times, nil
}

func (c *Collector) collectCaches() []model.CPUCache {
	entries, err := c.fs.ReadDir("sys/devices/system/cpu/cpu0/cache")
	if err != nil {
		return nil
	}
	caches := make([]model.CPUCache, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "index") {
			continue
		}
		prefix := "sys/devices/system/cpu/cpu0/cache/" + entry.Name() + "/"
		levelData, levelErr := c.fs.ReadFile(prefix + "level")
		typeData, typeErr := c.fs.ReadFile(prefix + "type")
		sizeData, sizeErr := c.fs.ReadFile(prefix + "size")
		if levelErr != nil || typeErr != nil || sizeErr != nil {
			continue
		}
		level, levelErr := strconv.Atoi(strings.TrimSpace(string(levelData)))
		size, sizeErr := parseCacheSize(strings.TrimSpace(string(sizeData)))
		if levelErr != nil || sizeErr != nil || level < 1 {
			continue
		}
		caches = append(caches, model.CPUCache{Level: level, Type: strings.TrimSpace(string(typeData)), SizeBytes: size})
	}
	sort.Slice(caches, func(i, j int) bool {
		if caches[i].Level != caches[j].Level {
			return caches[i].Level < caches[j].Level
		}
		return caches[i].Type < caches[j].Type
	})
	return caches
}

func parseCacheSize(value string) (uint64, error) {
	if value == "" {
		return 0, collector.ErrFieldMissing
	}
	multiplier := uint64(1)
	last := value[len(value)-1]
	if last == 'K' || last == 'k' {
		multiplier, value = 1024, value[:len(value)-1]
	} else if last == 'M' || last == 'm' {
		multiplier, value = 1024*1024, value[:len(value)-1]
	}
	n, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cache size %q: %w", value, collector.ErrParse)
	}
	return n * multiplier, nil
}

// ParseCPUInfo parses the processor blocks from /proc/cpuinfo. It returns the
// logical CPU IDs separately so the collector can probe corresponding sysfs
// frequency files without re-parsing raw data.
func ParseCPUInfo(data []byte) (*model.CPUInfo, []int, error) {
	blocks := strings.Split(strings.TrimSpace(string(data)), "\n\n")
	if len(blocks) == 0 || blocks[0] == "" {
		return nil, nil, fmt.Errorf("processor blocks: %w", collector.ErrFieldMissing)
	}
	type topology struct{ physical, core string }
	var ids []int
	topologies := map[topology]struct{}{}
	sockets := map[string]struct{}{}
	var modelName string
	flags := map[string]struct{}{}
	for _, block := range blocks {
		fields := map[string]string{}
		for _, line := range strings.Split(block, "\n") {
			key, value, ok := strings.Cut(line, ":")
			if ok {
				fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
			}
		}
		processor, ok := fields["processor"]
		if !ok {
			return nil, nil, fmt.Errorf("processor ID: %w", collector.ErrFieldMissing)
		}
		id, err := strconv.Atoi(processor)
		if err != nil || id < 0 {
			return nil, nil, fmt.Errorf("processor ID %q: %w", processor, collector.ErrParse)
		}
		ids = append(ids, id)
		if modelName == "" {
			modelName = fields["model name"]
			if modelName == "" {
				modelName = fields["Hardware"]
			}
		}
		for _, flag := range strings.Fields(fields["flags"]) {
			flags[flag] = struct{}{}
		}
		physical, hasPhysical := fields["physical id"]
		core, hasCore := fields["core id"]
		if hasPhysical {
			sockets[physical] = struct{}{}
		}
		if hasPhysical && hasCore {
			topologies[topology{physical, core}] = struct{}{}
		}
	}
	sort.Ints(ids)
	if len(ids) == 0 {
		return nil, nil, fmt.Errorf("logical cores: %w", collector.ErrFieldMissing)
	}
	result := &model.CPUInfo{LogicalCores: len(ids), Model: modelName, Flags: make([]string, 0, len(flags))}
	if len(topologies) > 0 {
		value := len(topologies)
		result.PhysicalCores = &value
	}
	if len(sockets) > 0 {
		value := len(sockets)
		result.Sockets = &value
	}
	for flag := range flags {
		result.Flags = append(result.Flags, flag)
	}
	sort.Strings(result.Flags)
	return result, ids, nil
}

func (c *Collector) collectFrequencies(ids []int) []model.CPUFrequency {
	frequencies := make([]model.CPUFrequency, 0, len(ids))
	for _, id := range ids {
		prefix := fmt.Sprintf("sys/devices/system/cpu/cpu%d/cpufreq/", id)
		frequency := model.CPUFrequency{CPUID: id}
		frequency.CurrentMHz = readMHz(c.fs, prefix+"scaling_cur_freq")
		frequency.MinimumMHz = readMHz(c.fs, prefix+"cpuinfo_min_freq")
		frequency.MaximumMHz = readMHz(c.fs, prefix+"cpuinfo_max_freq")
		if frequency.CurrentMHz != nil || frequency.MinimumMHz != nil || frequency.MaximumMHz != nil {
			frequencies = append(frequencies, frequency)
		}
	}
	return frequencies
}

func readMHz(fs platform.SysFS, path string) *float64 {
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil
	}
	kHz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil || kHz < 0 {
		return nil
	}
	mhz := kHz / 1000
	return &mhz
}

// IsOptionalFrequencyError reports whether an absent cpufreq interface is an
// expected capability gap. It is retained for callers and tests that need to
// distinguish absence from malformed values without string matching.
func IsOptionalFrequencyError(err error) bool { return errors.Is(err, platform.ErrNotFound) }
