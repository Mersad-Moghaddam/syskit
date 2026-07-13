package model

// ContainerInfo is the cgroup-derived view of a container observed through its
// associated processes. Runtime metadata is intentionally only a hint because
// SysKit does not query runtime sockets.
type ContainerInfo struct {
	ID      string            `json:"id"`
	Runtime string            `json:"runtime,omitempty"`
	PIDs    int               `json:"pids"`
	Metrics *ContainerMetrics `json:"metrics,omitempty"`
}

// ContainerMetrics contains cgroup counters normalized to bytes and
// nanoseconds. Nil fields indicate unavailable controller data.
type ContainerMetrics struct {
	MemoryCurrentBytes  *uint64 `json:"memory_current_bytes,omitempty"`
	CPUUsageNanoseconds *uint64 `json:"cpu_usage_nanoseconds,omitempty"`
	ReadBytes           *uint64 `json:"read_bytes,omitempty"`
	WrittenBytes        *uint64 `json:"written_bytes,omitempty"`
}

// ContainerList is the collection returned by container inspection commands.
type ContainerList struct {
	Containers []ContainerInfo `json:"containers"`
}

// ContainerDetail expands one cgroup-derived container with the processes
// currently associated with its recognized ID.
type ContainerDetail struct {
	ContainerInfo
	Processes []Process `json:"processes"`
}
