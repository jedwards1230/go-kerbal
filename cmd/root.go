package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/core/tui"
)

func Execute() {
	config.LoadConfig()
	cfg := config.GetConfig()

	// If logging is enabled, logs will be output to debug.log.
	if cfg.Settings.EnableLogging {
		// clear debug file
		if err := os.Truncate("debug.log", 0); err != nil {
			log.Printf("Failed to clear debug.log: %v", err)
		}

		f, err := tea.LogToFile("debug.log", "[debug]")
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err = f.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	log.Println("Starting program")

	if cfg.Settings.KerbalDir == "" {
		log.Printf("ERROR: No Kerbal Directory found!")
	} else {
		log.Printf("Found Kerbal dir: " + cfg.Settings.KerbalDir)
	}

	m := tui.InitialModel()
	var opts []tea.ProgramOption

	// Always append alt screen program option.
	opts = append(opts, tea.WithAltScreen())

	// If mousewheel is enabled, append it to the program options.
	if cfg.Settings.EnableMouseWheel {
		opts = append(opts, tea.WithMouseAllMotion())
	}
	p := tea.NewProgram(m, opts...)
	log.Println("Program initialized")
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting server: %v", err)
		os.Exit(1)
	}
}
