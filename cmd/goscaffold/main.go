package main

import (
	"goscaffold"
	"goscaffold/internal/cli"
)

func main() {
	cli.Execute(goscaffold.TemplatesFS)
}
