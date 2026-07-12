package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// Process exit codes assigned at the CLI boundary. This table is canonical in
// specs/error-handling.md and mirrored in specs/cli-conventions.md; the two
// must always match. Codes are derived from the sentinel errors that propagate
// up from the platform and collector layers — never assigned deeper than the
// CLI.
const (
	exitOK          = 0 // success
	exitError       = 1 // unspecified runtime error
	exitUsage       = 2 // invalid flags, arguments, or usage
	exitPermission  = 3 // insufficient privilege to read a kernel interface
	exitUnsupported = 4 // required kernel interface missing or unsupported
	exitPartial     = 5 // some data collected; one or more collectors failed
)

// usageError marks an error as a command-usage problem — an invalid flag, flag
// value, or argument — so the boundary maps it to exit status 2.
type usageError struct {
	err error
}

func (e *usageError) Error() string { return e.err.Error() }

func (e *usageError) Unwrap() error { return e.err }

// PartialError reports that a command aggregated several collectors and some,
// but not all, failed. The successfully collected data is still rendered to
// stdout; the joined diagnostics are surfaced to stderr and the process exits
// with exitPartial (5). It is the seam the service layer uses to signal partial
// failure without blanking the whole result over one missing field
// (specs/error-handling.md "Partial-Failure Handling").
type PartialError struct {
	// Err is the joined error (typically errors.Join of each collector
	// failure), so every underlying cause remains inspectable with errors.Is.
	Err error
}

func (e *PartialError) Error() string {
	if e.Err == nil {
		return "partial failure"
	}
	return e.Err.Error()
}

func (e *PartialError) Unwrap() error { return e.Err }

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

// present translates an internal error into a user-facing message and the exit
// code it maps to. It is the single boundary where wrapped, lowercase internal
// errors become actionable full-sentence diagnostics for stderr
// (specs/error-handling.md "Internal vs. User-Facing Errors"). The message is
// empty when err is nil, and empty for usage errors because Cobra prints those
// itself.
func present(err error) (message string, code int) {
	switch {
	case err == nil:
		return "", exitOK
	case errors.Is(err, platform.ErrPermission):
		return "Permission denied reading a kernel interface. Try running with elevated privileges (sudo).", exitPermission
	case errors.Is(err, platform.ErrUnsupported):
		return "This information is not available on your kernel.", exitUnsupported
	}

	var perr *PartialError
	if errors.As(err, &perr) {
		return fmt.Sprintf("Some data could not be collected: %v", err), exitPartial
	}

	var uerr *usageError
	if errors.As(err, &uerr) {
		// Usage problems are concise and actionable; the message names the bad
		// flag/argument. Exit 2 signals usage error to scripts.
		return fmt.Sprintf("Error: %v\nRun 'syskit --help' for usage.", err), exitUsage
	}

	return fmt.Sprintf("Error: %v", err), exitError
}

// exitCode maps an execution error to a process exit code using the same
// mapping as present, so the two never diverge.
func exitCode(err error) int {
	_, code := present(err)
	return code
}
