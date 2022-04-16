package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/spf13/viper"
)

// Do computations for TUI app
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Update mod list
	case ModListUpdatedMsg:
		b.nav.listSelected = 0
		b.registry.ModList = msg
		b.registry.SortModList(b.sortOptions)
		b.logs = append(b.logs, "Mod list updated")
		b.checkActiveViewPortBounds()
		b.primaryViewport.GotoTop()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
		// checks to not overwrite any required input screens
		if !b.inputRequested {
			b.activeBox = constants.PrimaryBoxActive
		}
	// Update KSP dir
	case UpdateKspDirMsg:
		if msg {
			cfg := config.GetConfig()
			log.Print("Kerbal directory updated")
			b.textInput.Reset()
			b.textInput.SetValue(fmt.Sprintf("Success!: %v", cfg.Settings.KerbalDir))
			b.textInput.Blur()
			b.inputRequested = false
		} else {
			log.Printf("Error updating ksp dir: %v", msg)
			b.textInput.Reset()
			b.textInput.Placeholder = "Try again..."
		}
		b.splashViewport.SetContent(b.inputKspView())
	case DownloadModMsg:
		if msg {
			b.logs = append(b.logs, "Mod downloaded and installed")
		} else {
			b.logs = append(b.logs, "Error downloading mod")
		}
	case ErrorMsg:
		log.Printf("Error message in update: %v", msg)
	// Window resize
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		b.help.Width = msg.Width

		b.splashViewport.Width = msg.Width - b.primaryViewport.Style.GetHorizontalFrameSize()
		b.splashViewport.Height = msg.Height - constants.StatusBarHeight - b.primaryViewport.Style.GetVerticalFrameSize()
		if b.inputRequested {
			b.splashViewport.SetContent(b.inputKspView())
		} else {
			b.splashViewport.SetContent(b.logView())
		}

		b.primaryViewport.Width = (msg.Width / 2) - b.primaryViewport.Style.GetHorizontalFrameSize()
		b.primaryViewport.Height = msg.Height - (constants.StatusBarHeight * 3) - b.primaryViewport.Style.GetVerticalFrameSize()
		b.secondaryViewport.Width = (msg.Width / 2) - b.secondaryViewport.Style.GetHorizontalFrameSize()
		b.secondaryViewport.Height = msg.Height - (constants.StatusBarHeight * 3) - b.secondaryViewport.Style.GetVerticalFrameSize()

		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())

		if !b.ready {
			b.ready = true
		}
	// Key pressed
	case tea.KeyMsg:
		cmds = append(cmds, b.handleKeys(msg))
	// Mouse input
	case tea.MouseMsg:
		// TODO: fix scrolling beyond page. breaks things.
		switch msg.Type {
		case tea.MouseWheelUp:
			b.scrollView("up")
		case tea.MouseWheelDown:
			b.scrollView("down")
		}
	}

	if b.inputRequested {
		b.activeBox = constants.SplashBoxActive
		b.textInput.Placeholder = "KSP Directory..."
		b.textInput.Focus()
		b.textInput.Reset()
		b.splashViewport.SetContent(b.inputKspView())
		cmds = append(cmds, textinput.Blink)
	}

	b.textInput, cmd = b.textInput.Update(msg)
	cmds = append(cmds, cmd)

	b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// Handles all key press events
func (b *Bubble) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	cfg := config.GetConfig()

	switch {
	// Quit
	case key.Matches(msg, b.keyMap.Quit):
		b.logs = append(b.logs, "Quitting")
		log.Print("Quitting")
		b.inputRequested = false
		return tea.Quit
	// Down
	case key.Matches(msg, b.keyMap.Down):
		b.scrollView("down")
	// Up
	case key.Matches(msg, b.keyMap.Up):
		b.scrollView("up")
	// Space
	case key.Matches(msg, b.keyMap.Space):
		if b.nav.listSelected == b.nav.listCursor {
			b.nav.listSelected = -1
		} else {
			b.nav.listSelected = b.nav.listCursor
		}
		b.checkActiveViewPortBounds()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	// Enter
	case key.Matches(msg, b.keyMap.Enter):
		if b.inputRequested {
			cmds = append(cmds, b.updateKspDirCmd(b.textInput.Value()))
		}
	// Escape
	case key.Matches(msg, b.keyMap.Esc):
		if !b.inputRequested {
			//b.inputRequested = false
			b.textInput.Reset()
			b.textInput.Blur()
			b.activeBox = constants.PrimaryBoxActive
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		}
	// Swap view
	case key.Matches(msg, b.keyMap.SwapView):
		switch b.activeBox {
		case constants.PrimaryBoxActive:
			b.activeBox = constants.SecondaryBoxActive
		case constants.SecondaryBoxActive:
			b.activeBox = constants.PrimaryBoxActive
		case constants.SplashBoxActive:
			b.activeBox = constants.PrimaryBoxActive
		}
	// Show logs
	case key.Matches(msg, b.keyMap.ShowLogs):
		if b.activeBox == constants.SplashBoxActive && !b.inputRequested {
			b.activeBox = constants.PrimaryBoxActive
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		} else {
			b.activeBox = constants.SplashBoxActive
			b.splashViewport.SetContent(b.logView())
			b.splashViewport.GotoBottom()
		}
	// Refresh list
	case key.Matches(msg, b.keyMap.RefreshList):
		b.logs = append(b.logs, "Getting mod list")
		cmds = append(cmds, b.getAvailableModsCmd())
	// Hide incompatible
	case key.Matches(msg, b.keyMap.HideIncompatible):
		b.logs = append(b.logs, "Toggling compatible mod view")
		viper.Set("settings.hide_incompatible", !cfg.Settings.HideIncompatibleMods)
		viper.WriteConfigAs(viper.ConfigFileUsed())
		cmds = append(cmds, b.getAvailableModsCmd())
	// Swap sort order
	case key.Matches(msg, b.keyMap.SwapSortOrder):
		if b.sortOptions.SortOrder == "ascend" {
			b.sortOptions.SortOrder = "descend"
		} else if b.sortOptions.SortOrder == "descend" {
			b.sortOptions.SortOrder = "ascend"
		}
		b.logs = append(b.logs, "Swapping sort order to "+b.sortOptions.SortOrder)
		log.Printf("Swapping sort order to %s", b.sortOptions.SortOrder)
		b.registry.SortModList(b.sortOptions)
		b.activeBox = constants.PrimaryBoxActive
		b.checkActiveViewPortBounds()
		b.primaryViewport.GotoTop()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	// Input KSP dir
	case key.Matches(msg, b.keyMap.EnterKspDir):
		if b.activeBox == constants.SplashBoxActive && !b.inputRequested {
			b.inputRequested = false
			b.activeBox = constants.PrimaryBoxActive
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		} else {
			b.activeBox = constants.SplashBoxActive
			b.inputRequested = true
			b.textInput.Placeholder = "KSP Directory..."
			b.textInput.Focus()
			b.textInput.Reset()
			return textinput.Blink
		}
	// Download selected mod
	case key.Matches(msg, b.keyMap.Download):
		var mod = b.registry.SortedModList[b.nav.listSelected]
		b.logs = append(b.logs, "Downloading mod")
		cmds = append(cmds, b.downloadModCmd(mod.Download))
	}

	return tea.Batch(cmds...)
}

// Handles wrapping and button
// scrolling in the viewport.
func (b *Bubble) checkActiveViewPortBounds() {
	switch b.activeBox {
	case constants.PrimaryBoxActive:
		top := b.primaryViewport.YOffset - 3
		bottom := b.primaryViewport.Height + b.primaryViewport.YOffset - 4

		if b.nav.listCursor < top {
			b.primaryViewport.LineUp(1)
		} else if b.nav.listCursor > bottom {
			b.primaryViewport.LineDown(1)
		}

		if b.nav.listCursor > len(b.registry.SortedModList)-1 {
			b.nav.listCursor = 0
			b.primaryViewport.GotoTop()
		} else if b.nav.listCursor < 0 {
			b.nav.listCursor = len(b.registry.SortedModList) - 1
			b.primaryViewport.GotoBottom()
		}
	case constants.SecondaryBoxActive:
		if b.secondaryViewport.AtBottom() {
			b.secondaryViewport.GotoBottom()
		} else if b.secondaryViewport.AtTop() {
			b.secondaryViewport.GotoTop()
		}
	case constants.SplashBoxActive:
		if b.splashViewport.AtBottom() {
			b.splashViewport.GotoBottom()
		} else if b.splashViewport.AtTop() {
			b.splashViewport.GotoTop()
		}
	}
}

// Handles mouse scrolling in the viewport
func (b *Bubble) scrollView(dir string) {
	switch dir {
	case "up":
		switch b.activeBox {
		case constants.PrimaryBoxActive:
			b.nav.listCursor--
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case constants.SecondaryBoxActive:
			b.secondaryViewport.LineUp(1)
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case constants.SplashBoxActive:
			b.splashViewport.LineUp(1)
			b.splashViewport.SetContent(b.logView())
		}
	case "down":
		switch b.activeBox {
		case constants.PrimaryBoxActive:
			b.nav.listCursor++
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case constants.SecondaryBoxActive:
			b.secondaryViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case constants.SplashBoxActive:
			b.splashViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.splashViewport.SetContent(b.logView())
		}
	default:
		log.Panic("Invalid scroll direction: " + dir)
	}
	b.checkActiveViewPortBounds()
}
