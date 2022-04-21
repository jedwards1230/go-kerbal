package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/jedwards1230/go-kerbal/registry/database"
	"github.com/spf13/viper"
)

type (
	UpdatedModMapMsg    map[string][]database.Ckan
	InstalledModListMsg map[string]bool
	UpdateKspDirMsg     bool
	ErrorMsg            error
)

type SearchMsg struct {
	registry.ModIndex
}

// Request the mod list from the database
func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		log.Print("Checking available mods")
		b.registry.DB.UpdateDB(false)
		updatedModMap := b.registry.GetTotalModMap()
		if len(updatedModMap) == 0 {
			b.registry.DB.UpdateDB(true)
			updatedModMap = b.registry.GetTotalModMap()
		}
		return UpdatedModMapMsg(updatedModMap)
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
		// todo: crashed here setting current ksp dir as new ksp dir in-app
		viper.Set("settings.kerbal_ver", kerbalVer.String())
		viper.WriteConfigAs(viper.ConfigFileUsed())
		log.Printf("Kerbal dir: " + kerbalDir + "/")
		log.Printf("Kerbal Version: %v", kerbalVer)
		return UpdateKspDirMsg(true)
	}
}

// Download selected mods
func (b Bubble) downloadModCmd() tea.Cmd {
	return func() tea.Msg {
		mods, err := b.registry.DownloadMods(b.nav.installSelected)
		if err != nil {
			log.Printf("Error downloading: %v", err)
			return ErrorMsg(err)
		}

		err = registry.InstallMods(mods)
		if err != nil {
			log.Printf("Error downloading: %v", err)
			return ErrorMsg(err)
		}

		installedModList, err := dirfs.CheckInstalledMods()
		if err != nil {
			return ErrorMsg(err)
		}
		return InstalledModListMsg(installedModList)
	}
}

// Download selected mods
func (b Bubble) searchCmd(s string) tea.Cmd {
	return func() tea.Msg {
		searchMapIndex, err := b.registry.BuildSearchMapIndex(s)
		if err != nil {
			log.Printf("Error building search index: %v", err)
			return SearchMsg{searchMapIndex}
		}
		return SearchMsg{searchMapIndex}
	}
}
