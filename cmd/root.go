package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/tui"
)

func Execute() {
	//config.LoadConfig()
	//cfg := config.GetConfig()

	m := tui.InitialModel()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	//registry.BuildRegistry()
}
