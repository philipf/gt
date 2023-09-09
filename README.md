[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/philipf/gt)
![Build workflow](https://github.com/philipf/gt/actions/workflows/go.yml/badge.svg)

# gt - Go Time

> **Warning:** Work in progress, will most likely only work on my machine.

GT (Go Time) is a CLI to improve my productivity. The current feature set includes:
- `gt gtd`:  Creates a new GTD action and saves it my Obsidian Kanban board
- `gt settings`: Prints the configured settings
- `gt version`: Prints the version of the CLI


## Install

```bash
go install github.com/philipf/gt@latest
```

Running `gt` for the first time will create a config file in the user's home directory.

The config file is located at `~/.gt.yaml` and can be edited manually. (Future versions of the CLI will include a `gt config` command to edit the config file.)
