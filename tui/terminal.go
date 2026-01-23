package tui

import (
	"bufio"
	"fmt"
	"os"
)

// ClearScreen clears the terminal screen.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ClearLine clears the current line.
func ClearLine() {
	fmt.Print("\r\033[K")
}

// WaitForEnter waits for the user to press Enter.
func WaitForEnter() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress Enter to continue...")
	reader.ReadString('\n')
}
