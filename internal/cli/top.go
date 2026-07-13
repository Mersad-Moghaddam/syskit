package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
)

type topProvider func(service.ProcessOptions) (*model.ProcessList, error)
type topModel struct {
	provider topProvider
	interval time.Duration
	options  service.ProcessOptions
	list     *model.ProcessList
	err      error
}
type topTick struct{}
type topData struct {
	list *model.ProcessList
	err  error
}

func newTopCmd(provider topProvider) *cobra.Command {
	var interval time.Duration
	var sort string
	var limit int
	var filters []string
	cmd := &cobra.Command{Use: "top", Short: "Monitor processes interactively", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		if interval < minWatchInterval || interval > maxWatchInterval {
			return fmt.Errorf("top interval must be between %s and %s", minWatchInterval, maxWatchInterval)
		}
		if !isInteractiveTerminal(os.Stdout) {
			return fmt.Errorf("top requires an interactive terminal")
		}
		parsed, err := service.ParseProcessFilters(filters)
		if err != nil {
			return err
		}
		model := topModel{provider: provider, interval: interval, options: service.ProcessOptions{Filters: parsed, Sort: sort, Reverse: true, Limit: limit}}
		_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
		return err
	}}
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "refresh interval (250ms to 1m)")
	cmd.Flags().StringVar(&sort, "sort", "memory", "initial sort: cpu, memory, name, or pid")
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum process rows")
	cmd.Flags().StringSliceVar(&filters, "filter", nil, "filter with field=value (repeatable)")
	return cmd
}
func (m topModel) Init() tea.Cmd { return tea.Batch(m.fetch(), m.tick()) }
func (m topModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.KeyMsg:
		switch value.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			m.options.Sort = "cpu"
		case "m":
			m.options.Sort = "memory"
		case "n":
			m.options.Sort = "name"
		case "p":
			m.options.Sort = "pid"
		default:
			return m, nil
		}
		return m, m.fetch()
	case topData:
		m.list, m.err = value.list, value.err
	case topTick:
		return m, tea.Batch(m.fetch(), m.tick())
	}
	return m, nil
}
func (m topModel) fetch() tea.Cmd {
	return func() tea.Msg { list, err := m.provider(m.options); return topData{list, err} }
}
func (m topModel) tick() tea.Cmd {
	return tea.Tick(m.interval, func(time.Time) tea.Msg { return topTick{} })
}
func (m topModel) View() string {
	if m.err != nil {
		return "SysKit top\n\ncollection error: " + m.err.Error() + "\n\nq: quit"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "SysKit top — sort: %s\n\nPID     USER       CPU%%    MEM%%    COMMAND\n", m.options.Sort)
	if m.list != nil {
		for _, p := range m.list.Processes {
			user := p.User
			if user == "" {
				user = fmt.Sprint(p.UID)
			}
			cpu, mem := 0.0, 0.0
			if p.CPUPercent != nil {
				cpu = *p.CPUPercent
			}
			if p.MemoryPercent != nil {
				mem = *p.MemoryPercent
			}
			fmt.Fprintf(&b, "%-7d %-10s %5.1f  %5.1f  %s\n", p.PID, user, cpu, mem, p.Command)
		}
	}
	b.WriteString("\nc/m/n/p: sort  •  q: quit")
	return b.String()
}
