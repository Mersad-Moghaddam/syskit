//go:build linux && integration

package port

import (
	"testing"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
	"github.com/stretchr/testify/require"
)

func TestCollectorReadsLiveSocketTables(t *testing.T) {
	info, err := NewCollector(platform.RealFS()).Collect()
	require.NoError(t, err)
	require.NotEmpty(t, info.Sockets)
}
