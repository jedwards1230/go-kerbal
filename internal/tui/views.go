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

func (b Bubble) modListView() string {
	modMap := b.registry.GetActiveModList()

	pageStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryViewport.Width).
		Height(b.bubbles.paginator.PerPage + 1).Render

	pagerStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryViewport.Width).
		Align(lipgloss.Center).Render

	page := ""
	if len(b.registry.ModMapIndex) > 0 {
		start, end := b.bubbles.paginator.GetSliceBounds(len(b.registry.ModMapIndex))
		//log.Printf("cursor: %v hide: %v start: %v end: %v perpage: %v", b.bubbles.paginator.Cursor, b.nav.listCursorHide, start, end, b.bubbles.paginator.PerPage)
		for i, id := range b.registry.ModMapIndex[start:end] {
			mod := modMap[id.Key]

			checked := " "
			if b.nav.installSelected[mod.Identifier].Identifier != "" {
				checked = "x"
			}

			line := truncate.StringWithTail(
				fmt.Sprintf("[%s] %s", checked, mod.Name),
				uint(b.bubbles.primaryViewport.Width-2),
				internal.EllipsisStyle)

			if b.bubbles.paginator.Cursor == i && !b.nav.listCursorHide {
				page += lipgloss.NewStyle().
					Background(b.theme.SelectedListItemColor).
					Foreground(b.theme.UnselectedListItemColor).
					Width(b.bubbles.primaryViewport.Width).
					Render(line)
			} else if mod.Install.Installed {
				page += lipgloss.NewStyle().
					Foreground(b.theme.InstalledColor).
					Render(line)
			} else if !mod.IsCompatible {
				page += lipgloss.NewStyle().
					Foreground(b.theme.Orange).
					Render(line)
			} else {
				page += lipgloss.NewStyle().
					Render(line)
			}
			page += "\n"
		}
	}

	page = connectVert(
		pageStyle(page),
		pagerStyle(b.bubbles.paginator.View()),
	)

	return lipgloss.NewStyle().
		Width(b.bubbles.primaryViewport.Width).
		Height(b.bubbles.primaryViewport.Height - 3).
		Render(page)
}

func (b Bubble) modInfoView() string {
	if !b.nav.listCursorHide {
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
			return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
		}

		drawKVColor := func(k, v string, color lipgloss.AdaptiveColor) string {
			return connectHorz(
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

		return connectVert(
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
	// default to home view
	return lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width).
		Height(b.bubbles.secondaryViewport.Height - 3).
		Render(b.homeView())
}

func (b Bubble) homeView() string {
	contentStyle := lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width - 5).
		Padding(2).
		Render

	content := "" +
		"To do:\n " +
		"- Display error messages\n " +
		"- Make install queue editable\n " +
		"- Window resizing on Windows\n " +
		"- Better mod info formatting\n " +
		"- Pagination\n " +
		"- Ensure mods install/uninstall properly\n " +
		"- Dynamic command view\n "

	return contentStyle(content)
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
			// logs
			line := truncate.StringWithTail(
				connectHorz(idx, strings.Join(lineWords, " ")),
				uint(b.bubbles.splashViewport.Width-2),
				internal.EllipsisStyle)

			bodyList = append(bodyList, line)
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	body := connectVert(bodyList...)
	return lipgloss.NewStyle().
		Width(b.bubbles.splashViewport.Width).
		Height(b.bubbles.splashViewport.Height - 3).
		Render(body)
}

func (b Bubble) queueView() string {
	var content = ""
	titleStyle := lipgloss.NewStyle().
		Padding(2, 0, 1, 2).
		Bold(true).
		Width(b.bubbles.secondaryViewport.Width)

	pageStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryViewport.Width).
		Height(b.bubbles.paginator.PerPage + 1).Render

	pagerStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryViewport.Width).
		Align(lipgloss.Center).Render

	entryStyle := lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width-1).
		Padding(0, 0, 0, 4)
		/*
			removeStyle := entryStyle.Copy().
				Foreground(b.theme.UnselectedListItemColor)
		*/
	downloadStyle := entryStyle.Copy().
		Foreground(b.theme.Blue)

	installStyle := entryStyle.Copy().
		Foreground(b.theme.Green)

	if len(b.nav.installSelected) > 0 {
		selectedStyle := entryStyle.Copy().
			Foreground(b.theme.UnselectedListItemColor).
			Background(b.theme.SelectedListItemColor)

		trimName := func(s string) string {
			return truncate.StringWithTail(
				s,
				uint(b.bubbles.primaryViewport.Width-6),
				internal.EllipsisStyle)
		}

		var removeList, installList, dependencyList []string
		start, end := b.bubbles.paginator.GetSliceBounds(len(b.registry.ModMapIndex))
		//log.Printf("cursor: %v hide: %v start: %v end: %v perpage: %v", b.bubbles.paginator.Cursor, b.nav.listCursorHide, start, end, b.bubbles.paginator.PerPage)
		for i, entry := range b.registry.ModMapIndex[start:end] {
			typeMap := b.registry.Queue[entry.SearchBy]
			switch entry.SearchBy {
			case "remove":
				mod := typeMap[entry.Key]
				if b.bubbles.paginator.Cursor == i && !b.nav.listCursorHide {
					removeList = append(removeList, selectedStyle.Render(trimName(mod.Name)))
				} else {
					removeList = append(removeList, entryStyle.Render(trimName(mod.Name)))
				}
			case "install":
				mod := typeMap[entry.Key]
				if mod.Install.Installed {
					installList = append(installList, installStyle.Render(trimName(mod.Name)))
				} else if mod.Download.Downloaded {
					installList = append(installList, downloadStyle.Render(trimName(mod.Name)))
				} else if b.bubbles.paginator.Cursor-len(b.registry.Queue["remove"]) == i && !b.nav.listCursorHide {
					installList = append(installList, selectedStyle.Render(trimName(mod.Name)))
				} else {
					installList = append(installList, entryStyle.Render(trimName(mod.Name)))
				}
			case "dependency":
				mod := typeMap[entry.Key]
				if mod.Install.Installed {
					dependencyList = append(dependencyList, installStyle.Render(trimName(mod.Name)))
				} else if mod.Download.Downloaded {
					dependencyList = append(dependencyList, downloadStyle.Render(trimName(mod.Name)))
				} else if b.bubbles.paginator.Cursor-len(b.registry.Queue["remove"])-len(b.registry.Queue["install"]) == i && !b.nav.listCursorHide {
					dependencyList = append(dependencyList, selectedStyle.Render(trimName(mod.Name)))
				} else {
					dependencyList = append(dependencyList, entryStyle.Render(trimName(mod.Name)))
				}
			}
		}

		// Display mods to remove
		if len(b.registry.Queue["remove"]) > 0 {
			removeContent := connectVert(removeList...)
			content = connectVert(
				titleStyle.Foreground(b.theme.Red).Render("To Remove"),
				removeContent,
			)
		}

		// Display mods to intall
		if len(b.registry.Queue["install"]) > 0 {
			installContent := connectVert(installList...)

			content = connectVert(
				content,
				titleStyle.Foreground(b.theme.Green).Render("To Install"),
				installContent,
			)
		}

		// Display mod dependencies to install
		if len(b.registry.Queue["dependency"]) > 0 {
			installContent := connectVert(dependencyList...)
			content = connectVert(
				content,
				titleStyle.Foreground(b.theme.Green).Render("Dependencies"),
				installContent,
			)
		}

		if content != "" {
			return connectVert(
				pageStyle(content),
				pagerStyle(b.bubbles.paginator.View()),
			)
		} else {
			b.LogErrorf("Unable to parse queue: %v", b.nav.installSelected)
		}
	}
	return lipgloss.NewStyle().
		Padding(2).
		Align(lipgloss.Center).
		Width(b.bubbles.primaryViewport.Width).
		Height(b.bubbles.primaryViewport.Height - 3).
		Render("No mods in queue")
}

func (b Bubble) settingsView() string {
	cfg := config.GetConfig()

	titleStyle := lipgloss.NewStyle().
		Padding(1, 5).
		Bold(true).
		Width(b.bubbles.secondaryViewport.Width).
		Align(lipgloss.Left)

	var sortLines []string
	var configLines []string

	sortOrder := b.drawKV("Sort Order", b.registry.SortOptions.SortOrder, false)
	sortBy := b.drawKV("Sort By", b.registry.SortOptions.SortTag, false)
	compat := b.drawKV("Hide Incompatible", fmt.Sprintf("%v", cfg.Settings.HideIncompatibleMods), false)

	if b.nav.menuCursor == internal.MenuKspDir {
		configLines = append(configLines, b.drawKV("Kerbal Directory", cfg.Settings.KerbalDir, true))
	} else {
		configLines = append(configLines, b.drawKV("Kerbal Directory", cfg.Settings.KerbalDir, false))
	}

	configLines = append(configLines, b.drawKV("Kerbal Version", cfg.Settings.KerbalVer, false))
	configLines = append(configLines, b.drawKV("Logging", fmt.Sprintf("%v", cfg.Settings.EnableLogging), false))
	configLines = append(configLines, b.drawKV("Mousewheel", fmt.Sprintf("%v", cfg.Settings.EnableMouseWheel), false))
	configLines = append(configLines, b.drawKV("Metadata Repo", cfg.Settings.MetaRepo, false))
	configLines = append(configLines, b.drawKV("Last Repo Hash", cfg.Settings.LastRepoHash, false))
	configLines = append(configLines, b.drawKV("Theme", cfg.AppTheme, false))

	switch b.nav.menuCursor {
	case internal.MenuSortOrder:
		sortOrder = b.drawKV("Sort Order", b.registry.SortOptions.SortOrder, true)
	case internal.MenuSortTag:
		sortBy = b.drawKV("Sort By", b.registry.SortOptions.SortTag, true)
	case internal.MenuCompatible:
		compat = b.drawKV("Hide Incompatible", fmt.Sprintf("%v", cfg.Settings.HideIncompatibleMods), true)
	}

	sortOrder = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(sortOrder)
	sortBy = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(sortBy)
	compat = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(compat)

	sortLines = append(sortLines, sortOrder, sortBy, compat)

	for i := range sortLines {
		sortLines[i] = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(sortLines[i])
	}

	for i := range configLines {
		configLines[i] = lipgloss.NewStyle().Width(b.bubbles.secondaryViewport.Width).Render(configLines[i])
	}

	sortContent := connectVert(sortLines...)

	sortContent = lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width).
		Render(sortContent)

	sortOptions := connectVert(
		titleStyle.Render("Sorting"),
		sortContent,
	)

	configContent := connectVert(configLines...)

	configContent = lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width).
		Render(configContent)

	configOptions := connectVert(
		titleStyle.Render("Config"),
		configContent,
	)

	body := connectVert(
		sortOptions,
		configOptions,
	)

	body = lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width).
		Height(b.bubbles.secondaryViewport.Height - 3).
		Render(body)

	return connectVert(
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

	return connectVert(
		question,
		inText,
	)
}

func (b Bubble) statusBarView() string {
	cfg := config.GetConfig()
	width := lipgloss.Width

	statusBarStyle := lipgloss.NewStyle().
		Height(internal.StatusBarHeight)

	fileCount := fmt.Sprintf("Mod: %d/%d", b.bubbles.paginator.GetCursorIndex()+1, len(b.registry.ModMapIndex))
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
		colorLegend = connectHorz(
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
				line := connectHorz(strings.Join(lineWords[2:], " "))
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

		status = connectHorz(
			spin,
			status,
		)
	}

	return connectHorz(
		status,
		colorLegend,
		sortOptions,
		fileCount,
	)
}

func (b Bubble) commandView() string {
	commandStyle := lipgloss.NewStyle().
		Align(lipgloss.Center)

	switch b.activeBox {
	case internal.QueueView:
		if b.ready {
			options := commandStyle.Render(b.getBoolOptionsView())
			return lipgloss.NewStyle().
				Width(b.bubbles.commandViewport.Width).
				Height(b.bubbles.commandViewport.Height).
				Render(options)
		} else {
			loading := lipgloss.NewStyle().
				Bold(true).
				Align(lipgloss.Center).
				Width(b.bubbles.commandViewport.Width-4).
				Padding(1, 2, 0).
				Render("Loading...")

			return lipgloss.NewStyle().
				Width(b.bubbles.commandViewport.Width).
				Height(b.bubbles.commandViewport.Height).
				Render(loading)
		}
	}

	return commandStyle.Render(b.helpView())
}

/* func (b Bubble) getMainButtonsView() string {
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
		leftColumn := connectHorz(
			refresh,
			search,
			apply,
			enter,
		)

		leftColumn = lipgloss.NewStyle().
			Width(b.width - lipgloss.Width(settings)).
			Render(leftColumn)

		buttonRow = connectHorz(
			leftColumn,
			settings,
		)
	case internal.SearchView:
		leftColumn := connectHorz(
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

		buttonRow = connectHorz(
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

		buttonRow = connectHorz(
			escape,
			enableInput,
		)
	case internal.SettingsView:
		escape = lipgloss.NewStyle().Width(b.width - lipgloss.Width(settings)).Render(escape)

		buttonRow = connectHorz(
			escape,
			settings,
		)
	case internal.LogView:
		buttonRow = lipgloss.NewStyle().Width(b.width - lipgloss.Width(settings)).Render(escape)
	}

	return buttonRow
} */

func (b Bubble) helpView() string {
	/* titleStyle := lipgloss.NewStyle().
	Bold(true).
	Align(lipgloss.Left).
	Width(b.bubbles.commandViewport.Width).
	Padding(1, 2, 0) */

	leftColumn := []string{
		b.drawHelpKV("up", "Move up"),
		b.drawHelpKV("down", "Move down"),
		b.drawHelpKV("space", "Toggle mod info"),
		b.drawHelpKV("enter", "Add to queue"),
		b.drawHelpKV("tab", "Swap active windows"),
	}

	rightColumn := []string{
		b.drawHelpKV("1", "Refresh"),
		b.drawHelpKV("2", "Search"),
		b.drawHelpKV("3", "Apply"),
		b.drawHelpKV("0", "Settings"),
		b.drawHelpKV("shift+o", "Logs"),
	}

	content := connectHorz(connectVert(leftColumn...), connectVert(rightColumn...))

	content = lipgloss.NewStyle().
		Padding(1).
		Margin(1, 0).
		Render(content)

	return lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width).
		Render(content)
}

func (b Bubble) getBoolOptionsView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.bubbles.commandViewport.Width-4).
		Padding(1, 2, 0)

	optionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		BorderForeground(b.theme.InactiveBoxBorderColor).
		Padding(0, 4).
		Faint(true).
		Margin(1, 1)

	cancel := optionStyle.Render("Cancel")
	confirm := optionStyle.Render("Confirm")

	if b.nav.listCursorHide {
		if b.nav.boolCursor {
			confirm = optionStyle.Copy().
				BorderForeground(b.theme.ActiveBoxBorderColor).
				Border(lipgloss.RoundedBorder()).
				Faint(false).
				Render("Confirm")
			cancel = optionStyle.Copy().
				Render("Cancel")
		} else {
			cancel = optionStyle.Copy().
				BorderForeground(b.theme.ActiveBoxBorderColor).
				Border(lipgloss.RoundedBorder()).
				Faint(false).
				Render("Cancel")
			confirm = optionStyle.Copy().
				Render("Confirm")
		}
	}

	options := connectHorz(cancel, "  ", confirm)
	options = lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width - 4).
		Align(lipgloss.Center).
		Render(options)

	content := connectVert(
		titleStyle.Render("Apply?"),
		options,
	)

	return lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width).
		//Padding(3).
		Render(content)
}
