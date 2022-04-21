package tui

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

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
	modMap := b.registry.GetActiveModMap()

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.primaryViewport.Width).
		Height(3).
		Padding(1)

	title := titleStyle.Render("Mod List")
	if b.searchInput {
		title = titleStyle.Render("Search Mods")
	}

	s := ""
	for i, id := range b.registry.ModMapIndex {
		mod := modMap[id.Key]

		checked := " "
		if b.nav.installSelected[mod.Identifier].Identifier != "" {
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
	modMap := b.registry.GetActiveModMap()

	var title, body string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.secondaryViewport.Width).
		Height(3).
		Padding(1)

	if b.nav.listSelected >= 0 {
		id := b.registry.ModMapIndex[b.nav.listSelected]
		mod := modMap[id.Key]

		keyStyle := lipgloss.NewStyle().
			Align(lipgloss.Left).
			Bold(true).
			Width(b.secondaryViewport.Width/4).
			Padding(0, 2)

		valueStyle := lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(b.secondaryViewport.Width*3/4).
			Padding(0, 5)

		title = titleStyle.Render("Mod")

		nameKey := keyStyle.Render("Name")
		namealue := valueStyle.Render(mod.Name)
		name := lipgloss.JoinHorizontal(lipgloss.Top, nameKey, namealue)

		identifierKey := keyStyle.Render("Identifier")
		identifierValue := valueStyle.Render(mod.Identifier)
		identifier := lipgloss.JoinHorizontal(lipgloss.Top, identifierKey, identifierValue)

		authorKey := keyStyle.Render("Author")
		authorValue := valueStyle.Render(mod.Author)
		author := lipgloss.JoinHorizontal(lipgloss.Top, authorKey, authorValue)

		versionKey := keyStyle.Render("Mod Version")
		versionValue := valueStyle.Render(mod.Versions.Mod)
		version := lipgloss.JoinHorizontal(lipgloss.Top, versionKey, versionValue)

		versionKspKey := keyStyle.Render("KSP Versions")
		versionKspValue := fmt.Sprintf("%v - %v", mod.Versions.KspMin, mod.Versions.KspMax)
		versionKspValue = valueStyle.Render(versionKspValue)
		versionKsp := lipgloss.JoinHorizontal(lipgloss.Top, versionKspKey, versionKspValue)

		installedKey := keyStyle.
			Foreground(lipgloss.NoColor{}).
			Render("Installed?")
		installedValue := ""
		if mod.Install.Installed {
			installedValue = valueStyle.Copy().
				Foreground(b.theme.InstalledListItemColor).
				Render("Installed")
		} else {
			installedValue = valueStyle.
				Render("Not Installed")
		}
		installed := lipgloss.JoinHorizontal(lipgloss.Top, installedKey, installedValue)

		installDirKey := keyStyle.Render("Install dir")
		installDirValue := valueStyle.Render(mod.Install.InstallTo)
		installDir := lipgloss.JoinHorizontal(lipgloss.Top, installDirKey, installDirValue)

		downloadKey := keyStyle.Render("License")
		downloadValue := valueStyle.Render(mod.Install.Download)
		download := lipgloss.JoinHorizontal(lipgloss.Top, downloadKey, downloadValue)

		abstractKey := keyStyle.Render("Abstract")
		abstractValue := valueStyle.Render(mod.Abstract)
		abstract := lipgloss.JoinHorizontal(lipgloss.Top, abstractKey, abstractValue)

		dependenciesKey := keyStyle.Render("Dependencies")
		dependenciesValue := valueStyle.Render("None")
		if len(mod.ModDepends) > 0 {
			dependenciesValue = valueStyle.Render(strings.Join(mod.ModDepends, ", "))
		}
		dependencies := lipgloss.JoinHorizontal(lipgloss.Top, dependenciesKey, dependenciesValue)

		conflictsKey := keyStyle.Render("Conflicts")
		conflictsValue := valueStyle.Render("None")
		if len(mod.ModConflicts) > 0 {
			conflictsValue = valueStyle.Render(strings.Join(mod.ModConflicts, ", "))
		}
		conflicts := lipgloss.JoinHorizontal(lipgloss.Top, conflictsKey, conflictsValue)

		licenseKey := keyStyle.Render("License")
		licenseValue := valueStyle.Render(mod.License)
		license := lipgloss.JoinHorizontal(lipgloss.Top, licenseKey, licenseValue)

		body = lipgloss.JoinVertical(
			lipgloss.Top,
			name,
			identifier,
			author,
			"\n",
			version,
			versionKsp,
			"\n",
			installed,
			installDir,
			download,
			"\n",
			abstract,
			"\n",
			dependencies,
			conflicts,
			"\n",
			license,
		)
	} else {
		title = titleStyle.Render("Help Menu")
		body = lipgloss.NewStyle().
			Width(b.secondaryViewport.Width).
			Height(b.secondaryViewport.Height - 3).
			Render(b.help.View())
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

	file, err := ioutil.ReadFile("./logs/debug.log")
	if err != nil {
		log.Fatal(err)
	}

	body := lipgloss.NewStyle().
		Width(b.splashViewport.Width).
		Height(b.splashViewport.Height - 3).
		Render(string(file))

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

	var status string
	statusWidth := b.width - width(fileCountColumn) - width(sortOptionsColumn) - width(showCompatibleColumn)
	if b.searchInput {
		status = statusBarStyle.
			Align(lipgloss.Left).
			Padding(0, 1).
			Width(statusWidth).
			Render(b.textInput.View())
	} else {
		status = "Status: " + b.logs[len(b.logs)-1]
		status = statusBarStyle.
			Align(lipgloss.Left).
			Padding(0, 1).
			Width(statusWidth).
			Render(truncate.StringWithTail(
				status,
				uint(statusWidth-3),
				"..."),
			)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		status,
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

	showCompatible := buttonStyle.Render("2. Hide incompatible mods")
	if cfg.Settings.HideIncompatibleMods {
		showCompatible = buttonStyle.Render("2. Show incompatible mods")
	}
	sortOrder := buttonStyle.Render("3. Sort Order")
	download := buttonStyle.Render("5. Download mod")

	if b.searchInput {
		escape := buttonStyle.
			Align(lipgloss.Left).
			Render("Esc to close")

		leftColumn := lipgloss.JoinHorizontal(lipgloss.Top,
			escape,
			showCompatible,
			sortOrder,
			download,
		)

		enableInput := buttonStyle.
			Align(lipgloss.Right).
			Render("6. Enable text input")
		if b.inputRequested {
			enableInput = buttonStyle.
				Align(lipgloss.Right).
				Render("6. Disable text input")
		}

		leftColumn = lipgloss.NewStyle().Width(b.width - lipgloss.Width(enableInput)).Render(leftColumn)

		return lipgloss.JoinHorizontal(lipgloss.Top,
			leftColumn,
			enableInput,
		)
	} else {
		refresh := buttonStyle.Render("1. Refresh")
		enterDir := buttonStyle.Render("4. Enter KSP Dir")
		search := buttonStyle.Render("6. Search")

		leftColumn := lipgloss.JoinHorizontal(lipgloss.Top,
			refresh,
			showCompatible,
			sortOrder,
			enterDir,
			download,
			search,
		)

		settings := buttonStyle.
			Align(lipgloss.Right).
			Render("0. Settings")

		leftColumn = lipgloss.NewStyle().
			Width(b.width - lipgloss.Width(settings)).
			Render(leftColumn)

		return lipgloss.JoinHorizontal(lipgloss.Top,
			leftColumn,
			settings,
		)
	}
}
