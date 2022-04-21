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
		b.primaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(primaryBoxBorderColor)
		primaryBox = b.primaryViewport.View()

		b.secondaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(lipgloss.NormalBorder()).
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

		b.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(lipgloss.NormalBorder()).
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
	title := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.primaryViewport.Width).
		Height(3).
		Padding(1).
		Render("Mod List")

	s := ""
	for i, id := range b.registry.ModMapIndex {
		mod := b.registry.SortedMap[id.Key]

		checked := " "
		if b.nav.installSelected[mod.Identifier] {
			checked = "x"
		}

		line := truncate.StringWithTail(
			fmt.Sprintf("[%s] %s", checked, mod.Name),
			uint(b.primaryViewport.Width-2),
			"...")

		if b.nav.listCursor == i {
			s += lipgloss.NewStyle().
				Background(b.theme.SelectedListItemColor).
				Foreground(b.theme.UnselectedListItemColor).
				Width(b.primaryViewport.Width).
				Render(line)
		} else {
			if mod.Install.Installed {
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

	body := lipgloss.NewStyle().
		Width(b.primaryViewport.Width).
		Height(b.primaryViewport.Height - 3).
		Render(s)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) modInfoView() string {
	var title, body string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.secondaryViewport.Width).
		Height(3).
		Padding(1)

	bodyStyle := lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height - 3)

	if b.nav.listSelected >= 0 {
		id := b.registry.ModMapIndex[b.nav.listSelected]
		mod := b.registry.SortedMap[id.Key]

		title = titleStyle.Render(mod.Name)

		s := fmt.Sprintf(
			""+
				"Author:           %v\n\n"+
				"Installed:        %v\n"+
				"Installed to:     %v\n\n"+
				"Abstract:         %v\n\n"+
				"Identifier:       %v\n"+
				"Version:          %v\n"+
				"KSP Version:      %v - %v\n\n"+
				"Dependencies:     %v\n"+
				"Conflicts:        %v\n\n"+
				"License:          %v\n\n"+
				"Download:         %v\n\n",
			mod.Author,
			mod.Install.Installed,
			mod.Install.InstallTo,
			mod.Abstract,
			mod.Identifier,
			mod.Versions.Mod,
			mod.Versions.KspMin,
			mod.Versions.KspMax,
			mod.ModConflicts,
			mod.ModDepends,
			mod.License,
			mod.Install.Download)

		body = bodyStyle.Render(s)
	} else {
		title = titleStyle.Render("Help Menu")
		body = bodyStyle.Render(b.help.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) logView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.splashViewport.Width).
		Height(3).
		Padding(1).
		Render("Logs")

	content, err := ioutil.ReadFile("./logs/debug.log")
	if err != nil {
		log.Fatal(err)
	}
	s := string(content)

	body := lipgloss.NewStyle().
		Width(b.splashViewport.Width).
		Height(b.splashViewport.Height - 3).
		Render(s)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) settingsView() string {
	cfg := config.GetConfig()
	title := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.splashViewport.Width).
		Height(3).
		Padding(1).
		Render("Settings")

	lineStyle := lipgloss.NewStyle().
		Padding(1).
		Align(lipgloss.Left).
		Width(b.splashViewport.Width)

	content := lineStyle.Render(fmt.Sprintf("Kerbal Directory: %v", cfg.Settings.KerbalDir))
	content += lineStyle.Render(fmt.Sprintf("Kerbal Version: %v", cfg.Settings.KerbalVer))
	content += lineStyle.Render(fmt.Sprintf("Logging: %v", cfg.Settings.EnableLogging))
	content += lineStyle.Render(fmt.Sprintf("Mousewheel: %v", cfg.Settings.EnableMouseWheel))
	content += lineStyle.Render(fmt.Sprintf("Hide incompatible: %v", cfg.Settings.HideIncompatibleMods))
	content += lineStyle.Render(fmt.Sprintf("Metadata Repo: %v", cfg.Settings.MetaRepo))
	content += lineStyle.Render(fmt.Sprintf("Last Repo Hash: %v", cfg.Settings.LastRepoHash))
	content += lineStyle.Render(fmt.Sprintf("Theme: %v", cfg.AppTheme))

	body := lipgloss.NewStyle().
		Width(b.splashViewport.Width).
		Height(b.splashViewport.Height - 3).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) inputKspView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.splashViewport.Width).
		Height(3).
		Padding(1).
		Render("Kerbal Space Program Directory")

	question := lipgloss.NewStyle().
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

	inText = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(b.width).
		Padding(1).
		Render(inText)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		question,
		inText,
	)
}

func (b Bubble) statusBarView() string {
	cfg := config.GetConfig()
	width := lipgloss.Width

	statusBarStyle := lipgloss.NewStyle().
		Height(constants.StatusBarHeight)

	fileCount := fmt.Sprintf("Mod: %d/%d", b.nav.listCursor+1, len(b.registry.ModMapIndex))
	fileCountColumn := statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Render(fileCount)

	sortOptions := fmt.Sprintf("Sort: %s by %s", b.registry.SortOptions.SortOrder, b.registry.SortOptions.SortTag)
	sortOptionsColumn := statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 3).
		Render(sortOptions)

	var showCompatible string
	if cfg.Settings.HideIncompatibleMods {
		showCompatible = "Hiding incompatible mods"
	} else {
		showCompatible = "Showing incompatible mods"
	}
	showCompatibleColumn := statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 3).
		Render(showCompatible)

	status := "Status: " + b.logs[len(b.logs)-1]
	statusColumn := statusBarStyle.
		Align(lipgloss.Left).
		Padding(0, 1).
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

func (b Bubble) getMainButtonsView() string {
	cfg := config.GetConfig()

	buttonStyle := lipgloss.NewStyle().
		Underline(true).
		Padding(0, 2).
		Height(constants.StatusBarHeight)

	refreshColumn := buttonStyle.Render("1. Refresh")

	showCompatible := "2. Hide incompatible mods"
	if !cfg.Settings.HideIncompatibleMods {
		showCompatible = "2. Hide incompatible mods"
	} else {
		showCompatible = "2. Show incompatible mods"
	}
	hideIncompatibleColumn := buttonStyle.Render(showCompatible)

	sortOrderColumn := buttonStyle.Render("3. Sort Order")

	enterDirColumn := buttonStyle.Render("4. Enter KSP Dir")

	downloadColumn := buttonStyle.Render("5. Download mod")

	settingsColumn := buttonStyle.Render("0. Settings")

	return lipgloss.JoinHorizontal(lipgloss.Top,
		refreshColumn,
		hideIncompatibleColumn,
		sortOrderColumn,
		enterDirColumn,
		downloadColumn,
		settingsColumn,
	)
}
