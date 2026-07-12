package model

// CPUInfo is the static CPU topology and identity snapshot. Optional topology
// values remain nil when a virtualized or restricted kernel does not expose
// them; zero is never used to mean unavailable.
type CPUInfo struct {
	LogicalCores  int            `json:"logical_cores"`
	PhysicalCores *int           `json:"physical_cores,omitempty"`
	Sockets       *int           `json:"sockets,omitempty"`
	Model         string         `json:"model"`
	Architecture  string         `json:"architecture"`
	Flags         []string       `json:"flags,omitempty"`
	Caches        []CPUCache     `json:"caches,omitempty"`
	Frequencies   []CPUFrequency `json:"frequencies,omitempty"`
}

// CPUCache is a cache description reported by the kernel for CPU zero. Cache
// entries are shared across logical CPUs, so reading one CPU avoids duplicates.
type CPUCache struct {
	Level     int    `json:"level"`
	Type      string `json:"type"`
	SizeBytes uint64 `json:"size_bytes"`
}

// CPUFrequency records optional cpufreq values for one logical CPU. Linux
// exposes these values in kHz; SysKit normalizes them to MHz for presentation.
type CPUFrequency struct {
	CPUID      int      `json:"cpu_id"`
	CurrentMHz *float64 `json:"current_mhz,omitempty"`
	MinimumMHz *float64 `json:"minimum_mhz,omitempty"`
	MaximumMHz *float64 `json:"maximum_mhz,omitempty"`
}
