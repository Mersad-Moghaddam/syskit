package cli

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestDashboardModelRendersSnapshotAndError(t *testing.T) {
	m := dashboardModel{interval: time.Second}
	updated, _ := m.Update(dashboardData{snapshot: dashboardSnapshot{Hostname: "fixture", Uptime: 60, MemoryUsed: 40, MemoryTotal: 100, DiskUsed: 20, DiskTotal: 80, Interfaces: 2, TopProcess: "worker"}})
	view := updated.(dashboardModel).View()
	assert.Contains(t, view, "host: fixture")
	assert.Contains(t, view, "memory: 40 / 100 bytes")
	assert.Contains(t, view, "top process: worker")

	updated, _ = m.Update(dashboardData{err: errors.New("fixture failure")})
	assert.Contains(t, updated.(dashboardModel).View(), "collection error: fixture failure")
}

func TestDashboardQuitKey(t *testing.T) {
	m := dashboardModel{interval: time.Second}
	_, command := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.NotNil(t, command)
}

func TestDashboardCommandRejectsUnsafeInterval(t *testing.T) {
	cmd := newDashboardCmd(func() (dashboardSnapshot, error) { return dashboardSnapshot{}, nil })
	cmd.SetArgs([]string{"--interval", "10ms"})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "dashboard interval"))
}
