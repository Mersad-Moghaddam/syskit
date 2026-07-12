package port

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSocketTable(t *testing.T) {
	sockets, err := ParseSocketTable([]byte("  sl local_address rem_address st tx rx tr tm->when retrnsmt uid timeout inode\n   0: 0100007F:1F90 00000000:0000 0A 0 0 0 0 0 0 1000 0 42\n"), "tcp")
	assert.NoError(t, err)
	assert.Len(t, sockets, 1)
	assert.Equal(t, "127.0.0.1", sockets[0].LocalAddress)
	assert.Equal(t, uint16(8080), sockets[0].LocalPort)
	assert.Equal(t, "LISTEN", sockets[0].State)
}
