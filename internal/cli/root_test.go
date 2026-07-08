package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// executeArgs builds a fresh root command, runs it with the given args, and
// returns captured stdout, stderr, and the execution error.
func executeArgs(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	root := newRootCmd()
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)

	err := root.Execute()
	return stdout.String(), stderr.String(), err
}

func TestRootFormatFlag(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		wantUsage bool
	}{
		{name: "table", format: "table"},
		{name: "json", format: "json"},
		{name: "yaml", format: "yaml"},
		{name: "invalid rejected", format: "xml", wantUsage: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// version is a real subcommand, so PersistentPreRunE runs and
			// validates --format.
			_, _, err := executeArgs(t, "--format", tt.format, "version")

			if tt.wantUsage {
				require.Error(t, err)
				assert.Equal(t, exitUsage, exitCode(err))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, exitOK, exitCode(err))
		})
	}
}

func TestRootDefaultFormatIsTable(t *testing.T) {
	root := newRootCmd()
	f := root.PersistentFlags().Lookup("format")
	require.NotNil(t, f)
	assert.Equal(t, formatTable, f.DefValue)
}

// TestRootInvalidFormatNoSubcommand guards the bare-root path: an invalid
// --format with no subcommand must still be a usage error (exit 2), not a
// silent help dump. This regresses a defect where the root command lacked a
// RunE, so Cobra skipped PersistentPreRunE and printed help with exit 0.
func TestRootInvalidFormatNoSubcommand(t *testing.T) {
	_, _, err := executeArgs(t, "--format", "xml")
	require.Error(t, err)
	assert.Equal(t, exitUsage, exitCode(err))
}

// TestBareRootPrintsHelp confirms the bare `syskit` invocation succeeds and
// prints help rather than erroring.
func TestBareRootPrintsHelp(t *testing.T) {
	stdout, _, err := executeArgs(t)
	require.NoError(t, err)
	assert.Equal(t, exitOK, exitCode(err))
	assert.Contains(t, stdout, "syskit")
	assert.Contains(t, stdout, "Available Commands")
}

func TestUnknownFlagIsUsageError(t *testing.T) {
	_, _, err := executeArgs(t, "--nope")
	require.Error(t, err)
	assert.Equal(t, exitUsage, exitCode(err))
}

func TestVersionSubcommand(t *testing.T) {
	stdout, _, err := executeArgs(t, "version")
	require.NoError(t, err)
	assert.Equal(t, exitOK, exitCode(err))
	assert.Equal(t, version+"\n", stdout)
}

func TestVersionRejectsArgs(t *testing.T) {
	_, _, err := executeArgs(t, "version", "extra")
	require.Error(t, err)
	assert.Equal(t, exitUsage, exitCode(err))
}

func TestHelpSucceeds(t *testing.T) {
	stdout, _, err := executeArgs(t, "--help")
	require.NoError(t, err)
	assert.Contains(t, stdout, "syskit")
	assert.Contains(t, stdout, "read-only")
}

func TestExitCodeMapping(t *testing.T) {
	assert.Equal(t, exitOK, exitCode(nil))
	assert.Equal(t, exitError, exitCode(assertErr("boom")))
	assert.Equal(t, exitUsage, exitCode(&usageError{err: assertErr("bad flag")}))
}

// assertErr is a tiny error helper for the mapping table above.
type assertErr string

func (e assertErr) Error() string { return string(e) }
