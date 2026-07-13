package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatHandlesSpacesAndParentheses(t *testing.T) {
	p, err := ParseStat([]byte("42 (worker (main)) S 1 0 0 0 0 0 0 0 0 0 7 3 0 0 0 0 4 0 99\n"))
	assert.NoError(t, err)
	assert.Equal(t, 42, p.PID)
	assert.Equal(t, 1, p.PPID)
	assert.Equal(t, "worker (main)", p.Command)
	assert.Equal(t, uint64(10), p.CPUTime)
	assert.Equal(t, uint64(4), p.Threads)
	assert.Equal(t, uint64(99), p.StartTimeTicks)
}

func TestParsePasswd(t *testing.T) {
	users := ParsePasswd([]byte("root:x:0:0:root:/root:/bin/sh\nmersad:x:1000:1000::/home/mersad:/bin/bash\nbad\n"))
	assert.Equal(t, map[uint64]string{0: "root", 1000: "mersad"}, users)
}
