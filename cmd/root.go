package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/jedwards1230/go-kerbal/internal/tui"
	"github.com/spf13/viper"
)

func Execute() {
	config.LoadConfig("./")
	cfg := config.GetConfig()

	// If logging is enabled, logs will be output to debug.log.
	if cfg.Settings.EnableLogging {
		logPath := "./logs/debug.log"
		// Create log dir
		err := os.MkdirAll("./logs", os.ModePerm)
		if err != nil {
			log.Fatalf("error creating tmp dir: %v", err)
		}

		// clear previous logs
		if _, err := os.Stat(logPath); err == nil {
			if err := os.Truncate(logPath, 0); err != nil {
				log.Printf("Failed to clear %s: %v", err, logPath)
			}
		}

		f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)

		var LstdFlags = log.Lmsgprefix | log.Ltime | log.Lmicroseconds | log.Lshortfile

		log.SetFlags(LstdFlags)
	}
	log.Println("Starting program")

	if cfg.Settings.KerbalDir == "" || cfg.Settings.KerbalVer == "" {
		kerbalDir, err := dirfs.FindKspPath("")
		if err != nil {
			log.Printf("Error finding KSP folder: %v", err)
		}
		if kerbalDir == "" {
			viper.Set("settings.kerbal_dir", "")
			viper.Set("settings.kerbal_ver", "1.12.3")
			err := viper.WriteConfigAs(viper.ConfigFileUsed())
			if err != nil {
				log.Printf("Error saving log: %v", err)
			}
		} else {
			kerbalVer := dirfs.FindKspVersion(kerbalDir)
			viper.Set("settings.kerbal_dir", kerbalDir)
			viper.Set("settings.kerbal_ver", kerbalVer.String())
			err := viper.WriteConfigAs(viper.ConfigFileUsed())
			if err != nil {
				log.Printf("Error saving log: %v", err)
			}
			log.Printf("Kerbal dir: " + kerbalDir + "/")
			log.Printf("Kerbal Version: %v", kerbalVer)
		}
	} else {
		log.Printf("Kerbal dir: " + cfg.Settings.KerbalDir + "/")
		log.Printf("Kerbal Version: %v", cfg.Settings.KerbalVer)
	}

	m := tui.InitialModel()
	var opts []tea.ProgramOption

	// Always append alt screen program option.
	opts = append(opts, tea.WithAltScreen())

	// If mousewheel is enabled, append it to the program options.
	if cfg.Settings.EnableMouseWheel {
		opts = append(opts, tea.WithMouseCellMotion())
	}

	p := tea.NewProgram(m, opts...)
	log.Println("Program initialized")
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting server: %v", err)
		os.Exit(1)
	}
}
