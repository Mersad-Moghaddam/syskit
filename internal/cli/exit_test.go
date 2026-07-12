package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

// TestPresentExitCodes verifies the CLI boundary maps each sentinel/type to the
// canonical exit code (specs/error-handling.md 0–5 table).
func TestPresentExitCodes(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
		wantMsg  bool // whether a user-facing message is expected
	}{
		{name: "nil is success", err: nil, wantCode: exitOK, wantMsg: false},
		{name: "generic is 1", err: errors.New("boom"), wantCode: exitError, wantMsg: true},
		{name: "usage is 2", err: &usageError{err: errors.New("bad flag")}, wantCode: exitUsage, wantMsg: true},
		{
			name:     "permission is 3",
			err:      fmt.Errorf("reading proc/1/status: %w", platform.ErrPermission),
			wantCode: exitPermission,
			wantMsg:  true,
		},
		{
			name:     "unsupported is 4",
			err:      fmt.Errorf("reading pressure: %w", platform.ErrUnsupported),
			wantCode: exitUnsupported,
			wantMsg:  true,
		},
		{
			name:     "partial is 5",
			err:      &PartialError{Err: errors.Join(errors.New("cpu failed"), errors.New("mem failed"))},
			wantCode: exitPartial,
			wantMsg:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, code := present(tt.err)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantCode, exitCode(tt.err), "exitCode must agree with present")
			if tt.wantMsg {
				assert.NotEmpty(t, msg)
			} else {
				assert.Empty(t, msg)
			}
		})
	}
}

// TestPresentPermissionWinsOverWrapping confirms errors.Is matching survives
// wrapping, so a deeply wrapped permission error still maps to 3.
func TestPresentPermissionWinsOverWrapping(t *testing.T) {
	err := fmt.Errorf("collecting cpu: %w", fmt.Errorf("reading proc/stat: %w", platform.ErrPermission))
	_, code := present(err)
	assert.Equal(t, exitPermission, code)
}
