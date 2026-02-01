// Package tui provides terminal UI components using Lipgloss.
package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for terminal UI.
var Theme = struct {
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Success     lipgloss.Color
	Error       lipgloss.Color
	Warning     lipgloss.Color
	Info        lipgloss.Color
	Muted       lipgloss.Color
	Text        lipgloss.Color
	ScrollTrack lipgloss.Color
}{
	Primary:     lipgloss.Color("6"),   // Cyan
	Secondary:   lipgloss.Color("5"),   // Magenta
	Success:     lipgloss.Color("2"),   // Green
	Error:       lipgloss.Color("1"),   // Red
	Warning:     lipgloss.Color("3"),   // Yellow
	Info:        lipgloss.Color("4"),   // Blue
	Muted:       lipgloss.Color("8"),   // Gray
	Text:        lipgloss.Color("252"), // Light gray text
	ScrollTrack: lipgloss.Color("238"), // Dark gray scrollbar track
}
