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
		b.textInput, cmd = b.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if b.activeBox == internal.ModInfoView {
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
			b.activeBox = internal.SearchView
			b.textInput.Focus()
			// only search when input is updated
			_, ok := msg.(tea.KeyMsg)
			if ok {
				cmds = append(cmds, b.searchCmd(b.textInput.Value()))
			}
		} else {
			b.activeBox = internal.EnterKspDirView
			b.textInput.Focus()
			b.splashViewport.SetContent(b.inputKspView())
		}
	}

	return b, tea.Batch(cmds...)
}

// Handles wrapping and button
// scrolling in the viewport.
func (b *Bubble) checkActiveViewPortBounds() {
	switch b.activeBox {
	case internal.ModListView, internal.SearchView:
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
	case internal.ModInfoView:
		if b.secondaryViewport.AtBottom() {
			b.secondaryViewport.GotoBottom()
		} else if b.secondaryViewport.AtTop() {
			b.secondaryViewport.GotoTop()
		}
	case internal.LogView:
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
		case internal.ModListView, internal.SearchView:
			b.nav.listCursor--
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.ModInfoView:
			b.secondaryViewport.LineUp(1)
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.LogView:
			b.splashViewport.LineUp(1)
			b.splashViewport.SetContent(b.logView())
		}
	case "down":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			b.nav.listCursor++
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.ModInfoView:
			b.secondaryViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.primaryViewport.SetContent(b.modListView())
		case internal.LogView:
			b.splashViewport.LineDown(1)
			b.checkActiveViewPortBounds()
			b.splashViewport.SetContent(b.logView())
		}
	default:
		log.Panic("Invalid scroll direction: " + dir)
	}
	b.checkActiveViewPortBounds()
}
