package command

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/golden"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type fakeSystemService struct {
	info *model.SystemInfo
	err  error
}

func (s fakeSystemService) Collect() (*model.SystemInfo, error) { return s.info, s.err }

func TestSystemCommandOutput(t *testing.T) {
	info := &model.SystemInfo{Hostname: "fixture-host", OSName: "Fixture Linux", OSVersion: "1.0", KernelRelease: "6.12", KernelVersion: "#1", Architecture: "amd64", UptimeSeconds: 93784, BootTime: time.Date(2026, 7, 11, 9, 56, 56, 0, time.UTC), LoadAverage1: .42, LoadAverage5: .35, LoadAverage15: .30}
	for _, tt := range []struct{ name, format, golden string }{{"table", "table", "system_table.golden"}, {"json", "json", "system_json.golden"}} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewSystemCmd(fakeSystemService{info: info}, SystemOptions{Format: func() string { return tt.format }, NoHeader: func() bool { return false }})
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetArgs([]string{})
			require.NoError(t, cmd.Execute())
			golden.Assert(t, out.Bytes(), tt.golden)
		})
	}
}

func TestFormatUptime(t *testing.T) {
	assert.Equal(t, "1d 02h 03m", formatUptime(float64(26*time.Hour+3*time.Minute)/float64(time.Second)))
}
