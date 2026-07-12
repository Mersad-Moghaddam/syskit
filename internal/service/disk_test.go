package service

import (
	"testing"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/stretchr/testify/assert"
)

type fakeDiskCollector struct{ info *model.DiskInfo }

func (c fakeDiskCollector) Collect() (*model.DiskInfo, error) { return c.info, nil }
func TestDiskCollect(t *testing.T) {
	want := &model.DiskInfo{}
	got, err := NewDisk(fakeDiskCollector{want}).Collect()
	assert.NoError(t, err)
	assert.Same(t, want, got)
}
