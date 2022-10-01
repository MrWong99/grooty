package main

import (
	"fmt"
	"os"

	"github.com/MrWong99/grooty/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		if f, err := tea.LogToFile("debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}

	if err := tea.NewProgram(tui.NewModel()).Start(); err != nil {
		fmt.Printf("Could not start program: %v\n", err)
		os.Exit(1)
	}
}
