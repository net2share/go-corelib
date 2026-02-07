package tui

import (
	"fmt"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	sessionMu sync.Mutex
	inSession bool
)

// BeginSession enters the alternate screen buffer once for the entire session.
// Individual programs run in inline mode to avoid flickering between transitions.
func BeginSession() {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	if inSession {
		return
	}
	inSession = true
	fmt.Fprint(os.Stdout, "\033[?1049h")   // Enter alt screen
	fmt.Fprint(os.Stdout, "\033[H\033[2J") // Clear + cursor home
}

// EndSession exits the alternate screen buffer at session end.
func EndSession() {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	if !inSession {
		return
	}
	inSession = false
	fmt.Fprint(os.Stdout, "\033[?25h")   // Show cursor
	fmt.Fprint(os.Stdout, "\033[?1049l") // Exit alt screen
}

// InSession reports whether the session is active.
func InSession() bool {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	return inSession
}

// newProgram creates a tea.Program. In session mode: inline (no alt screen).
// Outside session: original WithAltScreen behavior.
func newProgram(m tea.Model) *tea.Program {
	if InSession() {
		fmt.Fprint(os.Stdout, "\033[H\033[2J") // Clear between programs
		return tea.NewProgram(m)
	}
	return tea.NewProgram(m, tea.WithAltScreen())
}
