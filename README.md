# goScaffold

`goScaffold` is a robust, clean-architecture CLI application built in Go that generates the skeleton of new Go projects from reusable templates.

## Features
- **Project Generation**: Scaffold new projects easily (e.g., `goscaffold create gogrep`).
- **Interactive Git Lifecycle**: Automatically initialize Git, create your first commit, configure remotes, and push right out of the box!
- **Variable Replacement**: Supports dynamic variables like `{{PROJECT_NAME}}`, `{{MODULE_NAME}}`, `{{AUTHOR}}`, and `{{YEAR}}`.
- **Embedded Templates**: Templates are bundled directly into the executable via Go's `embed` package for single-binary portability.
- **Multiple Templates**: Choose from `basic`, `cli`, `api`, `grpc`, `worker`, or `library`.

## Installation
Ensure you have Go 1.26+ installed.

```bash
git clone https://github.com/arthurgray2k/goScaffold.git
cd goScaffold
go install ./cmd/goscaffold
```

For more details on commands and workflows, please see [USAGE.md](USAGE.md).