package cli

import (
	"fmt"
	"io/fs"

	"github.com/spf13/cobra"
	"goscaffold/internal/templates"
)

func newInfoCmd(manager templates.Manager) *cobra.Command {
	return &cobra.Command{
		Use:   "info [template]",
		Short: "Show information and files of a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			tmplFS, err := manager.Get(name)
			if err != nil {
				return fmt.Errorf("failed to get template %q: %w", name, err)
			}

			fmt.Printf("Template: %s\n", name)
			fmt.Println("Files:")
			
			err = fs.WalkDir(tmplFS, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if path == "." {
					return nil
				}
				
				prefix := "  ├── "
				if d.IsDir() {
					fmt.Printf("%s%s/\n", prefix, path)
				} else {
					fmt.Printf("%s%s\n", prefix, path)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to walk template files: %w", err)
			}

			return nil
		},
	}
}
