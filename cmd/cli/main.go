package main

import (
	"github.com/sujalamati/ArachneDB/cli"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"fmt"
)

func main() {
	p := tea.NewProgram(cli.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
