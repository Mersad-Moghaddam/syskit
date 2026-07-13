package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestPortListFiltersByOwnerPID(t *testing.T) {
	s := NewPort(portCollectorStub{info: &model.PortInfo{Sockets: []model.Socket{
		{LocalPort: 8080, Owners: []model.SocketOwner{{PID: 100}, {PID: 200}}},
		{LocalPort: 9090, Owners: []model.SocketOwner{{PID: 300}}},
	}}})

	info, err := s.List(PortOptions{PID: 200})
	require.NoError(t, err)
	require.Len(t, info.Sockets, 1)
	assert.Equal(t, uint16(8080), info.Sockets[0].LocalPort)
}

func TestPortListRejectsNegativePID(t *testing.T) {
	s := NewPort(portCollectorStub{})
	_, err := s.List(PortOptions{PID: -1})
	require.EqualError(t, err, "PID must not be negative")
}

type portCollectorStub struct {
	info *model.PortInfo
	err  error
}

func (s portCollectorStub) Collect() (*model.PortInfo, error) { return s.info, s.err }
