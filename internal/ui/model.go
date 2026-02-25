package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Joselay/portui/internal/process"
)

type state int

const (
	stateNormal state = iota
	stateSearch
	stateConfirmKill
)

type refreshMsg struct {
	processes []process.Info
	err       error
}

type statusMsg string

// Model is the main application model.
type Model struct {
	processes []process.Info
	filtered  []process.Info
	cursor    int
	state     state
	search    string
	width     int
	height    int
	keys      keyMap
	help      help.Model
	status    string
	err       error
}

// New creates a new Model.
func New() Model {
	h := help.New()
	h.ShowAll = false
	return Model{
		keys: newKeyMap(),
		help: h,
	}
}

func (m Model) Init() tea.Cmd {
	return refresh
}

func refresh() tea.Msg {
	procs, err := process.List()
	return refreshMsg{processes: procs, err: err}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case refreshMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.processes = msg.processes
			m.err = nil
		}
		m.applyFilter()
		return m, nil

	case statusMsg:
		m.status = string(msg)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Search mode: capture typing
	if m.state == stateSearch {
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc:
			m.state = stateNormal
			if msg.Type == tea.KeyEsc {
				m.search = ""
				m.applyFilter()
			}
			return m, nil
		case tea.KeyBackspace:
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
				m.applyFilter()
			}
			return m, nil
		case tea.KeyRunes:
			m.search += string(msg.Runes)
			m.applyFilter()
			return m, nil
		}
		return m, nil
	}

	// Confirm kill mode
	if m.state == stateConfirmKill {
		switch {
		case key.Matches(msg, m.keys.Confirm):
			return m.killSelected()
		default:
			m.state = stateNormal
			m.status = ""
			return m, nil
		}
	}

	// Normal mode
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}

	case key.Matches(msg, m.keys.Kill):
		if len(m.filtered) > 0 {
			p := m.filtered[m.cursor]
			m.state = stateConfirmKill
			m.status = fmt.Sprintf("Kill %s (PID %d) on port %d? (y/n)", p.Command, p.PID, p.Port)
		}

	case key.Matches(msg, m.keys.Refresh):
		m.status = "Refreshing..."
		return m, refresh

	case key.Matches(msg, m.keys.Search):
		m.state = stateSearch
		m.search = ""

	case key.Matches(msg, m.keys.Clear):
		m.search = ""
		m.applyFilter()

	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
	}

	return m, nil
}

func (m *Model) applyFilter() {
	if m.search == "" {
		m.filtered = m.processes
	} else {
		query := strings.ToLower(m.search)
		m.filtered = nil
		for _, p := range m.processes {
			text := fmt.Sprintf("%s %d %d %s %s", p.Command, p.Port, p.PID, p.Protocol, p.User)
			if strings.Contains(strings.ToLower(text), query) {
				m.filtered = append(m.filtered, p)
			}
		}
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m Model) killSelected() (tea.Model, tea.Cmd) {
	if m.cursor >= len(m.filtered) {
		m.state = stateNormal
		return m, nil
	}

	p := m.filtered[m.cursor]
	err := process.Kill(p.PID)

	m.state = stateNormal
	if err != nil {
		m.status = fmt.Sprintf("Failed to kill PID %d: %v", p.PID, err)
		return m, nil
	}

	m.status = fmt.Sprintf("Sent SIGTERM to %s (PID %d)", p.Command, p.PID)
	return m, refresh
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("⚡ portui"))
	b.WriteString("  ")
	b.WriteString(helpBarStyle.Render(fmt.Sprintf("%d processes listening", len(m.filtered))))
	b.WriteString("\n\n")

	// Search bar
	if m.state == stateSearch {
		b.WriteString(searchStyle.Render("/ " + m.search + "█"))
		b.WriteString("\n\n")
	} else if m.search != "" {
		b.WriteString(helpBarStyle.Render(fmt.Sprintf("filter: %s", m.search)))
		b.WriteString("\n\n")
	}

	// Error
	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	// Table header
	header := fmt.Sprintf("  %-8s %-8s %-8s %-20s %-12s %s",
		"PORT", "PID", "PROTO", "COMMAND", "USER", "STATE")
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(helpBarStyle.Render(strings.Repeat("─", min(m.width, 80))))
	b.WriteString("\n")

	// Rows
	visibleRows := m.height - 10
	if m.help.ShowAll {
		visibleRows -= 4
	}
	if visibleRows < 3 {
		visibleRows = 3
	}

	start := 0
	if m.cursor >= visibleRows {
		start = m.cursor - visibleRows + 1
	}
	end := start + visibleRows
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		p := m.filtered[i]

		cursor := "  "
		if i == m.cursor {
			cursor = lipgloss.NewStyle().Foreground(accent).Render("▸ ")
		}

		port := portStyle.Render(fmt.Sprintf("%-8d", p.Port))
		pid := pidStyle.Render(fmt.Sprintf("%-8d", p.PID))
		proto := protocolStyle.Render(fmt.Sprintf("%-8s", p.Protocol))
		cmd := commandStyle.Render(fmt.Sprintf("%-20s", truncate(p.Command, 20)))
		user := normalRowStyle.Render(fmt.Sprintf("%-12s", p.User))
		state := normalRowStyle.Render(p.State)

		row := fmt.Sprintf("%s%s%s%s%s%s%s", cursor, port, pid, proto, cmd, user, state)

		if i == m.cursor {
			row = selectedStyle.Render(row)
		}

		b.WriteString(row)
		b.WriteString("\n")
	}

	// Status / confirm
	if m.status != "" {
		if m.state == stateConfirmKill {
			b.WriteString(confirmStyle.Render(m.status))
		} else {
			b.WriteString(statusStyle.Render(m.status))
		}
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(m.help.View(m.keys))

	return b.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
