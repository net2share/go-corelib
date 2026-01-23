package tui

import (
	"fmt"
	"strings"
)

// PrintProgress displays a progress bar.
func PrintProgress(downloaded, total int64) {
	if total <= 0 {
		return
	}

	percent := float64(downloaded) / float64(total) * 100
	barWidth := 40
	filled := int(float64(barWidth) * float64(downloaded) / float64(total))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Printf("\r[%s] %.1f%%", bar, percent)

	if downloaded >= total {
		fmt.Println()
	}
}
