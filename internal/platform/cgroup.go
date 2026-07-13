package platform

import (
	"fmt"
	"strconv"
	"strings"
)

// CgroupVersion identifies the Linux cgroup hierarchy layout.
type CgroupVersion uint8

const (
	CgroupUnknown CgroupVersion = iota
	CgroupV1
	CgroupV2
)

// CgroupMembership is one /proc/<pid>/cgroup hierarchy assignment.
type CgroupMembership struct {
	Hierarchy   string
	Controllers []string
	Path        string
}

// CgroupInfo is the normalized cgroup layout and membership view used by
// higher layers. Paths are relative to their cgroup mount roots.
type CgroupInfo struct {
	Version     CgroupVersion
	Memberships []CgroupMembership
}

// CgroupMetrics contains normalized controller counters. Nil fields mean the
// relevant controller or statistic was unavailable on this host.
type CgroupMetrics struct {
	MemoryCurrentBytes  *uint64
	CPUUsageNanoseconds *uint64
	ReadBytes           *uint64
	WrittenBytes        *uint64
}

// DetectCgroup determines the mounted cgroup layout and reads a process's
// assignments. pidPath is a slash-relative procfs path such as
// "proc/self/cgroup" or "proc/123/cgroup".
func DetectCgroup(fs SysFS, pidPath string) (*CgroupInfo, error) {
	data, err := fs.ReadFile(pidPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", pidPath, err)
	}
	info := &CgroupInfo{Memberships: ParseCgroupMembership(data)}
	if _, err := fs.ReadFile("sys/fs/cgroup/cgroup.controllers"); err == nil {
		info.Version = CgroupV2
	} else if len(info.Memberships) > 0 {
		info.Version = CgroupV1
	}
	return info, nil
}

// ParseCgroupMembership parses the kernel's hierarchy:controllers:path rows.
// Malformed rows are skipped so a transient or unknown controller never hides
// otherwise usable membership data.
func ParseCgroupMembership(data []byte) []CgroupMembership {
	var result []CgroupMembership
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(line, ":", 3)
		if len(parts) != 3 || parts[2] == "" {
			continue
		}
		controllers := []string(nil)
		if parts[1] != "" {
			controllers = strings.Split(parts[1], ",")
		}
		result = append(result, CgroupMembership{Hierarchy: parts[0], Controllers: controllers, Path: parts[2]})
	}
	return result
}

// ReadCgroupMetrics reads the accounting files exposed for a normalized cgroup.
// Missing optional controllers are represented by nil fields rather than errors.
func ReadCgroupMetrics(fs SysFS, info *CgroupInfo) (*CgroupMetrics, error) {
	if info == nil || info.Version == CgroupUnknown || len(info.Memberships) == 0 {
		return nil, fmt.Errorf("cgroup layout: %w", ErrUnsupported)
	}
	metrics := &CgroupMetrics{}
	if info.Version == CgroupV2 {
		path := "sys/fs/cgroup/" + cgroupPath(info.Memberships[0].Path)
		metrics.MemoryCurrentBytes = readUint(fs, path+"/memory.current")
		if data, err := fs.ReadFile(path + "/cpu.stat"); err == nil {
			metrics.CPUUsageNanoseconds = cpuUsageV2(data)
		}
		if data, err := fs.ReadFile(path + "/io.stat"); err == nil {
			metrics.ReadBytes, metrics.WrittenBytes = ioUsage(data)
		}
		return metrics, nil
	}
	for _, membership := range info.Memberships {
		path := cgroupPath(membership.Path)
		for _, controller := range membership.Controllers {
			switch controller {
			case "memory":
				metrics.MemoryCurrentBytes = readUint(fs, "sys/fs/cgroup/memory/"+path+"/memory.usage_in_bytes")
			case "cpuacct":
				metrics.CPUUsageNanoseconds = readUint(fs, "sys/fs/cgroup/cpuacct/"+path+"/cpuacct.usage")
			case "blkio":
				if data, err := fs.ReadFile("sys/fs/cgroup/blkio/" + path + "/blkio.throttle.io_service_bytes"); err == nil {
					metrics.ReadBytes, metrics.WrittenBytes = ioUsage(data)
				}
			}
		}
	}
	return metrics, nil
}
func cgroupPath(path string) string { return strings.TrimPrefix(path, "/") }

// ContainerIDFromCgroupPath extracts a runtime-style hexadecimal container ID
// from a cgroup path. It is intentionally best-effort: runtime metadata remains
// optional and callers must treat an empty result as unknown, not host scope.
func ContainerIDFromCgroupPath(path string) string {
	for _, part := range strings.Split(path, "/") {
		part = strings.TrimSuffix(part, ".scope")
		for _, prefix := range []string{"docker-", "cri-containerd-", "crio-"} {
			part = strings.TrimPrefix(part, prefix)
		}
		if len(part) == 64 && isHex(part) {
			return part
		}
	}
	return ""
}

// ContainerRuntimeFromCgroupPath returns a conservative runtime hint inferred
// from common cgroup naming conventions. An empty result means unknown.
func ContainerRuntimeFromCgroupPath(path string) string {
	for _, part := range strings.Split(path, "/") {
		switch {
		case strings.HasPrefix(part, "docker-") || part == "docker":
			return "docker"
		case strings.HasPrefix(part, "cri-containerd-") || strings.Contains(part, "containerd"):
			return "containerd"
		case strings.HasPrefix(part, "crio-") || strings.Contains(part, "crio"):
			return "cri-o"
		}
	}
	return ""
}
func isHex(value string) bool {
	for _, r := range value {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') {
			return false
		}
	}
	return true
}
func readUint(fs SysFS, path string) *uint64 {
	data, err := fs.ReadFile(path)
	if err != nil {
		return nil
	}
	value, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return nil
	}
	return &value
}
func cpuUsageV2(data []byte) *uint64 {
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "usage_usec" {
			value, err := strconv.ParseUint(fields[1], 10, 64)
			if err == nil {
				value *= 1000
				return &value
			}
		}
	}
	return nil
}
func ioUsage(data []byte) (*uint64, *uint64) {
	var read, written uint64
	foundRead, foundWritten := false, false
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		for _, field := range fields {
			if value, ok := strings.CutPrefix(field, "rbytes="); ok {
				if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
					read += parsed
					foundRead = true
				}
			}
			if value, ok := strings.CutPrefix(field, "wbytes="); ok {
				if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
					written += parsed
					foundWritten = true
				}
			}
		}
		// cgroup v1 blkio.throttle.io_service_bytes uses rows such as
		// "8:0 Read 4096" rather than v2's key=value format.
		if len(fields) >= 3 {
			if value, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
				switch strings.ToLower(fields[1]) {
				case "read":
					read += value
					foundRead = true
				case "write":
					written += value
					foundWritten = true
				}
			}
		}
	}
	var readResult, writtenResult *uint64
	if foundRead {
		readResult = &read
	}
	if foundWritten {
		writtenResult = &written
	}
	return readResult, writtenResult
}
