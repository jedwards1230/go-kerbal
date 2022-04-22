package tui

import (
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/spf13/viper"
)

// Do computations for TUI app
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if b.inputRequested {
		b.textInput, cmd = b.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if b.activeBox == internal.ModInfoBox {
		b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	// Update mod list
	case UpdatedModMapMsg:
		b.registry.TotalModMap = msg
		b.logs = append(b.logs, "Mod list updated")
		cmds = append(cmds, b.sortModMapCmd())
	case SortedMsg:
		b.registry.SortModMap()
		b.logs = append(b.logs, "Mod list sorted")
		log.Print("Sorted mod map")
		b.ready = true
		b.checkActiveViewPortBounds()
		b.primaryViewport.GotoTop()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case InstalledModListMsg:
		b.ready = true
		if len(b.registry.InstalledModList) != len(msg) {
			b.logs = append(b.logs, "Installed mod list updated")
			log.Printf("Updated installed mod list")
			b.registry.InstalledModList = msg
		} else {
			b.logs = append(b.logs, "No changes made")
		}
		cmds = append(cmds, b.getAvailableModsCmd())
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	// Update KSP dir
	case UpdateKspDirMsg:
		b.ready = true
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
	case SearchMsg:
		if len(msg) >= 0 {
			b.nav.listSelected = -1
			b.registry.ModMapIndex = registry.ModIndex(msg)
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		} else {
			log.Print("Error searching")
		}
	case ErrorMsg:
		b.ready = true
		log.Printf("ErrorMsg: %v", msg)
	case spinner.TickMsg:
		if !b.ready {
			if b.activeBox == internal.LogView {
				b.splashViewport.SetContent(b.logView())
			}
			b.spinner, cmd = b.spinner.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			b.spinner.Finish()
		}
	// Window resize
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		b.help.Width = msg.Width

		b.splashViewport.Width = msg.Width - b.primaryViewport.Style.GetHorizontalFrameSize()
		b.splashViewport.Height = msg.Height - internal.StatusBarHeight - b.primaryViewport.Style.GetVerticalFrameSize()
		if b.inputRequested && !b.searchInput {
			b.splashViewport.SetContent(b.inputKspView())
		} else {
			b.splashViewport.SetContent(b.logView())
		}

		b.primaryViewport.Width = (msg.Width / 2) - b.primaryViewport.Style.GetHorizontalFrameSize()
		b.primaryViewport.Height = msg.Height - (internal.StatusBarHeight * 2) - b.primaryViewport.Style.GetVerticalFrameSize()
		b.secondaryViewport.Width = (msg.Width / 2) - b.secondaryViewport.Style.GetHorizontalFrameSize()
		b.secondaryViewport.Height = msg.Height - (internal.StatusBarHeight * 2) - b.secondaryViewport.Style.GetVerticalFrameSize()

		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())

		if !b.ready {
			b.ready = true
		}
	// Key pressed
	case tea.KeyMsg:
		//log.Printf("Msg: %v %T", msg, msg)
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
		/* default:
		log.Printf("Msg: %v %T", msg, msg) */
	case MyTickMsg:
		//log.Printf("my tick: %v", msg)
		cmds = append(cmds, b.MyTickCmd())

	}

	if b.inputRequested {
		if b.searchInput {
			b.activeBox = internal.PrimaryBoxActive
			b.textInput.Focus()
			// only search when input is updated
			_, ok := msg.(tea.KeyMsg)
			if ok {
				cmds = append(cmds, b.searchCmd(b.textInput.Value()))
			}
		} else {
			b.activeBox = internal.SplashBoxActive
			b.textInput.Focus()
			b.splashViewport.SetContent(b.inputKspView())
		}
	}

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
		return tea.Quit
	// Down
	case key.Matches(msg, b.keyMap.Down):
		b.scrollView("down")
		b.inputRequested = false
		b.textInput.Blur()
	// Up
	case key.Matches(msg, b.keyMap.Up):
		b.scrollView("up")
		b.inputRequested = false
		b.textInput.Blur()
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
			if b.searchInput {
				b.textInput.Blur()
				b.inputRequested = false
				log.Printf("UPDATE: Start searching: %v", b.textInput.Value())
			} else {
				cmds = append(cmds, b.updateKspDirCmd(b.textInput.Value()))
			}
		} else {
			id := b.registry.ModMapIndex[b.nav.listCursor]
			modMap := b.registry.GetActiveModMap()
			mod := modMap[id.Key]

			// toggle mod selection
			b.nav.listSelected = b.nav.listCursor
			if b.nav.installSelected[mod.Identifier].Identifier != "" {
				delete(b.nav.installSelected, mod.Identifier)
			} else {
				b.nav.installSelected[mod.Identifier] = mod
			}

			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		}
	// Escape
	case key.Matches(msg, b.keyMap.Esc):
		if !b.inputRequested || b.searchInput {
			b.nav.listCursor = -1
			b.nav.installSelected = make(map[string]registry.Ckan, 0)
			b.textInput.Blur()
			b.textInput.Reset()
			cmds = append(cmds, b.sortModMapCmd())
			b.inputRequested = false
			b.searchInput = false
			b.activeBox = internal.PrimaryBoxActive
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		}
		b.textInput.Reset()
	// Swap view
	case key.Matches(msg, b.keyMap.SwapView):
		switch b.activeBox {
		case internal.PrimaryBoxActive:
			b.activeBox = internal.ModInfoBox
		case internal.ModInfoBox:
			b.activeBox = internal.PrimaryBoxActive
		case internal.SplashBoxActive:
			b.activeBox = internal.PrimaryBoxActive
		}
	// Show logs
	case key.Matches(msg, b.keyMap.ShowLogs):
		if b.activeBox == internal.SplashBoxActive && !b.inputRequested {
			b.activeBox = internal.PrimaryBoxActive
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		} else {
			b.activeBox = internal.SplashBoxActive
			b.splashViewport.SetContent(b.logView())
			b.splashViewport.GotoBottom()
		}
	// Refresh list
	case key.Matches(msg, b.keyMap.RefreshList):
		if !b.searchInput && !b.inputRequested {
			b.ready = false
			b.logs = append(b.logs, "Getting mod list")
			cmds = append(cmds, b.getAvailableModsCmd(), b.spinner.Tick)
		}
	// Hide incompatible
	case key.Matches(msg, b.keyMap.HideIncompatible):
		if !b.inputRequested {
			b.logs = append(b.logs, "Toggling compatible mod view")
			viper.Set("settings.hide_incompatible", !cfg.Settings.HideIncompatibleMods)
			viper.WriteConfigAs(viper.ConfigFileUsed())
			b.ready = false
			cmds = append(cmds, b.getAvailableModsCmd(), b.spinner.Tick)
		}
	// Swap sort order
	case key.Matches(msg, b.keyMap.SwapSortOrder):
		if !b.inputRequested {
			switch b.registry.SortOptions.SortOrder {
			case "ascend":
				b.registry.SortOptions.SortOrder = "descend"
			case "descend":
				b.registry.SortOptions.SortOrder = "ascend"
			}
			b.logs = append(b.logs, "Swapping sort order to "+b.registry.SortOptions.SortOrder)
			log.Printf("Swapping sort order to %s", b.registry.SortOptions.SortOrder)

			b.registry.SortModMap()
			cmds = append(cmds, b.sortModMapCmd())
			b.activeBox = internal.PrimaryBoxActive
		}
	// Input KSP dir
	// TODO: This has been hanging/acting slow. Something is wrong.
	case key.Matches(msg, b.keyMap.EnterKspDir):
		if !b.searchInput {
			if b.activeBox == internal.SplashBoxActive && !b.inputRequested {
				b.inputRequested = false
				b.activeBox = internal.PrimaryBoxActive
				b.primaryViewport.SetContent(b.modListView())
				b.secondaryViewport.SetContent(b.modInfoView())
			} else {
				b.activeBox = internal.SplashBoxActive
				b.inputRequested = true
				b.textInput.Placeholder = "KSP Directory..."
				b.textInput.Focus()
				b.textInput.Reset()
				if b.appConfig.Settings.KerbalDir != "" {
					b.textInput.SetValue(b.appConfig.Settings.KerbalDir)
				}
				return textinput.Blink
			}
		}
	// Download selected mod
	case key.Matches(msg, b.keyMap.Download):
		b.logs = append(b.logs, "Downloading mod")
		b.ready = false
		cmds = append(cmds, b.downloadModCmd(), b.spinner.Tick)
	// Search mods
	case key.Matches(msg, b.keyMap.Search):
		if b.searchInput && b.inputRequested {
			val := trimLastChar(b.textInput.Value())
			b.textInput.SetValue(val)
			b.inputRequested = false
			b.textInput.Blur()
		} else if b.searchInput && !b.inputRequested {
			b.inputRequested = true
			b.textInput.Focus()
			return textinput.Blink
		} else {
			b.searchInput = true
			b.inputRequested = true
			b.textInput.Reset()
			b.textInput.Placeholder = "Search..."
			return textinput.Blink
		}
	// View settings
	case key.Matches(msg, b.keyMap.Settings):
		if b.activeBox == internal.SplashBoxActive && !b.inputRequested {
			b.activeBox = internal.PrimaryBoxActive
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		} else {
			b.activeBox = internal.SplashBoxActive
			b.splashViewport.SetContent(b.settingsView())
			b.splashViewport.GotoTop()
		}
	}
	return tea.Batch(cmds...)
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

// Handles wrapping and button
// scrolling in the viewport.
func (b *Bubble) checkActiveViewPortBounds() {
	switch b.activeBox {
	case internal.PrimaryBoxActive:
		top := b.primaryViewport.YOffset - 3
		bottom := b.primaryViewport.Height + b.primaryViewport.YOffset - 4

		if b.nav.listCursor < top {
			b.primaryViewport.LineUp(1)
		} else if b.nav.listCursor > bottom {
			b.primaryViewport.LineDown(1)
		}

		if b.nav.listCursor > len(b.registry.ModMapIndex)-1 {
			b.nav.listCursor = 0
			b.primaryViewport.GotoTop()
		} else if b.nav.listCursor < 0 {
			b.nav.listCursor = len(b.registry.ModMapIndex) - 1
			b.primaryViewport.GotoBottom()
		}
	case internal.ModInfoBox:
		if b.secondaryViewport.AtBottom() {
			b.secondaryViewport.GotoBottom()
		} else if b.secondaryViewport.AtTop() {
			b.secondaryViewport.GotoTop()
		}
	case internal.SplashBoxActive:
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
		case internal.PrimaryBoxActive:
			b.nav.listCursor--
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.ModInfoBox:
			b.secondaryViewport.LineUp(1)
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.SplashBoxActive:
			b.splashViewport.LineUp(1)
			b.splashViewport.SetContent(b.logView())
		}
	case "down":
		switch b.activeBox {
		case internal.PrimaryBoxActive:
			b.nav.listCursor++
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.ModInfoBox:
			b.secondaryViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.SplashBoxActive:
			b.splashViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.splashViewport.SetContent(b.logView())
		}
	default:
		log.Panic("Invalid scroll direction: " + dir)
	}
	b.checkActiveViewPortBounds()
}
