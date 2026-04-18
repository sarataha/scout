package ui

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mirageglobe/scout/internal/filesystem"
)

// View renders the entire application UI.
func (m Model) View() tea.View {
	if m.Width == 0 {
		return tea.NewView("Loading…")
	}

	// ── Colours & metrics ──────────────────────────────────────────────
	t := Themes[m.ThemeIdx]
	accentColor := lipgloss.Color(t.Accent)
	dimColor := lipgloss.Color(t.Dim)
	textColor := lipgloss.Color(t.Text)
	selectedBg := lipgloss.Color(t.SelectedBg)
	selectedFg := lipgloss.Color(t.SelectedFg)
	dirColor := accentColor

	// Reserve space for borders (2 chars each side) and internal padding
	usableWidth := m.Width - 5
	if usableWidth < 20 {
		usableWidth = 20
	}

	leftWidth := 40
	if leftWidth > usableWidth*2/5 {
		leftWidth = usableWidth * 2 / 5
	}
	rightWidth := usableWidth - leftWidth

	contentHeight := m.Height - 5
	if contentHeight < 1 {
		contentHeight = 1
	}

	// ── Left pane: file list ───────────────────────────────────────────
	var listLines []string
	cwdDisplay := m.Cwd
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(cwdDisplay, home) {
		cwdDisplay = "~" + cwdDisplay[len(home):]
	}

	headerStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	listLines = append(listLines, headerStyle.Render(filesystem.Truncate(cwdDisplay, leftWidth-4)))

	if m.Err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
		listLines = append(listLines, errStyle.Render("Error: "+filesystem.Truncate(m.Err.Error(), leftWidth-4)))
	}

	visibleRows := contentHeight - len(listLines)
	if visibleRows < 1 {
		visibleRows = 1
	}

	scrollOffset := 0
	if m.Cursor >= visibleRows {
		scrollOffset = m.Cursor - visibleRows + 1
	}

	normalItem := lipgloss.NewStyle().Foreground(textColor).Width(leftWidth - 4)
	selectedItem := lipgloss.NewStyle().
		Foreground(selectedFg).
		Background(selectedBg).
		Bold(true).
		Width(leftWidth - 4)
	dirStyle := lipgloss.NewStyle().Foreground(dirColor).Bold(true).Width(leftWidth - 4)

	for i := scrollOffset; i < len(m.Entries) && len(listLines) < contentHeight; i++ {
		e := m.Entries[i]
		name := e.Name

		var symbol string
		var symStyle lipgloss.Style

		if e.IsDir {
			symbol = "▸"
			symStyle = lipgloss.NewStyle().Foreground(dirColor)
			name = name + "/"
		} else {
			symbol = "•"
			symStyle = lipgloss.NewStyle().Foreground(textColor)
		}

		if status, ok := m.GitStatus[name]; ok {
			switch status {
			case "M", "A", "R", "C", "U", "MM", "AM":
				symbol = "●"
			case "?":
				symbol = "○"
			default:
				symbol = "◆"
			}
		}

		line := symStyle.Render(symbol) + " " + name
		line = filesystem.Truncate(line, leftWidth-4)

		if i == m.Cursor {
			listLines = append(listLines, selectedItem.Render(line))
		} else {
			if e.IsDir {
				listLines = append(listLines, dirStyle.Render(line))
			} else {
				listLines = append(listLines, normalItem.Render(line))
			}
		}
	}

	for len(listLines) < contentHeight {
		listLines = append(listLines, strings.Repeat(" ", leftWidth-4))
	}

	leftContent := strings.Join(listLines, "\n")

	// ── Right pane: preview ────────────────────────────────────────────
	previewLines := strings.Split(strings.TrimSuffix(m.Preview, "\n"), "\n")
	startIdx := m.PreviewScroll
	endIdx := startIdx + contentHeight
	if endIdx > len(previewLines) {
		endIdx = len(previewLines)
	}

	visiblePreview := make([]string, endIdx-startIdx)
	copy(visiblePreview, previewLines[startIdx:endIdx])
	for i, l := range visiblePreview {
		visiblePreview[i] = filesystem.Truncate(l, rightWidth-4)
	}
	for len(visiblePreview) < contentHeight {
		visiblePreview = append(visiblePreview, "")
	}
	rightContent := strings.Join(visiblePreview, "\n")

	// ── Pane styles ────────────────────────────────────────────────────
	leftBorderColor := dimColor
	rightBorderColor := dimColor
	if m.FocusRight {
		rightBorderColor = accentColor
	} else {
		leftBorderColor = accentColor
	}

	leftPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(contentHeight + 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(leftBorderColor).
		Padding(0, 1).
		Render(leftContent)

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(contentHeight + 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(rightBorderColor).
		Padding(0, 1).
		Render(rightContent)

	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	// ── Status bar ─────────────────────────────────────────────────────
	statusStyle := lipgloss.NewStyle().Foreground(dimColor)
	count := fmt.Sprintf(" %d items", len(m.Entries))
	pos := ""
	if len(m.Entries) > 0 {
		pos = fmt.Sprintf(" %d/%d", m.Cursor+1, len(m.Entries))
	}
	help := " q:quit  ←/→:focus  ↑/↓:nav/scroll  v/enter:vim  t:theme ?:help"

	statusBar := statusStyle.Render(
		filesystem.Truncate(count+pos+"  │"+help, m.Width),
	)

	var layout string
	if m.ShowHelp {
		helpScreen := m.RenderHelp()
		layout = lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, helpScreen)
	} else {
		header := m.RenderHeader()
		layout = lipgloss.JoinVertical(lipgloss.Left, header, panes, statusBar)
	}

	v := tea.NewView(layout)
	v.AltScreen = true
	return v
}
