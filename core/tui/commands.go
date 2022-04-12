package tui

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/registry/database"
)

type ModListUpdatedMsg []database.Ckan
type UpdateKspDirMsg error

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
		// update config file
		// validate directory
		log.Printf("Input received: %s", s)
		err := errors.New("validate directory plz")
		return UpdateKspDirMsg(err)
	}
}
