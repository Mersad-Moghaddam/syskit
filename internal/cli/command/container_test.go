package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestContainerTable(t *testing.T) {
	table := containerTable(&model.ContainerList{Containers: []model.ContainerInfo{{ID: "abc", Runtime: "docker", PIDs: 2}}})
	assert.Equal(t, []string{"CONTAINER", "RUNTIME", "PIDS"}, table.Headers)
	assert.Equal(t, [][]string{{"abc", "docker", "2"}}, table.Rows)
}
