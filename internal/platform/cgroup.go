package platform

import (
	"fmt"
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
