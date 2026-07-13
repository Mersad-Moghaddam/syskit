package cli

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWatchIterationClearsAndRenders(t *testing.T) {
	var out bytes.Buffer
	err := runWatchIteration([]string{"system"}, &out, func(args []string, output io.Writer) error {
		_, err := output.Write([]byte("fixture output\n"))
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "\x1b[H\x1b[2Jfixture output\n", out.String())
}

func TestWatchThemedIterationKeepsSelectedAccentShell(t *testing.T) {
	var out bytes.Buffer
	theme := tuiTheme{accent: paletteAccent(5), color: false}
	err := runWatchIterationThemed([]string{"network"}, &out, func(args []string, output io.Writer) error {
		_, err := output.Write([]byte("fixture output\n"))
		return err
	}, theme, 2*time.Second)
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "SYSKIT WATCH")
	assert.Contains(t, out.String(), "network")
	assert.Contains(t, out.String(), "refresh 2s")
	assert.Contains(t, out.String(), "fixture output")
}
