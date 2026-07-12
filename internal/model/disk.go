package model

// DiskInfo describes mounted filesystem capacity from the current mount namespace.
type DiskInfo struct {
	Mounts []MountInfo `json:"mounts"`
}
type MountInfo struct {
	Source         string   `json:"source"`
	FilesystemType string   `json:"filesystem_type"`
	MountPoint     string   `json:"mount_point"`
	Options        []string `json:"options"`
	TotalBytes     *uint64  `json:"total_bytes,omitempty"`
	UsedBytes      *uint64  `json:"used_bytes,omitempty"`
	AvailableBytes *uint64  `json:"available_bytes,omitempty"`
	UsePercent     *float64 `json:"use_percent,omitempty"`
	TotalInodes    *uint64  `json:"total_inodes,omitempty"`
	FreeInodes     *uint64  `json:"free_inodes,omitempty"`
}
