//go:build linux && integration

package disk

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

func TestCollectorReadsLiveMounts(t *testing.T) {
	info, err := NewCollector(platform.RealFS()).Collect()
	require.NoError(t, err)
	require.NotEmpty(t, info.Mounts)
}
