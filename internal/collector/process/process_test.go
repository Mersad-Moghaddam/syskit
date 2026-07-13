package process

import (
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
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

func TestParseCPUAndMemoryTotals(t *testing.T) {
	assert.Equal(t, uint64(21), ParseCPUTotal([]byte("cpu  1 2 3 4 5 6\ncpu0 1 2\n")))
	assert.Equal(t, uint64(4096), ParseMemoryTotal([]byte("MemTotal:       4 kB\n")))
}

func BenchmarkParseStat(b *testing.B) {
	data := []byte("4242 (database worker) S 1 0 0 0 0 0 0 0 0 0 700 300 0 0 0 0 12 0 99000\n")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseStat(data)
	}
}

func BenchmarkCollectorWalk1000Processes(b *testing.B) {
	files := fstest.MapFS{
		"proc/stat":    &fstest.MapFile{Data: []byte("cpu  1 2 3 4 5 6 7 8 9 10\n")},
		"proc/meminfo": &fstest.MapFile{Data: []byte("MemTotal: 16777216 kB\n")},
		"etc/passwd":   &fstest.MapFile{Data: []byte("root:x:0:0:root:/root:/bin/sh\n")},
	}
	for pid := 1; pid <= 1000; pid++ {
		base := fmt.Sprintf("proc/%d/", pid)
		files[base+"stat"] = &fstest.MapFile{Data: []byte(fmt.Sprintf("%d (worker) S 1 0 0 0 0 0 0 0 0 0 7 3 0 0 0 0 4 0 99\n", pid))}
		files[base+"status"] = &fstest.MapFile{Data: []byte("Uid:\t0\t0\t0\t0\nVmRSS:\t1024 kB\nThreads:\t4\n")}
		files[base+"cmdline"] = &fstest.MapFile{Data: []byte("worker\x00--serve\x00")}
	}
	collector := NewCollector(platform.TestFS(files))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = collector.Collect()
	}
}
