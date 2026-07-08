package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Supported values for the persistent --format flag. FND-06 wires these to the
// concrete renderers; FND-03 only parses and validates the selection.
const (
	formatTable = "table"
	formatJSON  = "json"
	formatYAML  = "yaml"
)

// globalOptions holds the parsed global flags shared by every subcommand. It is
// created per invocation in newRootCmd — there is no mutable package-level state
// — and is the seam later stories extend and read: FND-06 selects a renderer
// from format, FND-08 folds in --config, FND-09 folds in logging flags.
type globalOptions struct {
	format string
}

// newRootCmd builds the syskit root command. It returns a fresh command on every
// call so tests can exercise flag parsing without shared state.
func newRootCmd() *cobra.Command {
	opts := &globalOptions{}

	cmd := &cobra.Command{
		Use:   "syskit",
		Short: "Inspect Linux system state from native kernel interfaces",
		Long: `SysKit is a Linux-only, read-only system-inspection toolkit.

It reads native kernel interfaces (/proc, /sys, Netlink, cgroups) directly,
never shelling out to other utilities, and renders CPU, memory, disk,
filesystem, process, network, and port information as a table, JSON, or YAML.`,
		// PersistentPreRunE validates the global flags before any subcommand runs.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return validateFormat(opts.format)
		},
		// RunE handles the bare `syskit` invocation by printing help. Defining it
		// also makes Cobra run PersistentPreRunE for the root command itself, so
		// an invalid global flag (e.g. --format xml) with no subcommand is
		// reported as a usage error instead of silently printing help.
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&opts.format, "format", formatTable,
		"output format: table, json, or yaml")

	// Report Cobra's own flag-parsing failures (unknown flag, missing value) as
	// usage errors so the boundary maps them to exit code 2.
	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return &usageError{err: err}
	})

	cmd.AddCommand(newVersionCmd())

	return cmd
}

// validateFormat returns a usage error when the --format value is not one of the
// supported formats. FND-06 selects a renderer from the same set of values.
func validateFormat(format string) error {
	switch format {
	case formatTable, formatJSON, formatYAML:
		return nil
	default:
		return &usageError{err: fmt.Errorf("invalid --format %q: must be one of table, json, yaml", format)}
	}
}

// Main builds and executes the root command and returns the process exit code.
// cmd/syskit passes the result straight to os.Exit, keeping main tiny; all CLI
// logic lives here.
func Main() int {
	root := newRootCmd()
	return exitCode(root.Execute())
}
