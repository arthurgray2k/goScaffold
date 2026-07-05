package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "{{PROJECT_NAME}}",
	Short: "{{PROJECT_NAME}} is a CLI application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to {{PROJECT_NAME}}!")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
