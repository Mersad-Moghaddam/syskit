package cli

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	minDashboardInterval = 500 * time.Millisecond
	maxDashboardInterval = time.Minute
)

type dashboardSnapshot struct {
	Hostname    string
	Uptime      float64
	MemoryUsed  uint64
	MemoryTotal uint64
	Interfaces  int
}

type dashboardProvider func() (dashboardSnapshot, error)

type dashboardModel struct {
	provider dashboardProvider
	interval time.Duration
	snapshot dashboardSnapshot
	err      error
}

type dashboardTick struct{}
type dashboardData struct {
	snapshot dashboardSnapshot
	err      error
}

func newDashboardCmd(provider dashboardProvider) *cobra.Command {
	var interval time.Duration
	cmd := &cobra.Command{Use: "dashboard", Short: "Show a live system dashboard", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		if interval < minDashboardInterval || interval > maxDashboardInterval {
			return fmt.Errorf("dashboard interval must be between %s and %s", minDashboardInterval, maxDashboardInterval)
		}
		model := dashboardModel{provider: provider, interval: interval}
		_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
		return err
	}}
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "dashboard refresh interval (500ms to 1m)")
	return cmd
}

func (m dashboardModel) Init() tea.Cmd { return tea.Batch(m.fetch(), m.tick()) }
func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.KeyMsg:
		if value.String() == "q" || value.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case dashboardData:
		m.snapshot, m.err = value.snapshot, value.err
	case dashboardTick:
		return m, tea.Batch(m.fetch(), m.tick())
	}
	return m, nil
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
	return fmt.Sprintf("%s\n\nhost: %s\nuptime: %s\nmemory: %d / %d bytes\nnetwork interfaces: %d\n\nrefresh: %s  •  q: quit", title, m.snapshot.Hostname, time.Duration(m.snapshot.Uptime*float64(time.Second)).Truncate(time.Second), m.snapshot.MemoryUsed, m.snapshot.MemoryTotal, m.snapshot.Interfaces, m.interval)
}
