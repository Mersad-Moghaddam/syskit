package cli

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestTopSortKeysUpdateOptions(t *testing.T) {
	m := topModel{interval: time.Second}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	assert.Equal(t, "cpu", updated.(topModel).options.Sort)
	updated, _ = updated.(topModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	assert.Equal(t, "memory", updated.(topModel).options.Sort)
}
