# portui

A terminal UI for managing processes listening on network ports.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow)

## Features

- List all processes listening on TCP/UDP ports
- Search and filter by process name, port, PID, or user
- Kill processes with confirmation prompt
- Vim-style navigation

## Installation

```bash
go install github.com/smaetongmenglay/portui@latest
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
| `x` | Kill selected process |
| `r` | Refresh process list |
| `/` | Search |
| `esc` | Clear search |
| `?` | Toggle help |
| `q` / `ctrl+c` | Quit |

## License

[MIT](LICENSE)
