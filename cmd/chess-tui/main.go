package main

import (
	"fmt"
	"os"

	"github.com/gabrzb/chess-go/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to run TUI:", err)
		os.Exit(1)
	}
}
