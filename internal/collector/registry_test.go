package collector_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/collector"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// stubCollect is a trivial CollectFunc used to populate registrations without
// pulling in a real domain.
func stubCollect(_ platform.SysFS) (any, error) { return "ok", nil }

func TestRegistry_RegisterAndLookup(t *testing.T) {
	r := collector.NewRegistry()

	require.NoError(t, r.Register(collector.Registration{
		Name:    "cpu",
		Summary: "processor stats",
		Collect: stubCollect,
	}))

	got, ok := r.Lookup("cpu")
	require.True(t, ok, "registered name must be found")
	assert.Equal(t, "cpu", got.Name)
	assert.Equal(t, "processor stats", got.Summary)

	_, ok = r.Lookup("memory")
	assert.False(t, ok, "unregistered name must not be found")
}

func TestRegistry_NamesAndAllAreSorted(t *testing.T) {
	r := collector.NewRegistry()
	for _, name := range []string{"network", "cpu", "memory"} {
		require.NoError(t, r.Register(collector.Registration{Name: name, Collect: stubCollect}))
	}

	assert.Equal(t, []string{"cpu", "memory", "network"}, r.Names())

	all := r.All()
	require.Len(t, all, 3)
	assert.Equal(t, "cpu", all[0].Name)
	assert.Equal(t, "memory", all[1].Name)
	assert.Equal(t, "network", all[2].Name)
}

func TestRegistry_DuplicateRegistrationIsRejected(t *testing.T) {
	r := collector.NewRegistry()
	require.NoError(t, r.Register(collector.Registration{Name: "cpu", Collect: stubCollect}))

	err := r.Register(collector.Registration{Name: "cpu", Collect: stubCollect})
	require.Error(t, err)
	assert.True(t, errors.Is(err, collector.ErrAlreadyRegistered),
		"duplicate name must return ErrAlreadyRegistered, got %v", err)

	// The original registration is preserved (no silent overwrite).
	assert.Len(t, r.Names(), 1)
}

func TestRegistry_InvalidRegistrationIsRejected(t *testing.T) {
	r := collector.NewRegistry()

	err := r.Register(collector.Registration{Name: "", Collect: stubCollect})
	assert.True(t, errors.Is(err, collector.ErrInvalidRegistration),
		"empty name must return ErrInvalidRegistration, got %v", err)

	err = r.Register(collector.Registration{Name: "cpu", Collect: nil})
	assert.True(t, errors.Is(err, collector.ErrInvalidRegistration),
		"nil Collect must return ErrInvalidRegistration, got %v", err)

	assert.Empty(t, r.Names(), "no invalid registration should be stored")
}

func TestRegistry_IndependentInstances(t *testing.T) {
	a := collector.NewRegistry()
	b := collector.NewRegistry()

	require.NoError(t, a.Register(collector.Registration{Name: "cpu", Collect: stubCollect}))

	_, ok := b.Lookup("cpu")
	assert.False(t, ok, "separate registries must not share state (no globals)")
}

// typedReading is a stand-in domain snapshot for exercising Adapt.
type typedReading struct{ value int }

type typedCollector struct{ fs platform.SysFS }

func (typedCollector) Collect() (typedReading, error) { return typedReading{value: 7}, nil }

func TestAdapt_ErasesTypedCollector(t *testing.T) {
	fn := collector.Adapt(func(fs platform.SysFS) collector.Collector[typedReading] {
		return typedCollector{fs: fs}
	})

	got, err := fn(nil)
	require.NoError(t, err)

	reading, ok := got.(typedReading)
	require.True(t, ok, "erased value must recover to the typed snapshot")
	assert.Equal(t, 7, reading.value)
}
