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

	if b.inputRequested {
		b.bubbles.textInput, cmd = b.bubbles.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	// Update mod list
	case UpdatedModMapMsg:
		b.registry.TotalModMap = msg
		return b, b.sortModMapCmd()
	case SortedMsg:
		b.registry.SortModList()
		b.LogCommand("Sorted mod map")
		b.ready = true
		b.checkActiveViewPortBounds()
		b.bubbles.primaryViewport.GotoTop()
		if b.activeBox == internal.SearchView {
			cmds = append(cmds, b.searchCmd(b.bubbles.textInput.Value()))
		}
	case InstalledModListMsg:
		b.ready = true
		if len(b.registry.InstalledModList) != len(msg) {
			b.LogCommand("Updated installed mod list")
			b.registry.InstalledModList = msg
		}
		cmds = append(cmds, b.getAvailableModsCmd())
	// Update KSP dir
	case UpdateKspDirMsg:
		b.ready = true
		if msg {
			cfg := config.GetConfig()
			b.LogCommand("Kerbal directory updated")
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
			b.nav.listSelected = -1
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
	// Window resize
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		b.bubbles.help.Width = msg.Width

		b.bubbles.splashViewport.Width = msg.Width - b.bubbles.primaryViewport.Style.GetHorizontalFrameSize()
		b.bubbles.splashViewport.Height = msg.Height - (internal.StatusBarHeight * 2) - b.bubbles.primaryViewport.Style.GetVerticalFrameSize()

		b.bubbles.primaryViewport.Width = (msg.Width / 2) - b.bubbles.primaryViewport.Style.GetHorizontalFrameSize()
		b.bubbles.primaryViewport.Height = msg.Height - (internal.StatusBarHeight * 2) - b.bubbles.primaryViewport.Style.GetVerticalFrameSize()
		b.bubbles.secondaryViewport.Width = (msg.Width / 2) - b.bubbles.secondaryViewport.Style.GetHorizontalFrameSize()
		b.bubbles.secondaryViewport.Height = msg.Height - (internal.StatusBarHeight * 2) - b.bubbles.secondaryViewport.Style.GetVerticalFrameSize()
	// Key pressed
	case tea.KeyMsg:
		cmds = append(cmds, b.handleKeys(msg))
	// Mouse input
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

	if b.activeBox == internal.EnterKspDirView || b.activeBox == internal.SearchView {
		if b.inputRequested {
			b.bubbles.textInput.Focus()
		} else {
			b.bubbles.textInput.Blur()
		}
	}

	cmds = append(cmds, b.updateActiveView(msg))

	return b, tea.Batch(cmds...)
}

// Update content for active view
func (b *Bubble) updateActiveView(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch b.activeBox {
	case internal.ModListView, internal.SearchView:
		b.bubbles.primaryViewport.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	case internal.ModInfoView:
		b.bubbles.secondaryViewport, cmd = b.bubbles.secondaryViewport.Update(msg)
		b.bubbles.primaryViewport.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	case internal.LogView:
		b.bubbles.splashViewport.SetContent(b.logView())
	case internal.EnterKspDirView:
		b.bubbles.splashViewport.SetContent(b.inputKspView())
	case internal.SettingsView:
		b.bubbles.splashViewport.SetContent(b.settingsView())
	}

	return cmd
}

func (b *Bubble) switchActiveView(newView int) {
	b.lastActiveBox = b.activeBox
	b.activeBox = newView
}

// Handles wrapping and button scrolling in the viewport
func (b *Bubble) checkActiveViewPortBounds() {
	switch b.activeBox {
	case internal.ModListView, internal.SearchView:
		top := b.bubbles.primaryViewport.YOffset - 3
		bottom := b.bubbles.primaryViewport.Height + b.bubbles.primaryViewport.YOffset - 4

		if b.nav.listCursor < top {
			b.bubbles.primaryViewport.LineUp(1)
		} else if b.nav.listCursor > bottom {
			b.bubbles.primaryViewport.LineDown(1)
		}

		if b.nav.listCursor > len(b.registry.ModMapIndex)-1 {
			b.nav.listCursor = 0
			b.nav.listSelected = b.nav.listCursor
			b.bubbles.primaryViewport.GotoTop()
		} else if b.nav.listCursor < 0 {
			b.nav.listCursor = len(b.registry.ModMapIndex) - 1
			b.nav.listSelected = b.nav.listCursor
			b.bubbles.primaryViewport.GotoBottom()
		}
	case internal.ModInfoView:
		if b.bubbles.secondaryViewport.AtBottom() {
			b.bubbles.secondaryViewport.GotoBottom()
		} else if b.bubbles.secondaryViewport.AtTop() {
			b.bubbles.secondaryViewport.GotoTop()
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
// TODO: fix scrolling beyond page and in big lists. breaks app.
func (b *Bubble) scrollView(dir string) {
	switch dir {
	case "up":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			b.nav.listCursor--
			b.nav.listSelected = b.nav.listCursor
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineUp(1)
		case internal.LogView:
			b.bubbles.splashViewport.LineUp(1)
		}
	case "down":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			b.nav.listCursor++
			b.nav.listSelected = b.nav.listCursor
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineDown(1)
		case internal.LogView:
			b.bubbles.splashViewport.LineDown(1)
		}
	default:
		log.Panic("Invalid scroll direction: " + dir)
	}
	b.checkActiveViewPortBounds()
}
