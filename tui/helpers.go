package tui

import (
	"bufio"
	"fmt"
	"os"
)

// WaitForEnter waits for the user to press Enter.
func WaitForEnter() {
	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// ClearLine clears the current line in the terminal.
func ClearLine() {
	fmt.Print("\r\033[K")
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
