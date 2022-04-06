package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/registry/database"
)

type getAvailableModsMsg []database.Ckan

func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		b.registry.DB.UpdateDB(true)
		updatedModList := b.registry.DB.GetModList()
		return getAvailableModsMsg(updatedModList)
	}
}

func (b Bubble) updateDbCmd() {
	b.registry.DB.UpdateDB(true)
}
