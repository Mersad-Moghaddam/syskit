package platform_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// idleFS returns a TestFS rooted at the idle-host fixture set.
func idleFS(t *testing.T) platform.SysFS {
	t.Helper()
	return platform.TestFS(os.DirFS("testdata/fixtures/idle-host"))
}

func TestFS_ReadFileResolvesFixtureBytes(t *testing.T) {
	fsys := idleFS(t)

	got, err := fsys.ReadFile("proc/stat")
	require.NoError(t, err)

	want, err := os.ReadFile(filepath.Join("testdata", "fixtures", "idle-host", "proc", "stat"))
	require.NoError(t, err)

	assert.Equal(t, want, got, "ReadFile must return the exact fixture bytes")
}

func TestFS_ReadFileDistinctFixtureVariants(t *testing.T) {
	idle := platform.TestFS(os.DirFS("testdata/fixtures/idle-host"))
	busy := platform.TestFS(os.DirFS("testdata/fixtures/busy-host"))

	idleLoad, err := idle.ReadFile("proc/loadavg")
	require.NoError(t, err)
	busyLoad, err := busy.ReadFile("proc/loadavg")
	require.NoError(t, err)

	assert.NotEqual(t, idleLoad, busyLoad, "the two fixture variants must differ")
	assert.Equal(t, "0.00 0.01 0.05 1/234 5678\n", string(idleLoad))
	assert.Equal(t, "3.14 2.72 1.61 8/512 91234\n", string(busyLoad))
}

func TestFS_ReadDirEnumeratesFixtureDirectory(t *testing.T) {
	fsys := idleFS(t)

	entries, err := fsys.ReadDir("proc")
	require.NoError(t, err)

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	assert.Contains(t, names, "stat")
	assert.Contains(t, names, "loadavg")
}

func TestFS_OpenStreamsFile(t *testing.T) {
	fsys := idleFS(t)

	f, err := fsys.Open("proc/stat")
	require.NoError(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.NoError(t, err)
	assert.Contains(t, string(data), "cpu ")
}

func TestFS_MissingFileIsErrNotFound(t *testing.T) {
	fsys := idleFS(t)

	tests := []struct {
		name string
		call func() error
	}{
		{"ReadFile", func() error { _, err := fsys.ReadFile("proc/does-not-exist"); return err }},
		{"Open", func() error { _, err := fsys.Open("proc/does-not-exist"); return err }},
		{"ReadDir", func() error { _, err := fsys.ReadDir("proc/does-not-exist"); return err }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			require.Error(t, err)
			assert.ErrorIs(t, err, platform.ErrNotFound)
		})
	}
}

func TestFS_UnderlyingPathErrorStillInspectable(t *testing.T) {
	fsys := idleFS(t)

	_, err := fsys.ReadFile("proc/does-not-exist")
	require.Error(t, err)

	// The sentinel must be reachable...
	assert.ErrorIs(t, err, platform.ErrNotFound)
	// ...and so must the underlying *fs.PathError, per the documented
	// double-%w wrapping choice.
	var perr *fs.PathError
	assert.ErrorAs(t, err, &perr)
}

func TestFS_RejectsUnsafePaths(t *testing.T) {
	fsys := idleFS(t)

	unsafe := []string{
		"/proc/stat",             // absolute
		"../idle-host/proc/stat", // parent traversal
		"proc/../../etc/passwd",  // escaping traversal
		"proc//stat",             // non-canonical
		"./proc/stat",            // non-canonical
		"",                       // empty
	}
	for _, name := range unsafe {
		t.Run(name, func(t *testing.T) {
			_, rerr := fsys.ReadFile(name)
			require.Error(t, rerr)
			assert.ErrorIs(t, rerr, platform.ErrNotFound)

			_, oerr := fsys.Open(name)
			require.Error(t, oerr)
			assert.ErrorIs(t, oerr, platform.ErrNotFound)

			_, derr := fsys.ReadDir(name)
			require.Error(t, derr)
			assert.ErrorIs(t, derr, platform.ErrNotFound)
		})
	}
}

func TestRealFS_ReadsProcStat(t *testing.T) {
	// Sanity check that RealFS is rooted at "/" and reads a real zero-length
	// pseudo-file to EOF. /proc/stat exists on every Linux host, including CI.
	fsys := platform.RealFS()

	data, err := fsys.ReadFile("proc/stat")
	require.NoError(t, err)
	assert.NotEmpty(t, data, "zero-length pseudo-file must still be read to EOF")
	assert.Contains(t, string(data), "cpu")
}
