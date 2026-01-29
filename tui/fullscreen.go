package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppInfo holds application metadata displayed in fullscreen components.
type AppInfo struct {
	Name      string
	Version   string
	BuildTime string
}

// globalAppInfo is the application info displayed in fullscreen footers.
var globalAppInfo *AppInfo

// SetAppInfo sets the global application info for fullscreen components.
// This should be called once at application startup.
// Consumer projects (like dnstm) can call this to control the footer display.
func SetAppInfo(name, version, buildTime string) {
	globalAppInfo = &AppInfo{
		Name:      name,
		Version:   version,
		BuildTime: buildTime,
	}
}

// ClearAppInfo clears the global application info.
func ClearAppInfo() {
	globalAppInfo = nil
}

// GetAppInfo returns the current global app info (or nil if not set).
func GetAppInfo() *AppInfo {
	return globalAppInfo
}

// renderFooter renders the footer with app info.
func renderFooter(width int) string {
	if globalAppInfo == nil {
		return ""
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	var footer string
	if globalAppInfo.BuildTime != "" && globalAppInfo.BuildTime != "unknown" {
		footer = fmt.Sprintf("%s %s (%s)", globalAppInfo.Name, globalAppInfo.Version, globalAppInfo.BuildTime)
	} else {
		footer = fmt.Sprintf("%s %s", globalAppInfo.Name, globalAppInfo.Version)
	}

	return footerStyle.Render(footer)
}

// MenuOption represents a single menu option.
type MenuOption struct {
	Label string
	Value string
}

// MenuConfig configures a full-screen menu.
type MenuConfig struct {
	Header      string // Text displayed above the box (outside border)
	Title       string
	Description string
	Options     []MenuOption
	Selected    int // Initial selected index
}

// menuModel is the bubbletea model for full-screen menu.
type menuModel struct {
	config   MenuConfig
	header   string
	cursor   int
	selected string
	width    int
	height   int
	quitting bool
}

func newMenuModel(cfg MenuConfig) menuModel {
	cursor := cfg.Selected
	if cursor < 0 || cursor >= len(cfg.Options) {
		cursor = 0
	}
	return menuModel{
		config: cfg,
		header: cfg.Header,
		cursor: cursor,
	}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			m.selected = ""
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.config.Options) - 1 // Wrap to last
			}
		case "down", "j":
			if m.cursor < len(m.config.Options)-1 {
				m.cursor++
			} else {
				m.cursor = 0 // Wrap to first
			}
		case "enter", " ":
			m.selected = m.config.Options[m.cursor].Value
			return m, tea.Quit
		case "home":
			m.cursor = 0
		case "end":
			m.cursor = len(m.config.Options) - 1
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m menuModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(Theme.Primary).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	selectedStyle := lipgloss.NewStyle().
		Foreground(Theme.Primary).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	cursorStyle := lipgloss.NewStyle().
		Foreground(Theme.Primary)

	helpStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	// Title
	if m.config.Title != "" {
		b.WriteString(titleStyle.Render(m.config.Title))
		b.WriteString("\n\n")
	}

	// Description
	if m.config.Description != "" {
		b.WriteString(descStyle.Render(m.config.Description))
		b.WriteString("\n\n")
	}

	// Options
	for i, opt := range m.config.Options {
		cursor := "  "
		style := normalStyle
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
			style = selectedStyle
		}
		b.WriteString(cursor + style.Render(opt.Label) + "\n")
	}

	// Help
	b.WriteString(helpStyle.Render("\n↑/↓: navigate • enter: select • q/esc: back"))

	// Create a box with the content left-aligned inside
	boxWidth := 60
	if m.width > 0 && m.width < 80 {
		boxWidth = m.width - 10
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Muted).
		Padding(1, 2).
		Width(boxWidth)

	box := boxStyle.Render(b.String())

	// If header is set, render it above the box
	if m.header != "" {
		headerStyle := lipgloss.NewStyle().
			Foreground(Theme.Muted)
		headerText := headerStyle.Render(m.header)
		box = headerText + "\n\n" + box
	}

	// Center the box on screen with footer
	if m.width > 0 && m.height > 0 {
		footer := renderFooter(m.width)
		if footer != "" {
			// Reserve space for footer at bottom
			contentHeight := m.height - 2
			centered := lipgloss.Place(m.width, contentHeight,
				lipgloss.Center, lipgloss.Center,
				box)
			footerCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, footer)
			return centered + "\n" + footerCentered
		}
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			box)
	}

	return box
}

// RunMenu runs a full-screen menu and returns the selected value.
// Returns empty string if user cancels (esc/q).
func RunMenu(cfg MenuConfig) (string, error) {
	m := newMenuModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(menuModel)
	return result.selected, nil
}

// RunMenuSimple runs a menu with just labels (values = labels).
func RunMenuSimple(title string, options []string) (string, error) {
	opts := make([]MenuOption, len(options))
	for i, o := range options {
		opts[i] = MenuOption{Label: o, Value: o}
	}
	return RunMenu(MenuConfig{
		Title:   title,
		Options: opts,
	})
}

// ConfirmConfig configures a confirmation dialog.
type ConfirmConfig struct {
	Title       string
	Description string
	Affirmative string // Default: "Yes"
	Negative    string // Default: "No"
	Default     bool   // Default selection
}

// RunConfirm runs a full-screen confirmation dialog.
func RunConfirm(cfg ConfirmConfig) (bool, error) {
	if cfg.Affirmative == "" {
		cfg.Affirmative = "Yes"
	}
	if cfg.Negative == "" {
		cfg.Negative = "No"
	}

	selected := 0
	if !cfg.Default {
		selected = 1
	}

	result, err := RunMenu(MenuConfig{
		Title:       cfg.Title,
		Description: cfg.Description,
		Options: []MenuOption{
			{Label: cfg.Affirmative, Value: "yes"},
			{Label: cfg.Negative, Value: "no"},
		},
		Selected: selected,
	})
	if err != nil {
		return false, err
	}

	return result == "yes", nil
}

// InputConfig configures an input dialog.
type InputConfig struct {
	Title       string
	Description string
	Placeholder string
	Value       string // Initial value
	Password    bool   // Hide input
}

// inputModel is the bubbletea model for text input.
type inputModel struct {
	config   InputConfig
	value    string
	cursor   int
	width    int
	height   int
	quitting bool
	done     bool
}

func newInputModel(cfg InputConfig) inputModel {
	return inputModel{
		config: cfg,
		value:  cfg.Value,
		cursor: len(cfg.Value),
	}
}

func (m inputModel) Init() tea.Cmd {
	return nil
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "backspace":
			if m.cursor > 0 {
				m.value = m.value[:m.cursor-1] + m.value[m.cursor:]
				m.cursor--
			}
		case "delete":
			if m.cursor < len(m.value) {
				m.value = m.value[:m.cursor] + m.value[m.cursor+1:]
			}
		case "left":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right":
			if m.cursor < len(m.value) {
				m.cursor++
			}
		case "home", "ctrl+a":
			m.cursor = 0
		case "end", "ctrl+e":
			m.cursor = len(m.value)
		default:
			if len(msg.String()) == 1 {
				m.value = m.value[:m.cursor] + msg.String() + m.value[m.cursor:]
				m.cursor++
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m inputModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(Theme.Primary).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	inputFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Primary).
		Padding(0, 1).
		Width(50)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	helpStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	// Title
	if m.config.Title != "" {
		b.WriteString(titleStyle.Render(m.config.Title))
		b.WriteString("\n\n")
	}

	// Description
	if m.config.Description != "" {
		b.WriteString(descStyle.Render(m.config.Description))
		b.WriteString("\n\n")
	}

	// Input field
	displayValue := m.value
	if m.config.Password && len(m.value) > 0 {
		displayValue = strings.Repeat("•", len(m.value))
	}

	if displayValue == "" && m.config.Placeholder != "" {
		displayValue = placeholderStyle.Render(m.config.Placeholder)
	} else {
		// Show cursor
		if m.cursor < len(displayValue) {
			displayValue = displayValue[:m.cursor] + "█" + displayValue[m.cursor+1:]
		} else {
			displayValue += "█"
		}
	}

	b.WriteString(inputFieldStyle.Render(displayValue))

	// Help
	b.WriteString(helpStyle.Render("\n\nenter: confirm • esc: cancel"))

	// Create a box with the content left-aligned inside
	boxWidth := 60
	if m.width > 0 && m.width < 80 {
		boxWidth = m.width - 10
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Muted).
		Padding(1, 2).
		Width(boxWidth)

	box := boxStyle.Render(b.String())

	// Center the box on screen with footer
	if m.width > 0 && m.height > 0 {
		footer := renderFooter(m.width)
		if footer != "" {
			contentHeight := m.height - 2
			centered := lipgloss.Place(m.width, contentHeight,
				lipgloss.Center, lipgloss.Center,
				box)
			footerCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, footer)
			return centered + "\n" + footerCentered
		}
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			box)
	}

	return box
}

// RunInput runs a full-screen input dialog.
// Returns the entered value and whether it was confirmed (not cancelled).
func RunInput(cfg InputConfig) (string, bool, error) {
	m := newInputModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return "", false, err
	}

	result := finalModel.(inputModel)
	if result.quitting {
		return "", false, nil
	}
	return result.value, true, nil
}

// SelectConfig configures a selection menu (alias for MenuConfig).
type SelectConfig = MenuConfig

// SelectOption is an alias for MenuOption.
type SelectOption = MenuOption

// RunSelect runs a full-screen selection menu (alias for RunMenu).
func RunSelect(cfg SelectConfig) (string, error) {
	return RunMenu(cfg)
}

// ListConfig configures a full-screen list display.
type ListConfig struct {
	Title       string
	Description string
	Items       []string
	EmptyText   string // Text to show when list is empty
}

// listModel is the bubbletea model for full-screen list display.
type listModel struct {
	config   ListConfig
	scroll   int // Scroll offset for long lists
	width    int
	height   int
	quitting bool
}

func newListModel(cfg ListConfig) listModel {
	if cfg.EmptyText == "" {
		cfg.EmptyText = "No items to display."
	}
	return listModel{
		config: cfg,
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc", "enter", " ":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			maxScroll := m.maxScroll()
			if m.scroll < maxScroll {
				m.scroll++
			}
		case "home":
			m.scroll = 0
		case "end":
			m.scroll = m.maxScroll()
		case "pgup":
			m.scroll -= 10
			if m.scroll < 0 {
				m.scroll = 0
			}
		case "pgdown":
			m.scroll += 10
			maxScroll := m.maxScroll()
			if m.scroll > maxScroll {
				m.scroll = maxScroll
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m listModel) maxScroll() int {
	visibleItems := m.visibleItemCount()
	if len(m.config.Items) <= visibleItems {
		return 0
	}
	return len(m.config.Items) - visibleItems
}

func (m listModel) visibleItemCount() int {
	// Estimate how many items fit in the box
	// Box has padding, title, description, help text
	available := m.height - 15 // Reserve space for chrome
	if available < 5 {
		available = 5
	}
	return available
}

func (m listModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(Theme.Primary).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	emptyStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted).
		Italic(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	scrollStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	// Title
	if m.config.Title != "" {
		b.WriteString(titleStyle.Render(m.config.Title))
		b.WriteString("\n\n")
	}

	// Description
	if m.config.Description != "" {
		b.WriteString(descStyle.Render(m.config.Description))
		b.WriteString("\n\n")
	}

	// Items or empty message
	if len(m.config.Items) == 0 {
		b.WriteString(emptyStyle.Render(m.config.EmptyText))
		b.WriteString("\n")
	} else {
		visibleCount := m.visibleItemCount()
		endIdx := m.scroll + visibleCount
		if endIdx > len(m.config.Items) {
			endIdx = len(m.config.Items)
		}

		// Show scroll indicator at top if scrolled
		if m.scroll > 0 {
			b.WriteString(scrollStyle.Render("  ↑ more above"))
			b.WriteString("\n")
		}

		for i := m.scroll; i < endIdx; i++ {
			b.WriteString("  • ")
			b.WriteString(itemStyle.Render(m.config.Items[i]))
			b.WriteString("\n")
		}

		// Show scroll indicator at bottom if more items
		if endIdx < len(m.config.Items) {
			b.WriteString(scrollStyle.Render("  ↓ more below"))
			b.WriteString("\n")
		}
	}

	// Help
	if len(m.config.Items) > m.visibleItemCount() {
		b.WriteString(helpStyle.Render("\n↑/↓: scroll • enter/q/esc: close"))
	} else {
		b.WriteString(helpStyle.Render("\nenter/q/esc: close"))
	}

	// Create a box with the content left-aligned inside
	boxWidth := 60
	if m.width > 0 && m.width < 80 {
		boxWidth = m.width - 10
	} else if m.width >= 80 {
		boxWidth = 70
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Muted).
		Padding(1, 2).
		Width(boxWidth)

	box := boxStyle.Render(b.String())

	// Center the box on screen with footer
	if m.width > 0 && m.height > 0 {
		footer := renderFooter(m.width)
		if footer != "" {
			contentHeight := m.height - 2
			centered := lipgloss.Place(m.width, contentHeight,
				lipgloss.Center, lipgloss.Center,
				box)
			footerCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, footer)
			return centered + "\n" + footerCentered
		}
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			box)
	}

	return box
}

// ShowList displays a full-screen list and waits for user to dismiss it.
func ShowList(cfg ListConfig) error {
	m := newListModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	return err
}

// App represents a full-screen TUI application.
type App struct {
	Title     string
	Version   string
	BuildTime string
	Banner    string
}

// AppScreen represents a screen in the app.
type AppScreen interface {
	Run() (next string, err error)
}

// AppMessage is used for displaying messages.
type AppMessage struct {
	Type    string // "success", "error", "warning", "info"
	Message string
}

// ShowMessage displays a message and waits for any key.
func ShowMessage(msg AppMessage) error {
	style := lipgloss.NewStyle()
	switch msg.Type {
	case "success":
		style = style.Foreground(Theme.Success)
	case "error":
		style = style.Foreground(Theme.Error)
	case "warning":
		style = style.Foreground(Theme.Warning)
	default:
		style = style.Foreground(Theme.Info)
	}

	_, err := RunMenu(MenuConfig{
		Title:   style.Render(msg.Message),
		Options: []MenuOption{{Label: "OK", Value: "ok"}},
	})
	return err
}

// ProgressConfig configures a progress display.
type ProgressConfig struct {
	Title   string
	Message string
}

// progressModel displays progress in full-screen.
type progressModel struct {
	config ProgressConfig
	width  int
	height int
	done   bool
	doneCh chan struct{}
}

func (m progressModel) Init() tea.Cmd {
	return m.waitForDone
}

func (m progressModel) waitForDone() tea.Msg {
	<-m.doneCh
	return struct{}{}
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case struct{}:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m progressModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(Theme.Primary).
		Bold(true)

	msgStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	if m.config.Title != "" {
		b.WriteString(titleStyle.Render(m.config.Title))
		b.WriteString("\n\n")
	}

	b.WriteString(msgStyle.Render(m.config.Message))
	b.WriteString(" ")
	b.WriteString(spinner())

	content := b.String()

	if m.width > 0 && m.height > 0 {
		footer := renderFooter(m.width)
		if footer != "" {
			contentHeight := m.height - 2
			centered := lipgloss.Place(m.width, contentHeight,
				lipgloss.Center, lipgloss.Center,
				content)
			footerCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, footer)
			return centered + "\n" + footerCentered
		}
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content)
	}

	return content
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var spinnerIdx = 0

func spinner() string {
	s := spinnerFrames[spinnerIdx%len(spinnerFrames)]
	spinnerIdx++
	return s
}

// Progress represents a running progress indicator.
type Progress struct {
	program *tea.Program
	doneCh  chan struct{}
}

// StartProgress starts a full-screen progress indicator.
func StartProgress(cfg ProgressConfig) *Progress {
	doneCh := make(chan struct{})
	m := progressModel{
		config: cfg,
		doneCh: doneCh,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		p.Run()
	}()

	return &Progress{
		program: p,
		doneCh:  doneCh,
	}
}

// Done stops the progress indicator.
func (p *Progress) Done() {
	close(p.doneCh)
}

// Update updates the progress message.
func (p *Progress) Update(msg string) {
	// Note: This is a simplified implementation
	// For real updates, we'd need custom messages
	fmt.Print(msg)
}
