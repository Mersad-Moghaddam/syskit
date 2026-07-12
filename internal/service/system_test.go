package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type fakeSystemCollector struct {
	info *model.SystemInfo
	err  error
}

func (c fakeSystemCollector) Collect() (*model.SystemInfo, error) { return c.info, c.err }

func TestSystemCollect(t *testing.T) {
	want := &model.SystemInfo{Hostname: "fixture"}
	got, err := NewSystem(fakeSystemCollector{info: want}).Collect()
	assert.NoError(t, err)
	assert.Same(t, want, got)
}

func TestSystemCollectPropagatesError(t *testing.T) {
	want := errors.New("unavailable")
	_, err := NewSystem(fakeSystemCollector{err: want}).Collect()
	assert.ErrorIs(t, err, want)
}
