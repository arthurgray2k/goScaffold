package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/arthurgray2k/goScaffold/internal/config"
	"github.com/arthurgray2k/goScaffold/internal/filesystem"
	"github.com/arthurgray2k/goScaffold/internal/generator"
	"github.com/arthurgray2k/goScaffold/internal/templates"
	"github.com/arthurgray2k/goScaffold/internal/variables"
	"github.com/spf13/cobra"
)

type createFlags struct {
	templateName string
	dryRun       bool
	force        bool
}

type createInput struct {
	projectName  string
	moduleName   string
	templateName string
	author       string
	git          gitInput
}

type gitInput struct {
	initGit      bool
	createCommit bool
	configRemote bool
	pushGit      bool
	remoteURL    string
}

func newCreateCmd(manager templates.Manager) *cobra.Command {
	flags := createFlags{
		templateName: "basic",
	}

	cmd := &cobra.Command{
		Use:   "create [project-name]",
		Short: "Scaffold and initialize a new Go project",
		Long:  "create sets up a new Go project, optionally initializing git and configuring the repository.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(cmd, args, manager, flags)
		},
	}

	cmd.Flags().StringVarP(&flags.templateName, "template", "t", flags.templateName, "Template to use")
	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Print files that would be created without writing")
	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Overwrite existing files")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string, manager templates.Manager, flags createFlags) error {
	cfg := configFromContext(cmd)
	reader := bufio.NewReader(os.Stdin)

	input, err := promptCreateInput(reader, args, cfg, flags.templateName)
	if err != nil {
		return err
	}

	fmt.Println()

	destDir, err := generateProject(manager, input, flags)
	if err != nil {
		return err
	}

	if flags.dryRun {
		fmt.Println("Dry run completed. No files were written, skipping Git setup.")
		return nil
	}

	setupGit(destDir, input.git)

	fmt.Println("\nProject successfully created!")
	return nil
}

func configFromContext(cmd *cobra.Command) *config.Config {
	if cfg, ok := cmd.Context().Value("config").(*config.Config); ok {
		return cfg
	}
	return &config.Config{}
}

func promptCreateInput(reader *bufio.Reader, args []string, cfg *config.Config, defaultTemplate string) (createInput, error) {
	projectName := ""
	if len(args) > 0 {
		projectName = args[0]
	}

	if projectName == "" {
		input, err := prompt(reader, "Project Name     : ")
		if err != nil {
			return createInput{}, err
		}
		projectName = input
	}
	if projectName == "" {
		return createInput{}, fmt.Errorf("project name is required")
	}

	moduleName, err := promptModuleName(reader, cfg, projectName)
	if err != nil {
		return createInput{}, err
	}

	templateName, err := promptWithDefault(reader, "Project Type", defaultTemplate)
	if err != nil {
		return createInput{}, err
	}

	author, err := promptAuthor(reader, cfg)
	if err != nil {
		return createInput{}, err
	}

	git, err := promptGitInput(reader)
	if err != nil {
		return createInput{}, err
	}

	return createInput{
		projectName:  projectName,
		moduleName:   moduleName,
		templateName: templateName,
		author:       author,
		git:          git,
	}, nil
}

func promptModuleName(reader *bufio.Reader, cfg *config.Config, projectName string) (string, error) {
	defaultModule := projectName
	if cfg.DefaultModulePrefix != "" {
		defaultModule = fmt.Sprintf("%s/%s", cfg.DefaultModulePrefix, projectName)
	}

	moduleName, err := promptWithDefault(reader, "Module", defaultModule)
	if err != nil {
		return "", err
	}

	if err := variables.ValidateModuleName(moduleName); err != nil {
		return "", fmt.Errorf("invalid module name '%s': %w", moduleName, err)
	}

	return moduleName, nil
}

func promptAuthor(reader *bufio.Reader, cfg *config.Config) (string, error) {
	defaultAuthor := cfg.DefaultAuthor
	if defaultAuthor == "" {
		defaultAuthor = "Anonymous"
	}

	return promptWithDefault(reader, "Author", defaultAuthor)
}

func promptGitInput(reader *bufio.Reader) (gitInput, error) {
	fmt.Println()
	fmt.Print("Initialize Git?          [Y/n]: ")
	initGit, err := askYesNo(reader, true)
	if err != nil {
		return gitInput{}, err
	}

	git := gitInput{initGit: initGit}
	if !git.initGit {
		return git, nil
	}

	fmt.Print("Create initial commit?   [Y/n]: ")
	git.createCommit, err = askYesNo(reader, true)
	if err != nil {
		return gitInput{}, err
	}

	fmt.Print("Configure remote?        [Y/n]: ")
	git.configRemote, err = askYesNo(reader, true)
	if err != nil {
		return gitInput{}, err
	}

	if git.configRemote {
		git.remoteURL, err = prompt(reader, "Remote URL               : ")
		if err != nil {
			return gitInput{}, err
		}
	}

	if git.configRemote && git.remoteURL != "" {
		fmt.Print("Push to GitHub?          [Y/n]: ")
		git.pushGit, err = askYesNo(reader, true)
		if err != nil {
			return gitInput{}, err
		}
	}

	return git, nil
}

func promptWithDefault(reader *bufio.Reader, label, defaultValue string) (string, error) {
	input, err := prompt(reader, fmt.Sprintf("%-17s[%s]: ", label, defaultValue))
	if err != nil {
		return "", err
	}
	if input == "" {
		return defaultValue, nil
	}
	return input, nil
}

func prompt(reader *bufio.Reader, message string) (string, error) {
	fmt.Print(message)
	input, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func generateProject(manager templates.Manager, input createInput, flags createFlags) (string, error) {
	var fsys filesystem.FS = filesystem.OSFS{}
	if flags.dryRun {
		fsys = filesystem.DryRunFS{}
	}

	destDir := filepath.Join(".", input.projectName)
	opts := generator.Options{
		TemplateName: input.templateName,
		DestDir:      destDir,
		Force:        flags.force,
		Values: &variables.Values{
			ProjectName: input.projectName,
			ModuleName:  input.moduleName,
			Year:        fmt.Sprintf("%d", time.Now().Year()),
			Author:      input.author,
		},
	}

	fmt.Printf("Generating project '%s' using template '%s'...\n", input.projectName, input.templateName)
	if err := generator.New(manager, fsys).Generate(opts); err != nil {
		return "", fmt.Errorf("failed to generate project: %w", err)
	}

	return destDir, nil
}

func setupGit(destDir string, git gitInput) {
	if !git.initGit {
		return
	}

	fmt.Println("Initializing Git repository...")
	if err := runCmd(destDir, "git", "init"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to initialize git: %v\n", err)
		return
	}

	if git.createCommit {
		fmt.Println("Creating initial commit...")
		runCmd(destDir, "git", "add", ".")
		runCmd(destDir, "git", "commit", "-m", "Initial commit")
	}

	if git.configRemote && git.remoteURL != "" {
		fmt.Printf("Configuring remote '%s'...\n", git.remoteURL)
		runCmd(destDir, "git", "remote", "add", "origin", git.remoteURL)
	}

	if git.configRemote && git.remoteURL != "" && git.pushGit {
		fmt.Println("Pushing to remote...")
		// This can fail if authentication is not configured, so report it without failing creation.
		if err := runCmd(destDir, "git", "push", "-u", "origin", "main"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to push to remote: %v\n", err)
			fmt.Fprintf(os.Stderr, "You may need to push manually later.\n")
		}
	}
}

func askYesNo(reader *bufio.Reader, defaultYes bool) (bool, error) {
	input, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultYes, nil
	}
	return input == "y" || input == "yes", nil
}

func runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
