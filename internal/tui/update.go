package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/registry"
)

// Do computations for TUI app
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	//b.bubbles.paginator, cmd = b.bubbles.paginator.Update(msg)
	//cmds = append(cmds, cmd)

	if b.inputRequested {
		b.bubbles.textInput, cmd = b.bubbles.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case UpdatedModMapMsg:
		b.registry.TotalModMap = msg
		cmds = append(cmds, b.sortModMapCmd())

	case SortedMsg:
		b.registry.SortModList()
		b.nav.listCursorHide = true
		b.ready = true
		b.bubbles.paginator.GoToStart()
		if b.activeBox == internal.SearchView {
			cmds = append(cmds, b.searchCmd(b.bubbles.textInput.Value()))
		}

	case InstalledModListMsg:
		b.ready = true
		//cmds = append(cmds, b.getAvailableModsCmd())

	case UpdateKspDirMsg:
		b.ready = true
		if msg {
			cfg := config.GetConfig()
			b.LogSuccess("Kerbal directory updated")
			b.bubbles.textInput.Reset()
			b.bubbles.textInput.SetValue(fmt.Sprintf("Success!: %v", cfg.Settings.KerbalDir))
			b.inputRequested = false
		} else {
			b.LogErrorf("Error updating ksp dir: %v", msg)
			b.bubbles.textInput.Reset()
			b.bubbles.textInput.Placeholder = "Try again..."
		}

	case SearchMsg:
		if len(msg) >= 0 {
			b.nav.listCursorHide = true
			b.registry.ModMapIndex = registry.ModIndex(msg)
		} else {
			b.LogError("Error searching")
		}

	case ErrorMsg:
		b.ready = true
		b.LogErrorf("ErrorMsg: %v", msg)

	case spinner.TickMsg:
		if !b.ready {
			b.bubbles.spinner, cmd = b.bubbles.spinner.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			b.bubbles.spinner.Finish()
		}

	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height

		b.bubbles.paginator.PerPage = b.height - 11
		b.bubbles.paginator.Height = b.height - 11
		b.bubbles.paginator.Width = (b.width / 2) - 4

		b.bubbles.secondaryViewport.Width = (msg.Width / 2) - b.bubbles.secondaryViewport.Style.GetHorizontalFrameSize()
		b.bubbles.secondaryViewport.Height = (msg.Height * 2 / 3) - internal.StatusBarHeight - b.bubbles.secondaryViewport.Style.GetVerticalFrameSize() - 4

		b.bubbles.commandViewport.Width = (msg.Width / 2) - b.bubbles.commandViewport.Style.GetHorizontalFrameSize()
		b.bubbles.commandViewport.Height = (msg.Height / 3) - internal.StatusBarHeight - b.bubbles.commandViewport.Style.GetVerticalFrameSize() - 1

		b.bubbles.splashViewport.Width = msg.Width - b.bubbles.splashViewport.Style.GetHorizontalFrameSize()
		b.bubbles.splashViewport.Height = msg.Height - internal.StatusBarHeight - b.bubbles.splashViewport.Style.GetVerticalFrameSize() - 6

	case tea.KeyMsg:
		cmds = append(cmds, b.handleKeys(msg))

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			b.scrollView("up")
		case tea.MouseWheelDown:
			b.scrollView("down")
		}

	case MyTickMsg:
		cmds = append(cmds, b.MyTickCmd())

	default:
		log.Printf("%T", msg)
	}

	if b.inputRequested {
		b.bubbles.textInput.Focus()
	} else {
		b.bubbles.textInput.Blur()
	}

	cmds = append(cmds, b.updateActiveView(msg))

	return b, tea.Batch(cmds...)
}

// Update content for active view
func (b *Bubble) updateActiveView(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	b.bubbles.paginator.SetTotalPages(len(b.registry.ModMapIndex))
	b.checkActiveViewPortBounds()
	b.updateActiveMod()

	b.bubbles.commandViewport.SetContent(b.commandView())

	switch b.activeBox {
	case internal.ModListView, internal.SearchView:
		b.bubbles.paginator.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	case internal.ModInfoView:
		b.bubbles.paginator.SetContent(b.modListView())
		b.bubbles.secondaryViewport, cmd = b.bubbles.secondaryViewport.Update(msg)
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	case internal.EnterKspDirView:
		b.bubbles.splashViewport.SetContent(b.inputKspView())
	case internal.SettingsView:
		b.bubbles.paginator.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.settingsView())
	case internal.LogView:
		b.bubbles.splashViewport.SetContent(b.logView())
	case internal.QueueView:
		b.bubbles.paginator.SetContent(b.queueView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	}

	return cmd
}

func (b *Bubble) switchActiveView(newView int) {
	b.bubbles.paginator.Page = 0
	b.bubbles.paginator.Cursor = 0
	b.lastActiveBox = b.activeBox
	b.activeBox = newView
}

// Handles wrapping and button scrolling in the viewport
func (b *Bubble) checkActiveViewPortBounds() {
	switch b.activeBox {
	case internal.ModInfoView:
		if b.bubbles.secondaryViewport.AtBottom() {
			b.bubbles.secondaryViewport.GotoBottom()
		} else if b.bubbles.secondaryViewport.AtTop() {
			b.bubbles.secondaryViewport.GotoTop()
		}
	case internal.SettingsView:
		if b.nav.menuCursor >= internal.MenuInputs {
			b.nav.menuCursor = 0
		} else if b.nav.menuCursor < 0 {
			b.nav.menuCursor = internal.MenuInputs - 1
		}
	case internal.LogView:
		if b.bubbles.splashViewport.AtBottom() {
			b.bubbles.splashViewport.GotoBottom()
		} else if b.bubbles.splashViewport.AtTop() {
			b.bubbles.splashViewport.GotoTop()
		}
	}
}

// Handles mouse scrolling in the viewport
func (b *Bubble) scrollView(dir string) {
	switch dir {
	case "up":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			b.bubbles.paginator.LineUp()
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineUp(1)
		case internal.SettingsView:
			b.nav.menuCursor--
		case internal.LogView:
			b.bubbles.splashViewport.LineUp(1)
		}
	case "down":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			b.bubbles.paginator.LineDown()
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineDown(1)
		case internal.SettingsView:
			b.nav.menuCursor++
		case internal.LogView:
			b.bubbles.splashViewport.LineDown(1)
		}
	default:
		log.Panic("Invalid scroll direction: " + dir)
	}
}

func (b *Bubble) updateActiveMod() {
	if !b.nav.listCursorHide && len(b.registry.ModMapIndex) > 0 {
		cursor := b.bubbles.paginator.GetCursorIndex()

		// todo: check if all incompat map can be used for all cases
		if b.activeBox == internal.QueueView {
			//log.Printf("idx: %v, cur: %v", b.registry.ModMapIndex, cursor)
			id := b.registry.ModMapIndex[cursor]
			b.nav.activeMod = b.registry.SortedNonCompatibleMap[id.Key]
		} else {
			modMap := b.registry.GetActiveModList()
			//log.Printf("idx: %v, cur: %v, cur idx: %v", len(b.registry.ModMapIndex), cursor, b.bubbles.paginator.Index)
			id := b.registry.ModMapIndex[cursor]
			b.nav.activeMod = modMap[id.Key]
		}
	}
}
