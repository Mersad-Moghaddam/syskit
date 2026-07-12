package cli

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResolveVerbosityPrecedence checks the quiet > debug > verbose precedence
// (specs/logging-strategy.md "Verbosity Flags").
func TestResolveVerbosityPrecedence(t *testing.T) {
	tests := []struct {
		name                  string
		verbose, debug, quiet bool
		want                  verbosity
	}{
		{name: "none is normal", want: verbosityNormal},
		{name: "verbose", verbose: true, want: verbosityVerbose},
		{name: "debug", debug: true, want: verbosityDebug},
		{name: "quiet", quiet: true, want: verbosityQuiet},
		{name: "quiet beats debug", debug: true, quiet: true, want: verbosityQuiet},
		{name: "debug beats verbose", verbose: true, debug: true, want: verbosityDebug},
		{name: "quiet beats all", verbose: true, debug: true, quiet: true, want: verbosityQuiet},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, resolveVerbosity(tt.verbose, tt.debug, tt.quiet))
		})
	}
}

// TestLoggerEmitsByLevel verifies that the logger admits messages only at or
// above its configured level: silent by default, info at verbose, debug at
// debug, and nothing at quiet.
func TestLoggerEmitsByLevel(t *testing.T) {
	tests := []struct {
		name      string
		v         verbosity
		wantInfo  bool
		wantDebug bool
		wantError bool
	}{
		{name: "normal is silent", v: verbosityNormal},
		{name: "verbose emits info not debug", v: verbosityVerbose, wantInfo: true, wantError: true},
		{name: "debug emits debug and info", v: verbosityDebug, wantInfo: true, wantDebug: true, wantError: true},
		{name: "quiet silences even error", v: verbosityQuiet},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := newLogger(&buf, tt.v)

			log.Debug("dbg")
			assert.Equal(t, tt.wantDebug, bytes.Contains(buf.Bytes(), []byte("dbg")))

			buf.Reset()
			log.Info("nfo")
			assert.Equal(t, tt.wantInfo, bytes.Contains(buf.Bytes(), []byte("nfo")))

			buf.Reset()
			log.Error("err")
			assert.Equal(t, tt.wantError, bytes.Contains(buf.Bytes(), []byte("err")))
		})
	}
}

// TestLoggerLevelMapping pins the slog level chosen for each verbosity.
func TestLoggerLevelMapping(t *testing.T) {
	assert.Equal(t, silentLevel, levelFor(verbosityNormal))
	assert.Equal(t, silentLevel, levelFor(verbosityQuiet))
	assert.Equal(t, slog.LevelInfo, levelFor(verbosityVerbose))
	assert.Equal(t, slog.LevelDebug, levelFor(verbosityDebug))
}
