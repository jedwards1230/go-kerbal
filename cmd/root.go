package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/core/tui"
)

func Execute() {
	config.LoadConfig()
	cfg := config.GetConfig()

	m := tui.InitialModel()
	var opts []tea.ProgramOption

	// Always append alt screen program option.
	opts = append(opts, tea.WithAltScreen())

	// If mousewheel is enabled, append it to the program options.
	if cfg.Settings.EnableMouseWheel {
		opts = append(opts, tea.WithMouseAllMotion())
	}
	p := tea.NewProgram(m, opts...)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
