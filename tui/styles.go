package tui

import "github.com/charmbracelet/lipgloss"

// Pre-configured styles for common UI elements.
var (
	TitleStyle   = lipgloss.NewStyle().Foreground(Theme.Primary).Bold(true)
	SuccessStyle = lipgloss.NewStyle().Foreground(Theme.Success)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Theme.Error)
	WarnStyle    = lipgloss.NewStyle().Foreground(Theme.Warning)
	InfoStyle    = lipgloss.NewStyle().Foreground(Theme.Info)
	MutedStyle   = lipgloss.NewStyle().Foreground(Theme.Muted)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Theme.Primary).
			Padding(0, 1)

	// Content styles for box content
	KeyStyle    = lipgloss.NewStyle().Foreground(Theme.Muted)
	ValueStyle  = lipgloss.NewStyle().Foreground(Theme.Primary).Bold(true)
	HeaderStyle = lipgloss.NewStyle().Foreground(Theme.Warning).Bold(true)
	CodeStyle   = lipgloss.NewStyle().Foreground(Theme.Success)
)
