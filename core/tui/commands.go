package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/jedwards1230/go-kerbal/registry/database"
	"github.com/spf13/viper"
)

type ModListUpdatedMsg []database.Ckan
type UpdateKspDirMsg bool

// Request the mod list from the database
func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		b.registry.DB.UpdateDB(false)
		updatedModList := b.registry.GetModList()
		if len(updatedModList) == 0 {
			b.registry.DB.UpdateDB(true)
			updatedModList = b.registry.GetModList()
		}
		return ModListUpdatedMsg(updatedModList)
	}
}

func (b Bubble) updateKspDirCmd(s string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Input received: %s", s)
		kerbalDir, err := dirfs.FindKspPath(s)
		if err == nil {
			kerbalVer := dirfs.FindKspVersion(kerbalDir)
			viper.Set("settings.kerbal_dir", kerbalDir)
			viper.Set("settings.kerbal_ver", kerbalVer.String())
			viper.WriteConfigAs(viper.ConfigFileUsed())
			log.Printf("Kerbal dir: " + kerbalDir + "/")
			log.Printf("Kerbal Version: %v", kerbalVer)
			return UpdateKspDirMsg(true)
		}
		return UpdateKspDirMsg(false)
	}
}
