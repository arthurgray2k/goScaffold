package main

import (
	"os"
	"{{MODULE_NAME}}/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
