package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
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
func renderFooter(_ int) string {
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

// calcBoxWidth calculates the appropriate box width based on terminal width.
func calcBoxWidth(termWidth int) int {
	if termWidth > 0 && termWidth < 80 {
		return termWidth - 10
	}
	if termWidth >= 80 {
		return 70
	}
	return 60
}

// createBoxStyle creates a standard box style with rounded border.
func createBoxStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Theme.Muted).
		Padding(1, 2).
		Width(width)
}

// centerWithFooter centers content on screen and adds footer if app info is set.
func centerWithFooter(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return content
	}

	footer := renderFooter(width)
	if footer != "" {
		contentHeight := height - 2
		centered := lipgloss.Place(width, contentHeight,
			lipgloss.Center, lipgloss.Center,
			content)
		footerCentered := lipgloss.PlaceHorizontal(width, lipgloss.Center, footer)
		return centered + "\n" + footerCentered
	}
	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		content)
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
		Foreground(Theme.Text)

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
	boxWidth := calcBoxWidth(m.width)
	box := createBoxStyle(boxWidth).Render(b.String())

	// If header is set, render it above the box
	if m.header != "" {
		headerStyle := lipgloss.NewStyle().
			Foreground(Theme.Muted)
		headerText := headerStyle.Render(m.header)
		box = headerText + "\n\n" + box
	}

	return centerWithFooter(box, m.width, m.height)
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
	boxWidth := calcBoxWidth(m.width)
	box := createBoxStyle(boxWidth).Render(b.String())

	return centerWithFooter(box, m.width, m.height)
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
	return max(available, 5)
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
		Foreground(Theme.Text)

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
		endIdx := min(m.scroll+visibleCount, len(m.config.Items))

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
	boxWidth := calcBoxWidth(m.width)
	box := createBoxStyle(boxWidth).Render(b.String())

	return centerWithFooter(box, m.width, m.height)
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

	return centerWithFooter(content, m.width, m.height)
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

// InfoSection represents a section in the info view.
type InfoSection struct {
	Title string   // Section title (optional)
	Rows  []InfoRow // Key-value rows
}

// InfoRow represents a single row in an info section.
type InfoRow struct {
	Key     string
	Value   string
	Columns []string // If set, renders as aligned columns (ignores Key/Value)
}

// InfoConfig configures a full-screen info display.
type InfoConfig struct {
	Title       string
	Description string
	Sections    []InfoSection
}

// infoModel is the bubbletea model for full-screen info display.
type infoModel struct {
	config   InfoConfig
	scroll   int
	width    int
	height   int
	quitting bool
}

func newInfoModel(cfg InfoConfig) infoModel {
	return infoModel{
		config: cfg,
	}
}

func (m infoModel) Init() tea.Cmd {
	return nil
}

func (m infoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m infoModel) totalLines() int {
	lines := 0
	for i, section := range m.config.Sections {
		if section.Title != "" {
			lines++ // Section title
		}
		lines += len(section.Rows)
		if i < len(m.config.Sections)-1 {
			lines++ // Spacing between sections
		}
	}
	return lines
}

func (m infoModel) maxScroll() int {
	visibleLines := m.visibleLineCount()
	total := m.totalLines()
	if total <= visibleLines {
		return 0
	}
	return total - visibleLines
}

func (m infoModel) visibleLineCount() int {
	available := m.height - 15
	return max(available, 5)
}

func (m infoModel) View() string {
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

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(Theme.Warning).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(Theme.Muted)

	valueStyle := lipgloss.NewStyle().
		Foreground(Theme.Primary)

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

	// Build all lines
	var allLines []string
	for i, section := range m.config.Sections {
		if section.Title != "" {
			allLines = append(allLines, sectionTitleStyle.Render(section.Title))
		}

		// Calculate column widths for rows with Columns
		var colWidths []int
		for _, row := range section.Rows {
			if len(row.Columns) > 0 {
				for j, col := range row.Columns {
					colLen := len(col)
					if j >= len(colWidths) {
						colWidths = append(colWidths, colLen)
					} else if colLen > colWidths[j] {
						colWidths[j] = colLen
					}
				}
			}
		}

		for _, row := range section.Rows {
			if len(row.Columns) > 0 {
				// Render aligned columns
				var parts []string
				for j, col := range row.Columns {
					width := 0
					if j < len(colWidths) {
						width = colWidths[j]
					}
					// Last column doesn't need padding
					if j == len(row.Columns)-1 {
						parts = append(parts, valueStyle.Render(col))
					} else {
						parts = append(parts, valueStyle.Render(fmt.Sprintf("%-*s", width, col)))
					}
				}
				allLines = append(allLines, strings.Join(parts, "  "))
			} else if row.Key != "" {
				line := keyStyle.Render(row.Key+": ") + valueStyle.Render(row.Value)
				allLines = append(allLines, line)
			} else {
				// Value-only row (for items like list entries)
				allLines = append(allLines, valueStyle.Render(row.Value))
			}
		}
		if i < len(m.config.Sections)-1 {
			allLines = append(allLines, "") // Spacing between sections
		}
	}

	// Apply scrolling
	visibleCount := m.visibleLineCount()
	startIdx := m.scroll
	endIdx := min(startIdx+visibleCount, len(allLines))

	// Show scroll indicator at top if scrolled
	if m.scroll > 0 {
		b.WriteString(scrollStyle.Render("↑ more above"))
		b.WriteString("\n")
	}

	for i := startIdx; i < endIdx; i++ {
		b.WriteString(allLines[i])
		b.WriteString("\n")
	}

	// Show scroll indicator at bottom if more content
	if endIdx < len(allLines) {
		b.WriteString(scrollStyle.Render("↓ more below"))
		b.WriteString("\n")
	}

	// Help
	if len(allLines) > visibleCount {
		b.WriteString(helpStyle.Render("\n↑/↓: scroll • enter/q/esc: close"))
	} else {
		b.WriteString(helpStyle.Render("\nenter/q/esc: close"))
	}

	// Create a box with the content left-aligned inside
	boxWidth := calcBoxWidth(m.width)
	box := createBoxStyle(boxWidth).Render(b.String())

	return centerWithFooter(box, m.width, m.height)
}

// ShowInfo displays a full-screen info view and waits for user to dismiss it.
func ShowInfo(cfg InfoConfig) error {
	m := newInfoModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()
	return err
}

// ProgressLineType indicates the type of progress line.
type ProgressLineType int

const (
	ProgressLineText ProgressLineType = iota
	ProgressLineInfo
	ProgressLineStatus
	ProgressLineSuccess
	ProgressLineWarning
	ProgressLineError
)

// ProgressLine represents a single line in the progress view.
type ProgressLine struct {
	Type    ProgressLineType
	Message string
}

// progressViewMsg is sent to update the progress view.
type progressViewMsg struct {
	line *ProgressLine
	done bool
}

// progressViewModel is the bubbletea model for progress view using viewport.
type progressViewModel struct {
	title      string
	lines      []ProgressLine
	done       bool
	width      int
	height     int
	vpWidth    int // viewport width for text wrapping
	quitting   bool
	msgCh      chan progressViewMsg
	viewport   viewport.Model
	ready      bool
	autoScroll bool
}

func newProgressViewModel(title string, msgCh chan progressViewMsg) progressViewModel {
	return progressViewModel{
		title:      title,
		msgCh:      msgCh,
		autoScroll: true,
	}
}

func (m progressViewModel) Init() tea.Cmd {
	return m.waitForMsg
}

func (m progressViewModel) waitForMsg() tea.Msg {
	msg := <-m.msgCh
	return msg
}

// renderLine renders a single progress line with appropriate styling.
func (m progressViewModel) renderLine(line ProgressLine) string {
	textStyle := lipgloss.NewStyle().Foreground(Theme.Text)
	infoStyle := lipgloss.NewStyle().Foreground(Theme.Info)
	statusStyle := lipgloss.NewStyle().Foreground(Theme.Success)
	successStyle := lipgloss.NewStyle().Foreground(Theme.Success).Bold(true)
	warningStyle := lipgloss.NewStyle().Foreground(Theme.Warning)
	errorStyle := lipgloss.NewStyle().Foreground(Theme.Error)

	// Icon prefix is 2 chars wide (icon + space)
	const iconWidth = 2

	// wrapWithIcon wraps the message and prepends icon to first line, indents rest
	wrapWithIcon := func(icon string, msg string, style lipgloss.Style) string {
		if m.vpWidth <= iconWidth {
			return style.Render(icon + " " + msg)
		}
		wrapped := wordwrap.String(msg, m.vpWidth-iconWidth)
		lines := strings.Split(wrapped, "\n")
		for i, l := range lines {
			if i == 0 {
				lines[i] = style.Render(icon + " " + l)
			} else {
				lines[i] = style.Render("  " + l) // indent continuation lines
			}
		}
		return strings.Join(lines, "\n")
	}

	switch line.Type {
	case ProgressLineInfo:
		return wrapWithIcon("ℹ", line.Message, infoStyle)
	case ProgressLineStatus:
		return wrapWithIcon("✓", line.Message, statusStyle)
	case ProgressLineSuccess:
		return wrapWithIcon("✓", line.Message, successStyle)
	case ProgressLineWarning:
		return wrapWithIcon("⚠", line.Message, warningStyle)
	case ProgressLineError:
		return wrapWithIcon("✗", line.Message, errorStyle)
	default:
		if line.Message != "" {
			if m.vpWidth > 0 {
				wrapped := wordwrap.String(line.Message, m.vpWidth)
				return textStyle.Render(wrapped)
			}
			return textStyle.Render(line.Message)
		}
		return ""
	}
}

// updateViewportContent updates the viewport with the current lines.
func (m *progressViewModel) updateViewportContent() {
	var lines []string
	for _, line := range m.lines {
		lines = append(lines, m.renderLine(line))
	}
	content := strings.Join(lines, "\n")
	m.viewport.SetContent(content)

	// Auto-scroll to bottom if enabled
	if m.autoScroll {
		m.viewport.GotoBottom()
	}
}

func (m progressViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc", "enter", " ":
			if m.done {
				m.quitting = true
				return m, tea.Quit
			}
		case "up", "k", "pgup":
			m.autoScroll = false
		case "down", "j", "pgdown":
			// Re-enable auto-scroll if at bottom after this move
			if m.viewport.AtBottom() {
				m.autoScroll = true
			}
		case "home", "g":
			m.autoScroll = false
		case "end", "G":
			m.autoScroll = true
		}

		// Pass key events to viewport
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate viewport dimensions
		boxWidth := calcBoxWidth(m.width)

		// Viewport height (leave space for title, help, borders, padding)
		vpHeight := 15
		if m.height > 0 {
			vpHeight = min(max(m.height-16, 5), 25)
		}

		// Viewport width (box - padding - border - scrollbar space)
		vpWidth := boxWidth - 10
		m.vpWidth = vpWidth

		if !m.ready {
			m.viewport = viewport.New(vpWidth, vpHeight)
			m.viewport.MouseWheelEnabled = true
			m.viewport.MouseWheelDelta = 3
			m.ready = true
			m.updateViewportContent()
		} else {
			m.viewport.Width = vpWidth
			m.viewport.Height = vpHeight
			m.updateViewportContent() // re-wrap text for new width
		}

	case progressViewMsg:
		if msg.done {
			m.done = true
			return m, nil
		}
		if msg.line != nil {
			m.lines = append(m.lines, *msg.line)
			m.updateViewportContent()
		}
		return m, m.waitForMsg
	}

	return m, tea.Batch(cmds...)
}

// renderScrollbar renders a scrollbar based on viewport scroll position.
func (m progressViewModel) renderScrollbar() string {
	if !m.ready || m.viewport.TotalLineCount() <= m.viewport.Height {
		return ""
	}

	trackStyle := lipgloss.NewStyle().Foreground(Theme.ScrollTrack)
	thumbStyle := lipgloss.NewStyle().Foreground(Theme.Primary)

	height := m.viewport.Height
	totalLines := m.viewport.TotalLineCount()
	visibleLines := m.viewport.Height

	// Calculate thumb size (clamped between 1 and height)
	thumbSize := min(max(height*visibleLines/totalLines, 1), height)

	// Calculate thumb position
	scrollPercent := m.viewport.ScrollPercent()
	thumbPos := int(float64(height-thumbSize) * scrollPercent)

	var sb strings.Builder
	for i := range height {
		if i >= thumbPos && i < thumbPos+thumbSize {
			sb.WriteString(thumbStyle.Render("█"))
		} else {
			sb.WriteString(trackStyle.Render("│"))
		}
		if i < height-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m progressViewModel) View() string {
	if m.quitting {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Foreground(Theme.Primary).Bold(true)
	helpStyle := lipgloss.NewStyle().Foreground(Theme.Muted)

	// Calculate box dimensions
	boxWidth := calcBoxWidth(m.width)

	var b strings.Builder

	// Title
	if m.title != "" {
		b.WriteString(titleStyle.Render(m.title))
		b.WriteString("\n\n")
	}

	// Viewport content with scrollbar
	if m.ready {
		hasScrollbar := m.viewport.TotalLineCount() > m.viewport.Height
		if hasScrollbar {
			scrollbar := m.renderScrollbar()
			content := lipgloss.JoinHorizontal(lipgloss.Top, m.viewport.View(), "  ", scrollbar)
			b.WriteString(content)
		} else {
			b.WriteString(m.viewport.View())
		}
	}

	// Help text
	b.WriteString("\n\n")
	if m.done {
		hasScrollbar := m.ready && m.viewport.TotalLineCount() > m.viewport.Height
		if hasScrollbar {
			b.WriteString(helpStyle.Render("↑/↓: scroll • enter/q/esc: close"))
		} else {
			b.WriteString(helpStyle.Render("enter/q/esc: close"))
		}
	} else {
		b.WriteString(helpStyle.Render("..."))
	}

	// Create a box with fixed dimensions
	box := createBoxStyle(boxWidth).Render(b.String())

	return centerWithFooter(box, m.width, m.height)
}

// ProgressView manages a real-time progress display.
type ProgressView struct {
	program *tea.Program
	msgCh   chan progressViewMsg
	doneCh  chan struct{}
}

// NewProgressView creates and starts a new progress view.
func NewProgressView(title string) *ProgressView {
	msgCh := make(chan progressViewMsg, 100)
	doneCh := make(chan struct{})
	m := newProgressViewModel(title, msgCh)
	p := tea.NewProgram(m, tea.WithAltScreen())

	pv := &ProgressView{
		program: p,
		msgCh:   msgCh,
		doneCh:  doneCh,
	}

	go func() {
		p.Run()
		close(doneCh)
	}()

	return pv
}

// AddLine adds a line to the progress view.
func (pv *ProgressView) AddLine(lineType ProgressLineType, message string) {
	pv.msgCh <- progressViewMsg{
		line: &ProgressLine{Type: lineType, Message: message},
	}
}

// AddText adds a plain text line.
func (pv *ProgressView) AddText(message string) {
	pv.AddLine(ProgressLineText, message)
}

// AddInfo adds an info line.
func (pv *ProgressView) AddInfo(message string) {
	pv.AddLine(ProgressLineInfo, message)
}

// AddStatus adds a status line (checkmark).
func (pv *ProgressView) AddStatus(message string) {
	pv.AddLine(ProgressLineStatus, message)
}

// AddSuccess adds a success line.
func (pv *ProgressView) AddSuccess(message string) {
	pv.AddLine(ProgressLineSuccess, message)
}

// AddWarning adds a warning line.
func (pv *ProgressView) AddWarning(message string) {
	pv.AddLine(ProgressLineWarning, message)
}

// AddError adds an error line.
func (pv *ProgressView) AddError(message string) {
	pv.AddLine(ProgressLineError, message)
}

// Done signals completion and waits for user to dismiss.
func (pv *ProgressView) Done() {
	pv.msgCh <- progressViewMsg{done: true}
	<-pv.doneCh
}

// Dismiss closes the progress view immediately without waiting for user input.
func (pv *ProgressView) Dismiss() {
	pv.program.Quit()
	<-pv.doneCh
}
