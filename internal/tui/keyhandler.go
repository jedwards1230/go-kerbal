package tui

import (
	"log"

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
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = false
			} else {
				b.scrollView("down")
			}
		case internal.SettingsView:
			b.scrollView("down")
		}
		b.inputRequested = false

	// Up
	case key.Matches(msg, b.keyMap.Up):
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = false
			} else {
				b.scrollView("up")
			}
		case internal.SettingsView:
			b.scrollView("up")
		}
		b.inputRequested = false

	// Left
	case key.Matches(msg, b.keyMap.Left):
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.boolCursor = !b.nav.boolCursor
			} else {
				b.nav.boolCursor = false
				b.bubbles.paginator.PrevPage()
			}
		}

	// Right
	case key.Matches(msg, b.keyMap.Right):
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.boolCursor = !b.nav.boolCursor
			} else {
				b.nav.boolCursor = true
				b.bubbles.paginator.NextPage()
			}
		}

	// Space
	case key.Matches(msg, b.keyMap.Space):
		b.nav.listCursorHide = !b.nav.listCursorHide

	// Enter
	case key.Matches(msg, b.keyMap.Enter):
		cmds = append(cmds, b.handleEnterKey())

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
		case internal.QueueView:
			b.nav.listCursorHide = !b.nav.listCursorHide
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
		if b.activeBox == internal.QueueView {
			b.switchActiveView(internal.ModListView)
		} else {
			b.switchActiveView(internal.QueueView)
			b.prepareQueueView()
		}

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
	if len(b.registry.ModMapIndex) > 0 {
		cursor := b.bubbles.paginator.GetSliceStart() + b.bubbles.paginator.Cursor
		id := b.registry.ModMapIndex[cursor]
		modMap := b.registry.GetActiveModList()
		mod := modMap[id.Key]

		// toggle mod selection
		if b.nav.installSelected[mod.Identifier].Identifier != "" {
			delete(b.nav.installSelected, mod.Identifier)
		} else {
			b.nav.installSelected[mod.Identifier] = mod
		}
	}
}

func (b *Bubble) resetView() tea.Cmd {
	b.nav.boolCursor = false
	b.nav.listCursor = 0
	b.nav.listCursorHide = true
	b.nav.installSelected = make(map[string]registry.Ckan, 0)
	b.registry.Queue.List = make(map[string]map[string]registry.Ckan, 0)
	b.bubbles.textInput.Reset()
	b.inputRequested = false
	b.searchInput = false
	b.switchActiveView(internal.ModListView)
	b.lastActiveBox = internal.ModListView
	return b.getAvailableModsCmd()
}

// Handle inputs when Enter is pressed
func (b *Bubble) handleEnterKey() tea.Cmd {
	var cmds []tea.Cmd

	switch b.activeBox {
	case internal.ModListView, internal.SearchView:
		b.toggleSelectedItem()
	case internal.EnterKspDirView:
		cmds = append(cmds, b.updateKspDirCmd(b.bubbles.textInput.Value()))
	case internal.SettingsView:
		cmds = append(cmds, b.handleSettingsInput())
	case internal.QueueView:
		if b.nav.listCursorHide {
			if b.nav.boolCursor {
				// apply mods in queue
				b.ready = false
				cmds = append(cmds, b.applyModsCmd(), b.bubbles.spinner.Tick)
			} else {
				// cancel
				b.switchActiveView(internal.ModListView)
				b.nav.listCursor = 0
				b.ready = false
				cmds = append(cmds, b.sortModMapCmd())
			}
		}
	}

	b.inputRequested = false
	return tea.Batch(cmds...)
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
		b.nav.listCursorHide = true
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

func (b *Bubble) prepareQueueView() {
	removeQueue := make(map[string]registry.Ckan, 0)
	installQueue := make(map[string]registry.Ckan, 0)
	b.nav.listCursor = -1

	for _, mod := range b.nav.installSelected {
		if mod.Install.Installed {
			removeQueue[mod.Identifier] = mod
		} else {
			installQueue[mod.Identifier] = mod
		}
	}

	if len(removeQueue) > 0 {
		b.registry.Queue.SetRemovals(removeQueue)
	}
	if len(installQueue) > 0 {
		b.registry.Queue.SetSelections(installQueue)
	}

	// collect all mods and dependencies
	log.Print("Checking dependencies")
	if b.registry.Queue.InstallLen() > 0 {
		mods, err := b.registry.CheckDependencies()
		if err != nil {
			b.LogErrorf("%v", err)
		}
		b.registry.Queue.SetDependencies(mods)
	}
	idx, err := b.registry.BuildQueueIndex()
	if err != nil {
		b.LogErrorf("Cannot build queue index: %v", err)
	}
	b.registry.ModMapIndex = idx
}
