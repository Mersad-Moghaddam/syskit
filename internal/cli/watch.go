package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

const (
	minWatchInterval = 250 * time.Millisecond
	maxWatchInterval = time.Minute
)

type watchRunner func([]string, io.Writer) error

func newWatchCmd(run watchRunner) *cobra.Command {
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
		return runWatch(ctx, interval, args, cmd.OutOrStdout(), run)
	}}
	cmd.Flags().DurationVar(&interval, "interval", time.Second, "refresh interval (250ms to 1m)")
	return cmd
}

func runWatch(ctx context.Context, interval time.Duration, args []string, out io.Writer, run watchRunner) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := runWatchIteration(args, out, run); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func runWatchIteration(args []string, out io.Writer, run watchRunner) error {
	var buffer bytes.Buffer
	if err := run(append([]string(nil), args...), &buffer); err != nil {
		return err
	}
	_, err := fmt.Fprintf(out, "\x1b[H\x1b[2J%s", buffer.String())
	return err
}
