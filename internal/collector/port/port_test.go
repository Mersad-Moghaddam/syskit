package port

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

func TestParseSocketTable(t *testing.T) {
	sockets, err := ParseSocketTable([]byte("  sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode\n   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 42\n"), "tcp")
	assert.NoError(t, err)
	assert.Len(t, sockets, 1)
	assert.Equal(t, "127.0.0.1", sockets[0].LocalAddress)
	assert.Equal(t, uint16(8080), sockets[0].LocalPort)
	assert.Equal(t, "LISTEN", sockets[0].State)
}

func TestParseSocketTableIPv6(t *testing.T) {
	sockets, err := ParseSocketTable([]byte("  sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode\n   0: 00000000000000000000000001000000:1F90 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 42\n"), "tcp6")
	require.NoError(t, err)
	require.Len(t, sockets, 1)
	assert.Equal(t, "::1", sockets[0].LocalAddress)
	assert.Equal(t, "0A", sockets[0].RawState)
}

func TestParseUnixSocketTable(t *testing.T) {
	sockets, err := ParseUnixSocketTable([]byte("Num       RefCount Protocol Flags    Type St Inode Path\n0000000000000000: 00000002 00000000 00010000 0001 01 42 /run/example.sock\n"))
	require.NoError(t, err)
	require.Len(t, sockets, 1)
	assert.Equal(t, "unix", sockets[0].Protocol)
	assert.Equal(t, "LISTEN", sockets[0].State)
	assert.Equal(t, "/run/example.sock", sockets[0].LocalAddress)
}

func TestCollectorMapsSocketOwners(t *testing.T) {
	fsys := socketFS{MapFS: fstest.MapFS{
		"proc/net/tcp":     &fstest.MapFile{Data: []byte("  sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode\n   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 42\n")},
		"proc/123/cmdline": &fstest.MapFile{Data: []byte("server\x00--listen\x00")},
		"proc/123/fd/4":    &fstest.MapFile{},
		"proc/456/cmdline": &fstest.MapFile{Data: []byte("worker\x00")},
		"proc/456/fd/7":    &fstest.MapFile{},
	}, links: map[string]string{"proc/123/fd/4": "socket:[42]", "proc/456/fd/7": "socket:[42]"}}
	platformFS := platform.TestFS(fsys)
	entries, err := platformFS.ReadDir("proc")
	require.NoError(t, err)
	require.Len(t, entries, 3)
	fds, err := platformFS.ReadDir("proc/123/fd")
	require.NoError(t, err)
	require.Len(t, fds, 1)
	target, err := platformFS.ReadLink("proc/123/fd/4")
	require.NoError(t, err)
	require.Equal(t, "socket:[42]", target)
	info, err := NewCollector(platformFS).Collect()
	require.NoError(t, err)
	require.Len(t, info.Sockets, 1)
	require.Len(t, info.Sockets[0].Owners, 2)
	assert.Equal(t, []int{123, 456}, []int{info.Sockets[0].Owners[0].PID, info.Sockets[0].Owners[1].PID})
	assert.Equal(t, "server --listen", info.Sockets[0].Owners[0].Command)
}

func TestCollectorReportsPartialOwnerMapping(t *testing.T) {
	fsys := socketFS{MapFS: fstest.MapFS{
		"proc/net/tcp":  &fstest.MapFile{Data: []byte("  sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode\n   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000 0 0 42\n")},
		"proc/123/fd/4": &fstest.MapFile{},
	}, links: map[string]string{}, linkErrors: map[string]error{"proc/123/fd/4": fs.ErrPermission}}
	info, err := NewCollector(platform.TestFS(fsys)).Collect()
	require.NoError(t, err)
	assert.True(t, info.OwnerMappingPartial)
}

type socketFS struct {
	fstest.MapFS
	links      map[string]string
	linkErrors map[string]error
}

func (s socketFS) ReadLink(name string) (string, error) {
	if err, ok := s.linkErrors[name]; ok {
		return "", err
	}
	target, ok := s.links[name]
	if !ok {
		return "", fs.ErrNotExist
	}
	return target, nil
}
