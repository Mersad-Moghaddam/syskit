package command

import (
	"fmt"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/spf13/cobra"
)

type DiskService interface {
	Collect() (*model.DiskInfo, error)
}
type DiskOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewDiskCmd(s DiskService, o DiskOptions) *cobra.Command {
	return &cobra.Command{Use: "disk", Short: "Show mounted filesystem capacity", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.Collect()
		if err != nil {
			return fmt.Errorf("collecting disk information: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), diskTable(info))
		}
		return r.Render(c.OutOrStdout(), info)
	}}
}
func diskTable(info *model.DiskInfo) render.Table {
	t := render.Table{Headers: []string{"DEVICE", "TYPE", "SIZE", "USED", "AVAIL", "USE%", "MOUNT"}}
	for _, m := range info.Mounts {
		size, used, avail, pct := "unavailable", "unavailable", "unavailable", "unavailable"
		if m.TotalBytes != nil {
			size = fmt.Sprint(*m.TotalBytes)
		}
		if m.UsedBytes != nil {
			used = fmt.Sprint(*m.UsedBytes)
		}
		if m.AvailableBytes != nil {
			avail = fmt.Sprint(*m.AvailableBytes)
		}
		if m.UsePercent != nil {
			pct = fmt.Sprintf("%.0f%%", *m.UsePercent)
		}
		t.Rows = append(t.Rows, []string{m.Source, m.FilesystemType, size, used, avail, pct, m.MountPoint})
	}
	return t
}
