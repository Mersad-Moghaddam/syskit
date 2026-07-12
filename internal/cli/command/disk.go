package command

import (
	"fmt"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/spf13/cobra"
)

type DiskService interface {
	Collect() (*model.DiskInfo, error)
	Sample(time.Duration) (*model.DiskInfo, error)
}
type DiskOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewDiskCmd(s DiskService, o DiskOptions) *cobra.Command {
	var io bool
	var interval time.Duration
	var mount, filesystemType, device string
	cmd := &cobra.Command{Use: "disk", Short: "Show mounted filesystem capacity and I/O", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		var info *model.DiskInfo
		var err error
		if io {
			info, err = s.Sample(interval)
		} else {
			info, err = s.Collect()
		}
		if err != nil {
			return fmt.Errorf("collecting disk information: %w", err)
		}
		if !io {
			info = filterDiskMounts(info, mount, filesystemType, device)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			if io {
				return r.Render(c.OutOrStdout(), diskIOTable(info))
			}
			return r.Render(c.OutOrStdout(), diskTable(info))
		}
		return r.Render(c.OutOrStdout(), info)
	}}
	cmd.Flags().BoolVar(&io, "io", false, "show sampled disk I/O rates")
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "interval between disk I/O samples")
	cmd.Flags().StringVar(&mount, "mount", "", "filter by mount point")
	cmd.Flags().StringVar(&filesystemType, "type", "", "filter by filesystem type")
	cmd.Flags().StringVar(&device, "device", "", "filter by device source")
	return cmd
}
func filterDiskMounts(info *model.DiskInfo, mount, filesystemType, device string) *model.DiskInfo {
	if mount == "" && filesystemType == "" && device == "" {
		return info
	}
	result := &model.DiskInfo{Devices: info.Devices, CollectedAt: info.CollectedAt}
	for _, m := range info.Mounts {
		if mount != "" && m.MountPoint != mount {
			continue
		}
		if filesystemType != "" && m.FilesystemType != filesystemType {
			continue
		}
		if device != "" && m.Source != device {
			continue
		}
		result.Mounts = append(result.Mounts, m)
	}
	return result
}
func diskIOTable(info *model.DiskInfo) render.Table {
	t := render.Table{Headers: []string{"DEVICE", "READ OPS", "WRITE OPS", "READ B/s", "WRITE B/s"}}
	for _, d := range info.Devices {
		r, w := "unavailable", "unavailable"
		if d.ReadBytesPerSecond != nil {
			r = fmt.Sprintf("%.0f", *d.ReadBytesPerSecond)
		}
		if d.WrittenBytesPerSecond != nil {
			w = fmt.Sprintf("%.0f", *d.WrittenBytesPerSecond)
		}
		t.Rows = append(t.Rows, []string{d.Name, fmt.Sprint(d.ReadOperations), fmt.Sprint(d.WrittenOperations), r, w})
	}
	return t
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
