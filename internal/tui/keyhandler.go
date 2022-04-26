package tui

import (
	"log"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/spf13/viper"
)

// Handles all key press events
func (b *Bubble) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	switch {
	// Quit
	case key.Matches(msg, b.keyMap.Quit):
		log.Print("Quitting")
		return tea.Quit
	// Down
	case key.Matches(msg, b.keyMap.Down):
		b.scrollView("down")
		b.inputRequested = false
	// Up
	case key.Matches(msg, b.keyMap.Up):
		b.scrollView("up")
		b.inputRequested = false
	// Space
	case key.Matches(msg, b.keyMap.Space):
		if b.nav.listSelected == b.nav.listCursor {
			b.nav.listSelected = -1
		} else {
			b.nav.listSelected = b.nav.listCursor
		}
	// Enter
	case key.Matches(msg, b.keyMap.Enter):
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			b.toggleSelectedItem()
		case internal.EnterKspDirView:
			cmds = append(cmds, b.updateKspDirCmd(b.bubbles.textInput.Value()))
		case internal.SettingsView:
			cmds = append(cmds, b.handleSettingsInput())
		}
		b.inputRequested = false
	// Escape
	case key.Matches(msg, b.keyMap.Esc):
		cmds = append(cmds, b.resetView())
	// Swap view
	case key.Matches(msg, b.keyMap.SwapView):
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			b.switchActiveView(internal.ModInfoView)
		case internal.ModInfoView:
			b.switchActiveView(internal.ModListView)
		default:
			b.switchActiveView(internal.ModListView)
		}
	// Show logs
	case key.Matches(msg, b.keyMap.ShowLogs):
		b.prepareLogsView()
	// Refresh list
	case key.Matches(msg, b.keyMap.RefreshList):
		if b.activeBox != internal.EnterKspDirView && b.activeBox != internal.SearchView {
			b.ready = false
			cmds = append(cmds, b.getAvailableModsCmd(), b.bubbles.spinner.Tick)
		}
	// Search mods
	case key.Matches(msg, b.keyMap.Search):
		cmds = append(cmds, b.prepareSearchView())
	// Apply Changes
	case key.Matches(msg, b.keyMap.Apply):
		b.ready = false
		cmds = append(cmds, b.applyModsCmd(), b.bubbles.spinner.Tick)
	// View settings
	case key.Matches(msg, b.keyMap.Settings):
		b.prepareSettingsView()
	}

	// only perform search when input is updated
	if b.inputRequested && b.activeBox == internal.SearchView {
		cmds = append(cmds, b.searchCmd(b.bubbles.textInput.Value()))
	}

	return tea.Batch(cmds...)
}

func (b *Bubble) toggleSelectedItem() {
	id := b.registry.ModMapIndex[b.nav.listCursor]
	modMap := b.registry.GetActiveModList()
	mod := modMap[id.Key]

	// toggle mod selection
	b.nav.listSelected = b.nav.listCursor
	if b.nav.installSelected[mod.Identifier].Identifier != "" {
		delete(b.nav.installSelected, mod.Identifier)
	} else {
		b.nav.installSelected[mod.Identifier] = mod
	}
}

func (b *Bubble) resetView() tea.Cmd {
	b.nav.listCursor = 0
	b.nav.listSelected = -1
	b.nav.installSelected = make(map[string]registry.Ckan, 0)
	b.bubbles.textInput.Reset()
	b.inputRequested = false
	b.searchInput = false
	b.switchActiveView(internal.ModListView)
	b.lastActiveBox = internal.ModListView
	return b.sortModMapCmd()
}

// Handle the log view
func (b *Bubble) prepareLogsView() {
	if b.activeBox == internal.LogView {
		b.switchActiveView(internal.ModListView)
	} else {
		b.switchActiveView(internal.LogView)
		b.bubbles.splashViewport.GotoBottom()
	}
}

// Handle screen to input KSP dir
func (b *Bubble) prepareKspDirView() tea.Cmd {
	var cmd tea.Cmd
	if b.activeBox == internal.EnterKspDirView && !b.inputRequested {
		b.inputRequested = false
		b.switchActiveView(internal.ModListView)
	} else if b.activeBox != internal.EnterKspDirView {
		b.switchActiveView(internal.EnterKspDirView)
		b.inputRequested = true
		b.bubbles.textInput.Placeholder = "KSP Directory..."
		b.bubbles.textInput.Reset()
		if b.appConfig.Settings.KerbalDir != "" {
			b.bubbles.textInput.SetValue(b.appConfig.Settings.KerbalDir)
		}
		cmd = textinput.Blink
	}
	return cmd
}

// Handle search page
func (b *Bubble) prepareSearchView() tea.Cmd {
	var cmd tea.Cmd
	switch b.activeBox {
	case internal.SearchView:
		if b.inputRequested {
			val := trimLastChar(b.bubbles.textInput.Value())
			b.bubbles.textInput.SetValue(val)
			b.inputRequested = false
		} else {
			b.inputRequested = true
			cmd = textinput.Blink
		}
	default:
		b.switchActiveView(internal.SearchView)
		b.searchInput = true
		b.inputRequested = true
		b.bubbles.textInput.Reset()
		b.bubbles.textInput.Placeholder = "Search..."
		cmd = textinput.Blink
	}
	return cmd
}

// Handle settings page
func (b *Bubble) prepareSettingsView() tea.Cmd {
	var cmd tea.Cmd
	switch b.activeBox {
	case internal.SettingsView:
		b.switchActiveView(internal.ModListView)
	case internal.ModListView, internal.ModInfoView:
		b.switchActiveView(internal.SettingsView)
		b.bubbles.secondaryViewport.GotoTop()
	}
	return cmd
}

func (b *Bubble) handleSettingsInput() tea.Cmd {
	var cmds []tea.Cmd

	switch b.nav.menuCursor {
	case internal.MenuSortOrder:
		switch b.registry.SortOptions.SortOrder {
		case "ascend":
			b.registry.SortOptions.SortOrder = "descend"
		case "descend":
			b.registry.SortOptions.SortOrder = "ascend"
		}
		log.Printf("Swapping sort order to %s", b.registry.SortOptions.SortOrder)

		cmds = append(cmds, b.sortModMapCmd())
	case internal.MenuSortTag:
		// todo: filtering
	case internal.MenuCompatible:
		cfg := config.GetConfig()
		viper.Set("settings.hide_incompatible", !cfg.Settings.HideIncompatibleMods)
		viper.WriteConfigAs(viper.ConfigFileUsed())
		b.ready = false
		cmds = append(cmds, b.getAvailableModsCmd(), b.bubbles.spinner.Tick)
	case internal.MenuKspDir:
		cmds = append(cmds, b.prepareKspDirView())
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
