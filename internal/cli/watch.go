package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	minWatchInterval = 250 * time.Millisecond
	maxWatchInterval = time.Minute
)

type watchRunner func([]string, io.Writer) error

func newWatchCmd(run watchRunner) *cobra.Command {
	return newWatchCmdWithTheme(run, nil)
}

func newWatchCmdWithTheme(run watchRunner, selectedTheme *tuiTheme) *cobra.Command {
	var interval time.Duration
	cmd := &cobra.Command{Use: "watch <command>", Short: "Refresh a command continuously", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if interval < minWatchInterval || interval > maxWatchInterval {
			return fmt.Errorf("watch interval must be between %s and %s", minWatchInterval, maxWatchInterval)
		}
		if !isInteractiveTerminal(os.Stdout) {
			return fmt.Errorf("watch requires an interactive terminal")
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		return runWatchThemed(ctx, interval, args, cmd.OutOrStdout(), run, resolveTUITheme(selectedTheme))
	}}
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "refresh interval (250ms to 1m)")
	return cmd
}

func runWatch(ctx context.Context, interval time.Duration, args []string, out io.Writer, run watchRunner) error {
	return runWatchWithIteration(ctx, interval, args, out, func(args []string, out io.Writer, run watchRunner) error {
		return runWatchIteration(args, out, run)
	}, run)
}

func runWatchThemed(ctx context.Context, interval time.Duration, args []string, out io.Writer, run watchRunner, theme tuiTheme) error {
	return runWatchWithIteration(ctx, interval, args, out, func(args []string, out io.Writer, run watchRunner) error {
		return runWatchIterationThemed(args, out, run, theme, interval)
	}, run)
}

func runWatchWithIteration(ctx context.Context, interval time.Duration, args []string, out io.Writer, iteration func([]string, io.Writer, watchRunner) error, run watchRunner) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := iteration(args, out, run); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func runWatchIterationThemed(args []string, out io.Writer, run watchRunner, theme tuiTheme, interval time.Duration) error {
	var buffer bytes.Buffer
	if err := run(append([]string(nil), args...), &buffer); err != nil {
		return err
	}
	header := theme.badge("◉  SYSKIT WATCH") + "  " + theme.primaryStyle().Bold(theme.color).Render(strings.Join(args, " "))
	status := theme.primaryStyle().Render(fmt.Sprintf("● live  •  refresh %s  •  ctrl-c return", interval))
	_, err := fmt.Fprintf(out, "\x1b[H\x1b[2J%s\n%s\n\n%s", header, status, buffer.String())
	return err
}

func runWatchIteration(args []string, out io.Writer, run watchRunner) error {
	var buffer bytes.Buffer
	if err := run(append([]string(nil), args...), &buffer); err != nil {
		return err
	}
	_, err := fmt.Fprintf(out, "\x1b[H\x1b[2J%s", buffer.String())
	return err
}
