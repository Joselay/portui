# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Interactive TUI for viewing processes listening on network ports
- Search and filter processes by name, port, PID, or user
- Kill processes with confirmation prompt (SIGTERM)
- Vim-style navigation (j/k)
- Responsive layout adapting to terminal size
- Paginated process list

### Changed

- Redesigned UI with lazygit-style bordered panel layout
- Left panel shows process list, right panel shows selected process details
- Bottom status panel for confirmations and feedback
- Active panel highlighted with accent-colored borders, inactive panels dimmed
- Rounded border style (`╭╮╰╯`) across all panels
- Panel titles and scroll position indicator ("X of Y") embedded in borders
- Tab key to switch focus between panels
- `X` (shift+x) for force kill (SIGKILL) in addition to `x` for graceful kill (SIGTERM)
- State field now included in search filtering
- Dynamic column widths that adapt to terminal size
