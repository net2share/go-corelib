package tui

import "fmt"

// PrintSuccess displays a success message with a checkmark.
func PrintSuccess(msg string) {
	fmt.Println(SuccessStyle.Render("✓ " + msg))
}

// PrintError displays an error message with an X mark.
func PrintError(msg string) {
	fmt.Println(ErrorStyle.Render("✗ " + msg))
}

// PrintWarning displays a warning message with a warning symbol.
func PrintWarning(msg string) {
	fmt.Println(WarnStyle.Render("⚠ " + msg))
}

// PrintInfo displays an info message with an info symbol.
func PrintInfo(msg string) {
	fmt.Println(InfoStyle.Render("ℹ " + msg))
}

// PrintStep displays a step indicator with progress (e.g., [1/5] Installing...).
func PrintStep(step, total int, msg string) {
	fmt.Printf("[%d/%d] %s\n", step, total, msg)
}

// PrintStatus displays a status message with a bullet point.
func PrintStatus(msg string) {
	fmt.Println(SuccessStyle.Render("• " + msg))
}
