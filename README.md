# portui

A terminal UI for managing processes listening on network ports.

[![CI](https://github.com/Joselay/portui/actions/workflows/ci.yml/badge.svg)](https://github.com/Joselay/portui/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow)

## Features

- List all processes listening on TCP/UDP ports
- Search and filter by process name, port, PID, protocol, user, or state
- Panel-based UI with bordered panels (lazygit-style)
- Process list panel with detail panel side-by-side
- Kill processes with confirmation (SIGTERM or SIGKILL)
- Vim-style navigation
- Tab to switch focus between panels

## Installation

```bash
go install github.com/Joselay/portui@latest
```

Or build from source:

```bash
git clone https://github.com/Joselay/portui.git
cd portui
go build -o portui
```

## Usage

```bash
./portui
```

> Requires `lsof` (available by default on macOS and most Linux distributions). Some processes may require elevated permissions to view or kill.

## Keybindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `x` | Kill selected process (SIGTERM) |
| `X` | Force kill selected process (SIGKILL) |
| `r` | Refresh process list |
| `/` | Search |
| `esc` | Clear search |
| `tab` | Switch panel focus |
| `?` | Toggle help |
| `q` / `ctrl+c` | Quit |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
