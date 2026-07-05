# goScaffold Usage Guide

## Commands

### `create` (The Main Workflow)
The `create` command scaffolds a new project and handles Git initialization interactively.

```bash
goscaffold create [project-name] --template [template-name]
```

**Example:**
```bash
goscaffold create gogrep
```

**Interactive Prompts:**
If you don't provide the required arguments, `goscaffold` will ask for them:
```text
Project Name     : gogrep
Module           [gogrep]: github.com/arthurgray2k/gogrep
Project Type     [basic]: cli
Author           [Anonymous]: Arthur Gray

Initialize Git?          [Y/n]: y
Create initial commit?   [Y/n]: y
Configure remote?        [Y/n]: y
Remote URL               : https://github.com/arthurgray2k/gogrep.git
Push to GitHub?          [Y/n]: y
```

### `list`
List all available templates embedded in the binary.
```bash
goscaffold list
```

### `info`
View detailed information about a specific template.
```bash
goscaffold info [template-name]
```

### `version`
Print the version of `goScaffold`.
```bash
goscaffold version
```

## Configuration
You can define default values in `~/.goscaffold.yaml` so you don't have to type them every time.

**Example `~/.goscaffold.yaml`:**
```yaml
default_author: "Arthur Gray"
default_module_prefix: "github.com/arthurgray2k"
```
