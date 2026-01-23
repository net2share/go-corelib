package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// PrintTitle displays a boxed title.
func PrintTitle(title string) {
	fmt.Println()
	TitleColor.Println("╔" + strings.Repeat("═", utf8.RuneCountInString(title)+2) + "╗")
	TitleColor.Println("║ " + title + " ║")
	TitleColor.Println("╚" + strings.Repeat("═", utf8.RuneCountInString(title)+2) + "╝")
	fmt.Println()
}

// PrintSuccess displays a success message with a checkmark.
func PrintSuccess(msg string) {
	SuccessColor.Println("✓ " + msg)
}

// PrintStatus is an alias for PrintSuccess.
func PrintStatus(msg string) {
	PrintSuccess(msg)
}

// PrintError displays an error message with an X mark.
func PrintError(msg string) {
	ErrorColor.Println("✗ " + msg)
}

// PrintWarning displays a warning message with a warning symbol.
func PrintWarning(msg string) {
	WarnColor.Println("⚠ " + msg)
}

// PrintInfo displays an info message with an info symbol.
func PrintInfo(msg string) {
	InfoColor.Println("ℹ " + msg)
}

// PrintStep displays a step indicator with progress.
func PrintStep(step int, total int, msg string) {
	fmt.Printf("[%d/%d] %s\n", step, total, msg)
}
