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

		// set colors
		primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
		secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor

		// format active box
		switch b.activeBox {
		case constants.PrimaryBoxActive:
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		case constants.SecondaryBoxActive:
			secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor

		}

		// format views
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

		// organize views
		view = lipgloss.JoinVertical(
			lipgloss.Top,
			b.getMainButtonsView(),
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
	title := lipgloss.NewStyle().
		Bold(true).
		Width(b.primaryViewport.Width).
		Height(3).
		Padding(1).
		Render("Mod List")

	s := ""
	for i := range b.registry.SortedModList {
		checked := " "
		if i == b.nav.listSelected {
			checked = "x"
		}

		line := truncate.StringWithTail(
			fmt.Sprintf("[%s] %s", checked, b.registry.SortedModList[i].Name),
			uint(b.primaryViewport.Width-2),
			"...")

		if b.nav.listCursor == i {
			s += lipgloss.NewStyle().
				Background(b.theme.SelectedListItemColor).
				Foreground(b.theme.UnselectedListItemColor).
				Width(b.primaryViewport.Width).
				Render(line)
		} else {
			if b.registry.SortedModList[i].Installed {
				s += lipgloss.NewStyle().
					Foreground(b.theme.InstalledListItemColor).
					Render(line)
			} else {
				s += lipgloss.NewStyle().
					Render(line)
			}
		}
		s += "\n"
	}

	modList := lipgloss.NewStyle().
		Width(b.primaryViewport.Width).
		Height(b.primaryViewport.Height - 3).
		Render(s)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		modList,
	)
}

func (b Bubble) modInfoView() string {
	s := "\n"
	if b.nav.listSelected >= 0 {
		var mod = b.registry.SortedModList[b.nav.listSelected]

		s += "Mod\n\n"
		s += fmt.Sprintf(
			"Name:             %s\n"+
				"Author:           %s\n\n"+
				"Identifier:       %s\n"+
				"Installed:        %v\n\n"+
				"Version:          %s\n"+
				"KSP Max Version:  %s\n"+
				"KSP Min Version:  %s\n\n"+
				"Abstract:         %s\n\n"+
				"License:          %s\n\n"+
				"Download:         %s\n\n"+
				"Dependencies:         %s\n\n",
			mod.Name,
			mod.Author,
			mod.Identifier,
			mod.Installed,
			mod.Version,
			mod.VersionKspMax,
			mod.VersionKspMin,
			mod.Abstract,
			mod.License,
			mod.Download,
			mod.ModDepends)
	} else {
		s += b.help.View()
	}

	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(s)
}

func (b Bubble) logView() string {
	s := "Logs\n"

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

func (b Bubble) inputKspView() string {
	titleStyle := constants.BoldTextStyle.Copy()
	title := titleStyle.
		Align(lipgloss.Left).
		Width(b.width).
		Padding(1).
		Render("Please enter the path to your Kerbal Space Program directory:")

	inText := ""
	if b.inputRequested {
		inText = b.textInput.View()
	} else {
		inText = b.textInput.Value()
		inText += "\n\nPress Esc to close"
	}

	inTextStyle := constants.BoldTextStyle.Copy()
	inTextColumn := inTextStyle.
		Align(lipgloss.Left).
		Width(b.width).
		Padding(1).
		Render(inText)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		inTextColumn,
	)
}

func (b Bubble) statusBarView() string {
	cfg := config.GetConfig()
	var status = "Status: " + b.logs[len(b.logs)-1]

	width := lipgloss.Width
	var fileCount = fmt.Sprintf("Mod: %d/%d", b.nav.listCursor+1, len(b.registry.SortedModList))

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
		showCompatible = "Hiding incompatible mods"
	} else {
		showCompatible = "Showing incompatible mods"
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

func (b *Bubble) getMainButtonsView() string {
	cfg := config.GetConfig()
	refreshStyle := lipgloss.NewStyle().
		Underline(true)
	refreshColumn := refreshStyle.
		Align(lipgloss.Right).
		Padding(0, 2).
		Height(constants.StatusBarHeight).
		Render("1. Refresh")

	showCompatible := "2. Hide incompatible mods"
	if !cfg.Settings.HideIncompatibleMods {
		showCompatible = "2. Hide incompatible mods"
	} else {
		showCompatible = "2. Show incompatible mods"
	}

	hideIncompatibleStyle := lipgloss.NewStyle().
		Underline(true)
	hideIncompatibleColumn := hideIncompatibleStyle.
		Align(lipgloss.Right).
		Padding(0, 2).
		Height(constants.StatusBarHeight).
		Render(showCompatible)

	sortOrderStyle := lipgloss.NewStyle().
		Underline(true)
	sortOrderColumn := sortOrderStyle.
		Align(lipgloss.Right).
		Padding(0, 2).
		Height(constants.StatusBarHeight).
		Render("3. Sort Order")

	enterDirStyle := lipgloss.NewStyle().
		Underline(true)
	enterDirColumn := enterDirStyle.
		Align(lipgloss.Right).
		Padding(0, 2).
		Height(constants.StatusBarHeight).
		Render("4. Enter KSP Dir")

	downloadStyle := lipgloss.NewStyle().
		Underline(true)
	downloadColumn := downloadStyle.
		Align(lipgloss.Right).
		Padding(0, 2).
		Height(constants.StatusBarHeight).
		Render("5. Download mod")

	return lipgloss.JoinHorizontal(lipgloss.Top,
		refreshColumn,
		hideIncompatibleColumn,
		sortOrderColumn,
		enterDirColumn,
		downloadColumn,
	)
}
