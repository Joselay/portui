package ui

import "github.com/charmbracelet/lipgloss"

var (
	accent    = lipgloss.Color("#FF6AC1")
	subtle    = lipgloss.Color("#888888")
	highlight = lipgloss.Color("#7B61FF")
	green     = lipgloss.Color("#50FA7B")
	red       = lipgloss.Color("#FF5555")
	yellow    = lipgloss.Color("#F1FA8C")
	white     = lipgloss.Color("#F8F8F2")
	dimWhite  = lipgloss.Color("#BBBBBB")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			PaddingRight(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(highlight)

	normalRowStyle = lipgloss.NewStyle().
			Foreground(dimWhite)

	portStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green)

	pidStyle = lipgloss.NewStyle().
			Foreground(yellow)

	commandStyle = lipgloss.NewStyle().
			Foreground(white)

	protocolStyle = lipgloss.NewStyle().
			Foreground(subtle)

	confirmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(red).
			MarginTop(1)

	searchStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	helpBarStyle = lipgloss.NewStyle().
			Foreground(subtle)

	statusStyle = lipgloss.NewStyle().
			Foreground(green).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			MarginTop(1)
)
