package cli

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestDashboardModelRendersSnapshotAndError(t *testing.T) {
	m := dashboardModel{interval: time.Second, panel: overviewPanel}
	updated, _ := m.Update(dashboardData{snapshot: dashboardSnapshot{Hostname: "fixture", Uptime: 60, MemoryUsed: 40, MemoryTotal: 100, DiskUsed: 20, DiskTotal: 80, Interfaces: 2, TopProcess: "worker"}})
	view := updated.(dashboardModel).View()
	assert.Contains(t, view, "host: fixture")
	assert.Contains(t, view, "memory: 40 B / 100 B")
	assert.Contains(t, view, "OVERVIEW")

	updated, _ = m.Update(dashboardData{err: errors.New("fixture failure")})
	assert.Contains(t, updated.(dashboardModel).View(), "COLLECTION ERROR")
	assert.Contains(t, updated.(dashboardModel).View(), "fixture failure")
}

func TestDashboardDerivesCPUAndNetworkRates(t *testing.T) {
	beforeCPU := &model.CPUInfo{Times: []model.CPUTime{{CPUID: "all", Total: 100, Idle: 40, IOWait: 10}}}
	afterCPU := &model.CPUInfo{Times: []model.CPUTime{{CPUID: "all", Total: 200, Idle: 80, IOWait: 20}}}
	utilization := dashboardCPUUtilization(beforeCPU, afterCPU)
	assert.NotNil(t, utilization)
	assert.Equal(t, 50.0, *utilization)

	beforeNetwork := &model.NetworkInfo{CollectedAt: time.Unix(1, 0), Interfaces: []model.NetworkInterface{{Name: "eth0", RXBytes: 10, TXBytes: 20}}}
	afterNetwork := &model.NetworkInfo{CollectedAt: time.Unix(3, 0), Interfaces: []model.NetworkInterface{{Name: "eth0", RXBytes: 30, TXBytes: 50}}}
	rx, tx := dashboardNetworkRates(beforeNetwork, afterNetwork)
	assert.NotNil(t, rx)
	assert.NotNil(t, tx)
	assert.Equal(t, 10.0, *rx)
	assert.Equal(t, 15.0, *tx)
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

func TestDashboardUsesCompactLayoutAtMediumWidths(t *testing.T) {
	m := dashboardModel{interval: time.Second, panel: overviewPanel, width: 60, height: 16, snapshot: dashboardSnapshot{Hostname: "fixture", MemoryTotal: 1024, DiskTotal: 2048}}
	view := m.View()
	assert.Contains(t, view, "host: fixture")
	assert.NotContains(t, view, "HOST + CPU")
	for _, line := range strings.Split(view, "\n") {
		assert.LessOrEqual(t, lipgloss.Width(line), 68)
	}
}

func TestDashboardFormatsHumanReadableCapacity(t *testing.T) {
	assert.Equal(t, "1.0 GiB", formatTUIBytes(1<<30))
	assert.Equal(t, "2.0 MiB/s", formatTUIRate(2<<20))
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
