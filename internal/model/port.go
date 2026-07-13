package model

type PortInfo struct {
	Sockets []Socket `json:"sockets"`
}
type Socket struct {
	Protocol      string        `json:"protocol"`
	LocalAddress  string        `json:"local_address"`
	LocalPort     uint16        `json:"local_port"`
	RemoteAddress string        `json:"remote_address"`
	RemotePort    uint16        `json:"remote_port"`
	State         string        `json:"state"`
	RawState      string        `json:"raw_state"`
	Inode         uint64        `json:"inode"`
	Owners        []SocketOwner `json:"owners,omitempty"`
}

// SocketOwner is a best-effort association discovered from a process file
// descriptor. More than one process can own a socket after a fork.
type SocketOwner struct {
	PID     int    `json:"pid"`
	Command string `json:"command,omitempty"`
}
