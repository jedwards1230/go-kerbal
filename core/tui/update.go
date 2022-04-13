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
	case ModListUpdatedMsg:
		b.registry.ModList = msg
		b.registry.SortModList(b.sortOptions)
		b.logs = append(b.logs, "Mod list updated")
		b.activeBox = constants.PrimaryBoxActive
		b.checkActiveViewPortBounds()
		b.primaryViewport.GotoTop()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case UpdateKspDirMsg:
		if msg {
			log.Print("Kerbal directory updated")
			b.textInput.Reset()
			b.textInput.SetValue(fmt.Sprintf("Success!: %v", msg))
			b.textInput.Blur()
			b.inputRequested = false
		} else {
			log.Printf("Error updating ksp dir: %v", msg)
			b.textInput.Reset()
			b.textInput.Placeholder = "Try again..."
		}
		b.splashViewport.SetContent(b.inputKspView())
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
		b.primaryViewport.Height = msg.Height - constants.StatusBarHeight - b.primaryViewport.Style.GetVerticalFrameSize()
		b.secondaryViewport.Width = (msg.Width / 2) - b.secondaryViewport.Style.GetHorizontalFrameSize()
		b.secondaryViewport.Height = msg.Height - constants.StatusBarHeight - b.secondaryViewport.Style.GetVerticalFrameSize()

		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())

		if !b.ready {
			b.ready = true
		}
	case tea.KeyMsg:
		cmds = append(cmds, b.handleKeys(msg))
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
		b.splashViewport.SetContent(b.inputKspView())
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
	case key.Matches(msg, b.keyMap.Quit):
		b.logs = append(b.logs, "Quitting")
		return tea.Quit
	case key.Matches(msg, b.keyMap.Down):
		b.scrollView("down")
	case key.Matches(msg, b.keyMap.Up):
		b.scrollView("up")
	case key.Matches(msg, b.keyMap.Space):
		if b.selected == b.cursor {
			b.selected = -1
		} else {
			b.selected = b.cursor
		}
		b.checkActiveViewPortBounds()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case key.Matches(msg, b.keyMap.Enter):
		if b.inputRequested {
			cmds = append(cmds, b.updateKspDirCmd(b.textInput.Value()))
		}
	case key.Matches(msg, b.keyMap.Esc):
		b.inputRequested = false
		b.textInput.Reset()
		b.textInput.Blur()
		b.activeBox = constants.PrimaryBoxActive
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case key.Matches(msg, b.keyMap.SwapView):
		switch b.activeBox {
		case constants.PrimaryBoxActive:
			b.activeBox = constants.SecondaryBoxActive
		case constants.SecondaryBoxActive:
			b.activeBox = constants.PrimaryBoxActive
		case constants.SplashBoxActive:
			b.activeBox = constants.PrimaryBoxActive
		}
	case key.Matches(msg, b.keyMap.ShowLogs):
		if b.activeBox == constants.SplashBoxActive && !b.inputRequested {
			b.activeBox = constants.PrimaryBoxActive
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		} else {
			b.inputRequested = false
			b.activeBox = constants.SplashBoxActive
			b.splashViewport.SetContent(b.logView())
			b.splashViewport.GotoBottom()
		}
	case key.Matches(msg, b.keyMap.RefreshList):
		b.logs = append(b.logs, "Getting mod list")
		cmds = append(cmds, b.getAvailableModsCmd())
	case key.Matches(msg, b.keyMap.HideIncompatible):
		b.logs = append(b.logs, "Toggling compatible mod view")
		viper.Set("settings.hide_incompatible", !cfg.Settings.HideIncompatibleMods)
		viper.WriteConfigAs(viper.ConfigFileUsed())
		cmds = append(cmds, b.getAvailableModsCmd())
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
	case key.Matches(msg, b.keyMap.EnterKspDir):
		if b.activeBox == constants.SplashBoxActive && b.inputRequested {
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

		if b.cursor < top {
			b.primaryViewport.LineUp(1)
		} else if b.cursor > bottom {
			b.primaryViewport.LineDown(1)
		}

		if b.cursor > len(b.registry.SortedModList)-1 {
			b.cursor = 0
			b.primaryViewport.GotoTop()
		} else if b.cursor < 0 {
			b.cursor = len(b.registry.SortedModList) - 1
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
			b.cursor--
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
			b.cursor++
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
