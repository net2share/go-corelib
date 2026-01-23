package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ShowMenu displays a list of menu options.
func ShowMenu(options []MenuOption) {
	fmt.Println()
	for _, opt := range options {
		fmt.Printf("  %s. %s\n", opt.Key, opt.Label)
	}
	fmt.Println()
}

// PrintBox displays content in a bordered box with a title.
func PrintBox(title string, lines []string) {
	maxLen := utf8.RuneCountInString(title)
	for _, line := range lines {
		if utf8.RuneCountInString(line) > maxLen {
			maxLen = utf8.RuneCountInString(line)
		}
	}

	border := strings.Repeat("═", maxLen+2)

	fmt.Println()
	BoxColor.Printf("╔%s╗\n", border)
	BoxColor.Print("║ ")
	BoxTitleColor.Printf("%s", title)
	BoxColor.Printf("%s ║\n", strings.Repeat(" ", maxLen-utf8.RuneCountInString(title)))
	BoxColor.Printf("╠%s╣\n", border)

	for _, line := range lines {
		padding := maxLen - utf8.RuneCountInString(line)
		BoxColor.Print("║ ")
		// Color the values differently
		if strings.Contains(line, ":") && !strings.HasPrefix(line, " ") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				fmt.Print(parts[0] + ":")
				ValueColor.Print(parts[1])
			} else {
				fmt.Print(line)
			}
		} else {
			fmt.Print(line)
		}
		BoxColor.Printf("%s ║\n", strings.Repeat(" ", padding))
	}

	BoxColor.Printf("╚%s╝\n", border)
	fmt.Println()
}
