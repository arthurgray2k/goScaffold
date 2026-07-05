package filesystem

import (
	"path/filepath"
	"testing"
)

func TestOSFS(t *testing.T) {
	fsys := OSFS{}
	tempDir := t.TempDir()

	testPath := filepath.Join(tempDir, "testdir", "testfile.txt")
	testDirPath := filepath.Join(tempDir, "testdir")

	// Test Exists (should be false)
	exists, err := fsys.Exists(testPath)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	if exists {
		t.Fatalf("expected file to not exist")
	}

	// Test MkdirAll
	err = fsys.MkdirAll(testDirPath, 0755)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	// Test WriteFile
	err = fsys.WriteFile(testPath, []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Test Exists (should be true)
	exists, err = fsys.Exists(testPath)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	if !exists {
		t.Fatalf("expected file to exist")
	}
}

func TestDryRunFS(t *testing.T) {
	fsys := DryRunFS{}
	
	// Just verify it doesn't panic
	err := fsys.MkdirAll("some/path", 0755)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	err = fsys.WriteFile("some/path/file.txt", []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}
