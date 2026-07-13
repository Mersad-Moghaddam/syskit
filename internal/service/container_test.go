package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestContainerListGroupsAssociatedProcesses(t *testing.T) {
	s := NewContainer(processCollectorStub{list: &model.ProcessList{Processes: []model.Process{
		{PID: 1, ContainerID: "b", ContainerRuntime: "docker"},
		{PID: 2, ContainerID: "a", ContainerRuntime: "containerd"},
		{PID: 3, ContainerID: "b", ContainerRuntime: "docker"},
		{PID: 4},
	}}})

	list, err := s.List()
	require.NoError(t, err)
	require.Len(t, list.Containers, 2)
	assert.Equal(t, model.ContainerInfo{ID: "a", Runtime: "containerd", PIDs: 1}, list.Containers[0])
	assert.Equal(t, model.ContainerInfo{ID: "b", Runtime: "docker", PIDs: 2}, list.Containers[1])
}
