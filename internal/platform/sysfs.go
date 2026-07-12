package platform

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// SysFS is the narrow, read-only seam through which every collector reads
// procfs and sysfs. name is always a slash-relative path rooted at the mount
// root (e.g. "proc/stat", "sys/devices/system/cpu/present"), never absolute.
// Production code uses RealFS (rooted at "/"); tests use TestFS (rooted at a
// fixtures directory), so the same collector runs unchanged in both.
type SysFS interface {
	// ReadFile reads the entire named file relative to the mount root.
	ReadFile(name string) ([]byte, error)
	// Open opens the named file for streaming reads of large pseudo-files.
	Open(name string) (fs.File, error)
	// ReadDir lists a pseudo-directory (e.g. "proc" to enumerate PIDs).
	ReadDir(name string) ([]fs.DirEntry, error)
	StatFS(path string) (FSStats, error)
}

// FSStats is the filesystem capacity and inode subset used by collectors.
type FSStats struct{ TotalBytes, FreeBytes, AvailableBytes, TotalInodes, FreeInodes uint64 }

// osFS is a SysFS backed by the real operating-system filesystem, rooted at
// root. It is the only place in SysKit that touches the OS filesystem.
type osFS struct {
	root string
}

// RealFS returns a SysFS rooted at "/", mapping e.g. "proc/stat" to
// "/proc/stat".
func RealFS() SysFS { return osFS{root: "/"} }

// resolve maps a validated slash-relative name to an absolute host path.
func (o osFS) resolve(name string) string {
	return filepath.Join(o.root, filepath.FromSlash(name))
}

func (o osFS) ReadFile(name string) ([]byte, error) {
	if err := validate("reading", name); err != nil {
		return nil, err
	}
	// Critical: /proc and /sys pseudo-files routinely report a size of 0 in
	// their FileInfo even though reading them yields data. os.ReadFile trusts
	// that size hint to pre-size its buffer and can therefore under-read such
	// files. We instead open the file and drain it to EOF with io.ReadAll,
	// which grows its buffer dynamically and returns the true contents
	// regardless of the advertised size.
	f, err := os.Open(o.resolve(name))
	if err != nil {
		return nil, mapError("reading", name, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, mapError("reading", name, err)
	}
	return data, nil
}

func (o osFS) Open(name string) (fs.File, error) {
	if err := validate("opening", name); err != nil {
		return nil, err
	}
	f, err := os.Open(o.resolve(name))
	if err != nil {
		return nil, mapError("opening", name, err)
	}
	return f, nil
}

func (o osFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if err := validate("listing", name); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(o.resolve(name))
	if err != nil {
		return nil, mapError("listing", name, err)
	}
	return entries, nil
}

func (o osFS) StatFS(path string) (FSStats, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return FSStats{}, mapError("statting", path, err)
	}
	block := uint64(stat.Bsize)
	return FSStats{TotalBytes: stat.Blocks * block, FreeBytes: stat.Bfree * block, AvailableBytes: stat.Bavail * block, TotalInodes: stat.Files, FreeInodes: stat.Ffree}, nil
}

// fixtureFS is a SysFS backed by an arbitrary fs.FS, typically an
// os.DirFS rooted at a fixtures directory, so that ReadFile("proc/stat")
// resolves within the fixture root.
type fixtureFS struct {
	fsys fs.FS
}

// TestFS returns a SysFS backed by fsys (for example
// os.DirFS("testdata/fixtures/idle-host")), so collectors can be exercised
// against captured fixtures instead of the host kernel.
func TestFS(fsys fs.FS) SysFS { return fixtureFS{fsys: fsys} }

func (t fixtureFS) ReadFile(name string) ([]byte, error) {
	if err := validate("reading", name); err != nil {
		return nil, err
	}
	data, err := fs.ReadFile(t.fsys, name)
	if err != nil {
		return nil, mapError("reading", name, err)
	}
	return data, nil
}

func (t fixtureFS) Open(name string) (fs.File, error) {
	if err := validate("opening", name); err != nil {
		return nil, err
	}
	f, err := t.fsys.Open(name)
	if err != nil {
		return nil, mapError("opening", name, err)
	}
	return f, nil
}

func (t fixtureFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if err := validate("listing", name); err != nil {
		return nil, err
	}
	entries, err := fs.ReadDir(t.fsys, name)
	if err != nil {
		return nil, mapError("listing", name, err)
	}
	return entries, nil
}

func (t fixtureFS) StatFS(path string) (FSStats, error) {
	return FSStats{}, fmt.Errorf("statting %q: %w", path, ErrUnsupported)
}
