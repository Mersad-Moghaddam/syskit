package model

import "time"

// DiskInfo describes mounted filesystem capacity from the current mount namespace.
type DiskInfo struct {
	Mounts      []MountInfo  `json:"mounts"`
	Devices     []DiskDevice `json:"devices"`
	CollectedAt time.Time    `json:"collected_at"`
}
type DiskDevice struct {
	Name                  string   `json:"name"`
	ReadOperations        uint64   `json:"read_operations"`
	WrittenOperations     uint64   `json:"written_operations"`
	ReadBytes             uint64   `json:"read_bytes"`
	WrittenBytes          uint64   `json:"written_bytes"`
	ReadBytesPerSecond    *float64 `json:"read_bytes_per_second,omitempty"`
	WrittenBytesPerSecond *float64 `json:"written_bytes_per_second,omitempty"`
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
