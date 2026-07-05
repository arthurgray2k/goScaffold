# goScaffold

[![Build Status](https://github.com/arthurgray2k/goScaffold/actions/workflows/ci.yml/badge.svg)](https://github.com/arthurgray2k/goScaffold/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/arthurgray2k/goScaffold.svg)](https://pkg.go.dev/github.com/arthurgray2k/goScaffold)
[![GitHub Release](https://img.shields.io/github/v/release/arthurgray2k/goScaffold)](https://github.com/arthurgray2k/goScaffold/releases)

`goScaffold` is a robust, clean-architecture CLI application and reusable Go library that generates the skeleton of new Go projects from reusable templates.

## Features
- **Project Generation**: Scaffold new projects easily.
- **Interactive Git Lifecycle**: Automatically initialize Git, create your first commit, configure remotes, and push right out of the box!
- **Variable Replacement**: Supports dynamic variables like `{{PROJECT_NAME}}`, `{{MODULE_NAME}}`, `{{AUTHOR}}`, and `{{YEAR}}`.
- **Embedded Templates**: Templates are bundled directly into the executable via Go's `embed` package.
- **Multiple Templates**: Choose from `basic`, `cli`, `api`, `grpc`, `worker`, or `library`.

## Installation
Ensure you have Go 1.26+ installed.

```bash
# Using go install
go install github.com/arthurgray2k/goScaffold/cmd/goscaffold@latest
```

Alternatively, download a pre-compiled binary from our [Releases](https://github.com/arthurgray2k/goScaffold/releases) page.

## CLI Usage Examples

The `create` command scaffolds a new project and handles Git initialization interactively.

```bash
goscaffold create my-api --template api
```

List available templates:
```bash
goscaffold list
```

For full CLI documentation, see [USAGE.md](USAGE.md).

## Library Usage Examples

You can also use `goScaffold` programmatically in your own Go applications.

```go
package main

import (
    "os"
    "path/filepath"
    
    "goscaffold/internal/filesystem"
    "goscaffold/internal/generator"
    "goscaffold/internal/templates"
    "goscaffold/internal/variables"
)

func main() {
    // 1. Get the embedded templates via the templates package
    // (Assuming you expose the embedded fs somehow, or use your own)
    
    // 2. Setup the generator engine
    fsys := filesystem.OSFS{}
    gen := generator.New(templatesManager, fsys)
    
    // 3. Generate a project
    err := gen.Generate(generator.Options{
        TemplateName: "basic",
        DestDir:      filepath.Join(".", "my-project"),
        Force:        false,
        Values: &variables.Values{
            ProjectName: "my-project",
            ModuleName:  "github.com/user/my-project",
            Year:        "2026",
            Author:      "Alice",
        },
    })
    
    if err != nil {
        panic(err)
    }
}
```

## Development Instructions

To set up the project locally for development:

1. Clone the repository:
```bash
git clone https://github.com/arthurgray2k/goScaffold.git
cd goScaffold
```

2. Run the test suite:
```bash
go test -v -cover ./...
```

3. Build the CLI binary:
```bash
go build -o goscaffold ./cmd/goscaffold
```

## Release Process

Releases are completely automated via GitHub Actions native builds.

To trigger a new release:
1. Ensure your working tree is clean and you're on `main`.
2. Create and push an annotated Git tag matching `v*`:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```
3. The `release.yml` GitHub workflow will automatically build binaries for Windows, Linux, and macOS, generate SHA256 checksums, and publish a GitHub Release with the artifacts.

## Uninstallation
To remove `goScaffold` from your system, simply delete the executable from your Go bin directory:

```bash
rm $(go env GOPATH)/bin/goscaffold
```