package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/registry/datacollector"
)

type getAvailableModsMsg []datacollector.Ckan

func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		b.status = "Getting Mod List"
		updatedModList := datacollector.GetAvailableMods()
		return getAvailableModsMsg(updatedModList)
	}
}
