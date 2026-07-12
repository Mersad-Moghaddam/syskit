package command

import (
	"testing"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCPUTableMarksAbsentTopologyUnavailable(t *testing.T) {
	table := cpuTable(&model.CPUInfo{LogicalCores: 2, Model: "test", Architecture: "amd64"})
	assert.Equal(t, "unavailable", table.Rows[0][1])
	assert.Equal(t, "2", table.Rows[0][2])
}
