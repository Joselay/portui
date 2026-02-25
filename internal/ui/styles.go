package ui

import "github.com/charmbracelet/lipgloss"

// Dracula-inspired palette
var (
	accent    = lipgloss.Color("#FF6AC1")
	subtle    = lipgloss.Color("#555555")
	highlight = lipgloss.Color("#7B61FF")
	green     = lipgloss.Color("#50FA7B")
	red       = lipgloss.Color("#FF5555")
	yellow    = lipgloss.Color("#F1FA8C")
	white     = lipgloss.Color("#F8F8F2")
	dimWhite  = lipgloss.Color("#BBBBBB")
	bg        = lipgloss.Color("#282A36")
	panelBg   = lipgloss.Color("#21222C")
)

// Panel border styles
func activeBorder() lipgloss.Border {
	return lipgloss.RoundedBorder()
}

func inactiveBorder() lipgloss.Border {
	return lipgloss.RoundedBorder()
}

func activePanel(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(activeBorder()).
		BorderForeground(accent).
		Width(width).
		Height(height)
}

func inactivePanel(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(inactiveBorder()).
		BorderForeground(subtle).
		Width(width).
		Height(height)
}

// Row styles
var (
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

	cursorStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)
)

// Detail panel styles
var (
	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(highlight).
				Width(12)

	detailValueStyle = lipgloss.NewStyle().
				Foreground(white)

	detailValueGreen = lipgloss.NewStyle().
				Bold(true).
				Foreground(green)

	detailValueYellow = lipgloss.NewStyle().
				Foreground(yellow)

	detailValueDim = lipgloss.NewStyle().
			Foreground(dimWhite)
)

// Status styles
var (
	confirmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(red)

	searchStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(green)

	errorStyle = lipgloss.NewStyle().
			Foreground(red)

	helpBarStyle = lipgloss.NewStyle().
			Foreground(subtle)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent)
)
