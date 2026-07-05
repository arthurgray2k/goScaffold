package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of goscaffold",
		Long:  `All software has versions. This is goscaffold's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("goscaffold version 0.1.0")
		},
	}
}
