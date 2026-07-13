package cli

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestDashboardModelRendersSnapshotAndError(t *testing.T) {
	m := dashboardModel{interval: time.Second, panel: overviewPanel}
	updated, _ := m.Update(dashboardData{snapshot: dashboardSnapshot{Hostname: "fixture", Uptime: 60, MemoryUsed: 40, MemoryTotal: 100, DiskUsed: 20, DiskTotal: 80, Interfaces: 2, TopProcess: "worker"}})
	view := updated.(dashboardModel).View()
	assert.Contains(t, view, "host: fixture")
	assert.Contains(t, view, "memory: 40 / 100 bytes")
	assert.Contains(t, view, "overview")

	updated, _ = m.Update(dashboardData{err: errors.New("fixture failure")})
	assert.Contains(t, updated.(dashboardModel).View(), "collection error: fixture failure")
}

func TestDashboardNavigatesPanels(t *testing.T) {
	m := dashboardModel{interval: time.Second, panel: overviewPanel, snapshot: dashboardSnapshot{TopProcess: "worker"}}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	processes := updated.(dashboardModel)
	assert.Equal(t, processesPanel, processes.panel)
	assert.Contains(t, processes.View(), "top process: worker")
}

func TestDashboardHandlesSmallTerminal(t *testing.T) {
	m := dashboardModel{interval: time.Second}
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	assert.Contains(t, updated.(dashboardModel).View(), "terminal is too small")
}

func TestDashboardQuitKey(t *testing.T) {
	m := dashboardModel{interval: time.Second}
	_, command := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.NotNil(t, command)
}

func TestDashboardSkipsOverlappingCollection(t *testing.T) {
	m := dashboardModel{interval: time.Second, fetching: true}
	updated, command := m.Update(dashboardTick{})
	assert.True(t, updated.(dashboardModel).fetching)
	assert.NotNil(t, command)

	updated, _ = updated.(dashboardModel).Update(dashboardData{})
	assert.False(t, updated.(dashboardModel).fetching)
	updated, command = updated.(dashboardModel).Update(dashboardTick{})
	assert.True(t, updated.(dashboardModel).fetching)
	assert.NotNil(t, command)
}

func TestDashboardCommandRejectsUnsafeInterval(t *testing.T) {
	cmd := newDashboardCmd(func() (dashboardSnapshot, error) { return dashboardSnapshot{}, nil })
	cmd.SetArgs([]string{"--interval", "10ms"})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "dashboard interval"))
}

func TestInteractiveTerminalRejectsRegularFile(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "output")
	assert.NoError(t, err)
	t.Cleanup(func() { _ = file.Close() })
	assert.False(t, isInteractiveTerminal(file))
}
