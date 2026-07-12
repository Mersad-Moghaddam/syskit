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

func TestParseCacheSize(t *testing.T) {
	value, err := parseCacheSize("32K")
	assert.NoError(t, err)
	assert.Equal(t, uint64(32768), value)
	_, err = parseCacheSize("bad")
	assert.True(t, errors.Is(err, collector.ErrParse))
}

func TestParseCPUStat(t *testing.T) {
	times, err := ParseCPUStat([]byte("cpu 10 2 3 80 5 1 1 0 4 0\ncpu0 5 1 1 40 2 0 0 0 0 0\n"))
	assert.NoError(t, err)
	assert.Len(t, times, 2)
	assert.Equal(t, "all", times[0].CPUID)
	assert.Equal(t, uint64(102), times[0].Total)
	assert.Equal(t, uint64(4), times[0].Guest)
}

func BenchmarkParseCPUInfo(b *testing.B) {
	data := []byte("processor : 0\nmodel name : Test CPU\nphysical id : 0\ncore id : 0\nflags : sse avx\n")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ParseCPUInfo(data)
	}
}
