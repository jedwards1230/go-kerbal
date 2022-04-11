package tui

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/muesli/reflow/truncate"
)

func (b Bubble) View() string {
	var view string

	switch b.activeBox {
	case constants.PrimaryBoxActive, constants.SecondaryBoxActive:
		var primaryBox string
		var secondaryBox string

		primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
		secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor

		switch b.activeBox {
		case constants.PrimaryBoxActive:
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		case constants.SecondaryBoxActive:
			secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor

		}

		primaryBoxBorder := lipgloss.NormalBorder()
		b.primaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(primaryBoxBorder).
			BorderForeground(primaryBoxBorderColor)
		primaryBox = b.primaryViewport.View()

		secondaryBoxBorder := lipgloss.NormalBorder()
		b.secondaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(secondaryBoxBorder).
			BorderForeground(secondaryBoxBorderColor)
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
	case constants.SplashBoxActive:
		var splashBox string

		splashBoxBorder := lipgloss.NormalBorder()
		b.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(splashBoxBorder).
			BorderForeground(b.theme.ActiveBoxBorderColor)
		splashBox = b.splashViewport.View()

		view = lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				splashBox,
			),
			b.statusBarView(),
		)
	}

	return view
}

func (b Bubble) modListView() string {
	// construct list
	s := "\n  Mod List\n\n"
	for i := range b.registry.SortedModList {
		cursor := " "
		if b.cursor == i {
			cursor = ">"
		}
		checked := " "
		if i == b.selected {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, b.registry.SortedModList[i].Name)
	}

	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(s)
}

func (b Bubble) modInfoView() string {
	s := "\n"
	if b.selected >= 0 {
		var mod = b.registry.SortedModList[b.selected]

		s += "Mod\n\n"
		s += fmt.Sprintf(
			"Name:             %s\n"+
				"Identifier:       %s\n"+
				"Author:           %s\n\n"+
				"Version:          %s\n"+
				"KSP Max Version:  %s\n"+
				"KSP Min Version:  %s\n\n"+
				"Abstract:         %s\n\n"+
				"License:          %s\n\n"+
				"Download:         %s\n\n",
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
	cfg := config.GetConfig()
	var status = "Status: " + b.logs[len(b.logs)-1]

	width := lipgloss.Width
	var fileCount = fmt.Sprintf("Mod: %d/%d", b.cursor+1, len(b.registry.SortedModList))

	fileCountStyle := constants.BoldTextStyle.Copy()
	fileCountColumn := fileCountStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Render(fileCount)

	sortOptions := fmt.Sprintf("Sort: %s by %s", b.sortOptions.SortOrder, b.sortOptions.SortTag)

	sortOptionsStyle := constants.BoldTextStyle.Copy()
	sortOptionsColumn := sortOptionsStyle.
		Align(lipgloss.Right).
		Padding(0, 3).
		Height(constants.StatusBarHeight).
		Render(sortOptions)

	var showCompatible string
	if cfg.Settings.HideIncompatibleMods {
		showCompatible = "Hide incompatible mods"
	} else {
		showCompatible = "Show incompatible mods"
	}

	showCompatibleStyle := constants.BoldTextStyle.Copy()
	showCompatibleColumn := showCompatibleStyle.
		Align(lipgloss.Right).
		Padding(0, 3).
		Height(constants.StatusBarHeight).
		Render(showCompatible)

	statusStyle := constants.BoldTextStyle.Copy()
	statusColumn := statusStyle.
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Width(b.width - width(fileCountColumn) - width(sortOptionsColumn) - width(showCompatibleColumn)).
		Render(truncate.StringWithTail(
			status,
			uint(b.width-width(fileCountColumn)-width(sortOptionsColumn)-width(showCompatibleColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		statusColumn,
		sortOptionsColumn,
		showCompatibleColumn,
		fileCountColumn,
	)
}
