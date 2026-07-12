package cpu

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
)

func TestParseCPUInfo(t *testing.T) {
	data := []byte("processor : 0\nmodel name : Test CPU\nphysical id : 0\ncore id : 0\nflags : sse avx\n\nprocessor : 1\nmodel name : Test CPU\nphysical id : 0\ncore id : 1\nflags : sse avx2\n")
	info, ids, err := ParseCPUInfo(data)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1}, ids)
	assert.Equal(t, 2, info.LogicalCores)
	assert.EqualValues(t, 2, *info.PhysicalCores)
	assert.EqualValues(t, 1, *info.Sockets)
	assert.Equal(t, []string{"avx", "avx2", "sse"}, info.Flags)
}

func TestParseCPUInfoMalformed(t *testing.T) {
	_, _, err := ParseCPUInfo([]byte("model name : missing\n"))
	assert.True(t, errors.Is(err, collector.ErrFieldMissing))
}

func BenchmarkParseCPUInfo(b *testing.B) {
	data := []byte("processor : 0\nmodel name : Test CPU\nphysical id : 0\ncore id : 0\nflags : sse avx\n")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ParseCPUInfo(data)
	}
}
