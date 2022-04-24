package tui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/muesli/reflow/truncate"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Height(3).
		Padding(1)
)

func (b Bubble) View() string {
	var body string

	switch b.activeBox {
	case internal.ModListView, internal.ModInfoView, internal.SearchView:
		var primaryBox string
		var secondaryBox string

		// set colors
		primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
		secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor

		// format active box
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		case internal.ModInfoView:
			secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		case internal.EnterKspDirView, internal.LogView, internal.SettingsView:
			secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		}

		// format views
		b.bubbles.primaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(primaryBoxBorderColor)
		primaryBox = b.bubbles.primaryViewport.View()

		b.bubbles.secondaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(secondaryBoxBorderColor)
		secondaryBox = b.bubbles.secondaryViewport.View()

		// organize views
		body = lipgloss.JoinHorizontal(
			lipgloss.Top,
			primaryBox,
			secondaryBox,
		)
	case internal.LogView, internal.SettingsView, internal.EnterKspDirView:
		b.bubbles.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(b.theme.ActiveBoxBorderColor)
		body = b.bubbles.splashViewport.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		b.getMainButtonsView(),
		body,
		b.statusBarView(),
	)
}

func (b Bubble) modListView() string {
	modMap := b.registry.GetActiveModList()

	titleStyle := titleStyle.Copy().
		Width(b.bubbles.primaryViewport.Width)

	title := titleStyle.Render("Mod List")
	if b.activeBox == internal.SearchView {
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
			uint(b.bubbles.primaryViewport.Width-2),
			"...")

		if b.nav.listSelected == i {
			s += lipgloss.NewStyle().
				Background(b.theme.SelectedListItemColor).
				Foreground(b.theme.UnselectedListItemColor).
				Width(b.bubbles.primaryViewport.Width).
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
		Width(b.bubbles.primaryViewport.Width).
		Height(b.bubbles.primaryViewport.Height - 3).
		Render(s)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) modInfoView() string {
	modMap := b.registry.GetActiveModList()

	var title, body string

	titleStyle := titleStyle.Copy().
		Width(b.bubbles.primaryViewport.Width)

	if b.nav.listSelected >= 0 && b.nav.listSelected < len(b.registry.ModMapIndex) {
		id := b.registry.ModMapIndex[b.nav.listSelected]
		mod := modMap[id.Key]

		keyStyle := lipgloss.NewStyle().
			Align(lipgloss.Left).
			Bold(true).
			Width(b.bubbles.secondaryViewport.Width/4).
			Padding(0, 2)

		valueStyle := lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(b.bubbles.secondaryViewport.Width*3/4).
			Padding(0, 5)

		title = titleStyle.Render(mod.Name)

		abstract := titleStyle.
			Bold(false).
			Render(mod.Abstract)

		if mod.Description != "" {
			abstract = titleStyle.
				Bold(false).
				Render(mod.Description)
		}

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

		downloadKey := keyStyle.Render("Download")
		downloadValue := valueStyle.Render(mod.Download.URL)
		download := lipgloss.JoinHorizontal(lipgloss.Top, downloadKey, downloadValue)

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
			abstract,
			"\n",
			author,
			identifier,
			license,
			"\n",
			version,
			versionKsp,
			"\n",
			installed,
			installDir,
			download,
			"\n",
			dependencies,
			conflicts,
		)
	} else {
		title = titleStyle.Render("Help Menu")
		body = lipgloss.NewStyle().
			Width(b.bubbles.secondaryViewport.Width).
			Height(b.bubbles.secondaryViewport.Height - 3).
			Render(b.bubbles.help.View())
	}
	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) logView() string {
	file, err := os.Open("./logs/debug.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bodyList := make([]string, 0)
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		lineWords := strings.Fields(scanner.Text())
		if len(lineWords) > 1 {
			idx := lipgloss.NewStyle().
				Align(lipgloss.Left).
				Width(6).
				Padding(0, 1).
				Render(fmt.Sprint(i) + " ")

			lineWords[0] = lipgloss.NewStyle().
				Foreground(b.theme.Green).
				Render(lineWords[0])
			lineWords[1] = lipgloss.NewStyle().
				Foreground(b.theme.Orange).
				Width(16).
				Render(lineWords[1])
			line := lipgloss.JoinHorizontal(lipgloss.Left, idx, strings.Join(lineWords, " "))
			bodyList = append(bodyList, line)
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	title := titleStyle.Copy().
		Width(b.bubbles.splashViewport.Width).
		Render("Logs")

	body := lipgloss.JoinVertical(lipgloss.Top, bodyList...)
	body = lipgloss.NewStyle().
		Width(b.bubbles.splashViewport.Width).
		Height(b.bubbles.splashViewport.Height - 3).
		Render(string(body))

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) settingsView() string {
	cfg := config.GetConfig()
	title := titleStyle.Copy().
		Width(b.bubbles.splashViewport.Width).
		Render("Settings")

	lineStyle := lipgloss.NewStyle().
		Padding(1).
		Align(lipgloss.Left).
		Width(b.bubbles.splashViewport.Width)

	content := lineStyle.Render(fmt.Sprintf("Kerbal Directory: %v", cfg.Settings.KerbalDir))
	content += lineStyle.Render(fmt.Sprintf("Kerbal Version: %v", cfg.Settings.KerbalVer))
	content += lineStyle.Render(fmt.Sprintf("Logging: %v", cfg.Settings.EnableLogging))
	content += lineStyle.Render(fmt.Sprintf("Mousewheel: %v", cfg.Settings.EnableMouseWheel))
	content += lineStyle.Render(fmt.Sprintf("Hide incompatible: %v", cfg.Settings.HideIncompatibleMods))
	content += lineStyle.Render(fmt.Sprintf("Metadata Repo: %v", cfg.Settings.MetaRepo))
	content += lineStyle.Render(fmt.Sprintf("Last Repo Hash: %v", cfg.Settings.LastRepoHash))
	content += lineStyle.Render(fmt.Sprintf("Theme: %v", cfg.AppTheme))

	body := lipgloss.NewStyle().
		Width(b.bubbles.splashViewport.Width).
		Height(b.bubbles.splashViewport.Height - 3).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) inputKspView() string {
	title := titleStyle.Copy().
		Width(b.bubbles.splashViewport.Width).
		Render("Kerbal Space Program Directory")

	question := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(b.width).
		Padding(1).
		Render("Please enter the path to your Kerbal Space Program directory:")

	inText := ""
	if b.inputRequested {
		inText = b.bubbles.textInput.View()
	} else {
		inText = b.bubbles.textInput.Value()
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
		Height(internal.StatusBarHeight)

	fileCount := fmt.Sprintf("Mod: %d/%d", b.nav.listCursor+1, len(b.registry.ModMapIndex))
	fileCount = statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 3).
		Render(fileCount)

	sortOptions := fmt.Sprintf("Sort: %s by %s", b.registry.SortOptions.SortOrder, b.registry.SortOptions.SortTag)
	sortOptions = statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 3).
		Render(sortOptions)

	var showCompatible string
	if cfg.Settings.HideIncompatibleMods {
		showCompatible = "Hiding incompatible mods"
	} else {
		showCompatible = "Showing incompatible mods"
	}
	showCompatible = statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 3).
		Render(showCompatible)

	var status string
	statusWidth := b.width - width(fileCount) - width(sortOptions) - width(showCompatible)
	if b.searchInput {
		status = statusBarStyle.
			Align(lipgloss.Left).
			Padding(0, 2).
			Width(statusWidth).
			Render(b.bubbles.textInput.View())
	} else {
		// open log file
		file, err := os.Open("./logs/debug.log")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// read log file
		bodyList := make([]string, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineWords := strings.Fields(scanner.Text())
			if len(lineWords) > 1 {
				line := lipgloss.JoinHorizontal(lipgloss.Left, strings.Join(lineWords[2:], " "))
				bodyList = append(bodyList, line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		status = "Status: " + bodyList[len(bodyList)-1]

		// format status message
		status = statusBarStyle.
			Align(lipgloss.Left).
			Padding(0, 1).
			Width(statusWidth).
			Render(truncate.StringWithTail(
				status,
				uint(statusWidth-3),
				"..."),
			)

		spin := lipgloss.NewStyle().
			Width(b.bubbles.spinner.Style.GetWidth()).
			Padding(0, 0, 0, 1).
			Render("  ")

		if !b.ready {
			spin = lipgloss.NewStyle().
				Width(b.bubbles.spinner.Style.GetWidth()).
				Padding(0, 0, 0, 1).
				Render(b.bubbles.spinner.View())
		}

		status = lipgloss.JoinHorizontal(lipgloss.Top,
			spin,
			status,
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		status,
		sortOptions,
		showCompatible,
		fileCount,
	)
}

func (b Bubble) getMainButtonsView() string {
	var buttonRow string
	cfg := config.GetConfig()

	buttonStyle := lipgloss.NewStyle().
		Underline(true).
		Padding(0, 2).
		Height(internal.StatusBarHeight)

	escape := buttonStyle.
		Align(lipgloss.Left).
		Render("Esc. Home")

	refresh := buttonStyle.Render("1. Refresh")
	showCompatible := buttonStyle.Render("2. Hide incompatible mods")
	sortOrder := buttonStyle.Render("3. Sort Order")
	enterDir := buttonStyle.Render("4. Enter KSP Dir")
	download := buttonStyle.Render("5. Download mod")
	search := buttonStyle.Render("6. Search")

	if cfg.Settings.HideIncompatibleMods {
		showCompatible = buttonStyle.Render("2. Show incompatible mods")
	}

	settings := buttonStyle.
		Align(lipgloss.Right).
		Render("0. Settings")

	switch b.activeBox {
	case internal.ModInfoView, internal.ModListView:
		leftColumn := lipgloss.JoinHorizontal(lipgloss.Top,
			refresh,
			showCompatible,
			sortOrder,
			enterDir,
			download,
			search,
		)

		leftColumn = lipgloss.NewStyle().
			Width(b.width - lipgloss.Width(settings)).
			Render(leftColumn)

		buttonRow = lipgloss.JoinHorizontal(lipgloss.Top,
			leftColumn,
			settings,
		)
	case internal.SearchView:
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

		buttonRow = lipgloss.JoinHorizontal(lipgloss.Top,
			leftColumn,
			enableInput,
		)
	case internal.EnterKspDirView:
		enableInput := buttonStyle.
			Align(lipgloss.Right).
			Render("6. Enable text input")
		if b.inputRequested {
			enableInput = buttonStyle.
				Align(lipgloss.Right).
				Render("6. Disable text input")
		}

		escape = lipgloss.NewStyle().Width(b.width - lipgloss.Width(enableInput)).Render(escape)

		buttonRow = lipgloss.JoinHorizontal(lipgloss.Top,
			escape,
			enableInput,
		)
	case internal.SettingsView:
		escape = lipgloss.NewStyle().Width(b.width - lipgloss.Width(settings)).Render(escape)

		buttonRow = lipgloss.JoinHorizontal(lipgloss.Top,
			escape,
			settings,
		)
	case internal.LogView:
		buttonRow = lipgloss.NewStyle().Width(b.width - lipgloss.Width(settings)).Render(escape)
	}

	return buttonRow
}
