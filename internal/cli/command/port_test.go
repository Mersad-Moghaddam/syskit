package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestPortTableIncludesSocketOwnersAndInode(t *testing.T) {
	table := portTable(&model.PortInfo{Sockets: []model.Socket{{
		Protocol: "tcp", LocalAddress: "127.0.0.1", LocalPort: 8080,
		RemoteAddress: "0.0.0.0", State: "LISTEN", Inode: 42,
		Owners: []model.SocketOwner{{PID: 123, Command: "server --listen"}},
	}}})

	assert.Equal(t, []string{"PROTO", "LOCAL", "REMOTE", "STATE", "INODE", "PID", "COMMAND"}, table.Headers)
	assert.Equal(t, []string{"tcp", "127.0.0.1:8080", "0.0.0.0", "LISTEN", "42", "123", "server --listen"}, table.Rows[0])
}
