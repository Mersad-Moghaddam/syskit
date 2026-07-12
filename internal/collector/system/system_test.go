package system

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

func TestCollectorCollectsFixture(t *testing.T) {
	clock := func() time.Time { return time.Date(2026, 7, 12, 12, 0, 0, 0, time.UTC) }
	c := NewCollectorWithClock(platform.TestFS(os.DirFS("testdata/fixtures/standard")), clock)

	info, err := c.Collect()
	require.NoError(t, err)
	assert.Equal(t, "fixture-host", info.Hostname)
	assert.Equal(t, "Fixture Linux", info.OSName)
	assert.Equal(t, "1.0", info.OSVersion)
	assert.Equal(t, "6.12.0-fixture", info.KernelRelease)
	assert.Equal(t, 93784.5, info.UptimeSeconds)
	assert.Equal(t, time.Date(2026, 7, 11, 9, 56, 55, 500000000, time.UTC), info.BootTime)
	assert.Equal(t, 0.42, info.LoadAverage1)
	assert.Equal(t, 0.35, info.LoadAverage5)
	assert.Equal(t, 0.3, info.LoadAverage15)
}

func TestCollectorOptionalOSReleaseMissing(t *testing.T) {
	c := NewCollectorWithClock(platform.TestFS(os.DirFS("testdata/fixtures/no-os-release")), func() time.Time { return time.Unix(100, 0) })
	info, err := c.Collect()
	require.NoError(t, err)
	assert.Empty(t, info.OSName)
	assert.Empty(t, info.OSVersion)
}

func TestParsers(t *testing.T) {
	tests := []struct {
		name string
		data string
		want error
	}{
		{name: "uptime missing", data: "", want: collector.ErrFieldMissing},
		{name: "uptime invalid", data: "bad 2", want: collector.ErrParse},
		{name: "load missing", data: "0.1 0.2", want: collector.ErrFieldMissing},
		{name: "load invalid", data: "0.1 nope 0.3", want: collector.ErrParse},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name[:4] == "upti" {
				_, err := ParseUptime([]byte(tt.data))
				assert.True(t, errors.Is(err, tt.want))
				return
			}
			_, _, _, err := ParseLoadAverage([]byte(tt.data))
			assert.True(t, errors.Is(err, tt.want))
		})
	}
}

func TestParseOSRelease(t *testing.T) {
	name, version := ParseOSRelease([]byte("NAME=Fixture Linux\nVERSION_ID=\"1.0\"\n"))
	assert.Equal(t, "Fixture Linux", name)
	assert.Equal(t, "1.0", version)
}

func BenchmarkParseLoadAverage(b *testing.B) {
	data := []byte("0.42 0.35 0.30 1/123 456\n")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = ParseLoadAverage(data)
	}
}
