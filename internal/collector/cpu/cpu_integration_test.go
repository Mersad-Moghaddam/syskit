//go:build linux && integration

package cpu

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

func TestCollectorReadsLiveCPUInfo(t *testing.T) {
	info, err := NewCollector(platform.RealFS()).Collect()
	require.NoError(t, err)
	require.Positive(t, info.LogicalCores)
}
