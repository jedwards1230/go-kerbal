package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
)

func (b Bubble) View() string {
	b.primaryViewport.SetContent(b.modListView())
	b.secondaryViewport.SetContent(b.modInfoView())

	var primaryBox string
	var secondaryBox string
	primaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorder := lipgloss.NormalBorder()

	b.primaryViewport.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder)

	b.secondaryViewport.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(secondaryBoxBorder)

	primaryBox = b.primaryViewport.View()
	secondaryBox = b.secondaryViewport.View()

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			primaryBox,
			secondaryBox,
		),
	)

	return view
}

func (b Bubble) modListView() string {
	// construct list
	s := "\n  Mod List\n\n"
	for i, mod := range b.modList {
		cursor := " "
		if b.cursor == i {
			cursor = ">"
		}
		checked := " "
		if i == b.selected {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, mod)
	}

	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(s)
}

func (b Bubble) modInfoView() string {
	// construct list
	s := "\n  Mod \n\n"

	if b.selected >= 0 {
		mod := b.modList[b.selected]
		s += fmt.Sprintf("Name: %s\n\nAuthor: %s\n", mod.Name, mod.Author)

		s += "\nPress q to quit.\n\n"
	}

	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(s)
}
