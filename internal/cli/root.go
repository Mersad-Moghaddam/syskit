package cli

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/cli/command"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/cpu"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/disk"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/memory"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/network"
	"github.com/Mersad-Moghaddam/syskit/internal/collector/port"
	processcollector "github.com/Mersad-Moghaddam/syskit/internal/collector/process"
	systemcollector "github.com/Mersad-Moghaddam/syskit/internal/collector/system"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
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
	color      string
	noHeader   bool
	configPath string
	verbose    bool
	debug      bool
	quiet      bool

	// Derived at PersistentPreRunE time.
	cfg    *Config
	logger *slog.Logger
	level  verbosity
}

// newRootCmd builds the syskit root command. It returns a fresh command on every
// call so tests can exercise flag parsing without shared state.
func newRootCmd() *cobra.Command {
	cmd, _ := newRootCmdWithOptions()
	return cmd
}

func newRootCmdWithOptions() (*cobra.Command, *globalOptions) {
	return newRootCmdWithOptionsAndTheme(nil)
}

func newRootCmdWithOptionsAndTheme(theme *tuiTheme) (*cobra.Command, *globalOptions) {
	opts := &globalOptions{}
	activeTheme := resolveTUITheme(theme)
	theme = &activeTheme

	cmd := &cobra.Command{
		Use:   "syskit",
		Short: "Inspect Linux system state from native kernel interfaces",
		Long: `SysKit is a Linux-only, read-only system-inspection toolkit.

It reads native kernel interfaces (/proc, /sys, Netlink, cgroups) directly,
never shelling out to other utilities, and renders CPU, memory, disk,
filesystem, process, network, and port information as a table, JSON, or YAML.

Run syskit without a subcommand in an interactive terminal to open the
hierarchical control center with an animated, color-coded interface. Use arrow
keys or the mouse to browse every command family; one-shot results and live
views retain the selected accent and return to the same menu location.`,
		SilenceErrors: true, // Main presents errors at the boundary (present()).
		SilenceUsage:  true, // Main prints a concise usage hint; no full dump.
		// PersistentPreRunE loads configuration, resolves the effective global
		// options, builds the logger, and validates flags before any subcommand
		// runs. Configuration and logging are CLI-layer concerns resolved once
		// here and threaded down as plain values.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.resolve(cmd); err != nil {
				return err
			}
			theme.color = opts.colorEnabled(cmd.OutOrStdout())
			return nil
		},
		// RunE opens the discoverability menu for a bare interactive invocation
		// and preserves ordinary help for pipes and redirected output. Defining
		// it also makes Cobra run PersistentPreRunE for the root command itself,
		// so invalid global flags are still reported as usage errors.
		RunE: func(cmd *cobra.Command, args []string) error {
			stdin, inputIsFile := cmd.InOrStdin().(*os.File)
			stdout, outputIsFile := cmd.OutOrStdout().(*os.File)
			if inputIsFile && outputIsFile && isInteractiveTerminal(stdin) && isInteractiveTerminal(stdout) {
				stderr, ok := cmd.ErrOrStderr().(*os.File)
				if !ok {
					stderr = os.Stderr
				}
				return runInteractiveMenu(stdin, stdout, stderr, opts.colorEnabled(stdout), changedPersistentArgs(cmd))
			}
			return cmd.Help()
		},
	}

	pf := cmd.PersistentFlags()
	pf.StringVar(&opts.format, "format", formatTable, "output format: table, json, or yaml")
	pf.StringVar(&opts.color, "color", "auto", "color output: auto, always, or never")
	pf.BoolVar(&opts.noHeader, "no-header", false, "suppress table headers")
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
			NoHeader: func() bool { return opts.noHeader },
			Color:    func() bool { return opts.colorEnabled(cmd.OutOrStdout()) },
		},
	))
	cmd.AddCommand(command.NewCPUCmd(
		service.NewCPU(cpu.NewCollector(platform.RealFS())),
		command.CPUOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }},
	))
	cmd.AddCommand(command.NewMemoryCmd(service.NewMemory(memory.NewCollector(platform.RealFS())), command.MemoryOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }}))
	cmd.AddCommand(command.NewDiskCmd(service.NewDisk(disk.NewCollector(platform.RealFS())), command.DiskOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }}))
	cmd.AddCommand(command.NewFilesystemCmd(service.NewDisk(disk.NewCollector(platform.RealFS())), command.FilesystemOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }}))
	cmd.AddCommand(command.NewProcessCmd(service.NewProcess(processcollector.NewCollector(platform.RealFS())), command.ProcessOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }}))
	containerFS := platform.RealFS()
	cmd.AddCommand(command.NewContainerCmd(service.NewContainer(processcollector.NewCollector(containerFS), cgroupMetricsReader(containerFS)), command.ContainerOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }}))
	cmd.AddCommand(command.NewNetworkCmd(service.NewNetwork(network.NewCollectorWithAddresses(platform.RealFS(), platform.RealNetlink())), command.NetworkOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }}))
	cmd.AddCommand(command.NewPortCmd(service.NewPort(port.NewCollector(platform.RealFS())), command.PortOptions{Format: func() string { return opts.format }, NoHeader: func() bool { return opts.noHeader }, Color: func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }}))
	cmd.AddCommand(command.NewDiagnosticCmd(service.NewDiagnostic(systemcollector.NewCollector(platform.RealFS()), cpu.NewCollector(platform.RealFS()), memory.NewCollector(platform.RealFS()), disk.NewCollector(platform.RealFS()), processcollector.NewCollector(platform.RealFS()), network.NewCollectorWithAddresses(platform.RealFS(), platform.RealNetlink()), port.NewCollector(platform.RealFS())), func() string { return opts.format }, func() bool { return opts.noHeader }, func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }))
	cmd.AddCommand(command.NewPluginCmd(service.NewPlugin(), func() string { return opts.format }, func() bool { return opts.noHeader }, func() bool { return opts.colorEnabled(cmd.OutOrStdout()) }))
	var previousDashboardCPU *model.CPUInfo
	var previousDashboardNetwork *model.NetworkInfo
	cmd.AddCommand(newDashboardCmdWithTheme(func() (dashboardSnapshot, error) {
		system, err := service.NewSystem(systemcollector.NewCollector(platform.RealFS())).Collect()
		if err != nil {
			return dashboardSnapshot{}, err
		}
		memory, err := service.NewMemory(memory.NewCollector(platform.RealFS())).Collect()
		if err != nil {
			return dashboardSnapshot{}, err
		}
		cpu, err := service.NewCPU(cpu.NewCollector(platform.RealFS())).Collect()
		if err != nil {
			return dashboardSnapshot{}, err
		}
		network, err := service.NewNetwork(network.NewCollectorWithAddresses(platform.RealFS(), platform.RealNetlink())).Collect()
		if err != nil {
			return dashboardSnapshot{}, err
		}
		disks, err := service.NewDisk(disk.NewCollector(platform.RealFS())).Collect()
		if err != nil {
			return dashboardSnapshot{}, err
		}
		processes, err := service.NewProcess(processcollector.NewCollector(platform.RealFS())).List(service.ProcessOptions{Sort: "memory", Reverse: true, Limit: 1})
		if err != nil {
			return dashboardSnapshot{}, err
		}
		used := uint64(0)
		if memory.UsedBytes != nil {
			used = *memory.UsedBytes
		}
		diskUsed, diskTotal := uint64(0), uint64(0)
		for _, mount := range disks.Mounts {
			if mount.MountPoint == "/" {
				if mount.UsedBytes != nil {
					diskUsed = *mount.UsedBytes
				}
				if mount.TotalBytes != nil {
					diskTotal = *mount.TotalBytes
				}
				break
			}
		}
		top := "unavailable"
		if len(processes.Processes) > 0 {
			top = processes.Processes[0].Command
		}
		cpuPercent := dashboardCPUUtilization(previousDashboardCPU, cpu)
		rxRate, txRate := dashboardNetworkRates(previousDashboardNetwork, network)
		previousDashboardCPU, previousDashboardNetwork = cpu, network
		return dashboardSnapshot{Hostname: system.Hostname, Uptime: system.UptimeSeconds, CPUPercent: cpuPercent, MemoryUsed: used, MemoryTotal: memory.TotalBytes, SwapUsed: memory.SwapUsedBytes, SwapTotal: memory.SwapTotalBytes, DiskUsed: diskUsed, DiskTotal: diskTotal, Interfaces: len(network.Interfaces), NetworkRX: rxRate, NetworkTX: txRate, TopProcess: top}, nil
	}, theme))
	cmd.AddCommand(newWatchCmdWithTheme(func(args []string, out io.Writer) error {
		// A fresh root preserves the normal command/service construction path for
		// every refresh without sharing mutable Cobra invocation state.
		child := newRootCmd()
		child.SetArgs(append(args, "--format", "table"))
		child.SetOut(out)
		child.SetErr(out)
		return child.Execute()
	}, theme))
	cmd.AddCommand(newTopCmdWithTheme(func(options service.ProcessOptions) (*model.ProcessList, error) {
		return service.NewProcess(processcollector.NewCollector(platform.RealFS())).List(options)
	}, theme))

	return cmd, opts
}

func cgroupMetricsReader(fs platform.SysFS) service.ContainerMetricsReader {
	return func(process model.Process) (*model.ContainerMetrics, error) {
		info, err := platform.DetectCgroup(fs, fmt.Sprintf("proc/%d/cgroup", process.PID))
		if err != nil {
			return nil, err
		}
		metrics, err := platform.ReadCgroupMetrics(fs, info)
		if err != nil {
			return nil, err
		}
		return &model.ContainerMetrics{MemoryCurrentBytes: metrics.MemoryCurrentBytes, CPUUsageNanoseconds: metrics.CPUUsageNanoseconds, ReadBytes: metrics.ReadBytes, WrittenBytes: metrics.WrittenBytes}, nil
	}
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
	_, envColor := os.LookupEnv("SYSKIT_COLOR")
	o.color = cfg.resolveColor(cmd.Flags().Changed("color"), o.color, envColor, command)
	if os.Getenv("NO_COLOR") != "" {
		o.color = "never"
	}
	if err := validateColor(o.color); err != nil {
		return err
	}
	_, envNoHeader := os.LookupEnv("SYSKIT_NO_HEADER")
	o.noHeader = cfg.resolveNoHeader(cmd.Flags().Changed("no-header"), o.noHeader, envNoHeader, command)

	if (command == "watch" || command == "top" || command == "dashboard") && !cmd.Flags().Changed("interval") {
		_, envRefresh := os.LookupEnv("SYSKIT_REFRESH_INTERVAL")
		if err := cmd.Flags().Set("interval", cfg.resolveRefreshInterval(envRefresh, command).String()); err != nil {
			return &usageError{err: err}
		}
	}

	// Logger: verbosity from flags (quiet > debug > verbose), always on stderr.
	v := resolveVerbosity(o.verbose, o.debug, o.quiet)
	if !cmd.Flags().Changed("quiet") && !cmd.Flags().Changed("debug") && !cmd.Flags().Changed("verbose") {
		_, envVerbosity := os.LookupEnv("SYSKIT_VERBOSITY")
		v, err = parseVerbosity(cfg.resolveConfiguredVerbosity(envVerbosity, command))
		if err != nil {
			return &usageError{err: err}
		}
	}
	o.logger = newLogger(cmd.ErrOrStderr(), v)
	o.level = v
	o.logger.Debug("configuration resolved",
		"format", o.format,
		"config_path", path,
		"command", command,
	)
	return nil
}

func validateColor(color string) error {
	switch color {
	case "auto", "always", "never":
		return nil
	default:
		return &usageError{err: fmt.Errorf("invalid --color %q: must be one of auto, always, never", color)}
	}
}

func (o *globalOptions) colorEnabled(w io.Writer) bool {
	if o.color == "always" {
		return true
	}
	if o.color != "auto" {
		return false
	}
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// commandName returns the invoked subcommand's name for per-command config
// lookup. For the bare root it returns "".
func commandName(cmd *cobra.Command) string {
	if cmd.Name() == "syskit" {
		return ""
	}
	for cmd.Parent() != nil && cmd.Parent().Name() != "syskit" {
		cmd = cmd.Parent()
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
	root, opts := newRootCmdWithOptions()
	err := root.Execute()

	message, code := present(err)
	if message != "" && !opts.quiet && opts.level != verbosityQuiet {
		fmt.Fprintln(root.ErrOrStderr(), message)
	}
	return code
}
