package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

type CPUService interface {
	Collect() (*model.CPUInfo, error)
	Sample(time.Duration) (*model.CPUInfo, error)
}
type CPUOptions struct {
	Format   func() string
	NoHeader func() bool
}

// NewCPUCmd builds the static `syskit cpu` command. Per-core utilization is
// intentionally deferred to CPU-02 because it requires two timed snapshots.
func NewCPUCmd(service CPUService, options CPUOptions) *cobra.Command {
	var perCore bool
	var interval time.Duration
	cmd := &cobra.Command{Use: "cpu", Short: "Show CPU topology, model, and frequency information", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := service.Sample(interval)
			if err != nil {
				return fmt.Errorf("collecting CPU information: %w", err)
			}
			r, err := render.New(options.Format(), render.WithNoHeader(options.NoHeader()))
			if err != nil {
				return err
			}
			if options.Format() == "table" {
				return r.Render(cmd.OutOrStdout(), cpuTable(info, perCore))
			}
			return r.Render(cmd.OutOrStdout(), info)
		},
	}
	cmd.Flags().BoolVar(&perCore, "per-core", false, "show one row per logical CPU")
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "interval between CPU samples")
	return cmd
}

func cpuTable(info *model.CPUInfo, perCore bool) render.Table {
	physical, sockets := "unavailable", "unavailable"
	if info.PhysicalCores != nil {
		physical = fmt.Sprint(*info.PhysicalCores)
	}
	if info.Sockets != nil {
		sockets = fmt.Sprint(*info.Sockets)
	}
	_ = strings.Join(info.Flags, " ")
	table := render.Table{Headers: []string{"CPU", "MODEL", "PHYSICAL", "LOGICAL", "SOCKETS", "UTIL", "USER", "SYSTEM", "IDLE"}}
	for _, t := range info.Times {
		if !perCore && t.CPUID != "all" {
			continue
		}
		util := "unavailable"
		if t.Utilization != nil {
			util = fmt.Sprintf("%.1f%%", *t.Utilization)
		}
		table.Rows = append(table.Rows, []string{t.CPUID, info.Model, physical, fmt.Sprint(info.LogicalCores), sockets, util, fmt.Sprint(t.User), fmt.Sprint(t.System), fmt.Sprint(t.Idle)})
	}
	if len(table.Rows) == 0 {
		table.Rows = append(table.Rows, []string{"all", info.Model, physical, fmt.Sprint(info.LogicalCores), sockets, "unavailable", "", "", ""})
	}
	return table
}
