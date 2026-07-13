package service

import (
	"fmt"
	"sort"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type Diagnostic struct {
	memory MemoryCollector
	disk   DiskCollector
}

func NewDiagnostic(memory MemoryCollector, disk DiskCollector) *Diagnostic {
	return &Diagnostic{memory, disk}
}
func (s *Diagnostic) Collect(category, severity string) (*model.DiagnosticReport, error) {
	if category != "" && category != "memory" && category != "filesystem" {
		return nil, fmt.Errorf("unknown diagnostics category %q", category)
	}
	if severity != "" && severity != "info" && severity != "warning" && severity != "critical" {
		return nil, fmt.Errorf("unknown diagnostics severity %q", severity)
	}
	memory, err := s.memory.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting memory: %w", err)
	}
	disk, err := s.disk.Collect()
	if err != nil {
		return nil, fmt.Errorf("collecting disk: %w", err)
	}
	findings := EvaluateDiagnostics(memory, disk)
	filtered := findings[:0]
	for _, f := range findings {
		if (category == "" || f.Category == category) && (severity == "" || f.Severity == severity) {
			filtered = append(filtered, f)
		}
	}
	return &model.DiagnosticReport{Findings: filtered}, nil
}
func EvaluateDiagnostics(memory *model.MemoryInfo, disk *model.DiskInfo) []model.DiagnosticFinding {
	var findings []model.DiagnosticFinding
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
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	return findings
}
