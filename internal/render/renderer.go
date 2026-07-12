package render

import (
	"errors"
	"fmt"
	"io"
)

// Renderer converts an already-collected value into a user-facing
// representation and writes it to w. It is the single presentation seam for the
// CLI: commands build a value (a domain model for JSON, a Table for tables) and
// hand it to a Renderer selected from the --format flag.
//
// The epic refers to this contract as a "Formatter"; the canonical name here is
// Renderer to match ARCHITECTURE.md §4. The two terms describe the same seam.
//
// Renderers are deterministic (same input yields identical bytes), never
// collect system data, and never emit diagnostics — warnings and errors belong
// on stderr, outside this interface.
type Renderer interface {
	Render(w io.Writer, v any) error
}

// Sentinel errors returned by New and the renderers. Callers branch on these
// with errors.Is rather than string matching.
var (
	// ErrUnknownFormat indicates a --format value that names no renderer.
	ErrUnknownFormat = errors.New("unknown output format")
	// ErrFormatDeferred indicates a format that is planned but not yet
	// implemented. YAML is deferred to v0.2 (FND/OUT-03, ADR-009).
	ErrFormatDeferred = errors.New("output format not yet supported")
	// ErrUnsupportedValue indicates a renderer received a value shape it cannot
	// render (for example, a non-Table value handed to the table renderer).
	ErrUnsupportedValue = errors.New("value not supported by renderer")
)

// options holds the resolved, immutable configuration for a renderer. It is
// populated once by New from the functional Options and then never mutated, so
// renderers carry no mutable package-level or shared state.
type options struct {
	noHeader bool
	color    bool
}

// Option configures a Renderer at construction time.
type Option func(*options)

// WithNoHeader controls whether the table renderer emits its header row. It
// corresponds to the --no-header CLI flag and is ignored by structured
// renderers such as JSON.
func WithNoHeader(noHeader bool) Option {
	return func(o *options) { o.noHeader = noHeader }
}

// WithColor enables terminal color in table output. Color is OFF by default:
// Render writes to an arbitrary io.Writer, so the caller (which knows whether
// stdout is a TTY and whether NO_COLOR / --color apply) must opt in explicitly.
// Structured renderers never emit color regardless of this option.
func WithColor(enabled bool) Option {
	return func(o *options) { o.color = enabled }
}

// New returns the Renderer for the given output format.
//
//   - "json"  → machine-readable indented JSON.
//   - "table" → aligned, width-aware text table.
//   - "yaml"  → ErrFormatDeferred (planned for v0.2 via ADR-009).
//   - other   → ErrUnknownFormat.
//
// This is the seam the CLI's --format flag wires to. The "yaml" case is kept
// distinct from unknown formats so the deferred capability has a clean home to
// grow into.
func New(format string, opts ...Option) (Renderer, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	switch format {
	case "json":
		return jsonRenderer{}, nil
	case "table":
		return tableRenderer{noHeader: o.noHeader, color: o.color}, nil
	case "yaml":
		return nil, fmt.Errorf("format %q: %w", format, ErrFormatDeferred)
	default:
		return nil, fmt.Errorf("format %q: %w", format, ErrUnknownFormat)
	}
}
