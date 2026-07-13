package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

type FilesystemService interface {
	Collect() (*model.DiskInfo, error)
}
type FilesystemOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewFilesystemCmd(s FilesystemService, o FilesystemOptions) *cobra.Command {
	var filesystemType string
	var showPseudo bool
	cmd := &cobra.Command{Use: "filesystem", Short: "Show mounted filesystem inode and option details", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.Collect()
		if err != nil {
			return fmt.Errorf("collecting filesystem information: %w", err)
		}
		filtered := &model.DiskInfo{}
		for _, m := range info.Mounts {
			if filesystemType != "" && m.FilesystemType != filesystemType {
				continue
			}
			if !showPseudo && isPseudo(m.FilesystemType) {
				continue
			}
			filtered.Mounts = append(filtered.Mounts, m)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), filesystemTable(filtered))
		}
		return r.Render(c.OutOrStdout(), filtered)
	}}
	cmd.Flags().StringVar(&filesystemType, "type", "", "filter by filesystem type")
	cmd.Flags().BoolVar(&showPseudo, "show-pseudo", false, "include pseudo filesystems")
	return cmd
}
func isPseudo(t string) bool {
	switch t {
	case "proc", "sysfs", "tmpfs", "devtmpfs", "devpts", "cgroup", "cgroup2", "overlay", "mqueue", "securityfs", "tracefs", "debugfs", "pstore", "efivarfs", "bpf", "configfs", "fusectl", "autofs", "ramfs":
		return true
	}
	return false
}
func filesystemTable(info *model.DiskInfo) render.Table {
	t := render.Table{Headers: []string{"MOUNT", "TYPE", "SOURCE", "INODES", "FREE", "IUSE%", "OPTIONS"}}
	for _, m := range info.Mounts {
		total, free, pct := "unavailable", "unavailable", "unavailable"
		if m.TotalInodes != nil {
			total = fmt.Sprint(*m.TotalInodes)
		}
		if m.FreeInodes != nil {
			free = fmt.Sprint(*m.FreeInodes)
			if m.TotalInodes != nil && *m.TotalInodes > 0 {
				pct = fmt.Sprintf("%.0f%%", float64(*m.TotalInodes-*m.FreeInodes)*100/float64(*m.TotalInodes))
			}
		}
		t.Rows = append(t.Rows, []string{m.MountPoint, m.FilesystemType, m.Source, total, free, pct, strings.Join(m.Options, ",")})
	}
	return t
}
