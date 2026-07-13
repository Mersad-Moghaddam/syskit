package render

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// update regenerates golden files when the test binary is run with -update:
//
//	go test ./internal/render -update
var update = flag.Bool("update", false, "regenerate golden files")

// assertGolden compares got against the golden file at
// testdata/golden/<name>, regenerating it when -update is set.
func assertGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", "golden", name)
	if *update {
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, got, 0o644))
		return
	}
	want, err := os.ReadFile(path)
	require.NoError(t, err, "missing golden %s; run: go test ./internal/render -update", name)
	assert.Equal(t, string(want), string(got), "golden mismatch for %s", name)
}

// sampleLoad and sampleMemory are local models with snake_case json tags. The
// render package does not own domain models, so tests exercise JSON rendering
// against a stand-in that mirrors the model contract.
type sampleLoad struct {
	OneMinute      float64 `json:"one_minute"`
	FiveMinutes    float64 `json:"five_minutes"`
	FifteenMinutes float64 `json:"fifteen_minutes"`
}

type sampleMemory struct {
	TotalBytes  int64      `json:"total_bytes"`
	UsedBytes   int64      `json:"used_bytes"`
	SwapEnabled bool       `json:"swap_enabled"`
	LoadAverage sampleLoad `json:"load_average"`
}

func sampleValue() sampleMemory {
	return sampleMemory{
		TotalBytes:  16777216000,
		UsedBytes:   8388608000,
		SwapEnabled: true,
		LoadAverage: sampleLoad{OneMinute: 0.42, FiveMinutes: 0.35, FifteenMinutes: 0.3},
	}
}

func sampleTable() Table {
	return Table{
		Headers: []string{"NAME", "PID", "RSS_BYTES"},
		Rows: [][]string{
			{"systemd", "1", "12058624"},
			{"sshd", "812", "6410240"},
			{"the-long-daemon-name", "40231", "512"},
		},
	}
}

func TestJSONRendererGolden(t *testing.T) {
	r, err := New("json")
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, r.Render(&buf, sampleValue()))
	assertGolden(t, "sample.json.golden", buf.Bytes())
}

func TestJSONRendererDeterministic(t *testing.T) {
	r, err := New("json")
	require.NoError(t, err)

	var a, b bytes.Buffer
	require.NoError(t, r.Render(&a, sampleValue()))
	require.NoError(t, r.Render(&b, sampleValue()))
	assert.Equal(t, a.Bytes(), b.Bytes(), "same input must yield identical bytes")
}

func TestJSONRendererSnakeCaseAndShape(t *testing.T) {
	r, err := New("json")
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, r.Render(&buf, sampleValue()))
	out := buf.String()

	assert.Contains(t, out, `"total_bytes": 16777216000`)
	assert.Contains(t, out, `"load_average"`)
	assert.Contains(t, out, `"fifteen_minutes": 0.3`)
	assert.True(t, bytes.HasSuffix(buf.Bytes(), []byte("\n")), "must end with a trailing newline")
	// 2-space indent for the first nested level.
	assert.Contains(t, out, "\n  \"total_bytes\"")
	// No terminal control sequences in structured output.
	assert.NotContains(t, out, "\x1b")
}

func TestTableRendererGolden(t *testing.T) {
	r, err := New("table")
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, r.Render(&buf, sampleTable()))
	assertGolden(t, "sample_table.golden", buf.Bytes())
}

func TestTableRendererNoHeaderGolden(t *testing.T) {
	r, err := New("table", WithNoHeader(true))
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, r.Render(&buf, sampleTable()))
	assertGolden(t, "sample_table_noheader.golden", buf.Bytes())
}

func TestTableRendererRightAlignsNumericColumns(t *testing.T) {
	r, err := New("table")
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, r.Render(&buf, sampleTable()))
	lines := splitLines(buf.String())
	require.Len(t, lines, 4) // header + 3 rows

	// PID and RSS_BYTES are numeric → right-aligned: their values line up on the
	// right edge of the column. NAME is text → left-aligned.
	// header: "NAME                  PID  RSS_BYTES"
	assert.True(t, hasPrefixField(lines[0], "NAME"), "text header left-aligned")
	// "1" right-aligns under the 5-wide PID column ("40231").
	assert.Contains(t, lines[1], "    1")
	assert.Contains(t, lines[2], "  812")
	// RSS value "512" right-aligns under the widest RSS cell.
	assert.Contains(t, lines[3], "      512")
	// No trailing whitespace on any line.
	for _, ln := range lines {
		assert.Equal(t, ln, trimRight(ln), "no trailing whitespace")
	}
}

func TestTableRendererNoHeaderOmitsHeaderRow(t *testing.T) {
	withHeader, err := New("table")
	require.NoError(t, err)
	noHeader, err := New("table", WithNoHeader(true))
	require.NoError(t, err)

	var a, b bytes.Buffer
	require.NoError(t, withHeader.Render(&a, sampleTable()))
	require.NoError(t, noHeader.Render(&b, sampleTable()))

	assert.Equal(t, 4, len(splitLines(a.String())))
	assert.Equal(t, 3, len(splitLines(b.String())))
	assert.NotContains(t, b.String(), "RSS_BYTES")
}

func TestTableRendererColorNeverInBytes(t *testing.T) {
	plain, err := New("table")
	require.NoError(t, err)
	colored, err := New("table", WithColor(true))
	require.NoError(t, err)

	var a, b bytes.Buffer
	require.NoError(t, plain.Render(&a, sampleTable()))
	require.NoError(t, colored.Render(&b, sampleTable()))

	assert.NotContains(t, b.String(), "\x1b", "no escape sequences in table output")
	assert.Equal(t, a.String(), b.String(), "color option must not alter rendered bytes")
}

func TestTableRendererRejectsNonTable(t *testing.T) {
	r, err := New("table")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = r.Render(&buf, sampleValue())
	assert.ErrorIs(t, err, ErrUnsupportedValue)
	assert.Empty(t, buf.String())
}

func TestNewUnknownFormat(t *testing.T) {
	r, err := New("toml")
	assert.Nil(t, r)
	assert.ErrorIs(t, err, ErrUnknownFormat)
}

func TestYAMLRendererGolden(t *testing.T) {
	r, err := New("yaml")
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, r.Render(&buf, sampleValue()))
	assertGolden(t, "sample.yaml.golden", buf.Bytes())
	assert.Contains(t, buf.String(), "total_bytes:")
	assert.NotContains(t, buf.String(), "TotalBytes")
}

func TestNewReturnsRequestedRenderer(t *testing.T) {
	j, err := New("json")
	require.NoError(t, err)
	_, ok := j.(jsonRenderer)
	assert.True(t, ok)

	tbl, err := New("table")
	require.NoError(t, err)
	_, ok = tbl.(tableRenderer)
	assert.True(t, ok)
}

func BenchmarkTableRenderer1000Rows(b *testing.B) {
	rows := make([][]string, 1000)
	for i := range rows {
		rows[i] = []string{"worker-" + strconv.Itoa(i), strconv.Itoa(i + 1), strconv.FormatUint(uint64(i+1)*1048576, 10), "running"}
	}
	table := Table{Headers: []string{"NAME", "PID", "RSS_BYTES", "STATE"}, Rows: rows}
	r, err := New("table")
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := r.Render(io.Discard, table); err != nil {
			b.Fatal(err)
		}
	}
}

// --- small string helpers local to the test (no shared util package) ---

func splitLines(s string) []string {
	s = trimTrailingNewline(s)
	if s == "" {
		return nil
	}
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

func trimTrailingNewline(s string) string {
	for len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return s
}

func trimRight(s string) string {
	for len(s) > 0 && s[len(s)-1] == ' ' {
		s = s[:len(s)-1]
	}
	return s
}

func hasPrefixField(line, field string) bool {
	return len(line) >= len(field) && line[:len(field)] == field
}
