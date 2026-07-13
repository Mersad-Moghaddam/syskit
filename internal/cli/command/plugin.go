package command

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/plugin"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

type PluginService interface {
	List([]string) ([]plugin.Info, error)
	Inspect([]string, string) (*plugin.Info, error)
	Run(context.Context, []string, string) (any, error)
}

func NewPluginCmd(s PluginService, format func() string, noHeader, color func() bool) *cobra.Command {
	var dirs []string
	var timeout time.Duration
	cmd := &cobra.Command{Use: "plugins", Short: "Inspect discovered plugins", Args: cobra.NoArgs}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List plugin manifests without executing them", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		items, err := s.List(dirs)
		if err != nil {
			return fmt.Errorf("discovering plugins: %w", err)
		}
		r, err := render.New(format(), render.WithNoHeader(noHeader()), render.WithColor(color()))
		if err != nil {
			return err
		}
		if format() != "table" {
			return r.Render(c.OutOrStdout(), items)
		}
		t := render.Table{Headers: []string{"NAME", "VERSION", "API", "STATUS", "PERMISSIONS", "PATH"}}
		for _, item := range items {
			t.Rows = append(t.Rows, []string{item.Name, item.Version, item.APIVersion, item.Status, strings.Join(item.Permissions, ","), item.Path})
		}
		return r.Render(c.OutOrStdout(), t)
	}})
	cmd.AddCommand(&cobra.Command{Use: "inspect <name>", Short: "Show one plugin manifest without executing it", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		item, err := s.Inspect(dirs, args[0])
		if err != nil {
			return fmt.Errorf("inspecting plugin: %w", err)
		}
		r, err := render.New(format(), render.WithNoHeader(noHeader()), render.WithColor(color()))
		if err != nil {
			return err
		}
		if format() != "table" {
			return r.Render(c.OutOrStdout(), item)
		}
		t := render.Table{Headers: []string{"NAME", "VERSION", "API", "STATUS", "PERMISSIONS", "PATH"}, Rows: [][]string{{item.Name, item.Version, item.APIVersion, item.Status, strings.Join(item.Permissions, ","), item.Path}}}
		return r.Render(c.OutOrStdout(), t)
	}})
	cmd.AddCommand(&cobra.Command{Use: "run <name>", Short: "Execute one compatible plugin explicitly", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if timeout < 100*time.Millisecond || timeout > time.Minute {
			return fmt.Errorf("plugin timeout must be between 100ms and 1m")
		}
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()
		value, err := s.Run(ctx, dirs, args[0])
		if err != nil {
			return fmt.Errorf("running plugin: %w", err)
		}
		r, err := render.New(format(), render.WithNoHeader(noHeader()), render.WithColor(color()))
		if err != nil {
			return err
		}
		if format() == "table" {
			return r.Render(c.OutOrStdout(), pluginResultTable(value))
		}
		return r.Render(c.OutOrStdout(), value)
	}})
	cmd.PersistentFlags().StringSliceVar(&dirs, "plugin-dir", nil, "plugin directory to inspect (repeatable)")
	cmd.PersistentFlags().DurationVar(&timeout, "timeout", 5*time.Second, "plugin execution timeout (100ms to 1m)")
	return cmd
}

func pluginResultTable(value any) render.Table {
	table := render.Table{Headers: []string{"FIELD", "VALUE"}}
	object, ok := value.(map[string]any)
	if !ok {
		data, _ := json.Marshal(value)
		return render.Table{Headers: []string{"VALUE"}, Rows: [][]string{{string(data)}}}
	}
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		data, _ := json.Marshal(object[key])
		table.Rows = append(table.Rows, []string{key, string(data)})
	}
	return table
}
