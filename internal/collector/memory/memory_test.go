package memory

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
)

func TestParseMemInfo(t *testing.T) {
	info, err := ParseMemInfo([]byte("MemTotal: 100 kB\nMemFree: 10 kB\nMemAvailable: 60 kB\nBuffers: 2 kB\nCached: 3 kB\nSReclaimable: 1 kB\nSwapTotal: 8 kB\nSwapFree: 3 kB\n"))
	assert.NoError(t, err)
	assert.Equal(t, uint64(102400), info.TotalBytes)
	assert.Equal(t, uint64(40960), *info.UsedBytes)
	assert.Equal(t, uint64(5120), info.SwapUsedBytes)
}
func TestParseMemInfoMissingTotal(t *testing.T) {
	_, err := ParseMemInfo([]byte("MemFree: 1 kB\n"))
	assert.True(t, errors.Is(err, collector.ErrFieldMissing))
}
func TestParsePSI(t *testing.T) {
	psi, err := ParsePSI([]byte("some avg10=1.20 avg60=0.30 avg300=0 total=1\nfull avg10=0.20 avg60=0.10 avg300=0 total=1\n"))
	assert.NoError(t, err)
	assert.Equal(t, 1.2, psi.SomeAvg10)
	assert.Equal(t, .1, psi.FullAvg60)
}
func BenchmarkParseMemInfo(b *testing.B) {
	data := []byte("MemTotal: 100 kB\nMemFree: 10 kB\n")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseMemInfo(data)
	}
}
