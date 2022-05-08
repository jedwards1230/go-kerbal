package tui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/ckan"
	"github.com/jedwards1230/go-kerbal/internal/common"
	"github.com/jedwards1230/go-kerbal/internal/config"
	"github.com/jedwards1230/go-kerbal/internal/style"
	"github.com/jedwards1230/go-kerbal/internal/theme"
)

func (b Bubble) modListView() string {
	pageStyle := styleWidth(b.bubbles.primaryPaginator.Width).
		Height(b.bubbles.primaryPaginator.PerPage + 1).Render

	pagerStyle := styleWidth(b.bubbles.primaryPaginator.Width).
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

			line := trunc(fmt.Sprintf("[%s] %s", checked, mod.Name), b.bubbles.primaryPaginator.Width-2)

			if b.bubbles.primaryPaginator.Cursor == i && !b.nav.listCursorHide {
				page += style.ListSelected.
					Width(b.bubbles.primaryPaginator.Width).
					Render(line)
			} else if mod.Installed() {
				page += style.Installed.Render(line)
			} else if !mod.IsCompatible {
				page += style.Incompatible.Render(line)
			} else {
				page += line
			}
			page += "\n"
		}
	}

	page = connectVert(
		pageStyle(page),
		pagerStyle(b.bubbles.primaryPaginator.View()),
	)

	return styleWidth(b.bubbles.primaryPaginator.Width).
		Height(b.bubbles.primaryPaginator.Height - 3).
		Render(page)
}

func (b Bubble) modInfoView() string {
	if !b.nav.listCursorHide {
		mod := b.nav.activeMod

		keyStyle := style.KeyStyle.Width(b.bubbles.secondaryViewport.Width / 4)

		valueStyle := style.ValueStyle.Width(b.bubbles.secondaryViewport.Width * 3 / 4)

		abstractStyle := style.AbstractInfo.Width(b.bubbles.secondaryViewport.Width)

		abstract := abstractStyle.Render(mod.Abstract)

		if mod.Description != "" {
			abstract = abstractStyle.Render(mod.Description)
		}

		drawKV := func(k, v string) string {
			return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
		}

		drawKVColor := func(k, v string, color lipgloss.AdaptiveColor) string {
			return connectHorz(
				keyStyle.Render(k),
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
		if mod.Installed() {
			installed = drawKVColor("Installed", "Installed", theme.AppTheme.InstalledColor)
		}
		installDir := drawKV("Install dir", mod.Install.InstallTo)
		download := drawKV("Download", mod.Download.URL)
		dependencies := drawKVColor("Dependencies", "None", theme.AppTheme.Green)
		if len(mod.ModDepends) > 0 {
			dependencies = drawKVColor("Dependencies", strings.Join(mod.ModDepends, ", "), theme.AppTheme.Orange)
		}
		conflicts := drawKVColor("Conflicts", "None", theme.AppTheme.Green)
		if len(mod.ModConflicts) > 0 {
			conflicts = drawKVColor("Conflicts", strings.Join(mod.ModConflicts, ", "), theme.AppTheme.Red)
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
	return styleWidth(b.bubbles.secondaryViewport.Width).
		Height(b.bubbles.secondaryViewport.Height - 3).
		Render(b.homeView())
}

func (b Bubble) logView() string {
	pageStyle := styleWidth(b.bubbles.splashPaginator.Width).
		Height(b.bubbles.splashPaginator.PerPage + 1).Render

	pagerStyle := styleWidth(b.bubbles.splashPaginator.Width).
		Align(lipgloss.Center).Render

	var bodyList []string
	start, end := b.bubbles.splashPaginator.GetSliceBounds()
	//log.Printf("start %d, end %d", start, end)
	// todo: not happy with this approach. maybe switch to structured logging.
	for i := range b.logs[start:end] {
		lineWords := strings.Fields(b.logs[i+start])
		if len(lineWords) > 2 {
			idx := style.LogLineNumber.Render(fmt.Sprint(i + 1))

			// timestamp
			time := style.LogTimestamp.Render(lineWords[0])

			// file info
			file := strings.ReplaceAll(lineWords[1], ".go", "")
			file = style.LogFile.Render(file)
			line := strings.Join(lineWords[2:], " ")
			if b.bubbles.splashPaginator.Cursor == i {
				line = style.ListSelected.
					Render(line)
			}
			// logs
			line = trunc(
				connectHorz(idx, time, file, line),
				b.bubbles.splashPaginator.Width-2,
			)

			bodyList = append(bodyList, line)
		}
	}

	body := connectVert(bodyList...)

	body = connectVert(
		pageStyle(body),
		pagerStyle(b.bubbles.splashPaginator.View()),
	)

	return styleWidth(b.bubbles.splashPaginator.Width).
		Height(b.bubbles.splashPaginator.Height).
		Render(body)
}

func (b Bubble) queueView() string {
	var content = ""
	titleStyle := style.QueueTitle.
		Width(b.bubbles.secondaryViewport.Width)

	pageStyle := styleWidth(b.bubbles.primaryPaginator.Width).
		Height(b.bubbles.primaryPaginator.PerPage + 1).Render

	pagerStyle := styleWidth(b.bubbles.primaryPaginator.Width).
		Align(lipgloss.Center).Render

	entryStyle := styleWidth(b.bubbles.secondaryViewport.Width-1).
		Padding(0, 0, 0, 4)

	/* removeStyle := entryStyle.Copy().
	Foreground(theme.AppTheme.UnselectedListItemColor)

	downloadStyle := entryStyle.Copy().
		Foreground(theme.AppTheme.Blue) */

	installStyle := entryStyle.Copy().
		Foreground(theme.AppTheme.Green)

	if b.registry.Queue.Len() > 0 {
		selectedStyle := entryStyle.Copy().
			Foreground(theme.AppTheme.UnselectedListItemColor).
			Background(theme.AppTheme.SelectedListItemColor)

		trimName := func(s string) string {
			return trunc(
				s,
				b.bubbles.primaryPaginator.Width-6,
			)
		}

		applyLineStyle := func(i int, mod ckan.Ckan) string {
			if b.bubbles.primaryPaginator.GetCursorIndex() == i && !b.nav.listCursorHide {
				return selectedStyle.Render(trimName(mod.Name))
			} else if mod.Installed() {
				return installStyle.Render(trimName(mod.Name))
			} else {
				return entryStyle.Render(trimName(mod.Name))
			}
		}

		removeLineStyle := func(i int, mod ckan.Ckan) string {
			/* if mod.Installed() {
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
				removeList = append(removeList, removeLineStyle(i, mod))

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
				titleStyle.Foreground(theme.AppTheme.Red).Render("To Remove"),
				removeContent,
			)
		}

		// Display mods to intall
		if b.registry.Queue.InstallLen() > 0 {
			installContent := connectVert(installList...)

			content = connectVert(
				content,
				titleStyle.Foreground(theme.AppTheme.Green).Render("To Install"),
				installContent,
			)
		}

		// Display mod dependencies to install
		if len(b.registry.Queue.GetDependencies()) > 0 {
			deoendencyContent := connectVert(dependencyList...)
			content = connectVert(
				content,
				titleStyle.Foreground(theme.AppTheme.Green).Render("Dependencies"),
				deoendencyContent,
			)
		}

		if content != "" {
			return connectVert(
				pageStyle(content),
				pagerStyle(b.bubbles.primaryPaginator.View()),
			)
		} else {
			common.LogError("Unable to parse queue")
		}
	}
	return styleWidth(b.bubbles.primaryPaginator.Width).
		Padding(2).
		Align(lipgloss.Center).
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

	kspDir := trunc(cfg.Settings.KerbalDir, (b.bubbles.secondaryViewport.Width*2/3)-3)

	if b.nav.menuCursor == internal.MenuKspDir {
		configLines = append(configLines, b.drawKV("Kerbal Directory", kspDir, true))
	} else {
		configLines = append(configLines, b.drawKV("Kerbal Directory", kspDir, false))
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

	sortOrder = styleWidth(b.bubbles.secondaryViewport.Width).Render(sortOrder)
	sortBy = styleWidth(b.bubbles.secondaryViewport.Width).Render(sortBy)
	compat = styleWidth(b.bubbles.secondaryViewport.Width).Render(compat)

	sortLines = append(sortLines, sortOrder, sortBy, compat)

	for i := range sortLines {
		sortLines[i] = styleWidth(b.bubbles.secondaryViewport.Width).Render(sortLines[i])
	}

	for i := range configLines {
		configLines[i] = styleWidth(b.bubbles.secondaryViewport.Width).Render(configLines[i])
	}

	sortContent := connectVert(sortLines...)

	sortContent = styleWidth(b.bubbles.secondaryViewport.Width).
		Render(sortContent)

	sortOptions := connectVert(
		titleStyle.Render("Sorting"),
		sortContent,
	)

	configContent := connectVert(configLines...)

	configContent = styleWidth(b.bubbles.secondaryViewport.Width).
		Render(configContent)

	configOptions := connectVert(
		titleStyle.Render("Config"),
		configContent,
	)

	body := connectVert(
		sortOptions,
		configOptions,
	)

	body = styleWidth(b.bubbles.secondaryViewport.Width).
		Height(b.bubbles.secondaryViewport.Height - 3).
		Render(body)

	return connectVert(
		body,
	)
}

func (b Bubble) statusBarView() string {
	cfg := config.GetConfig()
	width := lipgloss.Width

	fileCount := fmt.Sprintf("Mod: %d/%d", b.bubbles.primaryPaginator.GetCursorIndex()+1, len(b.registry.ModMapIndex))
	if b.activeBox == internal.LogView {
		fileCount = fmt.Sprintf("Mod: %d/%d", b.bubbles.splashPaginator.GetCursorIndex()+1, len(b.registry.ModMapIndex))
	}
	fileCount = style.StatusBar.Copy().
		Align(lipgloss.Right).
		PaddingRight(6).
		PaddingLeft(2).
		Width(22).
		Render(fileCount)

	sortOptions := fmt.Sprintf("Sort: %s by %s", b.registry.SortOptions.SortOrder, b.registry.SortOptions.SortTag)
	sortOptions = style.StatusBar.Copy().
		Align(lipgloss.Right).
		Padding(0, 1).
		Width(22).
		Render(sortOptions)

	installedLegend := lipgloss.NewStyle().
		Foreground(theme.AppTheme.InstalledColor).
		Padding(0, 1).
		Render("Installed")
	incompatibleLegend := lipgloss.NewStyle().
		Foreground(theme.AppTheme.IncompatibleColor).
		Padding(0, 1).
		Render("Incompatible")

	colorLegend := installedLegend
	if !cfg.Settings.HideIncompatibleMods {
		colorLegend = connectHorz(
			incompatibleLegend,
			installedLegend)
	}

	colorLegend = style.StatusBar.Copy().
		Align(lipgloss.Right).
		Padding(0, 1).
		Width(30).
		Render(colorLegend)

	var status string
	statusWidth := b.width - width(fileCount) - width(sortOptions) - width(colorLegend) - 1
	if b.searchInput {
		status = style.StatusBar.Copy().
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
		status = style.StatusBar.
			Align(lipgloss.Left).
			Padding(0, 1).
			Width(statusWidth).
			Render(trunc(
				status,
				statusWidth-3,
			))

		spin := styleWidth(b.bubbles.spinner.Style.GetWidth()).
			Padding(0, 0, 0, 1).
			Render("  ")

		if !b.ready {
			spin = styleWidth(b.bubbles.spinner.Style.GetWidth()).
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
			loading := styleWidth(b.bubbles.commandViewport.Width).
				Bold(true).
				Align(lipgloss.Center).
				Padding(1, 2, 0).
				Render("Loading...")

			return styleWidth(b.bubbles.commandViewport.Width).
				Height(b.bubbles.commandViewport.Height).
				Render(loading)
		}
	default:
		return commandStyle.Render(b.helpView())
	}
}
