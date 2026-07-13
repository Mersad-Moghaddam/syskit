package command

import (
	"fmt"
	"github.com/Mersad-Moghaddam/syskit/internal/plugin"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/spf13/cobra"
	"strings"
)

type PluginService interface {
	List([]string) ([]plugin.Info, error)
	Inspect([]string, string) (*plugin.Info, error)
}

func NewPluginCmd(s PluginService, format func() string, noHeader func() bool) *cobra.Command {
	var dirs []string
	cmd := &cobra.Command{Use: "plugins", Short: "Inspect discovered plugins", Args: cobra.NoArgs}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List plugin manifests without executing them", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		items, err := s.List(dirs)
		if err != nil {
			return fmt.Errorf("discovering plugins: %w", err)
		}
		r, err := render.New(format(), render.WithNoHeader(noHeader()))
		if err != nil {
			return err
		}
		if format() != "table" {
			return r.Render(c.OutOrStdout(), items)
		}
		t := render.Table{Headers: []string{"NAME", "VERSION", "API", "STATUS", "PATH"}}
		for _, item := range items {
			t.Rows = append(t.Rows, []string{item.Name, item.Version, item.APIVersion, item.Status, item.Path})
		}
		return r.Render(c.OutOrStdout(), t)
	}})
	cmd.AddCommand(&cobra.Command{Use: "inspect <name>", Short: "Show one plugin manifest without executing it", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		item, err := s.Inspect(dirs, args[0])
		if err != nil {
			return fmt.Errorf("inspecting plugin: %w", err)
		}
		r, err := render.New(format(), render.WithNoHeader(noHeader()))
		if err != nil {
			return err
		}
		if format() != "table" {
			return r.Render(c.OutOrStdout(), item)
		}
		t := render.Table{Headers: []string{"NAME", "VERSION", "API", "STATUS", "PERMISSIONS", "PATH"}, Rows: [][]string{{item.Name, item.Version, item.APIVersion, item.Status, strings.Join(item.Permissions, ","), item.Path}}}
		return r.Render(c.OutOrStdout(), t)
	}})
	cmd.PersistentFlags().StringSliceVar(&dirs, "plugin-dir", nil, "plugin directory to inspect (repeatable)")
	return cmd
}
