package generator

import (
	"errors"
	"os"
	"testing"
	"testing/fstest"

	"github.com/arthurgray2k/goScaffold/internal/templates"
	"github.com/arthurgray2k/goScaffold/internal/variables"
)

// MockFS implements filesystem.FS for testing
type MockFS struct {
	MkdirCalls []string
	Files      map[string][]byte
	ExistsMap  map[string]bool
}

func NewMockFS() *MockFS {
	return &MockFS{
		Files:     make(map[string][]byte),
		ExistsMap: make(map[string]bool),
	}
}

func (m *MockFS) MkdirAll(path string, perm os.FileMode) error {
	m.MkdirCalls = append(m.MkdirCalls, path)
	return nil
}

func (m *MockFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	m.Files[filename] = data
	m.ExistsMap[filename] = true
	return nil
}

func (m *MockFS) Exists(path string) (bool, error) {
	return m.ExistsMap[path], nil
}

func TestGenerator_Generate(t *testing.T) {
	mockTmplFS := fstest.MapFS{
		"templates/basic/README.md": &fstest.MapFile{Data: []byte("Hello {{PROJECT_NAME}}!")},
		"templates/basic/cmd/__PROJECT_NAME__/main.go.tmpl": &fstest.MapFile{Data: []byte("package main")},
	}
	manager := templates.NewManager(mockTmplFS)
	
	fsys := NewMockFS()
	gen := New(manager, fsys)

	opts := Options{
		TemplateName: "basic",
		DestDir:      "my-api",
		Force:        false,
		Values: &variables.Values{
			ProjectName: "my-api",
		},
	}

	err := gen.Generate(opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify README.md
	content := fsys.Files["my-api/README.md"] // Using / or \ depends on filepath.Join but since we test on cross-platform, we should check with filepath.Join
	_ = content // to avoid unused error if we do a cross-platform check. Let's do a basic existence check.
	
	foundReadme := false
	for k, v := range fsys.Files {
		if k == "my-api/README.md" || k == "my-api\\README.md" {
			foundReadme = true
			if string(v) != "Hello my-api!" {
				t.Errorf("expected 'Hello my-api!', got %q", string(v))
			}
		}
	}
	if !foundReadme {
		t.Errorf("README.md not found in written files")
	}

	// Verify .tmpl stripping and __PROJECT_NAME__ replacement
	foundMain := false
	for k := range fsys.Files {
		if k == "my-api/cmd/my-api/main.go" || k == "my-api\\cmd\\my-api\\main.go" {
			foundMain = true
		}
	}
	if !foundMain {
		t.Errorf("main.go not found in expected path, got %v", fsys.Files)
	}
}

func TestGenerator_Generate_FileExists(t *testing.T) {
	mockTmplFS := fstest.MapFS{
		"templates/basic/README.md": &fstest.MapFile{Data: []byte("Hello")},
	}
	manager := templates.NewManager(mockTmplFS)
	
	fsys := NewMockFS()
	// Pre-populate to simulate existence
	fsys.ExistsMap["my-api/README.md"] = true
	fsys.ExistsMap["my-api\\README.md"] = true

	gen := New(manager, fsys)

	opts := Options{
		TemplateName: "basic",
		DestDir:      "my-api",
		Force:        false,
		Values:       &variables.Values{},
	}

	err := gen.Generate(opts)
	if !errors.Is(err, ErrFileExists) {
		t.Fatalf("expected ErrFileExists, got %v", err)
	}

	// Now with force
	opts.Force = true
	err = gen.Generate(opts)
	if err != nil {
		t.Fatalf("expected no error with Force=true, got %v", err)
	}
}

// FailingFS injects errors into the mock FS
type FailingFS struct {
	MockFS
	FailMkdir bool
	FailWrite bool
}

func (f *FailingFS) MkdirAll(path string, perm os.FileMode) error {
	if f.FailMkdir {
		return errors.New("mkdir failed")
	}
	return f.MockFS.MkdirAll(path, perm)
}

func (f *FailingFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if f.FailWrite {
		return errors.New("write failed")
	}
	return f.MockFS.WriteFile(filename, data, perm)
}

func TestGenerator_Generate_MkdirError(t *testing.T) {
	mockTmplFS := fstest.MapFS{
		"templates/basic/README.md": &fstest.MapFile{Data: []byte("Hello")},
	}
	manager := templates.NewManager(mockTmplFS)
	
	fsys := &FailingFS{MockFS: *NewMockFS(), FailMkdir: true}
	gen := New(manager, fsys)

	opts := Options{TemplateName: "basic", DestDir: "my-api", Values: &variables.Values{}}
	err := gen.Generate(opts)
	if err == nil || err.Error() != "mkdir failed" {
		t.Fatalf("expected mkdir failed, got %v", err)
	}
}

func TestGenerator_Generate_WriteError(t *testing.T) {
	mockTmplFS := fstest.MapFS{
		"templates/basic/README.md": &fstest.MapFile{Data: []byte("Hello")},
	}
	manager := templates.NewManager(mockTmplFS)
	
	fsys := &FailingFS{MockFS: *NewMockFS(), FailWrite: true}
	gen := New(manager, fsys)

	opts := Options{TemplateName: "basic", DestDir: "my-api", Values: &variables.Values{}}
	err := gen.Generate(opts)
	if err == nil || err.Error() != "write failed" {
		t.Fatalf("expected write failed, got %v", err)
	}
}
