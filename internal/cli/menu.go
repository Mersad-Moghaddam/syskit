package cli

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const menuIntroInterval = 55 * time.Millisecond

type menuIntroTick struct{}

type menuPrompt struct {
	label       string
	placeholder string
	prefix      []string
}

type menuItem struct {
	title       string
	description string
	children    []menuItem
	args        []string
	prompt      *menuPrompt
	live        bool
	accent      tuiAccent
}

type menuSelection struct {
	title  string
	args   []string
	live   bool
	accent tuiAccent
}

type menuModel struct {
	root       menuItem
	current    *menuItem
	parents    []*menuItem
	cursors    []int
	cursor     int
	offset     int
	width      int
	height     int
	input      string
	inputItem  *menuItem
	inputError string
	selection  *menuSelection
	color      bool
	introFrame int
}

func newMenuModel() menuModel {
	return newMenuModelWithColor(true)
}

func newMenuModelWithColor(color bool) menuModel {
	root := interactiveMenuTree()
	model := menuModel{root: root, width: 80, height: 24, color: color}
	model.current = &model.root
	return model
}

func (m menuModel) Init() tea.Cmd {
	if m.introFrame >= menuIntroFrames-1 {
		return nil
	}
	return menuIntroAnimationTick()
}

func menuIntroAnimationTick() tea.Cmd {
	return tea.Tick(menuIntroInterval, func(time.Time) tea.Msg { return menuIntroTick{} })
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case menuIntroTick:
		if m.introFrame < menuIntroFrames-1 {
			m.introFrame++
			return m, menuIntroAnimationTick()
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = value.Width, value.Height
		m.ensureVisible()
		return m, nil
	case tea.MouseMsg:
		m.introFrame = menuIntroFrames - 1
		if m.inputItem != nil {
			return m, nil
		}
		event := tea.MouseEvent(value)
		switch event.Button {
		case tea.MouseButtonWheelUp:
			m.move(-1)
		case tea.MouseButtonWheelDown:
			m.move(1)
		case tea.MouseButtonLeft:
			if event.Action == tea.MouseActionPress {
				index := m.offset + event.Y - m.itemStartRow()
				if index >= 0 && index < len(m.current.children) {
					m.cursor = index
					return m.activate()
				}
			}
		}
		return m, nil
	case tea.KeyMsg:
		m.introFrame = menuIntroFrames - 1
		if m.inputItem != nil {
			return m.updateInput(value)
		}
		switch value.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.move(-1)
		case "down", "j":
			m.move(1)
		case "enter", "right", "l":
			return m.activate()
		case "esc", "left", "h", "backspace":
			if !m.back() {
				return m, tea.Quit
			}
		case "home", "g":
			m.cursor = 0
			m.ensureVisible()
		case "end", "G":
			m.cursor = max(0, len(m.current.children)-1)
			m.ensureVisible()
		}
	}
	return m, nil
}

func (m menuModel) updateInput(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.inputItem, m.input, m.inputError = nil, "", ""
	case tea.KeyEnter:
		value := strings.TrimSpace(m.input)
		if value == "" {
			m.inputError = "A value is required."
			return m, nil
		}
		item := m.inputItem
		args := append(append([]string(nil), item.prompt.prefix...), value)
		m.selection = &menuSelection{title: item.title, args: args, live: item.live, accent: item.accent}
		return m, tea.Quit
	case tea.KeyBackspace, tea.KeyDelete:
		if m.input != "" {
			_, size := utf8.DecodeLastRuneInString(m.input)
			m.input = m.input[:len(m.input)-size]
		}
		m.inputError = ""
	case tea.KeyRunes:
		m.input += string(key.Runes)
		m.inputError = ""
	case tea.KeySpace:
		m.input += " "
		m.inputError = ""
	}
	return m, nil
}

func (m menuModel) activate() (tea.Model, tea.Cmd) {
	if len(m.current.children) == 0 || m.cursor >= len(m.current.children) {
		return m, nil
	}
	item := &m.current.children[m.cursor]
	if len(item.children) > 0 {
		m.parents = append(m.parents, m.current)
		m.cursors = append(m.cursors, m.cursor)
		m.current, m.cursor, m.offset = item, 0, 0
		return m, nil
	}
	if item.prompt != nil {
		m.inputItem, m.input, m.inputError = item, "", ""
		return m, nil
	}
	m.selection = &menuSelection{title: item.title, args: append([]string(nil), item.args...), live: item.live, accent: item.accent}
	return m, tea.Quit
}

func (m *menuModel) back() bool {
	if len(m.parents) == 0 {
		return false
	}
	last := len(m.parents) - 1
	m.current = m.parents[last]
	m.parents = m.parents[:last]
	m.cursor = m.cursors[last]
	m.cursors = m.cursors[:last]
	m.offset = 0
	m.ensureVisible()
	return true
}

func (m *menuModel) move(delta int) {
	count := len(m.current.children)
	if count == 0 {
		return
	}
	m.cursor = (m.cursor + delta + count) % count
	m.ensureVisible()
}

func (m *menuModel) ensureVisible() {
	visible := m.visibleRows()
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m menuModel) visibleRows() int {
	if m.height <= 0 {
		return len(m.current.children)
	}
	return max(1, m.height-m.itemStartRow()-4)
}

func (m menuModel) compact() bool {
	return m.width > 0 && m.width < 68 || m.height > 0 && m.height < 20
}

func (m menuModel) itemStartRow() int {
	if m.compact() {
		return 5
	}
	return len(syskitLogo) + 4
}

func (m menuModel) resumedAfterCommand() menuModel {
	m.selection = nil
	m.inputItem, m.input, m.inputError = nil, "", ""
	m.introFrame = menuIntroFrames - 1
	return m
}

func (m menuModel) View() string {
	styles := newMenuStyles(m.color)
	var view strings.Builder
	view.WriteString(renderSysKitLogo(m.introFrame, m.width, m.color, m.compact()))
	view.WriteString("\n")
	view.WriteString(renderTUITagline(m.introFrame, m.width, m.color))
	view.WriteString("\n\n")
	currentTheme := tuiTheme{accent: m.current.accent, color: m.color}
	view.WriteString(currentTheme.primaryStyle().Bold(true).Render("⌂  " + m.breadcrumb()))
	view.WriteString("\n\n")

	if m.inputItem != nil {
		theme := tuiTheme{accent: m.inputItem.accent, color: m.color}
		view.WriteString(theme.badge(m.inputItem.accent.icon + "  " + m.inputItem.title))
		view.WriteString("\n")
		view.WriteString(styles.description.Render("  " + m.inputItem.description))
		view.WriteString("\n\n")
		view.WriteString("  " + m.inputItem.prompt.label + "\n")
		value := m.input
		if value == "" {
			value = styles.placeholder.Render(m.inputItem.prompt.placeholder)
		}
		inputStyle := styles.input
		if m.color {
			inputStyle = inputStyle.BorderForeground(m.inputItem.accent.primary)
		}
		view.WriteString(inputStyle.Render(value + "█"))
		if m.inputError != "" {
			view.WriteString("\n" + styles.error.Render("  "+m.inputError))
		}
		view.WriteString("\n\n" + theme.primaryStyle().Render("  ↵ run  •  esc cancel  •  ctrl-c quit"))
		return view.String()
	}

	visible := m.visibleRows()
	end := min(len(m.current.children), m.offset+visible)
	for index := m.offset; index < end; index++ {
		item := m.current.children[index]
		marker, suffix := "   ", ""
		if len(item.children) > 0 {
			suffix = "   ›"
		}
		label := fitMenuLabel(item.accent.icon+"  "+item.title, item.description, suffix, m.width)
		line := marker + label
		theme := tuiTheme{accent: item.accent, color: m.color}
		if index == m.cursor {
			line = "▌  " + label
			view.WriteString(theme.selectedStyle().PaddingRight(1).Render(line))
		} else {
			view.WriteString(theme.primaryStyle().Render(line))
		}
		view.WriteString("\n")
	}
	if len(m.current.children) > visible {
		view.WriteString(styles.muted.Render(fmt.Sprintf("  %d–%d of %d", m.offset+1, end, len(m.current.children))))
		view.WriteString("\n")
	}
	view.WriteString("\n")
	view.WriteString(styles.description.Render("  " + m.selectedDescription()))
	view.WriteString("\n\n")
	footerTheme := currentTheme
	if len(m.current.children) > 0 && m.cursor < len(m.current.children) {
		footerTheme.accent = m.current.children[m.cursor].accent
	}
	view.WriteString(footerTheme.primaryStyle().Render("  ↑↓/jk move  •  ↵/→ open  •  esc/← back  •  q quit  •  ◉ mouse"))
	return view.String()
}

func (m menuModel) breadcrumb() string {
	parts := []string{"Home"}
	for _, parent := range m.parents {
		if parent.title != "Home" {
			parts = append(parts, parent.title)
		}
	}
	if m.current.title != "Home" {
		parts = append(parts, m.current.title)
	}
	return strings.Join(parts, "  /  ")
}

func (m menuModel) selectedDescription() string {
	if len(m.current.children) == 0 || m.cursor >= len(m.current.children) {
		return "No actions available"
	}
	return m.current.children[m.cursor].description
}

type menuStyles struct {
	muted, description, input, placeholder, error lipgloss.Style
}

func newMenuStyles(color bool) menuStyles {
	if !color {
		plain := lipgloss.NewStyle()
		return menuStyles{muted: plain, description: plain,
			input: plain.Border(lipgloss.RoundedBorder()).Padding(0, 1).MarginLeft(2), placeholder: plain, error: plain}
	}
	return menuStyles{
		muted:       lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
		description: lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		input:       lipgloss.NewStyle().Foreground(lipgloss.Color("230")).Background(lipgloss.Color("235")).Border(lipgloss.RoundedBorder()).Padding(0, 1).MarginLeft(2),
		placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color("242")),
		error:       lipgloss.NewStyle().Foreground(lipgloss.Color("203")),
	}
}

func fitMenuLabel(title, description, suffix string, width int) string {
	if width < 36 {
		return truncateMenuText(title+suffix, max(8, width-5))
	}
	labelWidth := min(24, max(16, width/3))
	title = truncateMenuText(title, labelWidth-1)
	line := title + strings.Repeat(" ", max(1, labelWidth-lipgloss.Width(title))) + " " + description + suffix
	return truncateMenuText(line, max(12, width-5))
}

func truncateMenuText(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	if limit <= 1 {
		return string(runes[:limit])
	}
	return string(runes[:limit-1]) + "…"
}

func interactiveMenuTree() menuItem {
	accentIndex := 0
	nextAccent := func() tuiAccent {
		accent := paletteAccent(accentIndex)
		accentIndex++
		return accent
	}
	leaf := func(title, description string, args ...string) menuItem {
		return menuItem{title: title, description: description, args: args, accent: nextAccent()}
	}
	live := func(title, description string, args ...string) menuItem {
		return menuItem{title: title, description: description, args: args, live: true, accent: nextAccent()}
	}
	group := func(title, description string, children ...menuItem) menuItem {
		return menuItem{title: title, description: description, children: children, accent: nextAccent()}
	}
	prompt := func(title, description, label, placeholder string, prefix ...string) menuItem {
		return menuItem{title: title, description: description, prompt: &menuPrompt{label: label, placeholder: placeholder, prefix: prefix}, accent: nextAccent()}
	}

	cpu := group("CPU", "Topology, model, frequency, and sampled utilization",
		leaf("CPU overview", "Aggregate topology and one-second utilization sample", "cpu"),
		leaf("Per-core utilization", "One utilization row for every logical CPU", "cpu", "--per-core"),
		leaf("Fast sample", "Aggregate CPU utilization over 250ms", "cpu", "--interval", "250ms"),
	)
	filesystem := group("Filesystem", "Mount capacity, inode usage, types, and options",
		leaf("Real filesystems", "Hide pseudo filesystems and show usable mounts", "filesystem"),
		leaf("All filesystems", "Include procfs, sysfs, tmpfs, and other pseudo mounts", "filesystem", "--show-pseudo"),
	)
	diagnosticCategories := group("By category", "Collect only one diagnostic domain",
		leaf("CPU", "Load and CPU topology findings", "diagnostics", "--category", "cpu"),
		leaf("Memory", "Memory PSI and swap findings", "diagnostics", "--category", "memory"),
		leaf("Disk", "Device saturation availability findings", "diagnostics", "--category", "disk"),
		leaf("Filesystem", "Filesystem capacity findings", "diagnostics", "--category", "filesystem"),
		leaf("Process", "Process memory concentration findings", "diagnostics", "--category", "process"),
		leaf("Network", "Interface error and drop findings", "diagnostics", "--category", "network"),
		leaf("Ports", "Wildcard listener findings", "diagnostics", "--category", "ports"),
	)

	return group("Home", "Choose a system-inspection domain",
		group("System & hardware", "Host identity and core resource summaries",
			leaf("System overview", "Host, OS, kernel, uptime, boot time, and load", "system"),
			cpu,
			leaf("Memory", "RAM, swap, cache, and kernel pressure information", "memory"),
		),
		group("Storage", "Capacity, device I/O, mounts, and inodes",
			leaf("Disk capacity", "Mounted storage capacity and usage", "disk"),
			leaf("Disk I/O", "Sample per-device read and write rates", "disk", "--io"),
			filesystem,
		),
		group("Processes", "Process snapshots, hierarchy, and resource views",
			leaf("Process list", "All readable processes sorted by PID", "process"),
			leaf("Top CPU snapshot", "Sample CPU usage and show the busiest 20", "process", "--sort", "cpu", "--reverse", "--limit", "20", "--interval", "1s"),
			leaf("Top memory snapshot", "Show the 20 largest resident-memory users", "process", "--sort", "memory", "--reverse", "--limit", "20"),
			leaf("Process tree", "Parent and child hierarchy", "process", "tree"),
			leaf("Container processes", "Only processes with recognized container cgroups", "process", "--containers"),
			live("Interactive top", "Live sortable process monitor", "top"),
		),
		group("Network & ports", "Interfaces, routes, DNS, and socket ownership",
			leaf("Network summary", "Counters and metadata for every interface", "network"),
			leaf("Interfaces", "Focused interface addresses and counters", "network", "interfaces"),
			leaf("Routes", "IPv4 routes and default gateway", "network", "routes"),
			leaf("DNS", "Configured resolver nameservers", "network", "dns"),
			leaf("All sockets", "TCP, UDP, IPv6, and Unix sockets with owners", "ports"),
			leaf("Listening sockets", "Only sockets accepting connections", "ports", "--listening"),
		),
		group("Live monitoring", "Continuous dashboards and refreshed commands",
			live("Dashboard overview", "Live CPU, memory, disk, network, and host summary", "dashboard"),
			live("Dashboard processes", "Open the dashboard on its process panel", "dashboard", "--panel", "processes"),
			live("Top by memory", "Interactive process monitor sorted by memory", "top", "--sort", "memory"),
			live("Top by CPU", "Interactive process monitor sorted by CPU", "top", "--sort", "cpu"),
			live("Watch system", "Continuously refresh the host overview", "watch", "system"),
			live("Watch CPU", "Continuously refresh CPU sampling", "watch", "cpu"),
			live("Watch memory", "Continuously refresh memory information", "watch", "memory"),
			live("Watch network", "Continuously refresh interface counters", "watch", "network"),
		),
		group("Diagnostics", "Explainable, read-only system health findings",
			leaf("All findings", "Run every diagnostic category", "diagnostics"),
			leaf("Warnings", "Show warning-level findings only", "diagnostics", "--severity", "warning"),
			leaf("Critical findings", "Show critical findings only", "diagnostics", "--severity", "critical"),
			leaf("Information", "Show informational and unavailable findings", "diagnostics", "--severity", "info"),
			diagnosticCategories,
		),
		group("Containers", "Best-effort cgroup-derived container inspection",
			leaf("Container list", "IDs, process counts, and available cgroup counters", "containers"),
			prompt("Inspect container", "Show processes associated with one container ID", "Container ID", "enter the full container ID", "containers", "inspect"),
		),
		group("Plugins", "Discover, inspect, and explicitly run extensions",
			leaf("Plugin list", "Discover plugin manifests without executing code", "plugins", "list"),
			prompt("Inspect plugin", "Show compatibility and requested permissions", "Plugin name", "example", "plugins", "inspect"),
			prompt("Run plugin", "Explicitly execute one compatible plugin", "Plugin name", "example", "plugins", "run"),
		),
		group("Output formats", "Preview stable human and automation formats",
			leaf("Table output", "Human-readable system overview", "system", "--format", "table"),
			leaf("JSON output", "Stable machine-readable system schema", "system", "--format", "json"),
			leaf("YAML output", "Stable human-editable system schema", "system", "--format", "yaml"),
		),
		group("Utilities", "Version, help, and shell integration",
			leaf("Version", "Print the embedded SysKit version", "version"),
			leaf("Command help", "Show the complete command and flag reference", "--help"),
			group("Shell completion", "Generate completion scripts for your shell",
				leaf("Bash", "Generate Bash completion source", "completion", "bash"),
				leaf("Zsh", "Generate Zsh completion source", "completion", "zsh"),
				leaf("Fish", "Generate Fish completion source", "completion", "fish"),
				leaf("PowerShell", "Generate PowerShell completion source", "completion", "powershell"),
			),
		),
	)
}

func runInteractiveMenu(stdin, stdout, stderr *os.File, color bool, globalArgs []string) error {
	model := newMenuModelWithColor(color)
	for {
		result, err := tea.NewProgram(model, tea.WithInput(stdin), tea.WithOutput(stdout), tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
		if err != nil {
			return fmt.Errorf("running interactive menu: %w", err)
		}
		finished, ok := result.(menuModel)
		if !ok || finished.selection == nil {
			return nil
		}

		selection := finished.selection
		model = finished.resumedAfterCommand()
		if !selection.live {
			if err := runMenuCommandResult(stdin, stdout, *selection, globalArgs, color); err != nil {
				return err
			}
			continue
		}

		theme := &tuiTheme{accent: selection.accent, color: color}
		child, options := newRootCmdWithOptionsAndTheme(theme)
		child.SetIn(stdin)
		child.SetOut(stdout)
		child.SetErr(stderr)
		child.SetArgs(append(append([]string(nil), globalArgs...), selection.args...))
		err = child.Execute()
		if err == nil {
			continue
		}
		message, code := present(err)
		if options.quiet {
			message = "Live view exited with an error."
		}
		if err := runMenuCompletedResult(stdin, stdout, *selection, menuCommandResult{stderr: message, err: err, code: code}, color); err != nil {
			return err
		}
	}
}

func changedPersistentArgs(cmd *cobra.Command) []string {
	var args []string
	cmd.PersistentFlags().Visit(func(flag *pflag.Flag) {
		args = append(args, "--"+flag.Name+"="+flag.Value.String())
	})
	return args
}
