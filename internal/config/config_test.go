package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_NoFile(t *testing.T) {
	// Change home dir so it doesn't find a real config
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome) // Windows
	t.Setenv("HOME", tempHome)        // Unix

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DefaultAuthor != "" {
		t.Errorf("expected empty config, got %+v", cfg)
	}
}

func TestLoad_WithFile(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	yamlContent := `
default_author: Alice
default_license: MIT
default_module_prefix: github.com/alice
`
	configPath := filepath.Join(tempHome, ".goscaffold.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DefaultAuthor != "Alice" {
		t.Errorf("expected Alice, got %s", cfg.DefaultAuthor)
	}
	if cfg.DefaultModulePrefix != "github.com/alice" {
		t.Errorf("expected github.com/alice, got %s", cfg.DefaultModulePrefix)
	}
}

func TestLoad_InvalidYaml(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	yamlContent := `
default_author: [invalid yaml
`
	configPath := filepath.Join(tempHome, ".goscaffold.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatalf("expected error for invalid yaml, got nil")
	}
}

func TestLoad_ReadFileError(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	// Make .goscaffold.yaml a directory to force ReadFile error (not ErrNotExist)
	configPath := filepath.Join(tempHome, ".goscaffold.yaml")
	if err := os.Mkdir(configPath, 0755); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatalf("expected error for unreadable file, got nil")
	}
}
