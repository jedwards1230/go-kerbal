package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/registry/database"
)

type ModListUpdatedMsg []database.Ckan

// Request the mod list from the database
func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		updatedModList := b.registry.GetModList()
		if len(updatedModList) == 0 {
			b.registry.DB.UpdateDB(true)
			updatedModList = b.registry.GetModList()
		}
		return ModListUpdatedMsg(updatedModList)
	}
}
