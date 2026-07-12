package model

type PortInfo struct {
	Sockets []Socket `json:"sockets"`
}
type Socket struct {
	Protocol      string `json:"protocol"`
	LocalAddress  string `json:"local_address"`
	LocalPort     uint16 `json:"local_port"`
	RemoteAddress string `json:"remote_address"`
	RemotePort    uint16 `json:"remote_port"`
	State         string `json:"state"`
	Inode         uint64 `json:"inode"`
}
