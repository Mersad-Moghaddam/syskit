package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type fakeMemoryCollector struct{ info *model.MemoryInfo }

func (c fakeMemoryCollector) Collect() (*model.MemoryInfo, error) { return c.info, nil }
func TestMemoryCollect(t *testing.T) {
	want := &model.MemoryInfo{TotalBytes: 1}
	got, err := NewMemory(fakeMemoryCollector{want}).Collect()
	assert.NoError(t, err)
	assert.Same(t, want, got)
}
