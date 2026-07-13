package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestEvaluateDiagnosticsReportsDocumentedThresholds(t *testing.T) {
	use := 96.0
	pressure := &model.MemoryPSI{FullAvg10: 12}
	findings := EvaluateDiagnostics(&model.MemoryInfo{Pressure: pressure, SwapTotalBytes: 100, SwapUsedBytes: 80}, &model.DiskInfo{Mounts: []model.MountInfo{{MountPoint: "/", UsePercent: &use}}})
	assert.Len(t, findings, 3)
	assert.Equal(t, "critical", findings[0].Severity)
}

func TestDiagnosticRejectsUnknownFilters(t *testing.T) {
	s := NewDiagnostic(memoryCollectorStub{}, diskCollectorStub{})
	_, err := s.Collect("network", "")
	assert.EqualError(t, err, "unknown diagnostics category \"network\"")
	_, err = s.Collect("", "high")
	assert.EqualError(t, err, "unknown diagnostics severity \"high\"")
}

type memoryCollectorStub struct{}

func (memoryCollectorStub) Collect() (*model.MemoryInfo, error) { return &model.MemoryInfo{}, nil }

type diskCollectorStub struct{}

func (diskCollectorStub) Collect() (*model.DiskInfo, error) { return &model.DiskInfo{}, nil }
