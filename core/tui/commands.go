package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/jedwards1230/go-kerbal/registry/database"
	"github.com/spf13/viper"
)

type (
	UpdatedModListMsg   []database.Ckan
	InstalledModListMsg map[string]bool
	UpdateKspDirMsg     bool
	DownloadModMsg      bool
	ErrorMsg            error
)

// Request the mod list from the database
func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		log.Print("Checking available mods")
		b.registry.DB.UpdateDB(false)
		updatedModList := b.registry.GetModList()
		if len(updatedModList) == 0 {
			b.registry.DB.UpdateDB(true)
			updatedModList = b.registry.GetModList()
		}
		return UpdatedModListMsg(updatedModList)
	}
}

// Check filesystem for installed mods
func (b Bubble) getInstalledModsCmd() tea.Cmd {
	return func() tea.Msg {
		log.Print("Checking installed mods")
		installedModList, err := dirfs.CheckInstalledMods()
		if err != nil {
			return ErrorMsg(err)
		}
		log.Printf("Found %d mods already installed", len(installedModList))
		return InstalledModListMsg(installedModList)
	}
}

// Manually input KSP directory
func (b Bubble) updateKspDirCmd(s string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Input received: %s", s)
		kerbalDir, err := dirfs.FindKspPath(s)
		if err != nil {
			return UpdateKspDirMsg(false)
		}
		kerbalVer := dirfs.FindKspVersion(kerbalDir)
		viper.Set("settings.kerbal_dir", kerbalDir)
		viper.Set("settings.kerbal_ver", kerbalVer.String())
		viper.WriteConfigAs(viper.ConfigFileUsed())
		log.Printf("Kerbal dir: " + kerbalDir + "/")
		log.Printf("Kerbal Version: %v", kerbalVer)
		return UpdateKspDirMsg(true)
	}
}

// Download selected mods
func (b Bubble) downloadModCmd(url string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Mod download requested: %s", url)
		err := dirfs.DownloadMod(url)
		if err != nil {
			log.Printf("Error downloading: %v", err)
			return DownloadModMsg(false)
		}
		return DownloadModMsg(true)
	}
}
