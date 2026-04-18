package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/mirageglobe/scout/internal/filesystem"
)

// RenderHeader generates the single-line top bar of the application.
func (m Model) RenderHeader() string {
	t := Themes[m.ThemeIdx]
	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Dim))

	left := accentStyle.Render(" scout ") + dimStyle.Render("│ github.com/mirageglobe/scout")

	memMB := float64(m.Stats.Mem) / 1024 / 1024
	dirSizeStr := filesystem.HumanSize(m.Stats.DirSize)

	// Build stats string
	statsStr := fmt.Sprintf("MEM: %.1fMB  DIR: %s", memMB, dirSizeStr)
	right := dimStyle.Render(statsStr + " ")

	// Calculate space between
	width := m.Width
	paddingCount := width - lipgloss.Width(left) - lipgloss.Width(right)
	if paddingCount < 0 {
		paddingCount = 0
	}

	return left + strings.Repeat(" ", paddingCount) + right
}
