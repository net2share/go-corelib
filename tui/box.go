package tui

import (
	"fmt"
	"strings"
)

// PrintBox displays content in a styled box with a title.
func PrintBox(title string, lines []string) {
	content := strings.Join(lines, "\n")
	box := BoxStyle.Render(content)
	fmt.Println()
	fmt.Println(TitleStyle.Render(title))
	fmt.Println(box)
}

// PrintBoxSimple displays content in a styled box without a title.
func PrintBoxSimple(lines []string) {
	content := strings.Join(lines, "\n")
	box := BoxStyle.Render(content)
	fmt.Println()
	fmt.Println(box)
}
