package render

import (
	"fmt"
	"io"
	"strings"
)

// Table is the general contract the table renderer consumes. Commands build a
// Table from a domain model — choosing columns, formatting units, and ordering
// rows — and hand it to the table Renderer, which owns only alignment, padding,
// and the header. This keeps the render layer free of domain knowledge.
//
// Every row should have the same number of cells as Headers. Missing trailing
// cells are treated as empty; a row longer than Headers still aligns, widening
// the table to fit.
type Table struct {
	// Headers are the column titles, emitted as the header row unless the
	// renderer is configured with WithNoHeader.
	Headers []string
	// Rows are the data rows, one []string of cells per row.
	Rows [][]string
}

// tableRenderer formats a Table as an aligned, width-aware text table. Text
// cells are left-aligned and numeric cells are right-aligned per
// specs/rendering.md. Columns are padded to the widest cell (header or body).
//
// Color, when enabled by the CLI for a terminal, emphasizes only the header.
// Structured renderers ignore the option and never emit terminal escapes.
type tableRenderer struct {
	noHeader bool
	color    bool
}

// Render writes t as an aligned table. v must be a Table or *Table; any other
// value yields ErrUnsupportedValue. Columns are separated by two spaces and each
// line is trimmed of trailing padding so the output has no dangling whitespace.
func (r tableRenderer) Render(w io.Writer, v any) error {
	t, ok := asTable(v)
	if !ok {
		return fmt.Errorf("render table: got %T, want render.Table: %w", v, ErrUnsupportedValue)
	}

	cols := columnCount(t)
	widths := columnWidths(t, cols, r.noHeader)
	rightAlign := numericColumns(t, cols)

	var b strings.Builder
	if !r.noHeader {
		if r.color {
			b.WriteString("\x1b[1m")
		}
		writeRow(&b, t.Headers, cols, widths, rightAlign)
		if r.color {
			b.WriteString("\x1b[0m")
		}
	}
	for _, row := range t.Rows {
		writeRow(&b, row, cols, widths, rightAlign)
	}

	if _, err := io.WriteString(w, b.String()); err != nil {
		return fmt.Errorf("writing table: %w", err)
	}
	return nil
}

// asTable normalizes the accepted value shapes into a Table.
func asTable(v any) (Table, bool) {
	switch t := v.(type) {
	case Table:
		return t, true
	case *Table:
		if t == nil {
			return Table{}, false
		}
		return *t, true
	default:
		return Table{}, false
	}
}

// columnCount returns the number of columns, taking the maximum across the
// header and every row so ragged input still renders fully.
func columnCount(t Table) int {
	n := len(t.Headers)
	for _, row := range t.Rows {
		if len(row) > n {
			n = len(row)
		}
	}
	return n
}

// columnWidths computes the display width of each column as the widest cell in
// that column. The header is excluded when noHeader is set so a suppressed
// header never pads the body.
func columnWidths(t Table, cols int, noHeader bool) []int {
	widths := make([]int, cols)
	if !noHeader {
		for i, h := range t.Headers {
			if w := len(h); w > widths[i] {
				widths[i] = w
			}
		}
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if w := len(cell); w > widths[i] {
				widths[i] = w
			}
		}
	}
	return widths
}

// numericColumns reports, per column, whether every non-empty body cell is
// numeric. Such columns are right-aligned per specs/rendering.md; headers are
// ignored so a text header over numeric data still right-aligns the numbers.
func numericColumns(t Table, cols int) []bool {
	right := make([]bool, cols)
	seen := make([]bool, cols)
	for i := range right {
		right[i] = true
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i >= cols {
				break
			}
			if strings.TrimSpace(cell) == "" {
				continue
			}
			seen[i] = true
			if !isNumeric(cell) {
				right[i] = false
			}
		}
	}
	// A column with no observed values is not treated as numeric.
	for i := range right {
		right[i] = right[i] && seen[i]
	}
	return right
}

// writeRow renders one row padded to widths and right-aligns the columns flagged
// in rightAlign. It stops at the final non-empty cell so lines have no trailing
// whitespace.
func writeRow(b *strings.Builder, cells []string, cols int, widths []int, rightAlign []bool) {
	last := -1
	for i := cols - 1; i >= 0; i-- {
		if i < len(cells) && cells[i] != "" {
			last = i
			break
		}
	}
	for i := 0; i <= last; i++ {
		if i > 0 {
			b.WriteString("  ")
		}
		var cell string
		if i < len(cells) {
			cell = cells[i]
		}
		pad := widths[i] - len(cell)
		if pad < 0 {
			pad = 0
		}
		if rightAlign[i] {
			writePadding(b, pad)
			b.WriteString(cell)
		} else {
			b.WriteString(cell)
			if i < last {
				writePadding(b, pad)
			}
		}
	}
	b.WriteByte('\n')
}

func writePadding(b *strings.Builder, count int) {
	for range count {
		b.WriteByte(' ')
	}
}

// isNumeric reports whether s is a plain decimal number: optional leading sign,
// digits, an optional single decimal point, and nothing else. It deliberately
// rejects formatted values like "1.2 GB" or "42%" so only true numeric columns
// right-align.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	i := 0
	if s[0] == '+' || s[0] == '-' {
		i++
	}
	digits := false
	dot := false
	for ; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			digits = true
		case c == '.' && !dot:
			dot = true
		default:
			return false
		}
	}
	return digits
}
