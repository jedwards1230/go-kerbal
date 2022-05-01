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
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/muesli/reflow/truncate"
)

func (b Bubble) modListView() string {
	pageStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryPaginator.Width).
		Height(b.bubbles.primaryPaginator.PerPage + 1).Render

	pagerStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryPaginator.Width).
		Align(lipgloss.Center).Render

	page := ""
	if len(b.registry.ModMapIndex) > 0 {
		start, end := b.bubbles.primaryPaginator.GetSliceBounds()
		for i, id := range b.registry.ModMapIndex[start:end] {
			mod := b.registry.SortedModMap[id.Key]

			checked := " "
			if b.registry.Queue.CheckQueue(mod.Identifier) {
				checked = "x"
			}

			line := truncate.StringWithTail(
				fmt.Sprintf("[%s] %s", checked, mod.Name),
				uint(b.bubbles.primaryPaginator.Width-2),
				internal.EllipsisStyle)

			if b.bubbles.primaryPaginator.Cursor == i && !b.nav.listCursorHide {
				page += lipgloss.NewStyle().
					Background(b.theme.SelectedListItemColor).
					Foreground(b.theme.UnselectedListItemColor).
					Width(b.bubbles.primaryPaginator.Width).
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
		pagerStyle(b.bubbles.primaryPaginator.View()),
	)

	return lipgloss.NewStyle().
		Width(b.bubbles.primaryPaginator.Width).
		Height(b.bubbles.primaryPaginator.Height - 3).
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

func (b Bubble) logView() string {
	pageStyle := lipgloss.NewStyle().
		Width(b.bubbles.splashPaginator.Width).
		Height(b.bubbles.splashPaginator.PerPage + 1).Render

	pagerStyle := lipgloss.NewStyle().
		Width(b.bubbles.splashPaginator.Width).
		Align(lipgloss.Center).Render

	var bodyList []string
	start, end := b.bubbles.splashPaginator.GetSliceBounds()
	//log.Printf("start %d, end %d", start, end)
	// todo: not happy with this approach. maybe switch to structured logging.
	for i := range b.logs[start:end] {
		lineWords := strings.Fields(b.logs[i+start])
		if len(lineWords) > 2 {
			idx := lipgloss.NewStyle().
				Align(lipgloss.Left).
				Width(5).
				Padding(0, 2, 0, 1).
				Render(fmt.Sprint(i + 1))

			// timestamp
			time := lipgloss.NewStyle().
				Foreground(b.theme.Green).
				MarginRight(1).
				Render(lineWords[0])
			// file info
			file := lipgloss.NewStyle().
				Foreground(b.theme.Blue).
				Width(20).
				MarginRight(1).
				Render(lineWords[1])
			line := strings.Join(lineWords[2:], " ")
			if b.bubbles.splashPaginator.Cursor == i {
				line = lipgloss.NewStyle().
					Background(b.theme.SelectedListItemColor).
					Foreground(b.theme.UnselectedListItemColor).
					Render(line)
			}
			// logs
			line = truncate.StringWithTail(
				connectHorz(idx, time, file, line),
				uint(b.bubbles.splashPaginator.Width-2),
				internal.EllipsisStyle)

			bodyList = append(bodyList, line)
		}
	}

	body := connectVert(bodyList...)

	body = connectVert(
		pageStyle(body),
		pagerStyle(b.bubbles.splashPaginator.View()),
	)

	return lipgloss.NewStyle().
		Width(b.bubbles.splashPaginator.Width).
		Height(b.bubbles.splashPaginator.Height).
		Render(body)
}

func (b Bubble) queueView() string {
	var content = ""
	titleStyle := lipgloss.NewStyle().
		Padding(2, 0, 1, 2).
		Bold(true).
		Width(b.bubbles.secondaryViewport.Width)

	pageStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryPaginator.Width).
		Height(b.bubbles.primaryPaginator.PerPage + 1).Render

	pagerStyle := lipgloss.NewStyle().
		Width(b.bubbles.primaryPaginator.Width).
		Align(lipgloss.Center).Render

	entryStyle := lipgloss.NewStyle().
		Width(b.bubbles.secondaryViewport.Width-1).
		Padding(0, 0, 0, 4)

	/* removeStyle := entryStyle.Copy().
	Foreground(b.theme.UnselectedListItemColor)

	downloadStyle := entryStyle.Copy().
		Foreground(b.theme.Blue)

	installStyle := entryStyle.Copy().
		Foreground(b.theme.Green) */

	if b.registry.Queue.Len() > 0 {
		selectedStyle := entryStyle.Copy().
			Foreground(b.theme.UnselectedListItemColor).
			Background(b.theme.SelectedListItemColor)

		trimName := func(s string) string {
			return truncate.StringWithTail(
				s,
				uint(b.bubbles.primaryPaginator.Width-6),
				internal.EllipsisStyle)
		}

		applyLineStyle := func(i int, mod registry.Ckan) string {
			/* if mod.Install.Installed {
				return installStyle.Render(trimName(mod.Name))
			} else if mod.Download.Downloaded {
				return downloadStyle.Render(trimName(mod.Name))
			} else if b.bubbles.paginator.GetCursorIndex() == i && !b.nav.listCursorHide {
				return selectedStyle.Render(trimName(mod.Name))
			} else {
				return entryStyle.Render(trimName(mod.Name))
			} */

			if b.bubbles.primaryPaginator.GetCursorIndex() == i && !b.nav.listCursorHide {
				return selectedStyle.Render(trimName(mod.Name))
			} else {
				return entryStyle.Render(trimName(mod.Name))
			}
		}

		var removeList, installList, dependencyList []string
		start, end := b.bubbles.primaryPaginator.GetSliceBounds()
		for i, entry := range b.registry.ModMapIndex[start:end] {
			mod := b.registry.Queue.List[entry.SearchBy][entry.Key]
			switch entry.SearchBy {
			case "remove":
				removeList = append(removeList, applyLineStyle(i, mod))

			case "install":
				installList = append(installList, applyLineStyle(i, mod))

			case "dependency":
				dependencyList = append(dependencyList, applyLineStyle(i, mod))
			}
		}

		// Display mods to remove
		if b.registry.Queue.RemoveLen() > 0 {
			removeContent := connectVert(removeList...)

			content = connectVert(
				titleStyle.Foreground(b.theme.Red).Render("To Remove"),
				removeContent,
			)
		}

		// Display mods to intall
		if b.registry.Queue.InstallLen() > 0 {
			installContent := connectVert(installList...)

			content = connectVert(
				content,
				titleStyle.Foreground(b.theme.Green).Render("To Install"),
				installContent,
			)
		}

		// Display mod dependencies to install
		if len(b.registry.Queue.GetDependencies()) > 0 {
			deoendencyContent := connectVert(dependencyList...)
			content = connectVert(
				content,
				titleStyle.Foreground(b.theme.Green).Render("Dependencies"),
				deoendencyContent,
			)
		}

		if content != "" {
			return connectVert(
				pageStyle(content),
				pagerStyle(b.bubbles.primaryPaginator.View()),
			)
		} else {
			b.LogError("Unable to parse queue")
		}
	}
	return lipgloss.NewStyle().
		Padding(2).
		Align(lipgloss.Center).
		Width(b.bubbles.primaryPaginator.Width).
		Height(b.bubbles.primaryPaginator.PerPage + 2).
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

func (b Bubble) statusBarView() string {
	cfg := config.GetConfig()
	width := lipgloss.Width

	statusBarStyle := lipgloss.NewStyle().
		Height(internal.StatusBarHeight)

	fileCount := fmt.Sprintf("Mod: %d/%d", b.bubbles.primaryPaginator.GetCursorIndex()+1, len(b.registry.ModMapIndex))
	if b.activeBox == internal.LogView {
		fileCount = fmt.Sprintf("Mod: %d/%d", b.bubbles.splashPaginator.GetCursorIndex()+1, len(b.registry.ModMapIndex))
	}
	fileCount = statusBarStyle.
		Align(lipgloss.Right).
		PaddingRight(6).
		PaddingLeft(2).
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
		file, err := os.Open(internal.LogPath)
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
			if b.nav.listCursorHide {
				return commandStyle.Render(b.getBoolOptionsView())
			} else {
				return commandStyle.Render(b.helpView())
			}
		} else {
			loading := lipgloss.NewStyle().
				Bold(true).
				Align(lipgloss.Center).
				Width(b.bubbles.commandViewport.Width).
				Padding(1, 2, 0).
				Render("Loading...")

			return lipgloss.NewStyle().
				Width(b.bubbles.commandViewport.Width).
				Height(b.bubbles.commandViewport.Height).
				Render(loading)
		}
	default:
		return commandStyle.Render(b.helpView())
	}
}
