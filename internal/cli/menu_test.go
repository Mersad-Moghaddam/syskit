package cli

import (
	"sort"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func updateMenu(t *testing.T, model menuModel, msg tea.Msg) (menuModel, tea.Cmd) {
	t.Helper()
	updated, cmd := model.Update(msg)
	result, ok := updated.(menuModel)
	require.True(t, ok)
	return result, cmd
}

func TestInteractiveMenuCoversEveryCommandFamily(t *testing.T) {
	root := interactiveMenuTree()
	var paths []string
	var walk func(menuItem)
	walk = func(item menuItem) {
		if len(item.children) > 0 {
			for _, child := range item.children {
				walk(child)
			}
			return
		}
		args := item.args
		if item.prompt != nil {
			args = item.prompt.prefix
		}
		if len(args) == 0 {
			return
		}
		path := args[0]
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			switch args[0] {
			case "completion", "containers", "network", "plugins", "process":
				path += " " + args[1]
			}
		}
		paths = append(paths, path)
	}
	walk(root)
	sort.Strings(paths)

	for _, path := range []string{
		"completion bash", "completion fish", "completion powershell", "completion zsh",
		"containers", "containers inspect", "cpu", "dashboard", "diagnostics", "disk",
		"filesystem", "memory", "network", "network dns", "network interfaces",
		"network routes", "plugins inspect", "plugins list", "plugins run", "ports",
		"process", "process tree", "system", "top", "version", "watch",
	} {
		assert.Contains(t, paths, path, "menu must expose %q", path)
	}
}

func TestInteractiveMenuAssignsAnAccentAndIconToEveryOption(t *testing.T) {
	root := interactiveMenuTree()
	colors := map[string]bool{}
	var walk func(menuItem)
	walk = func(item menuItem) {
		assert.NotEmpty(t, item.accent.primary, item.title)
		assert.NotEmpty(t, item.accent.icon, item.title)
		colors[string(item.accent.primary)] = true
		for _, child := range item.children {
			walk(child)
		}
	}
	walk(root)
	assert.GreaterOrEqual(t, len(colors), 10)
}

func TestInteractiveMenuEntranceAnimationAdvancesAndCanBeSkipped(t *testing.T) {
	model := newMenuModel()
	require.NotNil(t, model.Init())
	model, command := updateMenu(t, model, menuIntroTick{})
	assert.Equal(t, 1, model.introFrame)
	require.NotNil(t, command)

	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, menuIntroFrames-1, model.introFrame)
	assert.Equal(t, 1, model.cursor, "the key that skips animation also performs its action")
}

func TestInteractiveMenuAnimationKeepsClickableRowsStable(t *testing.T) {
	model := newMenuModel()
	fullRow := model.itemStartRow()
	model.introFrame = menuIntroFrames - 1
	assert.Equal(t, fullRow, model.itemStartRow())

	model.width, model.height = 60, 12
	compactRow := model.itemStartRow()
	model.introFrame = 0
	assert.Equal(t, compactRow, model.itemStartRow())
}

func TestInteractiveMenuNavigatesIntoCPUAndBack(t *testing.T) {
	model := newMenuModel()

	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	assert.Equal(t, "System & hardware", model.current.title)
	assert.Equal(t, "Home  /  System & hardware", model.breadcrumb())

	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	assert.Equal(t, "CPU", model.current.title)
	assert.Equal(t, "CPU overview", model.current.children[0].title)
	assert.Equal(t, "Per-core utilization", model.current.children[1].title)

	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEsc})
	assert.Equal(t, "System & hardware", model.current.title)
	assert.Equal(t, 1, model.cursor, "returning restores the parent cursor")

	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyLeft})
	assert.Equal(t, "Home", model.current.title)
	assert.Equal(t, 0, model.cursor)
}

func TestInteractiveMenuSelectsPerCoreCPU(t *testing.T) {
	model := newMenuModel()
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	model, cmd := updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})

	require.NotNil(t, cmd)
	require.NotNil(t, model.selection)
	assert.Equal(t, "Per-core utilization", model.selection.title)
	assert.Equal(t, []string{"cpu", "--per-core"}, model.selection.args)
	assert.Equal(t, model.current.children[1].accent, model.selection.accent)
}

func TestInteractiveMenuReturnsToSameLocationWithoutReplayingIntro(t *testing.T) {
	model := newMenuModel()
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, model.selection)

	resumed := model.resumedAfterCommand()
	assert.Equal(t, "CPU", resumed.current.title)
	assert.Equal(t, 1, resumed.cursor)
	assert.Nil(t, resumed.selection)
	assert.Equal(t, menuIntroFrames-1, resumed.introFrame)
	assert.Nil(t, resumed.Init(), "the entrance animation only runs once")
}

func TestInteractiveMenuMouseSelectsVisibleRow(t *testing.T) {
	model := newMenuModel()
	click := tea.MouseMsg{X: 4, Y: model.itemStartRow() + 2, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress}
	model, _ = updateMenu(t, model, click)

	assert.Equal(t, "Processes", model.current.title)
	assert.Nil(t, model.selection)
}

func TestInteractiveMenuPromptValidatesAndBuildsArguments(t *testing.T) {
	model := newMenuModel()
	for range 6 {
		model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	}
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	assert.Equal(t, "Containers", model.current.title)
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, model.inputItem)

	model, cmd := updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	assert.Nil(t, cmd)
	assert.Equal(t, "A value is required.", model.inputError)

	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("abc123")})
	model, cmd = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEnter})
	require.NotNil(t, cmd)
	require.NotNil(t, model.selection)
	assert.Equal(t, []string{"containers", "inspect", "abc123"}, model.selection.args)
}

func TestInteractiveMenuPromptSupportsUnicodeBackspaceAndCancel(t *testing.T) {
	model := newMenuModel()
	item := &model.current.children[6].children[1]
	model.inputItem = item
	model.input = "شناسه"

	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyBackspace})
	assert.Equal(t, "شناس", model.input)
	model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEsc})
	assert.Nil(t, model.inputItem)
	assert.Empty(t, model.input)
}

func TestInteractiveMenuResizeKeepsSelectionVisible(t *testing.T) {
	model := newMenuModel()
	model, _ = updateMenu(t, model, tea.WindowSizeMsg{Width: 60, Height: 12})
	for range 9 {
		model, _ = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyDown})
	}
	assert.Equal(t, 9, model.cursor)
	assert.Equal(t, 7, model.offset)
	assert.Contains(t, model.View(), "8–10 of 10")
}

func TestInteractiveMenuQuitAndRootBackQuit(t *testing.T) {
	model := newMenuModel()
	_, cmd := updateMenu(t, model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	require.NotNil(t, cmd)

	_, cmd = updateMenu(t, model, tea.KeyMsg{Type: tea.KeyEsc})
	require.NotNil(t, cmd)
}

func TestInteractiveMenuViewDescribesControlsAndSelection(t *testing.T) {
	model := newMenuModel()
	view := model.View()
	assert.Contains(t, view, "SYSKIT")
	assert.Contains(t, view, "CONTROL CENTER")
	assert.Contains(t, view, "System & hardware")
	assert.Contains(t, view, "◉ mouse")
	assert.Contains(t, view, "Host identity and core resource summaries")
}

func TestInteractiveMenuCanRenderWithoutColor(t *testing.T) {
	view := newMenuModelWithColor(false).View()
	assert.NotContains(t, view, "\x1b[")
	assert.Contains(t, view, "CONTROL CENTER")
}

func TestChangedPersistentArgsPreservesExplicitRootOptions(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.PersistentFlags().String("format", "table", "")
	cmd.PersistentFlags().Bool("quiet", false, "")
	require.NoError(t, cmd.PersistentFlags().Set("format", "json"))
	require.NoError(t, cmd.PersistentFlags().Set("quiet", "true"))

	assert.ElementsMatch(t, []string{"--format=json", "--quiet=true"}, changedPersistentArgs(cmd))
}
