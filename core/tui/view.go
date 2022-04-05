package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/muesli/reflow/truncate"
)

func (b Bubble) View() string {
	var primaryBox string
	var secondaryBox string
	var view string

	primaryBoxBorder := lipgloss.NormalBorder()
	b.primaryViewport.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder)
	primaryBox = b.primaryViewport.View()

	if !b.loading {
		secondaryBoxBorder := lipgloss.NormalBorder()
		b.secondaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(secondaryBoxBorder)
		secondaryBox = b.secondaryViewport.View()

		view = lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				primaryBox,
				secondaryBox,
			),
			b.statusBarView(),
		)
	} else {
		view = lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				primaryBox,
			),
			b.statusBarView(),
		)
	}

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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, mod.Name)
	}

	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(s)
}

func (b Bubble) modInfoView() string {
	s := "\n"
	if b.selected >= 0 {
		var mod = b.modList[b.selected]

		s += "Mod\n\n"
		s += fmt.Sprintf(
			"Name: %s\n\n"+
				"Identifier: %s\n\n"+
				"Author: %s\n\n"+
				"Version: %s\n\n"+
				"Abstract: %s\n\n"+
				"Download: %s\n\n",
			mod.Name,
			mod.Identifier,
			mod.Author,
			mod.Version,
			mod.Abstract,
			mod.Download)
	} else {
		s += b.help.View()
	}

	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(s)
}

func (b Bubble) loadingView() string {
	s := "loading screen?"
	return s
}

func (b Bubble) statusBarView() string {
	var status = "Status: " + b.status

	width := lipgloss.Width
	selectedFileName := b.modList[b.cursor].Name
	var fileCount = fmt.Sprintf("%d/%d", b.cursor+1, len(b.modList))

	selectedFileStyle := constants.BoldTextStyle.Copy()
	selectedFileColumn := selectedFileStyle.
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Render(truncate.StringWithTail(selectedFileName, 50, "..."))

	fileCountStyle := constants.BoldTextStyle.Copy()
	fileCountColumn := fileCountStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Render(fileCount)

	statusStyle := constants.BoldTextStyle.Copy()
	statusColumn := statusStyle.
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Width(b.width - width(selectedFileColumn) - width(fileCountColumn)).
		Render(truncate.StringWithTail(
			status,
			uint(b.width-width(selectedFileColumn)-width(fileCountColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileColumn,
		statusColumn,
		fileCountColumn,
	)
}
