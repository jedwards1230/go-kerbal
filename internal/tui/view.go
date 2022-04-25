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

func (b Bubble) View() string {
	var body string

	if b.activeBox == internal.LogView {
		b.bubbles.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(b.theme.ActiveBoxBorderColor)
		body = b.bubbles.splashViewport.View()
	} else {
		var primaryBox string
		var secondaryBox string

		// set colors
		primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
		secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor

		// format active box
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		case internal.ModInfoView, internal.EnterKspDirView, internal.SettingsView:
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
		body = connectSides(
			primaryBox,
			secondaryBox,
		)
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

	title := b.styleTitle("Mod List")
	if b.activeBox == internal.SearchView {
		title = b.styleTitle("Search Mods")
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
			internal.EllipsisStyle)

		if b.nav.listSelected == i {
			s += lipgloss.NewStyle().
				Background(b.theme.SelectedListItemColor).
				Foreground(b.theme.UnselectedListItemColor).
				Width(b.bubbles.primaryViewport.Width).
				Render(line)
		} else {
			if mod.Install.Installed {
				s += lipgloss.NewStyle().
					Foreground(b.theme.InstalledColor).
					Render(line)
			} else if !mod.IsCompatible {
				s += lipgloss.NewStyle().
					Foreground(b.theme.Orange).
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
	var title, body string

	if b.nav.listSelected >= 0 && b.nav.listSelected < len(b.registry.ModMapIndex) {
		mod := b.nav.activeMod

		keyStyle := lipgloss.NewStyle().
			Align(lipgloss.Left).
			Bold(true).
			Width(b.bubbles.secondaryViewport.Width/4).
			Padding(0, 2)

		valueStyle := lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(b.bubbles.secondaryViewport.Width*3/4).
			Padding(0, 5)

		title = b.styleTitle(mod.Name)

		abstractStyle := lipgloss.NewStyle().
			Bold(false).
			Align(lipgloss.Center).
			Width(b.bubbles.secondaryViewport.Width).
			Height(3).
			Padding(1, 2)

		abstract := abstractStyle.
			Render(mod.Abstract)

		if mod.Description != "" {
			abstract = abstractStyle.
				Render(mod.Description)
		}

		drawKV := func(k, v string) string {
			return connectSides(keyStyle.Render(k), valueStyle.Render(v))
		}

		identifier := drawKV("Identifier", mod.Identifier)
		author := drawKV("Author", mod.Author)
		version := drawKV("Mod Version", mod.Versions.Mod)
		versionKspValue := fmt.Sprintf("%v - %v", mod.Versions.KspMin, mod.Versions.KspMax)
		versionKsp := drawKV("KSP Versions", valueStyle.Render(versionKspValue))
		installDir := drawKV("Install dir", mod.Install.InstallTo)
		download := drawKV("Download", mod.Download.URL)
		license := drawKV("License", mod.License)

		installedValue := valueStyle.Render("Not Installed")
		if mod.Install.Installed {
			installedValue = valueStyle.Copy().
				Foreground(b.theme.InstalledColor).
				Render("Installed")
		}
		installed := drawKV("Installed", installedValue)

		dependenciesValue := valueStyle.Render("None")
		if len(mod.ModDepends) > 0 {
			dependenciesValue = valueStyle.Copy().
				Foreground(b.theme.Orange).
				Render(strings.Join(mod.ModDepends, ", "))
		}
		dependencies := drawKV("Dependencies", dependenciesValue)

		conflictsValue := valueStyle.Render("None")
		if len(mod.ModConflicts) > 0 {
			conflictsValue = valueStyle.Copy().
				Foreground(b.theme.Red).
				Render(strings.Join(mod.ModConflicts, ", "))
		}
		conflicts := drawKV("Conflicts", conflictsValue)

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
		title = b.styleTitle("Help Menu")
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
	width := lipgloss.Width
	i := 1
	for scanner.Scan() {
		lineWords := strings.Fields(scanner.Text())
		if len(lineWords) > 2 {
			idx := lipgloss.NewStyle().
				Align(lipgloss.Left).
				Width(5).
				Padding(0, 1).
				Render(fmt.Sprint(i) + " ")

			// timestamp
			lineWords[0] = lipgloss.NewStyle().
				Foreground(b.theme.Green).
				Render(lineWords[0])
			// file info
			lineWords[1] = lipgloss.NewStyle().
				Foreground(b.theme.Orange).
				Width(17).
				Render(lineWords[1])
			// log output
			line := lipgloss.NewStyle().
				Width(b.bubbles.splashViewport.Width - width(lineWords[0]) - width(lineWords[1])).
				Render(connectSides(idx, strings.Join(lineWords, " ")))

			bodyList = append(bodyList, line)
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	title := b.styleTitle("Logs")

	body := lipgloss.JoinVertical(lipgloss.Top, bodyList...)
	body = lipgloss.NewStyle().
		Width(b.bubbles.splashViewport.Width - 1).
		Height(b.bubbles.splashViewport.Height - 3).
		Render(string(body))

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) settingsView() string {
	cfg := config.GetConfig()

	title := b.styleTitle("Mod List")

	keyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Width(b.bubbles.secondaryViewport.Width/4).
		Padding(0, 2)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(b.bubbles.secondaryViewport.Width*3/4).
		Padding(0, 5)

	drawKV := func(k, v string) string {
		return connectSides(keyStyle.Render(k), valueStyle.Render(v))
	}

	var lines []string
	lines = append(lines, drawKV("Sort Order", b.registry.SortOptions.SortOrder))
	lines = append(lines, drawKV("Sort By", b.registry.SortOptions.SortTag))
	lines = append(lines, drawKV("Hide Compatible", fmt.Sprintf("%v", cfg.Settings.HideIncompatibleMods)))

	config := lipgloss.JoinVertical(lipgloss.Top, lines...)

	body := lipgloss.JoinVertical(lipgloss.Top,
		config,
		b.configView(),
	)

	body = lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width).
		Height(b.bubbles.secondaryViewport.Height - 3).
		Render(body)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) configView() string {
	cfg := config.GetConfig()
	title := b.styleTitle("Config")

	keyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Width(b.bubbles.secondaryViewport.Width/4).
		Padding(0, 2)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(b.bubbles.secondaryViewport.Width*3/4).
		Padding(0, 5)

	drawKV := func(k, v string) string {
		return connectSides(keyStyle.Render(k), valueStyle.Render(v))
	}

	var lines []string
	lines = append(lines, drawKV("Kerbal Directory", cfg.Settings.KerbalDir))
	lines = append(lines, drawKV("Kerbal Version", cfg.Settings.KerbalVer))
	lines = append(lines, drawKV("Logging", fmt.Sprintf("%v", cfg.Settings.EnableLogging)))
	lines = append(lines, drawKV("Mousewheel", fmt.Sprintf("%v", cfg.Settings.EnableMouseWheel)))
	lines = append(lines, drawKV("Metadata Repo", cfg.Settings.MetaRepo))
	lines = append(lines, drawKV("Last Repo Hash", cfg.Settings.LastRepoHash))
	lines = append(lines, drawKV("Theme", cfg.AppTheme))

	content := lipgloss.JoinVertical(lipgloss.Top, lines...)

	body := lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width).
		Height(b.bubbles.secondaryViewport.Height - 3).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		body,
	)
}

func (b Bubble) inputKspView() string {
	title := b.styleTitle("Kerbal Space Program Directory")

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
		Padding(0, 6, 0, 2).
		Render(fileCount)

	sortOptions := fmt.Sprintf("Sort: %s by %s", b.registry.SortOptions.SortOrder, b.registry.SortOptions.SortTag)
	sortOptions = statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Render(sortOptions)

	installedLegend := lipgloss.NewStyle().
		Foreground(b.theme.InstalledColor).
		Padding(0, 1).
		Render("Installed")
	incompatibleLegend := lipgloss.NewStyle().
		Foreground(b.theme.IncompatibleColor).
		Padding(0, 1).
		Render("Incompatible")

	colorLegend := installedLegend
	if !cfg.Settings.HideIncompatibleMods {
		colorLegend = connectSides(
			incompatibleLegend,
			installedLegend)
	}

	colorLegend = statusBarStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Render(colorLegend)

	var status string
	statusWidth := b.width - width(fileCount) - width(sortOptions) - width(colorLegend)
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
				line := connectSides(strings.Join(lineWords[2:], " "))
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
				internal.EllipsisStyle),
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

		status = connectSides(
			spin,
			status,
		)
	}

	return connectSides(
		status,
		colorLegend,
		sortOptions,
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
	apply := buttonStyle.Render("5. Apply mods")
	search := buttonStyle.Render("6. Search")
	enter := buttonStyle.Render("Ent. Install mod")

	if cfg.Settings.HideIncompatibleMods {
		showCompatible = buttonStyle.Render("2. Show incompatible mods")
	}

	if b.nav.activeMod.Install.Installed {
		enter = buttonStyle.Render("Ent. Remove mod")
	}

	settings := buttonStyle.
		Align(lipgloss.Right).
		Render("0. Settings")

	switch b.activeBox {
	case internal.ModInfoView, internal.ModListView:
		leftColumn := connectSides(
			refresh,
			showCompatible,
			sortOrder,
			enterDir,
			apply,
			search,
			enter,
		)

		leftColumn = lipgloss.NewStyle().
			Width(b.width - lipgloss.Width(settings)).
			Render(leftColumn)

		buttonRow = connectSides(
			leftColumn,
			settings,
		)
	case internal.SearchView:
		leftColumn := connectSides(
			escape,
			showCompatible,
			sortOrder,
			apply,
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

		buttonRow = connectSides(
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

		buttonRow = connectSides(
			escape,
			enableInput,
		)
	case internal.SettingsView:
		escape = lipgloss.NewStyle().Width(b.width - lipgloss.Width(settings)).Render(escape)

		buttonRow = connectSides(
			escape,
			settings,
		)
	case internal.LogView:
		buttonRow = lipgloss.NewStyle().Width(b.width - lipgloss.Width(settings)).Render(escape)
	}

	return buttonRow
}

func (b Bubble) styleTitle(s string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Height(3).
		Border(lipgloss.NormalBorder(), false, false, true).
		Padding(1)
	if b.activeBox == internal.LogView {
		return titleStyle.
			Width(b.bubbles.splashViewport.Width).
			Render(s)
	}
	return titleStyle.
		Width(b.bubbles.primaryViewport.Width).
		Render(s)
}

func connectSides(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}
