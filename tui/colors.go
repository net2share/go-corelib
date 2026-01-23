package tui

import "github.com/fatih/color"

// Exported color variables that can be customized by applications.
var (
	TitleColor    = color.New(color.FgCyan, color.Bold)
	SuccessColor  = color.New(color.FgGreen)
	ErrorColor    = color.New(color.FgRed)
	WarnColor     = color.New(color.FgYellow)
	InfoColor     = color.New(color.FgBlue)
	PromptColor   = color.New(color.FgYellow)
	DefaultColor  = color.New(color.FgHiBlack)
	ValueColor    = color.New(color.FgCyan)
	BoxColor      = color.New(color.FgCyan)
	BoxTitleColor = color.New(color.FgGreen, color.Bold)
)
