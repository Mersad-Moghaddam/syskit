package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatHandlesSpacesAndParentheses(t *testing.T) {
	p, err := ParseStat([]byte("42 (worker (main)) S 1 0 0 0 0 0 0 0 0 0 7 3 0 0 0 0 4 0 0\n"))
	assert.NoError(t, err)
	assert.Equal(t, 42, p.PID)
	assert.Equal(t, 1, p.PPID)
	assert.Equal(t, "worker (main)", p.Command)
	assert.Equal(t, uint64(10), p.CPUTime)
	assert.Equal(t, uint64(4), p.Threads)
}
