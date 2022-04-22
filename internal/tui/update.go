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

	if b.activeBox == internal.ModInfoView {
		b.bubbles.secondaryViewport, cmd = b.bubbles.secondaryViewport.Update(msg)
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
		b.bubbles.primaryViewport.GotoTop()
		//b.bubbles.primaryViewport.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
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
		b.bubbles.primaryViewport.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	// Update KSP dir
	case UpdateKspDirMsg:
		b.ready = true
		if msg {
			cfg := config.GetConfig()
			log.Print("Kerbal directory updated")
			b.bubbles.textInput.Reset()
			b.bubbles.textInput.SetValue(fmt.Sprintf("Success!: %v", cfg.Settings.KerbalDir))
			b.bubbles.textInput.Blur()
			b.inputRequested = false
		} else {
			log.Printf("Error updating ksp dir: %v", msg)
			b.bubbles.textInput.Reset()
			b.bubbles.textInput.Placeholder = "Try again..."
		}
		b.bubbles.splashViewport.SetContent(b.inputKspView())
	case SearchMsg:
		if len(msg) >= 0 {
			b.nav.listSelected = -1
			b.registry.ModMapIndex = registry.ModIndex(msg)
			b.bubbles.primaryViewport.SetContent(b.modListView())
			b.bubbles.secondaryViewport.SetContent(b.modInfoView())
		} else {
			log.Print("Error searching")
		}
	case ErrorMsg:
		b.ready = true
		log.Printf("ErrorMsg: %v", msg)
	case spinner.TickMsg:
		if !b.ready {
			if b.activeBox == internal.LogView {
				b.bubbles.splashViewport.SetContent(b.logView())
			}
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
		b.bubbles.splashViewport.Height = msg.Height - internal.StatusBarHeight - b.bubbles.primaryViewport.Style.GetVerticalFrameSize()
		if b.inputRequested && !b.searchInput {
			b.bubbles.splashViewport.SetContent(b.inputKspView())
		} else {
			b.bubbles.splashViewport.SetContent(b.logView())
		}

		b.bubbles.primaryViewport.Width = (msg.Width / 2) - b.bubbles.primaryViewport.Style.GetHorizontalFrameSize()
		b.bubbles.primaryViewport.Height = msg.Height - (internal.StatusBarHeight * 2) - b.bubbles.primaryViewport.Style.GetVerticalFrameSize()
		b.bubbles.secondaryViewport.Width = (msg.Width / 2) - b.bubbles.secondaryViewport.Style.GetHorizontalFrameSize()
		b.bubbles.secondaryViewport.Height = msg.Height - (internal.StatusBarHeight * 2) - b.bubbles.secondaryViewport.Style.GetVerticalFrameSize()

		b.bubbles.primaryViewport.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())

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
			b.activeBox = internal.SearchView
			b.bubbles.textInput.Focus()
			// only search when input is updated
			_, ok := msg.(tea.KeyMsg)
			if ok {
				cmds = append(cmds, b.searchCmd(b.bubbles.textInput.Value()))
			}
		} else {
			b.activeBox = internal.EnterKspDirView
			b.bubbles.textInput.Focus()
			b.bubbles.splashViewport.SetContent(b.inputKspView())
		}
	}

	return b, tea.Batch(cmds...)
}

// Handles wrapping and button
// scrolling in the viewport.
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
			b.bubbles.primaryViewport.GotoTop()
		} else if b.nav.listCursor < 0 {
			b.nav.listCursor = len(b.registry.ModMapIndex) - 1
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
func (b *Bubble) scrollView(dir string) {
	switch dir {
	case "up":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			b.nav.listCursor--
			b.checkActiveViewPortBounds()
			b.bubbles.primaryViewport.SetContent(b.modListView())
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineUp(1)
			b.checkActiveViewPortBounds()
			b.bubbles.primaryViewport.SetContent(b.modListView())
		case internal.LogView:
			b.bubbles.splashViewport.LineUp(1)
			b.bubbles.splashViewport.SetContent(b.logView())
		}
	case "down":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			b.nav.listCursor++
			b.checkActiveViewPortBounds()
			b.bubbles.primaryViewport.SetContent(b.modListView())
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.bubbles.primaryViewport.SetContent(b.modListView())
		case internal.LogView:
			b.bubbles.splashViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.bubbles.splashViewport.SetContent(b.logView())
		}
	default:
		log.Panic("Invalid scroll direction: " + dir)
	}
	b.checkActiveViewPortBounds()
}
