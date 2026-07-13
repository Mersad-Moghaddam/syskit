package cli

import (
	"bytes"
	"io"
	"testing"

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
