package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/registry/database"
)

type getAvailableModsMsg []database.Ckan

// Request the mod list from the database
func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		updatedModList := b.registry.DB.GetModList()
		return getAvailableModsMsg(updatedModList)
	}
}

// Request to ping the metadata repo and update if needed.
// Cloning the repo can be forced with a bool parameter
func (b Bubble) updateDbCmd(force bool) {
	b.registry.DB.UpdateDB(force)
}
