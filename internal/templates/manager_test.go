package templates

import (
	"errors"
	"testing"
	"testing/fstest"
)

func TestFSManager_List(t *testing.T) {
	mockFS := fstest.MapFS{
		"templates/basic/README.md": &fstest.MapFile{Data: []byte("basic content")},
		"templates/cli/main.go":     &fstest.MapFile{Data: []byte("cli content")},
		"templates/notadir":         &fstest.MapFile{Data: []byte("just a file")},
	}

	mgr := NewManager(mockFS)
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 templates, got %d", len(list))
	}

	expected := map[string]bool{"basic": true, "cli": true}
	for _, item := range list {
		if !expected[item] {
			t.Errorf("unexpected template %s in list", item)
		}
	}
}

func TestFSManager_Get(t *testing.T) {
	mockFS := fstest.MapFS{
		"templates/basic/README.md": &fstest.MapFile{Data: []byte("basic content")},
		"templates/notadir":         &fstest.MapFile{Data: []byte("file")},
	}

	mgr := NewManager(mockFS)

	t.Run("existing template", func(t *testing.T) {
		sub, err := mgr.Get("basic")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
		f, err := sub.Open("README.md")
		if err != nil {
			t.Fatalf("failed to open README.md in sub filesystem: %v", err)
		}
		f.Close()
	})

	t.Run("non-existent template", func(t *testing.T) {
		_, err := mgr.Get("doesnotexist")
		if !errors.Is(err, ErrTemplateNotFound) {
			t.Errorf("expected ErrTemplateNotFound, got %v", err)
		}
	})

	t.Run("not a directory", func(t *testing.T) {
		_, err := mgr.Get("notadir")
		if !errors.Is(err, ErrTemplateNotFound) {
			t.Errorf("expected ErrTemplateNotFound, got %v", err)
		}
	})
}
