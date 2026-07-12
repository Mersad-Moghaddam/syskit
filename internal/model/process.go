package model

// Process is a procfs process snapshot. CPU counters are raw clock ticks and
// resident memory stays in bytes so services can later derive percentages.
type Process struct {
	PID           int    `json:"pid"`
	PPID          int    `json:"ppid"`
	UID           uint64 `json:"uid"`
	State         string `json:"state"`
	Command       string `json:"command"`
	CPUTime       uint64 `json:"cpu_time"`
	ResidentBytes uint64 `json:"resident_bytes"`
	Threads       uint64 `json:"threads"`
}
type ProcessList struct {
	Processes []Process `json:"processes"`
}
