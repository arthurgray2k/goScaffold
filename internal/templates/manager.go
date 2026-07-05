package templates

import (
	"errors"
	"io/fs"
)

// ErrTemplateNotFound is returned when a requested template does not exist.
var ErrTemplateNotFound = errors.New("template not found")

// Manager defines the interface for interacting with embedded templates.
type Manager interface {
	// List returns the names of all available templates.
	List() ([]string, error)
	
	// Get returns an fs.FS restricted to the specific template directory.
	Get(name string) (fs.FS, error)
}

type fsManager struct {
	baseFS fs.FS
}

// NewManager creates a new Manager backed by the provided filesystem.
// The provided fs should have a "templates" directory at its root.
func NewManager(baseFS fs.FS) Manager {
	return &fsManager{
		baseFS: baseFS,
	}
}

func (m *fsManager) List() ([]string, error) {
	entries, err := fs.ReadDir(m.baseFS, "templates")
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}
	return templates, nil
}

func (m *fsManager) Get(name string) (fs.FS, error) {
	templatePath := "templates/" + name

	// Verify it exists and is a directory
	info, err := fs.Stat(m.baseFS, templatePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrTemplateNotFound
	}

	subFS, err := fs.Sub(m.baseFS, templatePath)
	if err != nil {
		return nil, err
	}

	return subFS, nil
}
