package model

// MemoryInfo is a byte-normalized memory snapshot. Optional fields are nil
// when the running kernel does not expose the corresponding interface.
type MemoryInfo struct {
	TotalBytes     uint64     `json:"total_bytes"`
	UsedBytes      *uint64    `json:"used_bytes,omitempty"`
	AvailableBytes *uint64    `json:"available_bytes,omitempty"`
	FreeBytes      uint64     `json:"free_bytes"`
	BuffersBytes   uint64     `json:"buffers_bytes"`
	CacheBytes     uint64     `json:"cache_bytes"`
	SwapTotalBytes uint64     `json:"swap_total_bytes"`
	SwapUsedBytes  uint64     `json:"swap_used_bytes"`
	SwapFreeBytes  uint64     `json:"swap_free_bytes"`
	Pressure       *MemoryPSI `json:"pressure,omitempty"`
}

// MemoryPSI contains kernel-provided memory stall percentages.
type MemoryPSI struct {
	SomeAvg10 float64 `json:"some_avg10"`
	SomeAvg60 float64 `json:"some_avg60"`
	FullAvg10 float64 `json:"full_avg10"`
	FullAvg60 float64 `json:"full_avg60"`
}
