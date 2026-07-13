package cli

import (
	"fmt"
	"os"
	"strings"
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
	theme    tuiTheme
}

type dashboardTick struct{}
type dashboardData struct {
	snapshot dashboardSnapshot
	err      error
}

func newDashboardCmd(provider dashboardProvider) *cobra.Command {
	return newDashboardCmdWithTheme(provider, nil)
}

func newDashboardCmdWithTheme(provider dashboardProvider, selectedTheme *tuiTheme) *cobra.Command {
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
		model := dashboardModel{provider: provider, interval: interval, panel: panel, fetching: true, theme: resolveTUITheme(selectedTheme)}
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
	theme := m.theme
	if theme.accent.primary == "" {
		theme = defaultTUITheme
	}
	title := theme.badge("◉  SYSKIT LIVE") + "  " + theme.primaryStyle().Bold(theme.color).Render("SYSTEM DASHBOARD")
	tabs := dashboardTabs(theme, m.panel)
	if m.err != nil {
		return title + "\n" + tabs + "\n\n" + dashboardCard(theme, "COLLECTION ERROR", m.err.Error(), max(30, m.width-4)) + "\n\nq: quit"
	}
	if m.panel == processesPanel {
		body := fmt.Sprintf("top process: %s\n\nLive process ranking is available in syskit top.", m.snapshot.TopProcess)
		return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", title, tabs, dashboardCard(theme, "PROCESS FOCUS", body, max(34, m.width-4)), dashboardFooter(theme, m.interval))
	}
	if m.width > 0 && (m.width < 48 || m.height > 0 && m.height < 12) {
		return fmt.Sprintf("%s\n\nterminal is too small (%dx%d)\nresize to at least 48x12\n\nq: quit", title, m.width, m.height)
	}
	host := fmt.Sprintf("host: %s\nuptime: %s\ncpu: %s %s", m.snapshot.Hostname, time.Duration(m.snapshot.Uptime*float64(time.Second)).Truncate(time.Second), dashboardPercent(m.snapshot.CPUPercent), dashboardOptionalBar(theme, m.snapshot.CPUPercent))
	memory := fmt.Sprintf("memory: %s / %s\n%s\nswap: %s / %s", formatTUIBytes(m.snapshot.MemoryUsed), formatTUIBytes(m.snapshot.MemoryTotal), dashboardBar(theme, m.snapshot.MemoryUsed, m.snapshot.MemoryTotal, 20), formatTUIBytes(m.snapshot.SwapUsed), formatTUIBytes(m.snapshot.SwapTotal))
	storage := fmt.Sprintf("disk: %s / %s\n%s", formatTUIBytes(m.snapshot.DiskUsed), formatTUIBytes(m.snapshot.DiskTotal), dashboardBar(theme, m.snapshot.DiskUsed, m.snapshot.DiskTotal, 20))
	network := fmt.Sprintf("network: %d interfaces%s", m.snapshot.Interfaces, dashboardNetworkRatesText(m.snapshot.NetworkRX, m.snapshot.NetworkTX))
	if m.width > 0 && m.width < 76 || m.height > 0 && m.height < 18 {
		compact := strings.Join([]string{
			fmt.Sprintf("host: %s  •  uptime: %s", m.snapshot.Hostname, time.Duration(m.snapshot.Uptime*float64(time.Second)).Truncate(time.Second)),
			fmt.Sprintf("cpu: %s %s", dashboardPercent(m.snapshot.CPUPercent), dashboardOptionalBar(theme, m.snapshot.CPUPercent)),
			strings.Split(memory, "\n")[0],
			strings.Split(storage, "\n")[0],
			network,
		}, "\n")
		return fmt.Sprintf("%s\n%s\n\n%s%s", title, tabs, compact, dashboardFooter(theme, m.interval))
	}
	cardWidth := max(32, (max(72, m.width)-3)/2)
	rowOne := lipgloss.JoinHorizontal(lipgloss.Top, dashboardCard(theme, "HOST + CPU", host, cardWidth), " ", dashboardCard(theme, "MEMORY", memory, cardWidth))
	rowTwo := lipgloss.JoinHorizontal(lipgloss.Top, dashboardCard(theme, "STORAGE", storage, cardWidth), " ", dashboardCard(theme, "NETWORK", network, cardWidth))
	return fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s", title, tabs, rowOne, rowTwo, dashboardFooter(theme, m.interval))
}

func dashboardTabs(theme tuiTheme, active string) string {
	overview, processes := "  OVERVIEW  ", "  PROCESSES  "
	if active == overviewPanel {
		overview = theme.badge("OVERVIEW")
	} else {
		processes = theme.badge("PROCESSES")
	}
	return overview + "  " + processes
}

func dashboardCard(theme tuiTheme, heading, body string, width int) string {
	headingStyle := theme.primaryStyle().Bold(theme.color)
	return theme.borderStyle().Width(max(20, width-4)).Render(headingStyle.Render(heading) + "\n" + body)
}

func dashboardFooter(theme tuiTheme, interval time.Duration) string {
	return "\n" + theme.primaryStyle().Render(fmt.Sprintf("● live  •  refresh %s  •  tab/←/→ switch  •  q quit", interval))
}

func dashboardOptionalBar(theme tuiTheme, value *float64) string {
	if value == nil {
		return ""
	}
	return dashboardBar(theme, uint64(*value), 100, 12)
}

func dashboardBar(theme tuiTheme, used, total uint64, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}
	filled := min(width, int(float64(used)/float64(total)*float64(width)))
	active := strings.Repeat("█", filled)
	if theme.color {
		active = theme.primaryStyle().Render(active)
	}
	return active + strings.Repeat("░", width-filled)
}

func formatTUIBytes(value uint64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	amount := float64(value)
	unit := 0
	for amount >= 1024 && unit < len(units)-1 {
		amount /= 1024
		unit++
	}
	if unit == 0 {
		return fmt.Sprintf("%d %s", value, units[unit])
	}
	return fmt.Sprintf("%.1f %s", amount, units[unit])
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
	return fmt.Sprintf(" (rx %s, tx %s)", formatTUIRate(*rx), formatTUIRate(*tx))
}

func formatTUIRate(value float64) string {
	units := []string{"B/s", "KiB/s", "MiB/s", "GiB/s"}
	unit := 0
	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	if unit == 0 {
		return fmt.Sprintf("%.0f %s", value, units[unit])
	}
	return fmt.Sprintf("%.1f %s", value, units[unit])
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
