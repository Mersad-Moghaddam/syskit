// Package render turns already-collected values into user-facing output behind
// a single Renderer interface: table and JSON now, YAML at v0.2, and the TUI
// later. The epic calls this seam a "Formatter"; the canonical name is Renderer
// to match ARCHITECTURE.md §4.
//
// Structured renderers (JSON) consume a domain model directly and rely on its
// json: struct tags for snake_case field names. The table renderer consumes the
// package-local Table contract (Headers + Rows of pre-formatted strings), which
// commands build from a model — keeping domain knowledge out of the render
// layer.
//
// Renderers are deterministic and golden-file-testable, write only to the
// provided writer (stdout at the CLI boundary; diagnostics go to stderr
// elsewhere), and never collect data. Color stays out of rendered bytes so
// golden files remain stable. Per ADR-004 render depends only on model and the
// standard library.
package render
