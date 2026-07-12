// Package golden provides the shared golden-file test helper used across
// SysKit's packages to pin user-facing output.
//
// A golden file is the recorded, reviewed expected output for a command or
// renderer. A test produces bytes, and [Assert] compares them against the
// stored golden; a mismatch fails the test with a readable diff. When the
// output contract intentionally changes, the golden is regenerated with the
// -update flag and the diff is reviewed like any other change
// (specs/testing-strategy.md, "Golden-File & End-to-End Tests"):
//
//	go test ./... -update
//
// Golden files live in each package's own testdata/golden/ directory. Go runs
// tests with the working directory set to the package directory, so a caller
// passing the name "cpu_json.golden" resolves to
// <caller-package>/testdata/golden/cpu_json.golden. Fixtures read with [Read]
// resolve relative to the caller's testdata/ in the same way.
package golden

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// update reports whether golden files should be rewritten from the test output
// instead of compared against it. It is wired to the -update test flag so that
// `go test ./... -update` regenerates every golden. It is read only after
// flag.Parse (which the testing package calls before tests run), so it
// introduces no mutable state during a test.
var update = flag.Bool("update", false, "rewrite golden files from test output")

// goldenDir is the conventional subdirectory, under a package's testdata/, that
// holds golden files.
const goldenDir = "golden"

// Assert compares got against the golden file named by name (resolved to
// testdata/golden/<name> in the calling package). On a byte mismatch it fails
// the test with a diff. When the -update flag is set, it rewrites the golden
// file with got and does not fail, so the recorded contract tracks the
// intentional change.
func Assert(t testing.TB, got []byte, name string) {
	t.Helper()
	path := filepath.Join("testdata", goldenDir, name)
	assert(t, got, path, *update)
}

// AssertString is the string convenience form of [Assert].
func AssertString(t testing.TB, got, name string) {
	t.Helper()
	Assert(t, []byte(got), name)
}

// assert is the flag-independent core of [Assert], separated so the
// compare-and-update logic is unit-testable without toggling the global -update
// flag. When update is true it writes got to path (creating the directory) and
// returns; otherwise it reads path and fails t if the bytes differ or the file
// is missing.
func assert(t testing.TB, got []byte, path string, update bool) {
	t.Helper()

	if update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("golden: creating directory for %s: %v", path, err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("golden: writing %s: %v", path, err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden: reading %s: %v\n(run `go test ./... -update` to create it)", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("golden mismatch for %s\n--- want ---\n%s\n--- got ---\n%s\n(run `go test ./... -update` to accept the new output)",
			path, want, got)
	}
}

// Read returns the contents of a fixture file located relative to the calling
// package's testdata/ directory (e.g. Read(t, "proc/stat") reads
// testdata/proc/stat). It fails the test if the file cannot be read, so callers
// use the returned bytes directly.
func Read(t testing.TB, path string) []byte {
	t.Helper()
	full := filepath.Join("testdata", filepath.FromSlash(path))
	data, err := os.ReadFile(full)
	if err != nil {
		t.Fatalf("golden: reading fixture %s: %v", full, err)
	}
	return data
}
