package model

// ContainerInfo is the cgroup-derived view of a container observed through its
// associated processes. Runtime metadata is intentionally only a hint because
// SysKit does not query runtime sockets.
type ContainerInfo struct {
	ID      string `json:"id"`
	Runtime string `json:"runtime,omitempty"`
	PIDs    int    `json:"pids"`
}

// ContainerList is the collection returned by container inspection commands.
type ContainerList struct {
	Containers []ContainerInfo `json:"containers"`
}
