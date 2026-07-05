package filesystem

import (
	"errors"
	"fmt"
	"os"
)

// FS abstracts the filesystem operations to enable dry-runs and easier testing.
type FS interface {
	MkdirAll(path string, perm os.FileMode) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Exists(path string) (bool, error)
}

// OSFS implements FS by interacting with the real disk.
type OSFS struct{}

func (OSFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OSFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (OSFS) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// DryRunFS implements FS by only printing what would happen, without making changes.
type DryRunFS struct {
	// We embed OSFS just for Exists so we know if a file exists on disk
	OSFS
}

func (DryRunFS) MkdirAll(path string, perm os.FileMode) error {
	fmt.Printf("[DRY-RUN] mkdir -p %s\n", path)
	return nil
}

func (DryRunFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	fmt.Printf("[DRY-RUN] write %s (%d bytes)\n", filename, len(data))
	return nil
}
