package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestProcessListFiltersByUserName(t *testing.T) {
	s := NewProcess(processCollectorStub{list: &model.ProcessList{Processes: []model.Process{
		{PID: 1, UID: 0, User: "root"},
		{PID: 2, UID: 1000, User: "mersad"},
	}}})
	filters, err := ParseProcessFilters([]string{"user=mersad"})
	require.NoError(t, err)
	list, err := s.List(ProcessOptions{Filters: filters})
	require.NoError(t, err)
	require.Len(t, list.Processes, 1)
	assert.Equal(t, 2, list.Processes[0].PID)
}

type processCollectorStub struct {
	list *model.ProcessList
	err  error
}

func (s processCollectorStub) Collect() (*model.ProcessList, error) { return s.list, s.err }
