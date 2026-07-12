package cli

import (
	"io"
	"log/slog"
)

// verbosity is the effective diagnostic level chosen from the --verbose,
// --debug, and --quiet flags. Logging is a CLI-layer concern: collectors,
// services, and the platform layer never log — they return errors, and the CLI
// decides whether and how to surface them (specs/logging-strategy.md).
type verbosity int

const (
	// verbosityNormal is the default: silent success. Nothing is emitted on a
	// successful run.
	verbosityNormal verbosity = iota
	// verbosityVerbose enables info-level diagnostics (--verbose/-v).
	verbosityVerbose
	// verbosityDebug enables debug-level diagnostics (--debug).
	verbosityDebug
	// verbosityQuiet suppresses all diagnostics, including errors (--quiet/-q).
	verbosityQuiet
)

// silentLevel is above slog.LevelError, so a handler set to it emits nothing —
// used for the default (silent success) and for --quiet.
const silentLevel = slog.LevelError + 4

// resolveVerbosity applies the precedence quiet > debug > verbose
// (specs/logging-strategy.md "Verbosity Flags"). --quiet is the strongest
// signal: the user asked for silence and SysKit honors it even for errors.
func resolveVerbosity(verbose, debug, quiet bool) verbosity {
	switch {
	case quiet:
		return verbosityQuiet
	case debug:
		return verbosityDebug
	case verbose:
		return verbosityVerbose
	default:
		return verbosityNormal
	}
}

// levelFor maps a verbosity to the slog level its handler should admit.
func levelFor(v verbosity) slog.Level {
	switch v {
	case verbosityDebug:
		return slog.LevelDebug
	case verbosityVerbose:
		return slog.LevelInfo
	case verbosityQuiet, verbosityNormal:
		return silentLevel
	default:
		return silentLevel
	}
}

// newLogger builds the diagnostic logger. It writes to w (always stderr in
// production — never stdout, so structured data on stdout is never corrupted)
// with a slog text handler whose level is set from the verbosity. At the
// default and quiet levels the handler admits nothing, so a successful run is
// silent (specs/logging-strategy.md "Structured Logging").
func newLogger(w io.Writer, v verbosity) *slog.Logger {
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{Level: levelFor(v)})
	return slog.New(handler)
}
