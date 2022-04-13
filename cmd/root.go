package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/core/tui"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/spf13/viper"
)

func Execute() {
	config.LoadConfig()
	cfg := config.GetConfig()

	// If logging is enabled, logs will be output to debug.log.
	if cfg.Settings.EnableLogging {
		// clear debug file
		if _, err := os.Stat("debug.log"); err == nil {
			if err := os.Truncate("debug.log", 0); err != nil {
				log.Printf("Failed to clear debug.log: %v", err)
			}
		}

		f, err := tea.LogToFile("debug.log", "")
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

	if cfg.Settings.KerbalDir == "" || cfg.Settings.KerbalVer == "" {
		kerbalDir, err := dirfs.FindKspPath("")
		if err != nil {
			log.Printf("Error finding KSP folder: %v", err)
		}
		if kerbalDir == "" {
			viper.Set("settings.kerbal_dir", "")
			viper.Set("settings.kerbal_ver", "1.12.3")
			viper.WriteConfigAs(viper.ConfigFileUsed())
			log.Printf("***** TODO: FIND KSP DIR FOR LINUX *****")
		} else {
			kerbalVer := dirfs.FindKspVersion(kerbalDir)
			viper.Set("settings.kerbal_dir", kerbalDir)
			viper.Set("settings.kerbal_ver", kerbalVer.String())
			viper.WriteConfigAs(viper.ConfigFileUsed())
			log.Printf("Kerbal dir: " + kerbalDir + "/")
			log.Printf("Kerbal Version: %v", kerbalVer)
		}
	} else {
		log.Printf("Kerbal dir: " + cfg.Settings.KerbalDir + "/")
		log.Printf("Kerbal Version: %v", cfg.Settings.KerbalVer)
	}

	// TODO: Handle custom folderpath input for game directory if not found

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
