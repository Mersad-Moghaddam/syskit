package cli

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

const (
	minDashboardInterval = 500 * time.Millisecond
	maxDashboardInterval = time.Minute
	overviewPanel        = "overview"
	processesPanel       = "processes"
)

type dashboardSnapshot struct {
	Hostname    string
	Uptime      float64
	MemoryUsed  uint64
	MemoryTotal uint64
	SwapUsed    uint64
	SwapTotal   uint64
	CPUPercent  *float64
	Interfaces  int
	NetworkRX   *float64
	NetworkTX   *float64
	DiskUsed    uint64
	DiskTotal   uint64
	TopProcess  string
}

type dashboardProvider func() (dashboardSnapshot, error)

type dashboardModel struct {
	provider dashboardProvider
	interval time.Duration
	snapshot dashboardSnapshot
	err      error
	panel    string
	width    int
	height   int
	fetching bool
}

type dashboardTick struct{}
type dashboardData struct {
	snapshot dashboardSnapshot
	err      error
}

func newDashboardCmd(provider dashboardProvider) *cobra.Command {
	var interval time.Duration
	var panel string
	cmd := &cobra.Command{Use: "dashboard", Short: "Show a live system dashboard", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		if interval < minDashboardInterval || interval > maxDashboardInterval {
			return fmt.Errorf("dashboard interval must be between %s and %s", minDashboardInterval, maxDashboardInterval)
		}
		if panel != overviewPanel && panel != processesPanel {
			return fmt.Errorf("dashboard panel must be %q or %q", overviewPanel, processesPanel)
		}
		if !isInteractiveTerminal(os.Stdout) {
			return fmt.Errorf("dashboard requires an interactive terminal")
		}
		model := dashboardModel{provider: provider, interval: interval, panel: panel, fetching: true}
		_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
		return err
	}}
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "dashboard refresh interval (500ms to 1m)")
	cmd.Flags().StringVar(&panel, "panel", overviewPanel, "initial panel: overview or processes")
	return cmd
}

func isInteractiveTerminal(file *os.File) bool {
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func (m dashboardModel) Init() tea.Cmd { return tea.Batch(m.fetch(), m.tick()) }
func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.KeyMsg:
		if value.String() == "q" || value.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if value.String() == "tab" || value.String() == "left" || value.String() == "right" {
			m.panel = nextDashboardPanel(m.panel)
		}
	case dashboardData:
		m.snapshot, m.err = value.snapshot, value.err
		m.fetching = false
	case tea.WindowSizeMsg:
		m.width, m.height = value.Width, value.Height
	case dashboardTick:
		// Never start a second collection while the previous one is still
		// running. A slow kernel interface should reduce refresh frequency,
		// rather than accumulate concurrent reads and stale updates.
		if m.fetching {
			return m, m.tick()
		}
		m.fetching = true
		return m, tea.Batch(m.fetch(), m.tick())
	}
	return m, nil
}
func nextDashboardPanel(panel string) string {
	if panel == processesPanel {
		return overviewPanel
	}
	return processesPanel
}
func (m dashboardModel) fetch() tea.Cmd {
	return func() tea.Msg { snapshot, err := m.provider(); return dashboardData{snapshot, err} }
}
func (m dashboardModel) tick() tea.Cmd {
	return tea.Tick(m.interval, func(time.Time) tea.Msg { return dashboardTick{} })
}
func (m dashboardModel) View() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("SysKit dashboard")
	if m.err != nil {
		return title + "\n\ncollection error: " + m.err.Error() + "\n\nq: quit"
	}
	if m.panel == processesPanel {
		return fmt.Sprintf("%s — processes\n\ntop process: %s\n\nrefresh: %s  •  tab: switch panel  •  q: quit", title, m.snapshot.TopProcess, m.interval)
	}
	if m.width > 0 && (m.width < 48 || m.height > 0 && m.height < 12) {
		return fmt.Sprintf("%s\n\nterminal is too small (%dx%d)\nresize to at least 48x12\n\nq: quit", title, m.width, m.height)
	}
	return fmt.Sprintf("%s — overview\n\nhost: %s\nuptime: %s\ncpu: %s\nmemory: %d / %d bytes\nswap: %d / %d bytes\ndisk: %d / %d bytes\nnetwork: %d interfaces%s\n\nrefresh: %s  •  tab: switch panel  •  q: quit", title, m.snapshot.Hostname, time.Duration(m.snapshot.Uptime*float64(time.Second)).Truncate(time.Second), dashboardPercent(m.snapshot.CPUPercent), m.snapshot.MemoryUsed, m.snapshot.MemoryTotal, m.snapshot.SwapUsed, m.snapshot.SwapTotal, m.snapshot.DiskUsed, m.snapshot.DiskTotal, m.snapshot.Interfaces, dashboardNetworkRatesText(m.snapshot.NetworkRX, m.snapshot.NetworkTX), m.interval)
}

func dashboardPercent(value *float64) string {
	if value == nil {
		return "sampling"
	}
	return fmt.Sprintf("%.1f%%", *value)
}

func dashboardNetworkRatesText(rx, tx *float64) string {
	if rx == nil || tx == nil {
		return " (sampling)"
	}
	return fmt.Sprintf(" (rx %.0f B/s, tx %.0f B/s)", *rx, *tx)
}

func dashboardCPUUtilization(before, after *model.CPUInfo) *float64 {
	if before == nil || after == nil {
		return nil
	}
	var old, current *model.CPUTime
	for i := range before.Times {
		if before.Times[i].CPUID == "all" {
			old = &before.Times[i]
		}
	}
	for i := range after.Times {
		if after.Times[i].CPUID == "all" {
			current = &after.Times[i]
		}
	}
	if old == nil || current == nil || current.Total <= old.Total {
		return nil
	}
	total := current.Total - old.Total
	idle := (current.Idle + current.IOWait) - (old.Idle + old.IOWait)
	if idle > total {
		return nil
	}
	value := float64(total-idle) * 100 / float64(total)
	return &value
}

func dashboardNetworkRates(before, after *model.NetworkInfo) (*float64, *float64) {
	if before == nil || after == nil {
		return nil, nil
	}
	elapsed := after.CollectedAt.Sub(before.CollectedAt).Seconds()
	if elapsed <= 0 {
		return nil, nil
	}
	previous := make(map[string]model.NetworkInterface, len(before.Interfaces))
	for _, iface := range before.Interfaces {
		previous[iface.Name] = iface
	}
	var rx, tx uint64
	for _, iface := range after.Interfaces {
		old, ok := previous[iface.Name]
		if !ok || iface.RXBytes < old.RXBytes || iface.TXBytes < old.TXBytes {
			continue
		}
		rx += iface.RXBytes - old.RXBytes
		tx += iface.TXBytes - old.TXBytes
	}
	rxRate, txRate := float64(rx)/elapsed, float64(tx)/elapsed
	return &rxRate, &txRate
}
