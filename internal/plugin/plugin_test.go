package plugin

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validManifest = `{"name":"example","version":"1.0.0","api_version":"v1","executable":"plugin","collectors":["example"],"permissions":[],"output_schemas":{"example":"object"},"author":"tester","license":"MIT"}`

func TestDiscoverValidatesCompatibility(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "example"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "example", "manifest.json"), []byte(validManifest), 0644))
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

func TestInspectFindsPlugin(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "example"), 0755))
	manifest := strings.Replace(validManifest, `"permissions":[]`, `"permissions":["procfs"]`, 1)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "example", "manifest.json"), []byte(manifest), 0644))
	info, err := Inspect([]string{dir}, "example")
	require.NoError(t, err)
	assert.Equal(t, []string{"procfs"}, info.Permissions)
	_, err = Inspect([]string{dir}, "missing")
	assert.EqualError(t, err, "plugin \"missing\" not found")
}

func TestRunExecutesCompatibleJSONPlugin(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "example")
	require.NoError(t, os.Mkdir(pluginDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "manifest.json"), []byte(validManifest), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "plugin"), []byte("#!/bin/sh\nread request\nprintf '{\"ok\":true}'\n"), 0755))
	value, err := Run(context.Background(), []string{dir}, "example")
	require.NoError(t, err)
	assert.Equal(t, true, value.(map[string]any)["ok"])
}

func TestDiscoverRejectsIncompleteManifest(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "example"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "example", "manifest.json"), []byte(`{"name":"example","version":"1.0.0","api_version":"v1"}`), 0644))
	_, err := Discover([]string{dir})
	assert.ErrorContains(t, err, "requires executable, collectors, permissions, output_schemas, author, and license")
}
