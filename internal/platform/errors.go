package platform

import (
	"errors"
	"fmt"
	"io/fs"
)

// Sentinel errors for well-known kernel-interface failure conditions. Callers
// branch on these with errors.Is rather than string matching (see
// specs/error-handling.md).
var (
	// ErrNotFound indicates a kernel interface file or directory is absent.
	ErrNotFound = errors.New("kernel interface not found")
	// ErrPermission indicates insufficient privilege to read an interface.
	ErrPermission = errors.New("permission denied reading kernel interface")
	// ErrUnsupported indicates the kernel does not provide this interface.
	ErrUnsupported = errors.New("kernel interface not supported")
)

// mapError wraps a low-level filesystem error with operation context and, for
// recognised conditions, the matching platform sentinel.
//
// Wrapping choice: for fs.ErrNotExist and fs.ErrPermission the returned error
// wraps BOTH the platform sentinel and the original error using two %w verbs
// (Go 1.20+). This keeps errors.Is(err, ErrNotFound)/ErrPermission working for
// callers that branch on sentinels, while errors.As(err, &pathErr) still
// recovers the underlying *fs.PathError for detailed diagnostics.
func mapError(op, name string, err error) error {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return fmt.Errorf("%s %q: %w: %w", op, name, ErrNotFound, err)
	case errors.Is(err, fs.ErrPermission):
		return fmt.Errorf("%s %q: %w: %w", op, name, ErrPermission, err)
	default:
		return fmt.Errorf("%s %q: %w", op, name, err)
	}
}

// validate ensures name is a safe, slash-relative path rooted at the mount root.
// It rejects absolute paths, ".." traversal, and other non-canonical forms via
// fs.ValidPath. An invalid path can never resolve to a real kernel interface,
// so it is reported as ErrNotFound.
func validate(op, name string) error {
	if !fs.ValidPath(name) {
		return fmt.Errorf("%s %q: invalid interface path: %w", op, name, ErrNotFound)
	}
	return nil
}
