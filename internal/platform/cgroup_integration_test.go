//go:build linux && integration

package platform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectCgroupOnLiveHost(t *testing.T) {
	info, err := DetectCgroup(RealFS(), "proc/self/cgroup")
	require.NoError(t, err)
	require.NotEqual(t, CgroupUnknown, info.Version)
	require.NotEmpty(t, info.Memberships)
}
