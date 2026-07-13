package model

// Process is a procfs process snapshot. CPU counters are raw clock ticks and
// resident memory stays in bytes so services can later derive percentages.
type Process struct {
	PID              int      `json:"pid"`
	PPID             int      `json:"ppid"`
	UID              uint64   `json:"uid"`
	User             string   `json:"user,omitempty"`
	State            string   `json:"state"`
	Command          string   `json:"command"`
	CPUTime          uint64   `json:"cpu_time"`
	CPUPercent       *float64 `json:"cpu_percent,omitempty"`
	StartTimeTicks   uint64   `json:"start_time_ticks"`
	ResidentBytes    uint64   `json:"resident_bytes"`
	MemoryPercent    *float64 `json:"memory_percent,omitempty"`
	Threads          uint64   `json:"threads"`
	ContainerID      string   `json:"container_id,omitempty"`
	ContainerRuntime string   `json:"container_runtime,omitempty"`
}
type ProcessList struct {
	Processes        []Process `json:"processes"`
	CPUTimeTotal     uint64    `json:"cpu_time_total"`
	TotalMemoryBytes uint64    `json:"total_memory_bytes"`
	Partial          bool      `json:"partial"`
}
