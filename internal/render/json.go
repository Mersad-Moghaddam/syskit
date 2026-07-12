package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// jsonRenderer emits v as indented JSON: pure machine-readable stdout with no
// terminal control sequences and no diagnostics. Field names come from the
// value's own `json:` struct tags (snake_case per the model contract), so the
// renderer stays agnostic to any specific domain model.
type jsonRenderer struct{}

// Render marshals v to indented JSON (2-space indent) with a trailing newline.
//
// It is deterministic: encoding/json emits struct fields in declaration order
// and map keys in sorted order, so identical input produces identical bytes.
// HTML escaping is disabled so characters such as <, > and & (common in process
// command lines) are written literally rather than as <-style escapes.
func (jsonRenderer) Render(w io.Writer, v any) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	// Encode already appends a trailing newline.
	if _, err := w.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("writing json: %w", err)
	}
	return nil
}
