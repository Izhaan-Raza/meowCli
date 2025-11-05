package main

import (
	"fmt"
	"os"

	"meowCli/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error aagay bc: %v ", err)
		os.Exit(1)
	}
}
