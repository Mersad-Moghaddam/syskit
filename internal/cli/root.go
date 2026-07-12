package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/cli/command"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/cpu"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/disk"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/memory"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/network"
	processcollector "github.com/Mersad-Moghaddam/syskit/internal/collector/process"
	systemcollector "github.com/Mersad-Moghaddam/syskit/internal/collector/system"
	"github.com/Mersad-Moghaddam/syskit/internal/platform"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
)

// Supported values for the persistent --format flag. FND-06 wires these to the
// concrete renderers; the CLI parses and validates the selection here.
const (
	formatTable = "table"
	formatJSON  = "json"
	formatYAML  = "yaml"
)

// globalOptions holds the parsed global flags and the values derived from them
// (loaded config, constructed logger) shared by every subcommand. It is created
// per invocation in newRootCmd — there is no mutable package-level state — and
// is the seam commands read: the effective format, the diagnostic logger, and
// the resolved configuration.
type globalOptions struct {
	// Raw flag values.
	format     string
	configPath string
	verbose    bool
	debug      bool
	quiet      bool

	// Derived at PersistentPreRunE time.
	cfg    *Config
	logger *slog.Logger
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
		SilenceErrors: true, // Main presents errors at the boundary (present()).
		SilenceUsage:  true, // Main prints a concise usage hint; no full dump.
		// PersistentPreRunE loads configuration, resolves the effective global
		// options, builds the logger, and validates flags before any subcommand
		// runs. Configuration and logging are CLI-layer concerns resolved once
		// here and threaded down as plain values.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.resolve(cmd)
		},
		// RunE handles the bare `syskit` invocation by printing help. Defining it
		// also makes Cobra run PersistentPreRunE for the root command itself, so
		// an invalid global flag (e.g. --format xml) with no subcommand is
		// reported as a usage error instead of silently printing help.
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	pf := cmd.PersistentFlags()
	pf.StringVar(&opts.format, "format", formatTable, "output format: table, json, or yaml")
	pf.StringVar(&opts.configPath, "config", "", "path to a specific config file")
	pf.BoolVarP(&opts.verbose, "verbose", "v", false, "show info-level diagnostics on stderr")
	pf.BoolVar(&opts.debug, "debug", false, "show debug-level diagnostics on stderr")
	pf.BoolVarP(&opts.quiet, "quiet", "q", false, "suppress all diagnostics, including errors")

	// Report Cobra's own flag-parsing failures (unknown flag, missing value) as
	// usage errors so the boundary maps them to exit code 2.
	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return &usageError{err: err}
	})

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(command.NewSystemCmd(
		service.NewSystem(systemcollector.NewCollector(platform.RealFS())),
		command.SystemOptions{
			Format:   func() string { return opts.format },
			NoHeader: func() bool { return opts.cfg != nil && opts.cfg.NoHeader },
		},
	))
	cmd.AddCommand(command.NewCPUCmd(
		service.NewCPU(cpu.NewCollector(platform.RealFS())),
		command.CPUOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.cfg != nil && opts.cfg.NoHeader }},
	))
	cmd.AddCommand(command.NewMemoryCmd(service.NewMemory(memory.NewCollector(platform.RealFS())), command.MemoryOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.cfg != nil && opts.cfg.NoHeader }}))
	cmd.AddCommand(command.NewDiskCmd(service.NewDisk(disk.NewCollector(platform.RealFS())), command.DiskOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.cfg != nil && opts.cfg.NoHeader }}))
	cmd.AddCommand(command.NewFilesystemCmd(service.NewDisk(disk.NewCollector(platform.RealFS())), command.FilesystemOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.cfg != nil && opts.cfg.NoHeader }}))
	cmd.AddCommand(command.NewProcessCmd(service.NewProcess(processcollector.NewCollector(platform.RealFS())), command.ProcessOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.cfg != nil && opts.cfg.NoHeader }}))
	cmd.AddCommand(command.NewNetworkCmd(service.NewNetwork(network.NewCollector(platform.RealFS())), command.NetworkOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.cfg != nil && opts.cfg.NoHeader }}))

	return cmd
}

// resolve loads configuration, applies flag>env>file>default precedence to the
// global options, builds the diagnostic logger, and validates the effective
// format. It runs in PersistentPreRunE so every subcommand sees resolved values.
func (o *globalOptions) resolve(cmd *cobra.Command) error {
	// Explicit --config flag wins over SYSKIT_CONFIG for discovery.
	path := o.configPath
	if path == "" {
		path = os.Getenv("SYSKIT_CONFIG")
	}

	cfg, err := Load(path)
	if err != nil {
		// A malformed config file is a real, user-actionable error.
		return err
	}
	o.cfg = cfg

	// Effective format: flag > env > per-command section > global > default.
	command := commandName(cmd)
	_, envFormat := os.LookupEnv("SYSKIT_FORMAT")
	flagChanged := cmd.Flags().Changed("format")
	o.format = cfg.resolveFormat(flagChanged, o.format, envFormat, command)
	if err := validateFormat(o.format); err != nil {
		return err
	}

	// Logger: verbosity from flags (quiet > debug > verbose), always on stderr.
	v := resolveVerbosity(o.verbose, o.debug, o.quiet)
	o.logger = newLogger(cmd.ErrOrStderr(), v)
	o.logger.Debug("configuration resolved",
		"format", o.format,
		"config_path", path,
		"command", command,
	)
	return nil
}

// commandName returns the invoked subcommand's name for per-command config
// lookup. For the bare root it returns "".
func commandName(cmd *cobra.Command) string {
	if cmd.Name() == "syskit" {
		return ""
	}
	return cmd.Name()
}

// validateFormat returns a usage error when the effective --format value is not
// one of the supported formats. FND-06 selects a renderer from the same set.
func validateFormat(format string) error {
	switch format {
	case formatTable, formatJSON, formatYAML:
		return nil
	default:
		return &usageError{err: fmt.Errorf("invalid --format %q: must be one of table, json, yaml", format)}
	}
}

// Main builds and executes the root command, presents any error at the CLI
// boundary (user-facing message to stderr, exit code from the sentinel), and
// returns the process exit code. cmd/syskit passes the result straight to
// os.Exit, keeping main tiny; all CLI logic lives here.
func Main() int {
	root := newRootCmd()
	err := root.Execute()

	message, code := present(err)
	if message != "" {
		fmt.Fprintln(root.ErrOrStderr(), message)
	}
	return code
}
