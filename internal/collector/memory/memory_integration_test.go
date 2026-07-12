//go:build linux && integration

package memory

import (
	"testing"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
	"github.com/stretchr/testify/require"
)

func TestCollectorReadsLiveMemory(t *testing.T) {
	info, err := NewCollector(platform.RealFS()).Collect()
	require.NoError(t, err)
	require.Positive(t, info.TotalBytes)
}
