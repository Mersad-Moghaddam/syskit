//go:build linux && integration

package network

import (
	"testing"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
	"github.com/stretchr/testify/require"
)

func TestCollectorReadsLiveNetworkInterfaces(t *testing.T) {
	info, err := NewCollector(platform.RealFS()).Collect()
	require.NoError(t, err)
	require.NotEmpty(t, info.Interfaces)
	metadataFound := false
	for _, iface := range info.Interfaces {
		require.NotEmpty(t, iface.Name)
		metadataFound = metadataFound || iface.MTU != nil
	}
	require.True(t, metadataFound)
}
