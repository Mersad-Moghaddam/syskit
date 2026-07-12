//go:build linux && integration

package system

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

func TestCollectorReadsLiveHost(t *testing.T) {
	info, err := NewCollector(platform.RealFS()).Collect()
	require.NoError(t, err)
	require.NotEmpty(t, info.Hostname)
	require.NotEmpty(t, info.KernelRelease)
	require.GreaterOrEqual(t, info.UptimeSeconds, float64(0))
}
