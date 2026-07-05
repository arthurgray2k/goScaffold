package cli

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/arthurgray2k/goScaffold/internal/config"
	"github.com/arthurgray2k/goScaffold/internal/templates"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(templatesFS fs.FS) {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize template manager
	tmplManager := templates.NewManager(templatesFS)

	rootCmd := &cobra.Command{
		Use:   "goscaffold",
		Short: "A tool to generate Go projects from reusable templates",
		Long:  `goscaffold helps you quickly generate new Go projects with standard layouts and structure.`,
		// Silence errors because we handle them gracefully
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Add subcommands
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newListCmd(tmplManager))
	rootCmd.AddCommand(newInfoCmd(tmplManager))
	rootCmd.AddCommand(newCreateCmd(tmplManager))


	// Store config in context for subcommands to use
	ctx := context.WithValue(context.Background(), "config", cfg)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
