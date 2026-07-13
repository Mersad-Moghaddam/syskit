package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

type MemoryService interface {
	Collect() (*model.MemoryInfo, error)
}
type MemoryOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewMemoryCmd(s MemoryService, o MemoryOptions) *cobra.Command {
	return &cobra.Command{Use: "memory", Short: "Show memory, swap, cache, and pressure information", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.Collect()
		if err != nil {
			return fmt.Errorf("collecting memory information: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), memoryTable(info))
		}
		return r.Render(c.OutOrStdout(), info)
	}}
}
func memoryTable(i *model.MemoryInfo) render.Table {
	used, available, pressure := "unavailable", "unavailable", "unavailable"
	if i.UsedBytes != nil {
		used = fmt.Sprint(*i.UsedBytes)
	}
	if i.AvailableBytes != nil {
		available = fmt.Sprint(*i.AvailableBytes)
	}
	if i.Pressure != nil {
		pressure = fmt.Sprintf("some %.2f%%", i.Pressure.SomeAvg10)
	}
	return render.Table{Headers: []string{"TOTAL", "USED", "AVAILABLE", "FREE", "BUFFERS", "CACHE", "SWAP USED", "PRESSURE"}, Rows: [][]string{{fmt.Sprint(i.TotalBytes), used, available, fmt.Sprint(i.FreeBytes), fmt.Sprint(i.BuffersBytes), fmt.Sprint(i.CacheBytes), fmt.Sprint(i.SwapUsedBytes), pressure}}}
}
