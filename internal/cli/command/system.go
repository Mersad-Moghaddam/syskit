// Package command contains thin Cobra commands that map user intent onto
// services and render their typed results.
package command

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

// SystemService is the service boundary consumed by the command.
type SystemService interface {
	Collect() (*model.SystemInfo, error)
}

// SystemOptions supplies already-resolved CLI values without making this
// command package depend on the parent cli package.
type SystemOptions struct {
	Format   func() string
	NoHeader func() bool
	Color    func() bool
}

// NewSystemCmd builds `syskit system`.
func NewSystemCmd(service SystemService, options SystemOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "Show host, kernel, uptime, and load information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := service.Collect()
			if err != nil {
				return fmt.Errorf("collecting system information: %w", err)
			}
			r, err := render.New(options.Format(), render.WithNoHeader(options.NoHeader()), render.WithColor(options.Color()))
			if err != nil {
				return err
			}
			if options.Format() == "table" {
				return r.Render(cmd.OutOrStdout(), systemTable(info))
			}
			return r.Render(cmd.OutOrStdout(), info)
		},
	}
}

func systemTable(info *model.SystemInfo) render.Table {
	os := info.OSName
	if info.OSVersion != "" {
		os += " " + info.OSVersion
	}
	return render.Table{
		Headers: []string{"HOST", "OS", "KERNEL", "ARCH", "UPTIME", "LOAD"},
		Rows: [][]string{{info.Hostname, os, info.KernelRelease, info.Architecture,
			formatUptime(info.UptimeSeconds), fmt.Sprintf("%.2f %.2f %.2f", info.LoadAverage1, info.LoadAverage5, info.LoadAverage15)}},
	}
}

func formatUptime(seconds float64) string {
	d := time.Duration(seconds * float64(time.Second)).Round(time.Second)
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	if days > 0 {
		return fmt.Sprintf("%dd %02dh %02dm", days, d/time.Hour, (d%time.Hour)/time.Minute)
	}
	return fmt.Sprintf("%02dh %02dm", d/time.Hour, (d%time.Hour)/time.Minute)
}
