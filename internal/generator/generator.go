package generator

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/arthurgray2k/goScaffold/internal/filesystem"
	"github.com/arthurgray2k/goScaffold/internal/templates"
	"github.com/arthurgray2k/goScaffold/internal/variables"
)

var ErrFileExists = errors.New("file already exists")

// Generator coordinates the reading of templates, applying variables, and writing to the filesystem.
type Generator struct {
	templates templates.Manager
	fsys      filesystem.FS
}

// New creates a new Generator.
func New(t templates.Manager, fsys filesystem.FS) *Generator {
	return &Generator{
		templates: t,
		fsys:      fsys,
	}
}

// Options defines the parameters for generating a project.
type Options struct {
	TemplateName string
	DestDir      string
	Force        bool
	Values       *variables.Values
}

// Generate executes the project scaffolding process.
func (g *Generator) Generate(opts Options) error {
	tmplFS, err := g.templates.Get(opts.TemplateName)
	if err != nil {
		return err
	}

	return fs.WalkDir(tmplFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}

		destPath := filepath.Join(opts.DestDir, filepath.FromSlash(path))
		destPath = strings.ReplaceAll(destPath, "__PROJECT_NAME__", opts.Values.ProjectName)
		
		// Strip .tmpl extension if present
		if filepath.Ext(destPath) == ".tmpl" {
			destPath = strings.TrimSuffix(destPath, ".tmpl")
		}

		if d.IsDir() {
			return g.fsys.MkdirAll(destPath, 0755)
		}

		// It's a file, check if we're allowed to write
		exists, err := g.fsys.Exists(destPath)
		if err != nil {
			return err
		}
		if exists && !opts.Force {
			return fmt.Errorf("%w: %s", ErrFileExists, destPath)
		}

		// Read file content
		f, err := tmplFS.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		// Replace variables
		replaced := opts.Values.Replace(string(content))

		// Ensure parent directory exists
		if err := g.fsys.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		return g.fsys.WriteFile(destPath, []byte(replaced), 0644)
	})
}
