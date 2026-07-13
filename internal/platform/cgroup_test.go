package platform

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectCgroupV2(t *testing.T) {
	fsys := TestFS(fstest.MapFS{
		"proc/self/cgroup":                 &fstest.MapFile{Data: []byte("0::/user.slice/test.scope\n")},
		"sys/fs/cgroup/cgroup.controllers": &fstest.MapFile{Data: []byte("cpu memory io\n")},
	})
	info, err := DetectCgroup(fsys, "proc/self/cgroup")
	require.NoError(t, err)
	assert.Equal(t, CgroupV2, info.Version)
	assert.Equal(t, "/user.slice/test.scope", info.Memberships[0].Path)
}

func TestDetectCgroupV1(t *testing.T) {
	fsys := TestFS(fstest.MapFS{"proc/self/cgroup": &fstest.MapFile{Data: []byte("5:memory,cpu:/docker/abc\n")}})
	info, err := DetectCgroup(fsys, "proc/self/cgroup")
	require.NoError(t, err)
	assert.Equal(t, CgroupV1, info.Version)
	assert.Equal(t, []string{"memory", "cpu"}, info.Memberships[0].Controllers)
}
