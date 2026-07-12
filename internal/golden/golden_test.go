package golden

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAssertMatchesEqualBytes confirms the public Assert passes when the test
// output equals the recorded golden file (internal/golden/testdata/golden/).
func TestAssertMatchesEqualBytes(t *testing.T) {
	Assert(t, []byte("hello golden\n"), "sample.golden")
}

// TestAssertStringMatches confirms the string convenience wrapper agrees with
// the same golden file.
func TestAssertStringMatches(t *testing.T) {
	AssertString(t, "hello golden\n", "sample.golden")
}

// TestAssertCoreCompare exercises the flag-independent core directly, so both
// the matching and mismatching compare paths are covered without touching the
// global -update flag. A temp dir holds the golden so real testdata is never
// modified.
func TestAssertCoreCompare(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "golden", "core.golden")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("expected\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Equal bytes: a real *testing.T must not fail.
	assert(t, []byte("expected\n"), path, false)

	// Different bytes: the core must record a failure. Use a spy TB so the
	// error does not fail this test.
	spy := &spyTB{TB: t}
	assert(spy, []byte("different\n"), path, false)
	if !spy.failed {
		t.Errorf("assert did not report a mismatch for differing bytes")
	}
}

// TestAssertCoreUpdateWrites confirms update=true creates the golden file (and
// its directory) and does not fail, then that a subsequent compare matches.
func TestAssertCoreUpdateWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "golden", "written.golden")

	assert(t, []byte("fresh output\n"), path, true) // writes, must not fail

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("update did not create the golden file: %v", err)
	}
	if string(got) != "fresh output\n" {
		t.Errorf("written golden = %q, want %q", got, "fresh output\n")
	}

	// The freshly written golden now compares equal.
	assert(t, []byte("fresh output\n"), path, false)
}

// TestReadFixture confirms Read returns fixture bytes relative to testdata/.
func TestReadFixture(t *testing.T) {
	got := Read(t, "fixtures/sample.txt")
	if string(got) != "fixture body\n" {
		t.Errorf("Read = %q, want %q", got, "fixture body\n")
	}
}

// spyTB is a testing.TB that records whether Error/Fatal was called instead of
// failing the enclosing test, so the mismatch path of assert can be asserted.
type spyTB struct {
	testing.TB
	failed bool
}

func (s *spyTB) Errorf(format string, args ...any) { s.failed = true }
func (s *spyTB) Error(args ...any)                 { s.failed = true }
func (s *spyTB) Fatalf(format string, args ...any) { s.failed = true }
func (s *spyTB) Fatal(args ...any)                 { s.failed = true }
func (s *spyTB) Helper()                           {}
