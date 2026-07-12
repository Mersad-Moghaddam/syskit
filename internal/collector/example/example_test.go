package example_test

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/example"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// fixtureFS returns a SysFS rooted at the named example fixture directory.
func fixtureFS(t *testing.T, variant string) platform.SysFS {
	t.Helper()
	return platform.TestFS(os.DirFS("testdata/" + variant))
}

func TestCollect_ValidFixture(t *testing.T) {
	c := example.NewCollector(fixtureFS(t, "valid"))

	got, err := c.Collect()
	require.NoError(t, err)

	assert.InDelta(t, 0.42, got.One, 1e-9)
	assert.InDelta(t, 0.35, got.Five, 1e-9)
	assert.InDelta(t, 0.30, got.Fifteen, 1e-9)
	assert.True(t, got.EntitiesKnown, "optional runnable/total token was present")
	assert.Equal(t, 1, got.Running)
	assert.Equal(t, 234, got.Total)
}

func TestCollect_OptionalMissingIsPartialNotError(t *testing.T) {
	c := example.NewCollector(fixtureFS(t, "no-entities"))

	got, err := c.Collect()
	require.NoError(t, err, "missing OPTIONAL data must not error")

	assert.InDelta(t, 0.10, got.One, 1e-9)
	assert.False(t, got.EntitiesKnown, "absent optional token is unavailable, not an error")
	assert.Zero(t, got.Running)
	assert.Zero(t, got.Total)
}

func TestCollect_MalformedFieldIsErrParse(t *testing.T) {
	c := example.NewCollector(fixtureFS(t, "malformed"))

	_, err := c.Collect()
	require.Error(t, err)
	assert.True(t, errors.Is(err, collector.ErrParse),
		"non-numeric load field must classify as ErrParse, got %v", err)
}

func TestCollect_MissingRequiredFieldIsErrFieldMissing(t *testing.T) {
	c := example.NewCollector(fixtureFS(t, "missing-field"))

	_, err := c.Collect()
	require.Error(t, err)
	assert.True(t, errors.Is(err, collector.ErrFieldMissing),
		"fewer than 3 load fields must classify as ErrFieldMissing, got %v", err)
}

func TestCollect_PlatformErrorPassesThrough(t *testing.T) {
	// A fixture directory with no proc/loadavg yields platform.ErrNotFound,
	// which the collector wraps for context but does not reclassify.
	c := example.NewCollector(platform.TestFS(os.DirFS(t.TempDir())))

	_, err := c.Collect()
	require.Error(t, err)
	assert.True(t, errors.Is(err, platform.ErrNotFound),
		"missing interface must surface platform.ErrNotFound, got %v", err)
}

func TestRegister_MakesCollectorDiscoverable(t *testing.T) {
	r := collector.NewRegistry()
	require.NoError(t, example.Register(r))

	reg, ok := r.Lookup("example")
	require.True(t, ok, "registered example collector must be discoverable by name")

	got, err := reg.Collect(fixtureFS(t, "valid"))
	require.NoError(t, err)

	reading, ok := got.(example.Reading)
	require.True(t, ok, "erased snapshot must recover to example.Reading")
	assert.InDelta(t, 0.42, reading.One, 1e-9)
}

// TestCollect_ConcurrentIsRaceFree proves the collector holds no shared mutable
// state: many goroutines share one collector and one SysFS and must all get the
// same result. Run under -race.
func TestCollect_ConcurrentIsRaceFree(t *testing.T) {
	c := example.NewCollector(fixtureFS(t, "valid"))

	const goroutines = 64
	var wg sync.WaitGroup
	wg.Add(goroutines)
	results := make([]example.Reading, goroutines)
	errs := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = c.Collect()
		}(i)
	}
	wg.Wait()

	for i := 0; i < goroutines; i++ {
		require.NoErrorf(t, errs[i], "goroutine %d", i)
		assert.Equal(t, results[0], results[i], "all goroutines must agree")
	}
}
