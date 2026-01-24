package tui

import (
	"fmt"
	"strings"
)

// PrintProgress displays a progress bar for downloads or long operations.
// current and total represent bytes or units of work.
func PrintProgress(current, total int64) {
	if total <= 0 {
		return
	}

	percent := float64(current) / float64(total) * 100
	barWidth := 40
	filled := int(float64(barWidth) * float64(current) / float64(total))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Printf("\r[%s] %.1f%%", bar, percent)
	if current >= total {
		fmt.Println()
	}
}

// PrintProgressWithLabel displays a progress bar with a custom label.
func PrintProgressWithLabel(label string, current, total int64) {
	if total <= 0 {
		return
	}

	percent := float64(current) / float64(total) * 100
	barWidth := 30
	filled := int(float64(barWidth) * float64(current) / float64(total))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Printf("\r%s [%s] %.1f%%", label, bar, percent)
	if current >= total {
		fmt.Println()
	}
}
