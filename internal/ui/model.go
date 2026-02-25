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

type panel int

const (
	panelList panel = iota
	panelDetail
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
	focus     panel
	search    string
	width     int
	height    int
	keys      keyMap
	help      help.Model
	status    string
	err       error
	forceKill bool
}

// New creates a new Model.
func New() Model {
	h := help.New()
	h.ShowAll = false
	return Model{
		keys:  newKeyMap(),
		help:  h,
		focus: panelList,
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

	case key.Matches(msg, m.keys.Tab):
		if m.focus == panelList {
			m.focus = panelDetail
		} else {
			m.focus = panelList
		}

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
			m.forceKill = false
			m.status = fmt.Sprintf("Kill %s (PID %d) on port %d? (y/n)", p.Command, p.PID, p.Port)
		}

	case key.Matches(msg, m.keys.ForceKill):
		if len(m.filtered) > 0 {
			p := m.filtered[m.cursor]
			m.state = stateConfirmKill
			m.forceKill = true
			m.status = fmt.Sprintf("Force kill %s (PID %d) on port %d? (y/n)", p.Command, p.PID, p.Port)
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
			text := fmt.Sprintf("%s %d %d %s %s %s", p.Command, p.Port, p.PID, p.Protocol, p.User, p.State)
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
	var err error
	if m.forceKill {
		err = process.ForceKill(p.PID)
	} else {
		err = process.Kill(p.PID)
	}

	m.state = stateNormal
	if err != nil {
		m.status = fmt.Sprintf("Failed to kill PID %d: %v", p.PID, err)
		return m, nil
	}

	sig := "SIGTERM"
	if m.forceKill {
		sig = "SIGKILL"
	}
	m.status = fmt.Sprintf("Sent %s to %s (PID %d)", sig, p.Command, p.PID)
	return m, refresh
}

// View renders the entire UI.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Calculate layout dimensions
	totalWidth := m.width
	helpHeight := 1
	if m.help.ShowAll {
		helpHeight = 5
	}
	statusPanelHeight := 1
	// 2 for top/bottom border of main panels, 2 for status panel borders, helpHeight, 1 for spacing
	topPanelHeight := m.height - statusPanelHeight - 2 - helpHeight - 3

	if topPanelHeight < 5 {
		topPanelHeight = 5
	}

	// Inner height = height minus 2 for borders
	innerTopHeight := topPanelHeight
	// Left panel takes ~50%, right panel takes the rest
	// Subtract 4 for the two panels' left+right borders (2 borders × 2 panels)
	leftInnerWidth := (totalWidth - 4) / 2
	rightInnerWidth := totalWidth - leftInnerWidth - 4

	if leftInnerWidth < 20 {
		leftInnerWidth = 20
	}
	if rightInnerWidth < 20 {
		rightInnerWidth = 20
	}

	// Render the panels
	listContent := m.renderListPanel(leftInnerWidth, innerTopHeight)
	detailContent := m.renderDetailPanel(rightInnerWidth, innerTopHeight)

	// Build list panel with border
	listTitle := fmt.Sprintf(" Processes (%d) ", len(m.filtered))
	if m.search != "" && m.state != stateSearch {
		listTitle = fmt.Sprintf(" Processes (%d) [filter: %s] ", len(m.filtered), m.search)
	}
	if m.state == stateSearch {
		listTitle = fmt.Sprintf(" / %s█ ", m.search)
	}

	listFooter := ""
	if len(m.filtered) > 0 {
		listFooter = fmt.Sprintf(" %d of %d ", m.cursor+1, len(m.filtered))
	}

	var leftPanel string
	if m.focus == panelList {
		leftPanel = renderBorderedPanel(listContent, listTitle, listFooter, leftInnerWidth, innerTopHeight, true)
	} else {
		leftPanel = renderBorderedPanel(listContent, listTitle, listFooter, leftInnerWidth, innerTopHeight, false)
	}

	detailTitle := " Details "
	var rightPanel string
	if m.focus == panelDetail {
		rightPanel = renderBorderedPanel(detailContent, detailTitle, "", rightInnerWidth, innerTopHeight, true)
	} else {
		rightPanel = renderBorderedPanel(detailContent, detailTitle, "", rightInnerWidth, innerTopHeight, false)
	}

	// Join top panels side by side
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	b.WriteString(topRow)
	b.WriteString("\n")

	// Status panel (inner width = total top row width - 2 for borders)
	statusInnerWidth := leftInnerWidth + rightInnerWidth + 2
	statusContent := m.renderStatusContent(statusInnerWidth)
	statusPanel := renderBorderedPanel(statusContent, " Status ", "", statusInnerWidth, statusPanelHeight, false)
	b.WriteString(statusPanel)
	b.WriteString("\n")

	// Help bar (no border, just text)
	b.WriteString(" ")
	b.WriteString(m.help.View(m.keys))

	return b.String()
}

func (m Model) renderListPanel(width, height int) string {
	var b strings.Builder

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf(" Error: %v", m.err)))
		content := b.String()
		return padToHeight(content, height)
	}

	// Column widths: cursor(1) + port + space + cmd + space + proto + space + state
	// Total overhead = 1 (cursor) + 3 (spaces) = 4
	portW := 7
	protoW := 5
	stateW := 8
	cmdW := width - portW - protoW - stateW - 4
	if cmdW < 8 {
		cmdW = 8
	}
	if cmdW > 24 {
		cmdW = 24
	}

	// Header — pad to exact width
	header := fmt.Sprintf(" %-*s %-*s %-*s %-*s",
		portW, "PORT", cmdW, "COMMAND", protoW, "PROTO", stateW, "STATE")
	b.WriteString(padRight(helpBarStyle.Render(truncateLine(header, width)), width))
	b.WriteString("\n")
	b.WriteString(helpBarStyle.Render(strings.Repeat("─", width)))
	b.WriteString("\n")

	// Visible rows
	visibleRows := height - 2 // header + separator
	if visibleRows < 1 {
		visibleRows = 1
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

		cursor := " "
		if i == m.cursor {
			cursor = cursorStyle.Render("▸")
		}

		port := portStyle.Render(fmt.Sprintf("%-*d", portW, p.Port))
		cmd := commandStyle.Render(fmt.Sprintf("%-*s", cmdW, truncate(p.Command, cmdW)))
		proto := protocolStyle.Render(fmt.Sprintf("%-*s", protoW, p.Protocol))
		st := normalRowStyle.Render(fmt.Sprintf("%-*s", stateW, truncate(p.State, stateW)))

		row := fmt.Sprintf("%s%s %s %s %s", cursor, port, cmd, proto, st)

		// Always pad to exact panel width before styling
		row = padRight(row, width)

		if i == m.cursor {
			row = selectedStyle.Render(row)
		}

		b.WriteString(row)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	// Fill remaining rows
	renderedLines := end - start
	linesUsed := 2 + renderedLines // header + sep + data rows
	for linesUsed < height {
		b.WriteString("\n")
		linesUsed++
	}

	return b.String()
}

func (m Model) renderDetailPanel(width, height int) string {
	var b strings.Builder

	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		b.WriteString(helpBarStyle.Render(" No process selected"))
		return padToHeight(b.String(), height)
	}

	p := m.filtered[m.cursor]

	rows := []struct {
		label string
		value string
		style lipgloss.Style
	}{
		{"Command", p.Command, detailValueStyle},
		{"PID", fmt.Sprintf("%d", p.PID), detailValueYellow},
		{"Port", fmt.Sprintf("%d", p.Port), detailValueGreen},
		{"Protocol", p.Protocol, detailValueDim},
		{"User", p.User, detailValueStyle},
		{"State", p.State, detailValueDim},
	}

	for i, r := range rows {
		label := detailLabelStyle.Render(r.label)
		value := r.style.Render(r.value)
		b.WriteString(fmt.Sprintf(" %s %s", label, value))
		if i < len(rows)-1 {
			b.WriteString("\n")
		}
	}

	return padToHeight(b.String(), height)
}

func (m Model) renderStatusContent(width int) string {
	if m.status != "" {
		if m.state == stateConfirmKill {
			return " " + confirmStyle.Render(m.status)
		}
		return " " + statusStyle.Render(m.status)
	}
	if m.err != nil {
		return " " + errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}
	return " " + helpBarStyle.Render("Ready")
}

// renderBorderedPanel draws a panel with rounded borders, title, and optional footer.
func renderBorderedPanel(content, title, footer string, width, height int, active bool) string {
	borderColor := subtle
	titleColor := subtle
	if active {
		borderColor = accent
		titleColor = accent
	}

	// Border characters (rounded)
	tl, tr, bl, br := "╭", "╮", "╰", "╯"
	h, v := "─", "│"

	bc := lipgloss.NewStyle().Foreground(borderColor)
	tc := lipgloss.NewStyle().Foreground(titleColor).Bold(true)

	// Top border with title
	titleRendered := tc.Render(title)
	titleLen := lipgloss.Width(title)
	topFill := width - titleLen
	if topFill < 0 {
		topFill = 0
	}
	topLine := bc.Render(tl) + titleRendered + bc.Render(strings.Repeat(h, topFill)+tr)

	// Bottom border with optional footer
	var bottomLine string
	if footer != "" {
		footerRendered := tc.Render(footer)
		footerLen := lipgloss.Width(footer)
		bottomFill := width - footerLen
		if bottomFill < 0 {
			bottomFill = 0
		}
		bottomLine = bc.Render(bl+strings.Repeat(h, bottomFill)) + footerRendered + bc.Render(br)
	} else {
		bottomLine = bc.Render(bl + strings.Repeat(h, width) + br)
	}

	// Split content into lines and pad/truncate each to width
	contentLines := strings.Split(content, "\n")
	var middle strings.Builder
	for i := 0; i < height; i++ {
		line := ""
		if i < len(contentLines) {
			line = contentLines[i]
		}
		// Pad the visual width
		lineWidth := lipgloss.Width(line)
		padding := width - lineWidth
		if padding < 0 {
			padding = 0
		}
		middle.WriteString(bc.Render(v) + line + strings.Repeat(" ", padding) + bc.Render(v))
		if i < height-1 {
			middle.WriteString("\n")
		}
	}

	return topLine + "\n" + middle.String() + "\n" + bottomLine
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "…"
}

func truncateLine(s string, maxLen int) string {
	if lipgloss.Width(s) <= maxLen {
		return s
	}
	// Simple byte truncation as fallback
	if len(s) > maxLen {
		return s[:maxLen]
	}
	return s
}

func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func padToHeight(content string, height int) string {
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines[:height], "\n")
}
