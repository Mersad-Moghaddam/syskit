package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestCPUTableMarksAbsentTopologyUnavailable(t *testing.T) {
	table := cpuTable(&model.CPUInfo{LogicalCores: 2, Model: "test", Architecture: "amd64"}, false)
	assert.Equal(t, "unavailable", table.Rows[0][2])
	assert.Equal(t, "2", table.Rows[0][3])
}
