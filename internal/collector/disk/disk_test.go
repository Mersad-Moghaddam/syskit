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
func BenchmarkParseMountInfo(b *testing.B) {
	data := []byte("24 1 8:1 / / rw - ext4 /dev/sda1 rw\n")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseMountInfo(data)
	}
}
