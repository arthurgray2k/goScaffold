package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/arthurgray2k/goScaffold/internal/config"
	"github.com/arthurgray2k/goScaffold/internal/filesystem"
	"github.com/arthurgray2k/goScaffold/internal/generator"
	"github.com/arthurgray2k/goScaffold/internal/templates"
	"github.com/arthurgray2k/goScaffold/internal/variables"
)

func newCreateCmd(manager templates.Manager) *cobra.Command {
	var templateName string
	var dryRun bool
	var force bool

	cmd := &cobra.Command{
		Use:   "create [project-name]",
		Short: "Scaffold and initialize a new Go project",
		Long:  "create sets up a new Go project, optionally initializing git and configuring the repository.",
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
				fmt.Print("Project Name     : ")
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

			fmt.Printf("Module           [%s]: ", defaultModule)
			input, _ := reader.ReadString('\n')
			moduleName := strings.TrimSpace(input)
			if moduleName == "" {
				moduleName = defaultModule
			}

			// Validate module name
			if err := variables.ValidateModuleName(moduleName); err != nil {
				return fmt.Errorf("invalid module name '%s': %w", moduleName, err)
			}

			fmt.Printf("Project Type     [%s]: ", templateName)
			input, _ = reader.ReadString('\n')
			tType := strings.TrimSpace(input)
			if tType != "" {
				templateName = tType
			}

			defaultAuthor := cfg.DefaultAuthor
			if defaultAuthor == "" {
				defaultAuthor = "Anonymous"
			}
			fmt.Printf("Author           [%s]: ", defaultAuthor)
			input, _ = reader.ReadString('\n')
			author := strings.TrimSpace(input)
			if author == "" {
				author = defaultAuthor
			}

			// Git Prompts
			fmt.Println()
			fmt.Print("Initialize Git?          [Y/n]: ")
			initGit := askYesNo(reader, true)

			createCommit := false
			configRemote := false
			pushGit := false
			remoteURL := ""

			if initGit {
				fmt.Print("Create initial commit?   [Y/n]: ")
				createCommit = askYesNo(reader, true)

				fmt.Print("Configure remote?        [Y/n]: ")
				configRemote = askYesNo(reader, true)
				if configRemote {
					fmt.Print("Remote URL               : ")
					input, _ = reader.ReadString('\n')
					remoteURL = strings.TrimSpace(input)
				}

				if configRemote && remoteURL != "" {
					fmt.Print("Push to GitHub?          [Y/n]: ")
					pushGit = askYesNo(reader, true)
				}
			}

			fmt.Println()

			var fsys filesystem.FS = filesystem.OSFS{}
			if dryRun {
				fsys = filesystem.DryRunFS{}
			}

			gen := generator.New(manager, fsys)
			destDir := filepath.Join(".", projectName)

			opts := generator.Options{
				TemplateName: templateName,
				DestDir:      destDir,
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
				fmt.Println("Dry run completed. No files were written, skipping Git setup.")
				return nil
			}

			// Git Operations
			if initGit {
				fmt.Println("Initializing Git repository...")
				if err := runCmd(destDir, "git", "init"); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to initialize git: %v\n", err)
				} else {
					if createCommit {
						fmt.Println("Creating initial commit...")
						runCmd(destDir, "git", "add", ".")
						runCmd(destDir, "git", "commit", "-m", "Initial commit")
					}
					if configRemote && remoteURL != "" {
						fmt.Printf("Configuring remote '%s'...\n", remoteURL)
						runCmd(destDir, "git", "remote", "add", "origin", remoteURL)
						
						if pushGit {
							fmt.Println("Pushing to remote...")
							// Note: might fail if user has no auth configured, but that's expected
							if err := runCmd(destDir, "git", "push", "-u", "origin", "main"); err != nil {
								fmt.Fprintf(os.Stderr, "Warning: failed to push to remote: %v\n", err)
								fmt.Fprintf(os.Stderr, "You may need to push manually later.\n")
							}
						}
					}
				}
			}

			fmt.Println("\nProject successfully created!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&templateName, "template", "t", "basic", "Template to use")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print files that would be created without writing")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing files")

	return cmd
}

func askYesNo(reader *bufio.Reader, defaultYes bool) bool {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}

func runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
