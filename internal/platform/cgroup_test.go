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

func TestReadCgroupMetricsV2(t *testing.T) {
	fsys := TestFS(fstest.MapFS{
		"sys/fs/cgroup/workload/memory.current": &fstest.MapFile{Data: []byte("4096\n")},
		"sys/fs/cgroup/workload/cpu.stat":       &fstest.MapFile{Data: []byte("usage_usec 12\n")},
		"sys/fs/cgroup/workload/io.stat":        &fstest.MapFile{Data: []byte("8:0 rbytes=10 wbytes=20\n")},
	})
	metrics, err := ReadCgroupMetrics(fsys, &CgroupInfo{Version: CgroupV2, Memberships: []CgroupMembership{{Path: "/workload"}}})
	require.NoError(t, err)
	assert.Equal(t, uint64(4096), *metrics.MemoryCurrentBytes)
	assert.Equal(t, uint64(12000), *metrics.CPUUsageNanoseconds)
	assert.Equal(t, uint64(10), *metrics.ReadBytes)
	assert.Equal(t, uint64(20), *metrics.WrittenBytes)
}

func TestReadCgroupMetricsV1(t *testing.T) {
	fsys := TestFS(fstest.MapFS{
		"sys/fs/cgroup/memory/docker/abc/memory.usage_in_bytes":          &fstest.MapFile{Data: []byte("2048\n")},
		"sys/fs/cgroup/cpuacct/docker/abc/cpuacct.usage":                 &fstest.MapFile{Data: []byte("9000\n")},
		"sys/fs/cgroup/blkio/docker/abc/blkio.throttle.io_service_bytes": &fstest.MapFile{Data: []byte("8:0 Read 12\n8:0 Write 34\n")},
	})
	metrics, err := ReadCgroupMetrics(fsys, &CgroupInfo{Version: CgroupV1, Memberships: []CgroupMembership{{Controllers: []string{"memory", "cpuacct", "blkio"}, Path: "/docker/abc"}}})
	require.NoError(t, err)
	assert.Equal(t, uint64(2048), *metrics.MemoryCurrentBytes)
	assert.Equal(t, uint64(9000), *metrics.CPUUsageNanoseconds)
	assert.Equal(t, uint64(12), *metrics.ReadBytes)
	assert.Equal(t, uint64(34), *metrics.WrittenBytes)
}

func TestContainerIDFromCgroupPath(t *testing.T) {
	id := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	assert.Equal(t, id, ContainerIDFromCgroupPath("/system.slice/docker-"+id+".scope"))
	assert.Equal(t, id, ContainerIDFromCgroupPath("/kubepods/cri-containerd-"+id+".scope"))
	assert.Empty(t, ContainerIDFromCgroupPath("/user.slice/not-a-container"))
}

func TestContainerRuntimeFromCgroupPath(t *testing.T) {
	id := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	assert.Equal(t, "docker", ContainerRuntimeFromCgroupPath("/system.slice/docker-"+id+".scope"))
	assert.Equal(t, "containerd", ContainerRuntimeFromCgroupPath("/kubepods/cri-containerd-"+id+".scope"))
	assert.Equal(t, "cri-o", ContainerRuntimeFromCgroupPath("/crio-"+id+".scope"))
	assert.Empty(t, ContainerRuntimeFromCgroupPath("/user.slice/not-a-container"))
}
