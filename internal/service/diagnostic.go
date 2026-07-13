package service

import (
	"fmt"
	"sort"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type Diagnostic struct {
	system  SystemCollector
	cpu     CPUCollector
	memory  MemoryCollector
	disk    DiskCollector
	process ProcessCollector
	network NetworkCollector
	port    PortCollector
}

func NewDiagnostic(system SystemCollector, cpu CPUCollector, memory MemoryCollector, disk DiskCollector, process ProcessCollector, network NetworkCollector, port PortCollector) *Diagnostic {
	return &Diagnostic{system, cpu, memory, disk, process, network, port}
}
func (s *Diagnostic) Collect(category, severity string) (*model.DiagnosticReport, error) {
	validCategories := map[string]bool{"": true, "cpu": true, "memory": true, "disk": true, "filesystem": true, "process": true, "network": true, "ports": true}
	if !validCategories[category] {
		return nil, fmt.Errorf("unknown diagnostics category %q", category)
	}
	if severity != "" && severity != "info" && severity != "warning" && severity != "critical" {
		return nil, fmt.Errorf("unknown diagnostics severity %q", severity)
	}
	system, err := s.system.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting system: %w", err)
	}
	cpu, err := s.cpu.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting CPU: %w", err)
	}
	memory, err := s.memory.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting memory: %w", err)
	}
	disk, err := s.disk.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting disk: %w", err)
	}
	processes, err := s.process.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting processes: %w", err)
	}
	network, err := s.network.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting network: %w", err)
	}
	ports, err := s.port.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting ports: %w", err)
	}
	findings := EvaluateDiagnostics(system, cpu, memory, disk, processes, network, ports)
	filtered := findings[:0]
	for _, f := range findings {
		if (category == "" || f.Category == category) && (severity == "" || f.Severity == severity) {
			filtered = append(filtered, f)
		}
	}
	return &model.DiagnosticReport{Findings: filtered}, nil
}
func EvaluateDiagnostics(system *model.SystemInfo, cpu *model.CPUInfo, memory *model.MemoryInfo, disk *model.DiskInfo, processes *model.ProcessList, network *model.NetworkInfo, ports *model.PortInfo) []model.DiagnosticFinding {
	var findings []model.DiagnosticFinding
	if cpu.LogicalCores > 0 && system.LoadAverage1 > float64(cpu.LogicalCores) {
		severity := "warning"
		if system.LoadAverage1 > float64(cpu.LogicalCores*2) {
			severity = "critical"
		}
		findings = append(findings, model.DiagnosticFinding{ID: "cpu-load", Severity: severity, Category: "cpu", Summary: "Load exceeds logical CPU capacity", Evidence: fmt.Sprintf("load1 %.2f across %d logical CPUs", system.LoadAverage1, cpu.LogicalCores), Sources: []string{"/proc/loadavg", "/proc/cpuinfo"}, Recommendation: "inspect runnable processes and CPU utilization"})
	}
	if memory.Pressure != nil && memory.Pressure.FullAvg10 >= 10 {
		findings = append(findings, model.DiagnosticFinding{ID: "memory-pressure", Severity: "warning", Category: "memory", Summary: "Memory pressure is elevated", Evidence: fmt.Sprintf("full PSI avg10 is %.2f%%", memory.Pressure.FullAvg10), Sources: []string{"/proc/pressure/memory"}, Recommendation: "inspect memory-heavy processes and swap activity"})
	}
	if memory.SwapTotalBytes > 0 && memory.SwapUsedBytes*100/memory.SwapTotalBytes >= 80 {
		findings = append(findings, model.DiagnosticFinding{ID: "swap-usage", Severity: "warning", Category: "memory", Summary: "Swap usage is high", Evidence: fmt.Sprintf("%d of %d bytes used", memory.SwapUsedBytes, memory.SwapTotalBytes), Sources: []string{"/proc/meminfo"}, Recommendation: "inspect memory pressure and working-set size"})
	}
	for _, mount := range disk.Mounts {
		if mount.UsePercent != nil && *mount.UsePercent >= 85 {
			severity := "warning"
			if *mount.UsePercent >= 95 {
				severity = "critical"
			}
			findings = append(findings, model.DiagnosticFinding{ID: "filesystem-capacity-" + mount.MountPoint, Severity: severity, Category: "filesystem", Summary: "Filesystem capacity is high", Evidence: fmt.Sprintf("%s is %.1f%% used", mount.MountPoint, *mount.UsePercent), Sources: []string{"/proc/self/mountinfo", "statfs"}, Recommendation: "free space or extend the filesystem"})
		}
	}
	findings = append(findings, model.DiagnosticFinding{ID: "disk-saturation-unavailable", Severity: "info", Category: "disk", Summary: "Disk saturation check is unavailable", Evidence: "device busy-time utilization is not collected", Sources: []string{"/proc/diskstats"}, Recommendation: "use sampled disk throughput and latency tooling when saturation is suspected"})
	if processes.TotalMemoryBytes > 0 {
		for _, process := range processes.Processes {
			percent := float64(process.ResidentBytes) * 100 / float64(processes.TotalMemoryBytes)
			if percent >= 50 {
				findings = append(findings, model.DiagnosticFinding{ID: fmt.Sprintf("process-memory-%d", process.PID), Severity: "warning", Category: "process", Summary: "One process holds a large memory share", Evidence: fmt.Sprintf("PID %d (%s) uses %.1f%% of memory", process.PID, process.Command, percent), Sources: []string{"/proc/<pid>/stat", "/proc/meminfo"}, Recommendation: "inspect the process working set and workload"})
			}
		}
	}
	var networkFaults uint64
	for _, iface := range network.Interfaces {
		networkFaults += iface.RXErrors + iface.TXErrors + iface.RXDrops + iface.TXDrops
	}
	if networkFaults > 0 {
		findings = append(findings, model.DiagnosticFinding{ID: "network-errors-drops", Severity: "warning", Category: "network", Summary: "Network interfaces report errors or drops", Evidence: fmt.Sprintf("%d cumulative errors and drops", networkFaults), Sources: []string{"/proc/net/dev"}, Recommendation: "inspect interface counters, links, and queue pressure"})
	}
	wildcardListeners := 0
	for _, socket := range ports.Sockets {
		if socket.State == "LISTEN" && (socket.LocalAddress == "0.0.0.0" || socket.LocalAddress == "::") {
			wildcardListeners++
		}
	}
	if wildcardListeners > 0 {
		findings = append(findings, model.DiagnosticFinding{ID: "ports-wildcard-listeners", Severity: "info", Category: "ports", Summary: "Services listen on wildcard addresses", Evidence: fmt.Sprintf("%d wildcard listening sockets", wildcardListeners), Sources: []string{"/proc/net/tcp", "/proc/net/tcp6"}, Recommendation: "confirm each exposed service is intentional and firewalled"})
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	return findings
}
