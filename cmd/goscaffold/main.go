package main

import (
	"github.com/arthurgray2k/goScaffold"
	"github.com/arthurgray2k/goScaffold/internal/cli"
)

func main() {
	cli.Execute(goscaffold.TemplatesFS)
}
