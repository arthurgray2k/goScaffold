package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"goscaffold/internal/config"
	"goscaffold/internal/filesystem"
	"goscaffold/internal/generator"
	"goscaffold/internal/templates"
	"goscaffold/internal/variables"
)

func newNewCmd(manager templates.Manager) *cobra.Command {
	var templateName string
	var dryRun bool
	var force bool

	cmd := &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new Go project",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract config from context
			var cfg *config.Config
			if c, ok := cmd.Context().Value("config").(*config.Config); ok {
				cfg = c
			} else {
				cfg = &config.Config{}
			}

			projectName := ""
			if len(args) > 0 {
				projectName = args[0]
			}

			reader := bufio.NewReader(os.Stdin)
			if projectName == "" {
				fmt.Print("Project name: ")
				input, _ := reader.ReadString('\n')
				projectName = strings.TrimSpace(input)
			}
			if projectName == "" {
				return fmt.Errorf("project name is required")
			}

			// Propose default module name based on config
			defaultModule := projectName
			if cfg.DefaultModulePrefix != "" {
				defaultModule = fmt.Sprintf("%s/%s", cfg.DefaultModulePrefix, projectName)
			}

			fmt.Printf("Module name [%s]: ", defaultModule)
			input, _ := reader.ReadString('\n')
			moduleName := strings.TrimSpace(input)
			if moduleName == "" {
				moduleName = defaultModule
			}

			// Validate module name
			if err := variables.ValidateModuleName(moduleName); err != nil {
				return fmt.Errorf("invalid module name '%s': %w", moduleName, err)
			}

			defaultAuthor := cfg.DefaultAuthor
			if defaultAuthor == "" {
				defaultAuthor = "Anonymous"
			}
			fmt.Printf("Author [%s]: ", defaultAuthor)
			input, _ = reader.ReadString('\n')
			author := strings.TrimSpace(input)
			if author == "" {
				author = defaultAuthor
			}

			var fsys filesystem.FS = filesystem.OSFS{}
			if dryRun {
				fsys = filesystem.DryRunFS{}
			}

			gen := generator.New(manager, fsys)

			opts := generator.Options{
				TemplateName: templateName,
				DestDir:      filepath.Join(".", projectName),
				Force:        force,
				Values: &variables.Values{
					ProjectName: projectName,
					ModuleName:  moduleName,
					Year:        fmt.Sprintf("%d", time.Now().Year()),
					Author:      author,
				},
			}

			fmt.Printf("Generating project '%s' using template '%s'...\n", projectName, templateName)
			if err := gen.Generate(opts); err != nil {
				return fmt.Errorf("failed to generate project: %w", err)
			}

			if dryRun {
				fmt.Println("Dry run completed. No files were written.")
			} else {
				fmt.Println("Project successfully generated!")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&templateName, "template", "t", "basic", "Template to use")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print files that would be created without writing")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing files")

	return cmd
}
