package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

// RenderHeader generates the single-line top bar of the application.
func (m Model) RenderHeader() string {
	t := Themes[m.ThemeIdx]
	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Dim))

	left := accentStyle.Render(" scout") + dimStyle.Render(" v"+Version+" │ github.com/mirageglobe/scout")

	memMB := float64(m.Stats.Mem) / 1024 / 1024

	// Build stats string
	now := time.Now()
	statsStr := fmt.Sprintf("%s  CPU: %.1f%%  MEM: %.1fMB", now.Format("Mon 2006-01-02 15:04"), m.Stats.CPU, memMB)
	right := dimStyle.Render(statsStr + " ")

	// Calculate space between
	width := m.Width
	paddingCount := width - lipgloss.Width(left) - lipgloss.Width(right)
	if paddingCount < 0 {
		paddingCount = 0
	}

	return left + strings.Repeat(" ", paddingCount) + right
}
