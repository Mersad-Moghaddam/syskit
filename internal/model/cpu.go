package model

import "time"

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
	Times         []CPUTime      `json:"times,omitempty"`
	CollectedAt   time.Time      `json:"collected_at"`
}

// CPUTime preserves raw cumulative Linux CPU counters and carries a service-
// derived utilization only when two snapshots have been compared.
type CPUTime struct {
	CPUID       string   `json:"cpu_id"`
	User        uint64   `json:"user"`
	Nice        uint64   `json:"nice"`
	System      uint64   `json:"system"`
	Idle        uint64   `json:"idle"`
	IOWait      uint64   `json:"iowait"`
	IRQ         uint64   `json:"irq"`
	SoftIRQ     uint64   `json:"softirq"`
	Steal       uint64   `json:"steal"`
	Guest       uint64   `json:"guest"`
	GuestNice   uint64   `json:"guest_nice"`
	Total       uint64   `json:"total"`
	Utilization *float64 `json:"utilization_percent,omitempty"`
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
