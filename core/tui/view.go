package tui

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/muesli/reflow/truncate"
)

func (b Bubble) View() string {
	var view string

	if b.activeBox == constants.SplashBoxActive {
		var splashBox string

		splashBoxBorder := lipgloss.NormalBorder()
		b.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(splashBoxBorder)
		splashBox = b.splashViewport.View()

		view = lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				splashBox,
			),
			b.statusBarView(),
		)
	} else {
		var primaryBox string
		var secondaryBox string

		primaryBoxBorder := lipgloss.NormalBorder()
		b.primaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(primaryBoxBorder)
		primaryBox = b.primaryViewport.View()
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
				"KSP Max Version: %s\n\n"+
				"KSP Min Version: %s\n\n"+
				"Abstract: %s\n\n"+
				"License: %s\n\n"+
				"Download: %s\n\n",
			mod.Name,
			mod.Identifier,
			mod.Author,
			mod.Version,
			mod.VersionKspMax,
			mod.VersionKspMin,
			mod.Abstract,
			mod.License,
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
	s := "Loading Screen\n"
	for i := range b.logs {
		s += b.logs[i] + "\n"
	}
	return lipgloss.NewStyle().
		Width(b.splashViewport.Width).
		Height(b.splashViewport.Height).
		Render(s)
}

func (b Bubble) logView() string {
	s := "Logs\n"
	/* for i := range b.logs {
		s += b.logs[i] + "\n"
	} */

	content, err := ioutil.ReadFile("debug.log")
	if err != nil {
		log.Fatal(err)
	}
	s += string(content)

	return lipgloss.NewStyle().
		Width(b.splashViewport.Width).
		Height(b.splashViewport.Height).
		Render(s)
}

func (b Bubble) statusBarView() string {
	var status = "Status: " + b.logs[len(b.logs)-1]

	width := lipgloss.Width
	var fileCount = fmt.Sprintf("%d/%d", b.cursor+1, len(b.modList))

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
		Width(b.width - width(fileCountColumn)).
		Render(truncate.StringWithTail(
			status,
			uint(b.width-width(fileCountColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		statusColumn,
		fileCountColumn,
	)
}
