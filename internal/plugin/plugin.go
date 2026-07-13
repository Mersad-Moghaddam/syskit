// Package plugin discovers declarative, out-of-process SysKit plugins. It
// never executes plugin binaries; execution remains an explicit later action.
package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const APIVersion = "v1"

type Manifest struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	APIVersion  string   `json:"api_version"`
	Permissions []string `json:"permissions,omitempty"`
}

type Info struct {
	Manifest
	Path   string `json:"path"`
	Status string `json:"status"`
}

func Discover(dirs []string) ([]Info, error) {
	if len(dirs) == 0 {
		dirs = DefaultDirs()
	}
	var result []Info
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("stating plugin directory %s: %w", dir, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("plugin path %s is not a directory", dir)
		}
		if info.Mode().Perm()&0002 != 0 {
			return nil, fmt.Errorf("plugin directory %s is world-writable", dir)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("reading plugin directory %s: %w", dir, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			path := filepath.Join(dir, entry.Name(), "manifest.json")
			data, err := os.ReadFile(path)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, fmt.Errorf("reading plugin manifest %s: %w", path, err)
			}
			var manifest Manifest
			if err := json.Unmarshal(data, &manifest); err != nil {
				return nil, fmt.Errorf("parsing plugin manifest %s: %w", path, err)
			}
			if manifest.Name == "" || manifest.Version == "" || manifest.APIVersion == "" {
				return nil, fmt.Errorf("plugin manifest %s requires name, version, and api_version", path)
			}
			status := "compatible"
			if manifest.APIVersion != APIVersion {
				status = "incompatible"
			}
			result = append(result, Info{Manifest: manifest, Path: filepath.Dir(path), Status: status})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

// DefaultDirs returns the documented discovery locations. Explicit command
// directories take precedence because Discover only uses this list when none
// were supplied.
func DefaultDirs() []string {
	var dirs []string
	if value := os.Getenv("SYSKIT_PLUGIN_DIR"); value != "" {
		dirs = append(dirs, filepath.SplitList(value)...)
	}
	if base := os.Getenv("XDG_DATA_HOME"); base != "" {
		dirs = append(dirs, filepath.Join(base, "syskit", "plugins"))
	} else if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, ".local", "share", "syskit", "plugins"))
	}
	return dirs
}
