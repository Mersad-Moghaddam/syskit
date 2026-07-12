package disk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMountInfoUnescapesMountPoint(t *testing.T) {
	mounts, err := ParseMountInfo([]byte("24 1 8:1 / /mnt/my\\040disk rw,relatime - ext4 /dev/sda1 rw\n"))
	assert.NoError(t, err)
	assert.Len(t, mounts, 1)
	assert.Equal(t, "/mnt/my disk", mounts[0].MountPoint)
	assert.Equal(t, "ext4", mounts[0].FilesystemType)
	assert.Equal(t, []string{"rw", "relatime"}, mounts[0].Options)
}

func TestParseDiskStats(t *testing.T) {
	devices, err := ParseDiskStats([]byte("8 0 sda 10 0 20 0 30 0 40 0 0 0 0\n"))
	assert.NoError(t, err)
	assert.Len(t, devices, 1)
	assert.Equal(t, uint64(20*512), devices[0].ReadBytes)
	assert.Equal(t, uint64(40*512), devices[0].WrittenBytes)
}
func BenchmarkParseMountInfo(b *testing.B) {
	data := []byte("24 1 8:1 / / rw - ext4 /dev/sda1 rw\n")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseMountInfo(data)
	}
}
