package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverValidatesCompatibility(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "example"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "example", "manifest.json"), []byte(`{"name":"example","version":"1.0.0","api_version":"v1"}`), 0644))
	items, err := Discover([]string{dir})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "compatible", items[0].Status)
}

func TestDefaultDirsUsesEnvironment(t *testing.T) {
	t.Setenv("SYSKIT_PLUGIN_DIR", "/one"+string(os.PathListSeparator)+"/two")
	t.Setenv("XDG_DATA_HOME", "/data")
	assert.Equal(t, []string{"/one", "/two", filepath.Join("/data", "syskit", "plugins")}, DefaultDirs())
}
