package command

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type containerServiceStub struct {
	list   *model.ContainerList
	detail *model.ContainerDetail
	err    error
}

func (s containerServiceStub) List() (*model.ContainerList, error)            { return s.list, s.err }
func (s containerServiceStub) Inspect(string) (*model.ContainerDetail, error) { return s.detail, s.err }

func TestContainerTable(t *testing.T) {
	table := containerTable(&model.ContainerList{Containers: []model.ContainerInfo{{ID: "abc", Runtime: "docker", PIDs: 2}}})
	assert.Equal(t, []string{"CONTAINER", "RUNTIME", "PIDS", "MEMORY", "CPU NS", "READ", "WRITE"}, table.Headers)
	assert.Equal(t, [][]string{{"abc", "docker", "2", "-", "-", "-", "-"}}, table.Rows)
}

func TestContainerInspectRendersTable(t *testing.T) {
	cmd := NewContainerCmd(containerServiceStub{detail: &model.ContainerDetail{ContainerInfo: model.ContainerInfo{ID: "abc", Runtime: "docker", PIDs: 1}, Processes: []model.Process{{PID: 7, Command: "worker"}}}}, ContainerOptions{Format: func() string { return "table" }, NoHeader: func() bool { return false }})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"inspect", "abc"})
	assert.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "CONTAINER")
	assert.Contains(t, out.String(), "worker")
}
