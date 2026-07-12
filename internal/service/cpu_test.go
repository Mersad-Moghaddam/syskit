package service

import (
	"testing"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/stretchr/testify/assert"
)

type fakeCPUCollector struct{ info *model.CPUInfo }

func (c fakeCPUCollector) Collect() (*model.CPUInfo, error) { return c.info, nil }
func TestCPUCollect(t *testing.T) {
	want := &model.CPUInfo{LogicalCores: 1}
	got, err := NewCPU(fakeCPUCollector{want}).Collect()
	assert.NoError(t, err)
	assert.Same(t, want, got)
}
