package tui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Prompt asks for user input with a label.
func Prompt(label string) string {
	reader := bufio.NewReader(os.Stdin)
	PromptColor.Printf("%s: ", label)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// PromptWithDefault asks for user input with a default value.
func PromptWithDefault(prompt, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)

	if defaultVal != "" {
		PromptColor.Printf("%s ", prompt)
		DefaultColor.Printf("[%s]", defaultVal)
		fmt.Print(": ")
	} else {
		PromptColor.Printf("%s: ", prompt)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}

// PromptInt asks for an integer input within a range.
func PromptInt(prompt string, defaultVal, min, max int) int {
	reader := bufio.NewReader(os.Stdin)

	for {
		PromptColor.Printf("%s ", prompt)
		DefaultColor.Printf("[%d]", defaultVal)
		fmt.Print(": ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return defaultVal
		}

		val, err := strconv.Atoi(input)
		if err != nil || val < min || val > max {
			PrintError(fmt.Sprintf("Please enter a number between %d and %d", min, max))
			continue
		}

		return val
	}
}

// PromptChoice asks the user to choose from a list of options.
func PromptChoice(prompt string, options []string, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)

	optStr := strings.Join(options, "/")
	PromptColor.Printf("%s (%s) ", prompt, optStr)
	DefaultColor.Printf("[%s]", defaultVal)
	fmt.Print(": ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}

	for _, opt := range options {
		if strings.EqualFold(input, opt) {
			return opt
		}
	}

	return defaultVal
}

// Confirm asks for a yes/no confirmation.
func Confirm(prompt string, defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)

	defaultStr := "y/N"
	if defaultYes {
		defaultStr = "Y/n"
	}

	PromptColor.Printf("%s ", prompt)
	DefaultColor.Printf("[%s]", defaultStr)
	fmt.Print(": ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes
	}

	return input == "y" || input == "yes"
}

// PromptPassword asks for password input (does not hide input - use for non-sensitive prompts).
// For actual password input, consider using golang.org/x/term.
func PromptPassword(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	PromptColor.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
