package cli

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
)

func TestTopSortKeysUpdateOptions(t *testing.T) {
	m := topModel{interval: time.Second}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	assert.Equal(t, "cpu", updated.(topModel).options.Sort)
	updated, _ = updated.(topModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	assert.Equal(t, "memory", updated.(topModel).options.Sort)
}

func TestTopScrollAndBackpressure(t *testing.T) {
	m := topModel{interval: time.Second, fetching: true, list: &model.ProcessList{Processes: []model.Process{{PID: 1}, {PID: 2}}}}
	updated, command := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, updated.(topModel).offset)
	assert.Nil(t, command)

	updated, command = updated.(topModel).Update(topTick{})
	assert.True(t, updated.(topModel).fetching)
	assert.NotNil(t, command)

	updated, _ = updated.(topModel).Update(topData{list: &model.ProcessList{Processes: []model.Process{{PID: 1}}}})
	assert.False(t, updated.(topModel).fetching)
	assert.Equal(t, 0, updated.(topModel).offset)
}

func TestTopViewUsesModernLiveLayout(t *testing.T) {
	memory := 12.5
	m := topModel{
		interval: time.Second,
		options:  service.ProcessOptions{Sort: "memory"},
		theme:    tuiTheme{accent: paletteAccent(3), color: false},
		list:     &model.ProcessList{Processes: []model.Process{{PID: 42, User: "fixture", Command: "worker", MemoryPercent: &memory}}},
	}
	view := m.View()
	assert.Contains(t, view, "SYSKIT TOP")
	assert.Contains(t, view, "LIVE PROCESS INTELLIGENCE")
	assert.Contains(t, view, "worker")
	assert.Contains(t, view, "refresh 1s")
}
