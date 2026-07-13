package cli

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const resultSpinnerInterval = 80 * time.Millisecond

type menuCommandResult struct {
	stdout string
	stderr string
	err    error
	code   int
}

type menuResultReady struct{ result menuCommandResult }
type menuResultTick struct{}

type menuResultModel struct {
	selection menuSelection
	run       func() menuCommandResult
	result    menuCommandResult
	loading   bool
	spinner   int
	width     int
	height    int
	yOffset   int
	xOffset   int
}

func newMenuResultModel(selection menuSelection, run func() menuCommandResult) menuResultModel {
	return menuResultModel{selection: selection, run: run, loading: true, width: 80, height: 24}
}

func (m menuResultModel) Init() tea.Cmd {
	return tea.Batch(m.execute(), menuResultSpinnerTick())
}

func (m menuResultModel) execute() tea.Cmd {
	return func() tea.Msg { return menuResultReady{result: m.run()} }
}

func menuResultSpinnerTick() tea.Cmd {
	return tea.Tick(resultSpinnerInterval, func(time.Time) tea.Msg { return menuResultTick{} })
}

func (m menuResultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch value := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = value.Width, value.Height
		m.clampOffsets()
	case menuResultReady:
		m.result, m.loading = value.result, false
		m.clampOffsets()
	case menuResultTick:
		if m.loading {
			m.spinner++
			return m, menuResultSpinnerTick()
		}
	case tea.MouseMsg:
		event := tea.MouseEvent(value)
		switch event.Button {
		case tea.MouseButtonWheelUp:
			m.yOffset -= 3
		case tea.MouseButtonWheelDown:
			m.yOffset += 3
		}
		m.clampOffsets()
	case tea.KeyMsg:
		switch value.String() {
		case "ctrl+c", "q", "esc", "enter":
			return m, tea.Quit
		case "up", "k":
			m.yOffset--
		case "down", "j":
			m.yOffset++
		case "pgup":
			m.yOffset -= m.resultRows()
		case "pgdown", " ":
			m.yOffset += m.resultRows()
		case "home", "g":
			m.yOffset = 0
		case "end", "G":
			m.yOffset = len(m.outputLines())
		case "left":
			m.xOffset -= 4
		case "right":
			m.xOffset += 4
		case "shift+left":
			m.xOffset = 0
		case "shift+right":
			m.xOffset = m.longestLine()
		}
		m.clampOffsets()
	}
	return m, nil
}

func (m *menuResultModel) clampOffsets() {
	lines := m.outputLines()
	maxY := max(0, len(lines)-m.resultRows())
	m.yOffset = min(max(0, m.yOffset), maxY)
	maxX := max(0, m.longestLine()-m.contentWidth())
	m.xOffset = min(max(0, m.xOffset), maxX)
}

func (m menuResultModel) outputLines() []string {
	if m.loading {
		return nil
	}
	output := strings.TrimRight(m.result.stdout, "\n")
	if m.result.stderr != "" {
		if output != "" {
			output += "\n\n"
		}
		output += strings.TrimRight(m.result.stderr, "\n")
	}
	if output == "" {
		output = "Command completed without output."
	}
	return strings.Split(stripTerminalControls(output), "\n")
}

func (m menuResultModel) resultRows() int { return max(1, m.height-9) }
func (m menuResultModel) contentWidth() int {
	return max(12, m.width-6)
}

func (m menuResultModel) longestLine() int {
	longest := 0
	for _, line := range m.outputLines() {
		longest = max(longest, lipgloss.Width(line))
	}
	return longest
}

func (m menuResultModel) View() string {
	theme := tuiTheme{accent: m.selection.accent, color: m.selection.accent.primary != ""}
	var view strings.Builder
	view.WriteString(theme.badge(m.selection.accent.icon + "  " + truncateMenuText(m.selection.title, max(8, m.width-12))))
	if m.width <= 0 || m.width >= 52 {
		view.WriteString("  ")
		view.WriteString(theme.primaryStyle().Bold(theme.color).Render("SYSKIT / RESULT"))
	}
	view.WriteString("\n")
	view.WriteString(theme.mutedStyle().Render(truncateMenuText("$ "+formatMenuCommand(m.selection.args), max(12, m.width-2))))
	view.WriteString("\n\n")

	if m.loading {
		spinner := []string{"◜", "◠", "◝", "◞", "◡", "◟"}[m.spinner%6]
		view.WriteString(theme.borderStyle().Width(max(20, m.width-6)).Render(
			theme.primaryStyle().Bold(theme.color).Render(spinner+"  collecting live Linux data") +
				"\n" + theme.mutedStyle().Render("Reading native kernel interfaces…"),
		))
		view.WriteString("\n\n")
		view.WriteString(theme.primaryStyle().Render("q/esc leave view"))
		return view.String()
	}

	status := "✓ SUCCESS"
	if m.result.err != nil {
		status = fmt.Sprintf("× FAILED · EXIT %d", m.result.code)
	}
	view.WriteString(theme.badge(status))
	position := fmt.Sprintf("lines %d–%d/%d  •  column %d", m.yOffset+1, min(len(m.outputLines()), m.yOffset+m.resultRows()), len(m.outputLines()), m.xOffset+1)
	if m.width > 0 && m.width < 52 {
		view.WriteString("\n")
	} else {
		view.WriteString("  ")
	}
	view.WriteString(theme.mutedStyle().Render(position))
	view.WriteString("\n")
	view.WriteString(m.renderOutputPanel(theme))
	view.WriteString("\n")
	footer := "↑↓/jk scroll  •  ←→ horizontal  •  pgup/pgdn page  •  g/G edges  •  ↵/esc return"
	if m.width > 0 && m.width < 64 {
		footer = "↑↓ scroll  •  ←→ pan  •  ↵ return"
	}
	view.WriteString(theme.primaryStyle().Render(footer))
	return view.String()
}

func (m menuResultModel) renderOutputPanel(theme tuiTheme) string {
	lines := m.outputLines()
	end := min(len(lines), m.yOffset+m.resultRows())
	contentWidth := m.contentWidth()
	border := "│"
	if theme.color {
		border = theme.primaryStyle().Render(border)
	}
	var panel strings.Builder
	panel.WriteString(theme.primaryStyle().Render("╭─ OUTPUT " + strings.Repeat("─", max(1, contentWidth-7)) + "╮"))
	panel.WriteByte('\n')
	for index := m.yOffset; index < end; index++ {
		line := sliceResultLine(lines[index], m.xOffset, contentWidth)
		lineStyle := lipgloss.NewStyle()
		if index == 0 && theme.color {
			lineStyle = lineStyle.Bold(true).Foreground(theme.accent.primary)
		}
		panel.WriteString(border + " " + lineStyle.Render(line) + strings.Repeat(" ", max(0, contentWidth-lipgloss.Width(line))) + " " + border)
		panel.WriteByte('\n')
	}
	for index := end - m.yOffset; index < m.resultRows(); index++ {
		panel.WriteString(border + strings.Repeat(" ", contentWidth+2) + border)
		panel.WriteByte('\n')
	}
	panel.WriteString(theme.primaryStyle().Render("╰" + strings.Repeat("─", max(1, contentWidth+2)) + "╯"))
	return panel.String()
}

func sliceResultLine(line string, offset, width int) string {
	var result strings.Builder
	position, written := 0, 0
	for _, current := range line {
		runeWidth := lipgloss.Width(string(current))
		if position+runeWidth <= offset {
			position += runeWidth
			continue
		}
		if written+runeWidth > width {
			break
		}
		result.WriteRune(current)
		position += runeWidth
		written += runeWidth
	}
	return result.String()
}

func stripTerminalControls(value string) string {
	var clean strings.Builder
	for index := 0; index < len(value); {
		if value[index] == 0x1b {
			index++
			if index < len(value) && value[index] == '[' {
				index++
				for index < len(value) {
					current := value[index]
					index++
					if current >= 0x40 && current <= 0x7e {
						break
					}
				}
			} else if index < len(value) && value[index] == ']' {
				index++
				for index < len(value) {
					if value[index] == 0x07 {
						index++
						break
					}
					if value[index] == 0x1b && index+1 < len(value) && value[index+1] == '\\' {
						index += 2
						break
					}
					index++
				}
			} else if index < len(value) {
				index++
			}
			continue
		}
		current := value[index]
		index++
		if current == '\n' || current == '\t' || current >= 0x20 {
			clean.WriteByte(current)
		}
	}
	return clean.String()
}

func formatMenuCommand(args []string) string {
	parts := []string{"syskit"}
	for _, arg := range args {
		if strings.ContainsAny(arg, " \t") {
			parts = append(parts, strconv.Quote(arg))
		} else {
			parts = append(parts, arg)
		}
	}
	return strings.Join(parts, " ")
}

func executeMenuCommand(selection menuSelection, globalArgs []string) menuCommandResult {
	var stdout, stderr bytes.Buffer
	child, options := newRootCmdWithOptions()
	child.SetOut(&stdout)
	child.SetErr(&stderr)
	args := make([]string, 0, len(globalArgs)+len(selection.args)+1)
	for _, arg := range globalArgs {
		if !strings.HasPrefix(arg, "--color=") {
			args = append(args, arg)
		}
	}
	args = append(args, "--color=never")
	args = append(args, selection.args...)
	child.SetArgs(args)
	err := child.Execute()
	if message, _ := present(err); message != "" && !options.quiet {
		fmt.Fprintln(&stderr, message)
	}
	return menuCommandResult{stdout: stdout.String(), stderr: stderr.String(), err: err, code: exitCode(err)}
}

func runMenuCommandResult(stdin, stdout *os.File, selection menuSelection, globalArgs []string, color bool) error {
	return runMenuResult(stdin, stdout, selection, func() menuCommandResult {
		return executeMenuCommand(selection, globalArgs)
	}, color)
}

func runMenuCompletedResult(stdin, stdout *os.File, selection menuSelection, result menuCommandResult, color bool) error {
	return runMenuResult(stdin, stdout, selection, func() menuCommandResult { return result }, color)
}

func runMenuResult(stdin, stdout *os.File, selection menuSelection, run func() menuCommandResult, color bool) error {
	if !color {
		selection.accent = tuiAccent{icon: selection.accent.icon}
	}
	model := newMenuResultModel(selection, run)
	_, err := tea.NewProgram(model, tea.WithInput(stdin), tea.WithOutput(stdout), tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	if err != nil {
		return fmt.Errorf("showing command result: %w", err)
	}
	return nil
}
