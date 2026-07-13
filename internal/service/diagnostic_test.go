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
