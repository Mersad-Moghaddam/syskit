package model

import "time"

// SystemInfo is a point-in-time summary of host identity and load. Numeric
// values retain their base units so JSON remains suitable for automation.
type SystemInfo struct {
	Hostname      string    `json:"hostname"`
	OSName        string    `json:"os_name,omitempty"`
	OSVersion     string    `json:"os_version,omitempty"`
	KernelRelease string    `json:"kernel_release"`
	KernelVersion string    `json:"kernel_version"`
	Architecture  string    `json:"architecture"`
	UptimeSeconds float64   `json:"uptime_seconds"`
	BootTime      time.Time `json:"boot_time"`
	LoadAverage1  float64   `json:"load_average_1"`
	LoadAverage5  float64   `json:"load_average_5"`
	LoadAverage15 float64   `json:"load_average_15"`
}
