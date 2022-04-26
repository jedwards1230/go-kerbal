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

	switch b.activeBox {
	case internal.LogView:
		b.bubbles.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(b.theme.ActiveBoxBorderColor)
		body = lipgloss.JoinVertical(lipgloss.Top,
			b.styleTitle("Logs"),
			b.bubbles.splashViewport.View(),
		)
	case internal.EnterKspDirView:
		b.bubbles.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(b.theme.ActiveBoxBorderColor)
		body = lipgloss.JoinVertical(lipgloss.Top,
			b.styleTitle("Enter Kerbal Space Program Directory"),
			b.bubbles.splashViewport.View(),
		)
	default:
		var primaryBox string
		var secondaryBox string

		primaryTitle := b.styleTitle("Mod List")
		secondaryTitle := b.styleTitle("Help Menu")
		switch b.activeBox {
		case internal.ModListView, internal.ModInfoView:
			if b.nav.listSelected >= 0 && b.nav.listSelected < len(b.registry.ModMapIndex) {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.SettingsView:
			secondaryTitle = b.styleTitle("Options")
		case internal.SearchView:
			primaryTitle = b.styleTitle("Search Mods")
			if b.nav.listSelected >= 0 && b.nav.listSelected < len(b.registry.ModMapIndex) {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		}

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
			BorderForeground(primaryBoxBorderColor).
			Align(lipgloss.Center)
		primaryBox = lipgloss.JoinVertical(lipgloss.Top,
			primaryTitle,
			b.bubbles.primaryViewport.View(),
		)

		b.bubbles.secondaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(secondaryBoxBorderColor)
		secondaryBox = lipgloss.JoinVertical(lipgloss.Top,
			secondaryTitle,
			b.bubbles.secondaryViewport.View(),
		)

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

	return lipgloss.NewStyle().
		Width(b.bubbles.primaryViewport.Width).
		Height(b.bubbles.primaryViewport.Height - 3).
		Render(s)
}

func (b Bubble) modInfoView() string {
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
			Padding(0, 2)

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

		drawKVColor := func(k, v string, color lipgloss.AdaptiveColor) string {
			return connectSides(
				keyStyle.
					Render(k),
				valueStyle.Copy().
					Foreground(color).
					Render(v))
		}

		identifier := drawKV("Identifier", mod.Identifier)
		license := drawKV("License", mod.License)
		author := drawKV("Author", mod.Author)
		version := drawKV("Mod Version", mod.Versions.Mod)
		versionKsp := drawKV("KSP Versions", fmt.Sprintf("%v - %v", mod.Versions.KspMin, mod.Versions.KspMax))
		installed := drawKV("Installed", "Not Installed")
		if mod.Install.Installed {
			installed = drawKVColor("Installed", "Installed", b.theme.InstalledColor)
		}
		installDir := drawKV("Install dir", mod.Install.InstallTo)
		download := drawKV("Download", mod.Download.URL)
		dependencies := drawKVColor("Dependencies", "None", b.theme.Green)
		if len(mod.ModDepends) > 0 {
			dependencies = drawKVColor("Dependencies", strings.Join(mod.ModDepends, ", "), b.theme.Orange)
		}
		conflicts := drawKVColor("Conflicts", "None", b.theme.Green)
		if len(mod.ModConflicts) > 0 {
			conflicts = drawKVColor("Conflicts", strings.Join(mod.ModConflicts, ", "), b.theme.Red)
		}

		return lipgloss.JoinVertical(
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
	}
	// default to help view
	return lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width).
		Height(b.bubbles.secondaryViewport.Height - 3).
		Render(b.bubbles.help.View())
}

func (b Bubble) logView() string {
	file, err := os.Open("./logs/debug.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bodyList := make([]string, 0)
	scanner := bufio.NewScanner(file)
	//width := lipgloss.Width
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
				Foreground(b.theme.Blue).
				Width(17).
				Render(lineWords[1])
			// log output
			/* line := lipgloss.NewStyle().
			//Width(b.bubbles.splashViewport.Width - width(lineWords[0]) - width(lineWords[1])).
			Render(connectSides(idx, strings.Join(lineWords, " "))) */

			line := truncate.StringWithTail(
				connectSides(idx, strings.Join(lineWords, " ")),
				uint(b.bubbles.splashViewport.Width-2),
				internal.EllipsisStyle)

			bodyList = append(bodyList, line)
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	body := lipgloss.JoinVertical(lipgloss.Top, bodyList...)
	return lipgloss.NewStyle().
		Width(b.bubbles.splashViewport.Width - 1).
		Height(b.bubbles.splashViewport.Height - 3).
		Render(body)
}

func (b Bubble) settingsView() string {
	cfg := config.GetConfig()

	title := lipgloss.NewStyle().
		Padding(1, 5).
		Width(b.bubbles.secondaryViewport.Width).
		Align(lipgloss.Left).
		Render("Sorting")

	keyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Width((b.bubbles.secondaryViewport.Width/4)+3).
		Padding(0, 2)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(0, 2)

	drawKV := func(k, v string) string {
		return connectSides(keyStyle.Render(k), valueStyle.Render(v))
	}

	drawKVColor := func(k, v string) string {
		return connectSides(
			keyStyle.
				Render(k),
			valueStyle.
				Foreground(b.theme.UnselectedListItemColor).
				Background(b.theme.SelectedListItemColor).
				Render(v))
	}

	var lines []string
	sortOrder := drawKV("Sort Order", b.registry.SortOptions.SortOrder)
	sortBy := drawKV("Sort By", b.registry.SortOptions.SortTag)
	compat := drawKV("Hide Incompatible", fmt.Sprintf("%v", cfg.Settings.HideIncompatibleMods))

	switch b.nav.menuCursor {
	case internal.MenuSortOrder:
		sortOrder = drawKVColor("Sort Order", b.registry.SortOptions.SortOrder)
	case internal.MenuSortTag:
		sortBy = drawKVColor("Sort By", b.registry.SortOptions.SortTag)
	case internal.MenuCompatible:
		compat = drawKVColor("Hide Incompatible", fmt.Sprintf("%v", cfg.Settings.HideIncompatibleMods))
	}

	sortOrder = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(sortOrder)
	sortBy = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(sortBy)
	compat = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(compat)

	lines = append(lines, sortOrder, sortBy, compat)

	content := lipgloss.JoinVertical(lipgloss.Top, lines...)

	body := lipgloss.JoinVertical(lipgloss.Top,
		content,
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
	title := lipgloss.NewStyle().
		Padding(1, 5).
		Width(b.bubbles.secondaryViewport.Width).
		Align(lipgloss.Left).
		Render("Config")

	keyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Width(b.bubbles.secondaryViewport.Width/4).
		Padding(0, 2)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(b.bubbles.secondaryViewport.Width*3/4).
		Padding(0, 2)

	drawKV := func(k, v string) string {
		return connectSides(keyStyle.Render(k), valueStyle.Render(v))
	}

	drawKVColor := func(k, v string) string {
		return connectSides(
			keyStyle.
				Render(k),
			valueStyle.Copy().
				Foreground(b.theme.UnselectedListItemColor).
				Background(b.theme.SelectedListItemColor).
				Render(v))
	}

	var lines []string

	if b.nav.menuCursor == internal.MenuKspDir {
		lines = append(lines, drawKVColor("Kerbal Directory", cfg.Settings.KerbalDir))
	} else {
		lines = append(lines, drawKV("Kerbal Directory", cfg.Settings.KerbalDir))
	}
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

	buttonStyle := lipgloss.NewStyle().
		Underline(true).
		Padding(0, 2).
		Align(lipgloss.Right).
		//Foreground(b.theme.LightGray).
		Height(internal.StatusBarHeight)

	escape := buttonStyle.
		Align(lipgloss.Left).
		Render("Esc. Home")

	refresh := buttonStyle.Render("1. Refresh")
	search := buttonStyle.Render("2. Search")
	apply := buttonStyle.Render("3. Apply mods")
	enter := buttonStyle.Render("Ent. Install mod")
	settings := buttonStyle.Render("0. Options")

	if b.nav.activeMod.Install.Installed {
		enter = buttonStyle.Render("Ent. Remove mod")
	}

	switch b.activeBox {
	case internal.ModInfoView, internal.ModListView:
		leftColumn := connectSides(
			refresh,
			search,
			apply,
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
	switch b.activeBox {
	case internal.LogView, internal.EnterKspDirView:
		return lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Height(3).
			Border(lipgloss.NormalBorder()).
			Padding(1).
			Width(b.bubbles.splashViewport.Width).
			Render(s)
	default:
		return lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Height(3).
			Border(lipgloss.NormalBorder()).
			Padding(1).
			Width(b.bubbles.primaryViewport.Width + 2).
			Render(s)
	}
}

func connectSides(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}
