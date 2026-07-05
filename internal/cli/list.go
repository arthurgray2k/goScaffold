package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"goscaffold/internal/templates"
)

func newListCmd(manager templates.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available project templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			tmplList, err := manager.List()
			if err != nil {
				return fmt.Errorf("failed to list templates: %w", err)
			}

			fmt.Println("Available templates:")
			for _, t := range tmplList {
				fmt.Printf("  - %s\n", t)
			}
			return nil
		},
	}
}
