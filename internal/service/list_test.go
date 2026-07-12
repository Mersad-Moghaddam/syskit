package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type listItem struct {
	name string
	id   int
}

func TestListPrimitives(t *testing.T) {
	filter, err := ParseFilter("name=beta")
	require.NoError(t, err)
	items := []listItem{{"beta", 2}, {"alpha", 1}}
	filtered, err := FilterItems(items, []Filter{filter}, map[string]func(listItem) string{"name": func(v listItem) string { return v.name }})
	require.NoError(t, err)
	require.Len(t, filtered, 1)
	sorted, err := SortItems(items, "id", map[string]func(listItem, listItem) bool{"id": func(a, b listItem) bool { return a.id < b.id }}, false)
	require.NoError(t, err)
	assert.Equal(t, "alpha", sorted[0].name)
	limited, err := LimitItems(sorted, 1)
	require.NoError(t, err)
	assert.Len(t, limited, 1)
}
func TestParseFilterRejectsBadInput(t *testing.T) {
	_, err := ParseFilter("name")
	assert.Error(t, err)
	_, err = LimitItems([]int{1}, -1)
	assert.Error(t, err)
}
