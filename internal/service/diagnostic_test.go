package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestEvaluateDiagnosticsReportsDocumentedThresholds(t *testing.T) {
	use := 96.0
	pressure := &model.MemoryPSI{FullAvg10: 12}
	findings := EvaluateDiagnostics(&model.SystemInfo{}, &model.CPUInfo{LogicalCores: 1}, &model.MemoryInfo{Pressure: pressure, SwapTotalBytes: 100, SwapUsedBytes: 80}, &model.DiskInfo{Mounts: []model.MountInfo{{MountPoint: "/", UsePercent: &use}}}, &model.ProcessList{TotalMemoryBytes: 1}, &model.NetworkInfo{Interfaces: []model.NetworkInterface{{Name: "lo"}}}, &model.PortInfo{})
	assert.Len(t, findings, 4)
	assert.Contains(t, findings, model.DiagnosticFinding{ID: "filesystem-capacity-/", Severity: "critical", Category: "filesystem", Summary: "Filesystem capacity is high", Evidence: "/ is 96.0% used", Sources: []string{"/proc/self/mountinfo", "statfs"}, Recommendation: "free space or extend the filesystem"})
}

func TestDiagnosticRejectsUnknownFilters(t *testing.T) {
	s := NewDiagnostic(nil, nil, nil, nil, nil, nil, nil)
	_, err := s.Collect("storage", "")
	assert.EqualError(t, err, "unknown diagnostics category \"storage\"")
	_, err = s.Collect("", "high")
	assert.EqualError(t, err, "unknown diagnostics severity \"high\"")
}

func TestEvaluateDiagnosticsCoversHostSignals(t *testing.T) {
	use := 10.0
	findings := EvaluateDiagnostics(
		&model.SystemInfo{LoadAverage1: 5}, &model.CPUInfo{LogicalCores: 2},
		&model.MemoryInfo{Pressure: &model.MemoryPSI{}}, &model.DiskInfo{Mounts: []model.MountInfo{{MountPoint: "/", UsePercent: &use}}},
		&model.ProcessList{TotalMemoryBytes: 100, Processes: []model.Process{{PID: 7, Command: "worker", ResidentBytes: 60}}},
		&model.NetworkInfo{Interfaces: []model.NetworkInterface{{Name: "eth0", RXErrors: 1}}},
		&model.PortInfo{Sockets: []model.Socket{{State: "LISTEN", LocalAddress: "0.0.0.0"}}},
	)
	assert.Len(t, findings, 5)
	ids := make([]string, 0, len(findings))
	for _, finding := range findings {
		ids = append(ids, finding.ID)
	}
	assert.Contains(t, ids, "cpu-load")
	assert.Contains(t, ids, "process-memory-7")
	assert.Contains(t, ids, "network-errors-drops")
	assert.Contains(t, ids, "ports-wildcard-listeners")
}

func TestEvaluateDiagnosticsReportsUnavailableSignals(t *testing.T) {
	findings := EvaluateDiagnostics(&model.SystemInfo{}, &model.CPUInfo{}, &model.MemoryInfo{}, &model.DiskInfo{}, &model.ProcessList{}, &model.NetworkInfo{}, &model.PortInfo{})
	ids := map[string]bool{}
	for _, finding := range findings {
		ids[finding.ID] = true
		assert.NotEmpty(t, finding.Evidence)
		assert.NotEmpty(t, finding.Sources)
		assert.NotEmpty(t, finding.Recommendation)
	}
	for _, id := range []string{"cpu-load-unavailable", "memory-pressure-unavailable", "filesystem-capacity-unavailable", "disk-saturation-unavailable", "process-memory-unavailable", "network-errors-unavailable"} {
		assert.True(t, ids[id], "missing unavailable finding %s", id)
	}
}

func TestDiagnosticCategoryCollectsOnlyRequiredDomain(t *testing.T) {
	service := NewDiagnostic(nil, nil, diagnosticMemoryCollector{}, nil, nil, nil, nil)
	report, err := service.Collect("memory", "info")
	assert.NoError(t, err)
	assert.Len(t, report.Findings, 1)
	assert.Equal(t, "memory-pressure-unavailable", report.Findings[0].ID)
}

type diagnosticMemoryCollector struct{}

func (diagnosticMemoryCollector) Collect() (*model.MemoryInfo, error) {
	return &model.MemoryInfo{}, nil
}
