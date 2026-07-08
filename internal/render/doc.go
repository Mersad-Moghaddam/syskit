// Package render turns typed model values into user-facing output behind a
// single Formatter/Renderer interface: table and JSON now, YAML at v0.2, and
// the TUI at v0.3.
//
// Renderers are deterministic and golden-file-testable, write only to the
// provided writer (stdout at the CLI boundary; diagnostics go to stderr
// elsewhere), and never collect data. Per ADR-004 render depends only on model
// and the standard library.
package render
