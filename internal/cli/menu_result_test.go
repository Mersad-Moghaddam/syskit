package cli

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func updateMenuResult(t *testing.T, model menuResultModel, msg tea.Msg) (menuResultModel, tea.Cmd) {
	t.Helper()
	updated, cmd := model.Update(msg)
	result, ok := updated.(menuResultModel)
	require.True(t, ok)
	return result, cmd
}

func TestMenuResultShowsThemedLoadingAndSuccess(t *testing.T) {
	selection := menuSelection{title: "CPU overview", args: []string{"cpu"}, accent: paletteAccent(2)}
	assert.Equal(t, selection.accent.primary, (tuiTheme{accent: selection.accent, color: true}).primaryStyle().GetForeground())
	model := newMenuResultModel(selection, func() menuCommandResult { return menuCommandResult{} })
	loading := model.View()
	assert.Contains(t, loading, "CPU overview")
	assert.Contains(t, loading, "collecting live Linux data")
	assert.Contains(t, loading, "syskit cpu")

	model, _ = updateMenuResult(t, model, menuResultReady{result: menuCommandResult{stdout: "CPU  UTIL\nall  12%\n"}})
	view := model.View()
	assert.Contains(t, view, "SUCCESS")
	assert.Contains(t, view, "CPU  UTIL")
	assert.Contains(t, view, "↵/esc return")
}

func TestMenuResultShowsFailureAndDiagnostics(t *testing.T) {
	selection := menuSelection{title: "Inspect plugin", args: []string{"plugins", "inspect", "missing"}, accent: paletteAccent(6)}
	model := newMenuResultModel(selection, func() menuCommandResult { return menuCommandResult{} })
	model, _ = updateMenuResult(t, model, menuResultReady{result: menuCommandResult{stderr: "plugin not found\n", err: errors.New("missing"), code: 1}})
	view := model.View()
	assert.Contains(t, view, "FAILED · EXIT 1")
	assert.Contains(t, view, "plugin not found")
}

func TestMenuResultScrollsVerticallyAndHorizontally(t *testing.T) {
	selection := menuSelection{title: "Processes", args: []string{"process"}, accent: paletteAccent(4)}
	model := newMenuResultModel(selection, func() menuCommandResult { return menuCommandResult{} })
	model.height, model.width = 13, 24
	model, _ = updateMenuResult(t, model, menuResultReady{result: menuCommandResult{stdout: "header-very-wide-column\nrow-one-very-wide-column\nrow-two-very-wide-column\nrow-three-very-wide-column\nrow-four-very-wide-column\n"}})
	model, _ = updateMenuResult(t, model, tea.KeyMsg{Type: tea.KeyPgDown})
	assert.Greater(t, model.yOffset, 0)
	model, _ = updateMenuResult(t, model, tea.KeyMsg{Type: tea.KeyRight})
	assert.Equal(t, 4, model.xOffset)
	assert.Contains(t, model.View(), "column 5")
}

func TestMenuResultMouseWheelScrolls(t *testing.T) {
	model := newMenuResultModel(menuSelection{accent: paletteAccent(1)}, func() menuCommandResult { return menuCommandResult{} })
	model.height = 12
	model, _ = updateMenuResult(t, model, menuResultReady{result: menuCommandResult{stdout: strings.Repeat("line\n", 12)}})
	model, _ = updateMenuResult(t, model, tea.MouseMsg{Button: tea.MouseButtonWheelDown, Action: tea.MouseActionPress})
	assert.Equal(t, 3, model.yOffset)
}

func TestMenuResultStripsTerminalControlSequences(t *testing.T) {
	assert.Equal(t, "red plain", stripTerminalControls("\x1b[31mred\x1b[0m\r plain\x07"))
	assert.Equal(t, "safe", stripTerminalControls("\x1b]0;malicious title\x07safe"))
	assert.Equal(t, "界", sliceResultLine("ab界cd", 2, 2))
}

func TestMenuResultCanRenderWithoutColor(t *testing.T) {
	selection := menuSelection{title: "System", args: []string{"system"}, accent: tuiAccent{icon: "◈"}}
	model := newMenuResultModel(selection, func() menuCommandResult { return menuCommandResult{} })
	model, _ = updateMenuResult(t, model, menuResultReady{result: menuCommandResult{stdout: "HOST\nfixture\n"}})
	assert.NotContains(t, model.View(), "\x1b[")
}

func TestMenuResultUsesCompactControlsAtNarrowWidths(t *testing.T) {
	selection := menuSelection{title: "A very long process command", args: []string{"process"}, accent: tuiAccent{icon: "◈"}}
	model := newMenuResultModel(selection, func() menuCommandResult { return menuCommandResult{} })
	model.width, model.height = 36, 12
	model, _ = updateMenuResult(t, model, menuResultReady{result: menuCommandResult{stdout: "HEADER\nvalue\n"}})
	view := model.View()
	assert.Contains(t, view, "←→ pan")
	assert.NotContains(t, view, "SYSKIT / RESULT")
	for _, line := range strings.Split(view, "\n") {
		assert.LessOrEqual(t, lipgloss.Width(line), 36)
	}
}

func TestExecuteMenuCommandUsesExistingCobraTree(t *testing.T) {
	result := executeMenuCommand(menuSelection{args: []string{"version"}}, nil)
	require.NoError(t, result.err)
	assert.Equal(t, exitOK, result.code)
	assert.Equal(t, version+"\n", result.stdout)
}

func TestFormatMenuCommandQuotesValuesWithSpaces(t *testing.T) {
	assert.Equal(t, `syskit plugins inspect "two words"`, formatMenuCommand([]string{"plugins", "inspect", "two words"}))
}
