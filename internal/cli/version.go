package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is the build version string. It defaults to "dev" and is overridden at
// link time, e.g.:
//
//	go build -ldflags "-X github.com/Mersad-Moghaddam/syskit/internal/cli.version=v0.1.0" ./cmd/syskit
var version = "dev"

// newVersionCmd returns the `syskit version` subcommand, which prints the build
// version to stdout and exits 0.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the syskit version",
		Args:  usageArgs(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), version)
			return nil
		},
	}
}
