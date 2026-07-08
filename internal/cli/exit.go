package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

// Process exit codes assigned at the CLI boundary. The full canonical table
// (3 permission, 4 unsupported, 5 partial) is defined in
// specs/cli-conventions.md and wired by FND-07; FND-03 establishes the seam with
// success, general error, and usage error.
const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

// usageError marks an error as a command-usage problem — an invalid flag, flag
// value, or argument — so exitCode maps it to exit status 2.
type usageError struct {
	err error
}

func (e *usageError) Error() string { return e.err.Error() }

func (e *usageError) Unwrap() error { return e.err }

// usageArgs wraps a positional-argument validator so a validation failure is
// reported as a usage error (exit 2) rather than a general error. Subcommands
// use it for their Args field so argument-count problems map consistently.
func usageArgs(validate cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := validate(cmd, args); err != nil {
			return &usageError{err: err}
		}
		return nil
	}
}

// exitCode maps an execution error to a process exit code: nil is success (0),
// usage errors are 2, and everything else defaults to a general error (1).
// FND-07 extends the default branch with errors.Is checks for the
// permission/unsupported/partial sentinels.
func exitCode(err error) int {
	if err == nil {
		return exitOK
	}
	var ue *usageError
	if errors.As(err, &ue) {
		return exitUsage
	}
	return exitError
}
