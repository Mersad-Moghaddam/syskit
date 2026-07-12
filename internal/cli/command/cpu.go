package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

type CPUService interface {
	Collect() (*model.CPUInfo, error)
}
type CPUOptions struct {
	Format   func() string
	NoHeader func() bool
}

// NewCPUCmd builds the static `syskit cpu` command. Per-core utilization is
// intentionally deferred to CPU-02 because it requires two timed snapshots.
func NewCPUCmd(service CPUService, options CPUOptions) *cobra.Command {
	return &cobra.Command{Use: "cpu", Short: "Show CPU topology, model, and frequency information", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := service.Collect()
			if err != nil {
				return fmt.Errorf("collecting CPU information: %w", err)
			}
			r, err := render.New(options.Format(), render.WithNoHeader(options.NoHeader()))
			if err != nil {
				return err
			}
			if options.Format() == "table" {
				return r.Render(cmd.OutOrStdout(), cpuTable(info))
			}
			return r.Render(cmd.OutOrStdout(), info)
		},
	}
}

func cpuTable(info *model.CPUInfo) render.Table {
	physical, sockets := "unavailable", "unavailable"
	if info.PhysicalCores != nil {
		physical = fmt.Sprint(*info.PhysicalCores)
	}
	if info.Sockets != nil {
		sockets = fmt.Sprint(*info.Sockets)
	}
	flags := strings.Join(info.Flags, " ")
	return render.Table{Headers: []string{"MODEL", "PHYSICAL", "LOGICAL", "SOCKETS", "ARCH", "FLAGS"}, Rows: [][]string{{info.Model, physical, fmt.Sprint(info.LogicalCores), sockets, info.Architecture, flags}}}
}
