// Package disk collects mount and filesystem-statistics data from Linux.
package disk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

type Collector struct{ fs platform.SysFS }

var _ collector.Collector[*model.DiskInfo] = (*Collector)(nil)

func NewCollector(fs platform.SysFS) *Collector { return &Collector{fs} }
func (c *Collector) Collect() (*model.DiskInfo, error) {
	data, err := c.fs.ReadFile("proc/self/mountinfo")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/self/mountinfo: %w", err)
	}
	mounts, err := ParseMountInfo(data)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/self/mountinfo: %w", err)
	}
	for i := range mounts {
		stats, err := c.fs.StatFS(mounts[i].MountPoint)
		if err != nil {
			continue
		}
		total, avail := stats.TotalBytes, stats.AvailableBytes
		used := total - stats.FreeBytes
		mounts[i].TotalBytes, mounts[i].AvailableBytes, mounts[i].UsedBytes = &total, &avail, &used
		if total > 0 {
			pct := float64(used) * 100 / float64(total)
			mounts[i].UsePercent = &pct
		}
		inodes, free := stats.TotalInodes, stats.FreeInodes
		mounts[i].TotalInodes, mounts[i].FreeInodes = &inodes, &free
	}
	statsData, err := c.fs.ReadFile("proc/diskstats")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/diskstats: %w", err)
	}
	devices, err := ParseDiskStats(statsData)
	if err != nil {
		return nil, fmt.Errorf("parsing /proc/diskstats: %w", err)
	}
	return &model.DiskInfo{Mounts: mounts, Devices: devices, CollectedAt: time.Now().UTC()}, nil
}
func ParseDiskStats(data []byte) ([]model.DiskDevice, error) {
	var devices []model.DiskDevice
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		f := strings.Fields(line)
		if len(f) < 10 {
			return nil, fmt.Errorf("diskstats fields: %w", collector.ErrParse)
		}
		values := make([]uint64, 0, 7)
		for _, idx := range []int{3, 5, 7, 9} {
			v, err := strconv.ParseUint(f[idx], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("diskstats %q: %w", f[idx], collector.ErrParse)
			}
			values = append(values, v)
		}
		devices = append(devices, model.DiskDevice{Name: f[2], ReadOperations: values[0], ReadBytes: values[1] * 512, WrittenOperations: values[2], WrittenBytes: values[3] * 512})
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("diskstats entries: %w", collector.ErrFieldMissing)
	}
	return devices, nil
}
func ParseMountInfo(data []byte) ([]model.MountInfo, error) {
	var mounts []model.MountInfo
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		parts := strings.Fields(line)
		separator := -1
		for i, p := range parts {
			if p == "-" {
				separator = i
				break
			}
		}
		if separator < 0 || separator+3 >= len(parts) || separator < 6 {
			return nil, fmt.Errorf("mountinfo line: %w", collector.ErrParse)
		}
		mounts = append(mounts, model.MountInfo{MountPoint: unescape(parts[4]), Options: strings.Split(parts[5], ","), FilesystemType: parts[separator+1], Source: unescape(parts[separator+2])})
	}
	if len(mounts) == 0 {
		return nil, fmt.Errorf("mountinfo entries: %w", collector.ErrFieldMissing)
	}
	return mounts, nil
}
func unescape(v string) string {
	return strings.NewReplacer("\\040", " ", "\\011", "\t", "\\012", "\n", "\\134", "\\").Replace(v)
}
