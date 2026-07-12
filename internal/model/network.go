package model

import "time"

type NetworkInfo struct {
	Interfaces  []NetworkInterface `json:"interfaces"`
	Routes      []Route            `json:"routes,omitempty"`
	Nameservers []string           `json:"nameservers,omitempty"`
	CollectedAt time.Time          `json:"collected_at"`
}
type Route struct {
	Interface   string `json:"interface"`
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Default     bool   `json:"default"`
}
type NetworkInterface struct {
	Name             string   `json:"name"`
	RXBytes          uint64   `json:"rx_bytes"`
	TXBytes          uint64   `json:"tx_bytes"`
	RXPackets        uint64   `json:"rx_packets"`
	TXPackets        uint64   `json:"tx_packets"`
	RXErrors         uint64   `json:"rx_errors"`
	TXErrors         uint64   `json:"tx_errors"`
	RXDrops          uint64   `json:"rx_drops"`
	TXDrops          uint64   `json:"tx_drops"`
	RXBytesPerSecond *float64 `json:"rx_bytes_per_second,omitempty"`
	TXBytesPerSecond *float64 `json:"tx_bytes_per_second,omitempty"`
}
