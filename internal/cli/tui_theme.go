package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// tuiAccent is a semantic color pair used by one menu action and every view it
// launches. Values are ANSI-256 colors so the design remains useful over SSH
// and on terminals without true-color support.
type tuiAccent struct {
	name      string
	primary   lipgloss.Color
	secondary lipgloss.Color
	icon      string
}

type tuiTheme struct {
	accent tuiAccent
	color  bool
}

var tuiPalette = []tuiAccent{
	{name: "electric cyan", primary: lipgloss.Color("45"), secondary: lipgloss.Color("23"), icon: "◈"},
	{name: "violet", primary: lipgloss.Color("141"), secondary: lipgloss.Color("54"), icon: "◆"},
	{name: "coral", primary: lipgloss.Color("203"), secondary: lipgloss.Color("52"), icon: "●"},
	{name: "amber", primary: lipgloss.Color("214"), secondary: lipgloss.Color("58"), icon: "◉"},
	{name: "mint", primary: lipgloss.Color("120"), secondary: lipgloss.Color("22"), icon: "✦"},
	{name: "sky", primary: lipgloss.Color("81"), secondary: lipgloss.Color("24"), icon: "◇"},
	{name: "pink", primary: lipgloss.Color("213"), secondary: lipgloss.Color("53"), icon: "⬢"},
	{name: "lime", primary: lipgloss.Color("155"), secondary: lipgloss.Color("28"), icon: "◐"},
	{name: "blue", primary: lipgloss.Color("75"), secondary: lipgloss.Color("18"), icon: "▣"},
	{name: "magenta", primary: lipgloss.Color("207"), secondary: lipgloss.Color("55"), icon: "✺"},
	{name: "gold", primary: lipgloss.Color("227"), secondary: lipgloss.Color("94"), icon: "◒"},
	{name: "orange", primary: lipgloss.Color("208"), secondary: lipgloss.Color("58"), icon: "⬡"},
}

var defaultTUITheme = tuiTheme{accent: tuiPalette[1], color: true}

func resolveTUITheme(theme *tuiTheme) tuiTheme {
	if theme == nil || theme.accent.primary == "" {
		return defaultTUITheme
	}
	return *theme
}

func paletteAccent(index int) tuiAccent {
	if index < 0 {
		index = -index
	}
	return tuiPalette[index%len(tuiPalette)]
}

func (t tuiTheme) primaryStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if t.color {
		style = style.Foreground(t.accent.primary)
	}
	return style
}

func (t tuiTheme) selectedStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if t.color {
		style = style.Bold(true).Foreground(lipgloss.Color("230")).Background(t.accent.secondary)
	}
	return style
}

func (t tuiTheme) badge(label string) string {
	style := lipgloss.NewStyle()
	label = "▰ " + label + " ▰"
	if t.color {
		style = style.Bold(true).Foreground(lipgloss.Color("230")).Background(t.accent.secondary)
	}
	return style.Render(label)
}

func (t tuiTheme) borderStyle() lipgloss.Style {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	if t.color {
		style = style.BorderForeground(t.accent.primary)
	}
	return style
}

func (t tuiTheme) mutedStyle() lipgloss.Style {
	style := lipgloss.NewStyle()
	if t.color {
		style = style.Foreground(lipgloss.Color("245"))
	}
	return style
}

const menuIntroFrames = 9

var syskitLogo = []string{
	"███████╗██╗   ██╗███████╗██╗  ██╗██╗████████╗",
	"██╔════╝╚██╗ ██╔╝██╔════╝██║ ██╔╝██║╚══██╔══╝",
	"███████╗ ╚████╔╝ ███████╗█████╔╝ ██║   ██║   ",
	"╚════██║  ╚██╔╝  ╚════██║██╔═██╗ ██║   ██║   ",
	"███████║   ██║   ███████║██║  ██╗██║   ██║   ",
	"╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝   ╚═╝   ",
}

func renderSysKitLogo(frame, width int, color, compact bool) string {
	if compact {
		label := "◢  S Y S K I T  ◣"
		if frame < 2 {
			label = "◢  S Y S · ·  ◣"
		}
		style := lipgloss.NewStyle()
		if color {
			style = style.Bold(true).Foreground(paletteAccent(frame).primary)
		}
		return style.Render(label)
	}

	var logo strings.Builder
	revealed := min(len(syskitLogo), frame+1)
	for index, line := range syskitLogo {
		if index < revealed {
			style := lipgloss.NewStyle()
			if color {
				style = style.Bold(true).Foreground(paletteAccent(index + frame).primary)
			}
			logo.WriteString(style.Render(line))
		} else {
			logo.WriteString(strings.Repeat(" ", min(width, lipgloss.Width(line))))
		}
		if index < len(syskitLogo)-1 {
			logo.WriteByte('\n')
		}
	}
	return logo.String()
}

func renderTUITagline(frame, width int, color bool) string {
	spinner := []string{"◜", "◠", "◝", "◞", "◡", "◟"}[frame%6]
	if frame >= menuIntroFrames-1 {
		spinner = "●"
	}
	left := lipgloss.NewStyle()
	right := lipgloss.NewStyle()
	if color {
		left = left.Bold(true).Foreground(paletteAccent(frame).primary)
		right = right.Foreground(lipgloss.Color("245"))
	}
	description := "native Linux intelligence  •  read-only  •  zero shell-outs"
	if width > 0 && width < 100 {
		description = "native Linux  •  read-only"
	}
	return fmt.Sprintf("%s  %s", left.Render(spinner+" SYSKIT // CONTROL CENTER"), right.Render(description))
}
