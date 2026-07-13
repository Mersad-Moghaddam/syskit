package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type fakeCPUCollector struct{ info *model.CPUInfo }

func (c fakeCPUCollector) Collect() (*model.CPUInfo, error) { return c.info, nil }
func TestCPUCollect(t *testing.T) {
	want := &model.CPUInfo{LogicalCores: 1}
	got, err := NewCPU(fakeCPUCollector{want}).Collect()
	assert.NoError(t, err)
	assert.Same(t, want, got)
}

type sequenceCPUCollector struct {
	infos []*model.CPUInfo
	next  int
}

func (c *sequenceCPUCollector) Collect() (*model.CPUInfo, error) {
	result := c.infos[c.next]
	c.next++
	return result, nil
}
func TestCPUSampleDerivesUtilization(t *testing.T) {
	c := &sequenceCPUCollector{infos: []*model.CPUInfo{{Times: []model.CPUTime{{CPUID: "all", Idle: 50, IOWait: 5, Total: 100}}}, {Times: []model.CPUTime{{CPUID: "all", Idle: 60, IOWait: 5, Total: 130}}}}}
	info, err := NewCPUWithSleep(c, func(time.Duration) {}).Sample(time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, info.Times[0].Utilization)
	assert.InDelta(t, 66.666, *info.Times[0].Utilization, .01)
}
