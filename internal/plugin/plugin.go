// Package plugin discovers declarative, out-of-process SysKit plugins. It
// never executes plugin binaries; execution remains an explicit later action.
package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const APIVersion = "v1"

type Manifest struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	APIVersion  string            `json:"api_version"`
	Permissions []string          `json:"permissions,omitempty"`
	Executable  string            `json:"executable,omitempty"`
	Collectors  []string          `json:"collectors,omitempty"`
	Schemas     map[string]string `json:"output_schemas,omitempty"`
	Author      string            `json:"author,omitempty"`
	License     string            `json:"license,omitempty"`
}

type Request struct {
	APIVersion string `json:"api_version"`
	Action     string `json:"action"`
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
			if manifest.Executable == "" || len(manifest.Collectors) == 0 || manifest.Permissions == nil || len(manifest.Schemas) == 0 || manifest.Author == "" || manifest.License == "" {
				return nil, fmt.Errorf("plugin manifest %s requires executable, collectors, permissions, output_schemas, author, and license", path)
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

func Inspect(dirs []string, name string) (*Info, error) {
	items, err := Discover(dirs)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].Name == name {
			return &items[i], nil
		}
	}
	return nil, fmt.Errorf("plugin %q not found", name)
}

// Run executes one explicitly selected compatible plugin using the v1 JSON
// stdin/stdout protocol. Discovery alone never calls this function.
func Run(ctx context.Context, dirs []string, name string) (any, error) {
	info, err := Inspect(dirs, name)
	if err != nil {
		return nil, err
	}
	if info.Status != "compatible" {
		return nil, fmt.Errorf("plugin %q uses incompatible API %q", name, info.APIVersion)
	}
	if info.Executable == "" || filepath.IsAbs(info.Executable) || filepath.Clean(info.Executable) != info.Executable || info.Executable == "." || info.Executable == ".." || strings.HasPrefix(info.Executable, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("plugin %q has invalid executable path", name)
	}
	path := filepath.Join(info.Path, info.Executable)
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stating plugin executable %s: %w", path, err)
	}
	if !stat.Mode().IsRegular() || stat.Mode().Perm()&0111 == 0 {
		return nil, fmt.Errorf("plugin executable %s is not an executable regular file", path)
	}
	request, _ := json.Marshal(Request{APIVersion: APIVersion, Action: "collect"})
	command := exec.CommandContext(ctx, path)
	command.Stdin = bytes.NewReader(append(request, '\n'))
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("plugin %q timed out: %w", name, ctx.Err())
		}
		return nil, fmt.Errorf("plugin %q failed: %w: %s", name, err, strings.TrimSpace(stderr.String()))
	}
	var result any
	decoder := json.NewDecoder(&stdout)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("plugin %q returned invalid JSON: %w", name, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("plugin %q returned multiple JSON values", name)
		}
		return nil, fmt.Errorf("plugin %q returned trailing invalid JSON: %w", name, err)
	}
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
