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
	fetching bool
	offset   int
	width    int
	height   int
	theme    tuiTheme
}
type topTick struct{}
type topData struct {
	list *model.ProcessList
	err  error
}

func newTopCmd(provider topProvider) *cobra.Command {
	return newTopCmdWithTheme(provider, nil)
}

func newTopCmdWithTheme(provider topProvider, selectedTheme *tuiTheme) *cobra.Command {
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
		model := topModel{provider: provider, interval: interval, options: service.ProcessOptions{Filters: parsed, Sort: sort, Reverse: true, Limit: limit}, fetching: true, theme: resolveTUITheme(selectedTheme)}
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
	case tea.WindowSizeMsg:
		m.width, m.height = value.Width, value.Height
		return m, nil
	case tea.KeyMsg:
		switch value.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.offset = m.nextOffset(1)
			return m, nil
		case "k", "up":
			m.offset = m.nextOffset(-1)
			return m, nil
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
		m.offset = 0
		if m.fetching {
			return m, nil
		}
		m.fetching = true
		return m, m.fetch()
	case topData:
		m.list, m.err = value.list, value.err
		m.fetching = false
		m.offset = m.nextOffset(0)
	case topTick:
		if m.fetching {
			return m, m.tick()
		}
		m.fetching = true
		return m, tea.Batch(m.fetch(), m.tick())
	}
	return m, nil
}

func (m topModel) nextOffset(delta int) int {
	if m.list == nil || len(m.list.Processes) == 0 {
		return 0
	}
	offset := m.offset + delta
	if offset < 0 {
		return 0
	}
	max := len(m.list.Processes) - 1
	if offset > max {
		return max
	}
	return offset
}
func (m topModel) fetch() tea.Cmd {
	return func() tea.Msg { list, err := m.provider(m.options); return topData{list, err} }
}
func (m topModel) tick() tea.Cmd {
	return tea.Tick(m.interval, func(time.Time) tea.Msg { return topTick{} })
}
func (m topModel) View() string {
	theme := m.theme
	if theme.accent.primary == "" {
		theme = defaultTUITheme
	}
	title := theme.badge("▲  SYSKIT TOP") + "  " + theme.primaryStyle().Bold(theme.color).Render("LIVE PROCESS INTELLIGENCE")
	if m.err != nil {
		return title + "\n\n" + theme.borderStyle().Render("collection error: "+m.err.Error()) + "\n\nq: quit"
	}
	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(renderTopSortTabs(theme, m.options.Sort))
	b.WriteString("\n\n")
	b.WriteString(theme.primaryStyle().Bold(theme.color).Render("PID     USER       CPU%    MEM%    COMMAND"))
	b.WriteString("\n")
	if m.list != nil {
		end := len(m.list.Processes)
		if m.height > 0 {
			end = min(end, m.offset+max(1, m.height-8))
		}
		for index, p := range m.list.Processes[m.offset:end] {
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
			row := fmt.Sprintf("%-7d %-10s %5.1f  %5.1f  %s", p.PID, user, cpu, mem, p.Command)
			if m.width > 0 {
				row = truncateMenuText(row, max(20, m.width-3))
			}
			if index == 0 {
				row = theme.primaryStyle().Render("▌ " + row)
			} else {
				row = "  " + row
			}
			b.WriteString(row + "\n")
		}
	}
	b.WriteString("\n" + theme.primaryStyle().Render(fmt.Sprintf("● live  •  refresh %s  •  c/m/n/p sort  •  j/k scroll  •  q quit", m.interval)))
	return b.String()
}

func renderTopSortTabs(theme tuiTheme, selected string) string {
	keys := []struct{ key, label string }{{"c", "CPU"}, {"m", "MEMORY"}, {"n", "NAME"}, {"p", "PID"}}
	parts := make([]string, 0, len(keys))
	for _, item := range keys {
		label := item.key + " " + item.label
		if strings.EqualFold(selected, item.label) {
			label = theme.badge(label)
		}
		parts = append(parts, label)
	}
	return strings.Join(parts, "  ")
}
