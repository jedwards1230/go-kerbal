package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/common"
	"github.com/jedwards1230/go-kerbal/internal/config"
	"github.com/jedwards1230/go-kerbal/internal/registry"
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
		b.bubbles.primaryPaginator.GoToStart()
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
			common.LogSuccess("Kerbal directory updated")
			b.bubbles.textInput.Reset()
			b.bubbles.textInput.SetValue(fmt.Sprintf("Success!: %v", cfg.Settings.KerbalDir))
			b.inputRequested = false
		} else {
			common.LogErrorf("Error updating ksp dir: %v", msg)
			b.bubbles.textInput.Reset()
			b.bubbles.textInput.Placeholder = "Try again..."
		}

	case SearchMsg:
		if len(msg) >= 0 {
			b.nav.listCursorHide = true
			b.registry.ModMapIndex = registry.ModIndex(msg)
		} else {
			common.LogError("Error searching")
		}

	case ErrorMsg:
		b.ready = true
		common.LogErrorf("ErrorMsg: %v", msg)

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

		b.bubbles.commandViewport.Width = (msg.Width / 3) - b.bubbles.commandViewport.Style.GetHorizontalFrameSize() + 8
		b.bubbles.commandViewport.Height = 7

		b.bubbles.secondaryViewport.Width = (msg.Width / 3) - b.bubbles.secondaryViewport.Style.GetHorizontalFrameSize() + 8
		b.bubbles.secondaryViewport.Height = msg.Height - b.bubbles.commandViewport.Height - internal.StatusBarHeight - b.bubbles.secondaryViewport.Style.GetVerticalFrameSize() - 7

		b.bubbles.primaryPaginator.SetWidth(msg.Width - b.bubbles.secondaryViewport.Width - 8)
		b.bubbles.primaryPaginator.SetHeight(b.bubbles.secondaryViewport.Height + b.bubbles.commandViewport.Height)

		b.bubbles.splashPaginator.SetWidth(msg.Width - 4)
		b.bubbles.splashPaginator.SetHeight(msg.Height - internal.StatusBarHeight - 9)

	case tea.KeyMsg:
		cmds = append(cmds, b.handleKeys(msg))

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			b.scrollView("up")
		case tea.MouseWheelDown:
			b.scrollView("down")
		}

	case TickMsg:
		cmds = append(cmds, b.TickCmd())
		//log.Print("tick")

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
	b.updateActiveMod()

	b.bubbles.primaryPaginator.SetTotalPages(len(b.registry.ModMapIndex))
	b.checkActiveViewPortBounds()

	b.bubbles.commandViewport.SetContent(b.commandView())

	switch b.activeBox {
	case internal.ModListView, internal.SearchView:
		b.bubbles.primaryPaginator.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	case internal.ModInfoView:
		b.bubbles.primaryPaginator.SetContent(b.modListView())
		b.bubbles.secondaryViewport, cmd = b.bubbles.secondaryViewport.Update(msg)
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	case internal.EnterKspDirView:
		b.bubbles.splashPaginator.SetTotalPages(1)
		b.bubbles.splashPaginator.SetContent(b.inputKspView())
	case internal.SettingsView:
		b.bubbles.primaryPaginator.SetContent(b.modListView())
		b.bubbles.secondaryViewport.SetContent(b.settingsView())
	case internal.QueueView:
		b.bubbles.primaryPaginator.SetContent(b.queueView())
		b.bubbles.secondaryViewport.SetContent(b.modInfoView())
	case internal.LogView:
		b.logs = b.checkLogs()
		b.bubbles.splashPaginator.SetTotalPages(len(b.logs))
		b.bubbles.splashPaginator.SetContent(b.logView())
	}

	return cmd
}

func (b *Bubble) switchActiveView(newView int) {
	b.bubbles.primaryPaginator.GoToStart()
	b.bubbles.splashPaginator.GoToStart()
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
	}
}

// Handles mouse scrolling in the viewport
func (b *Bubble) scrollView(dir string) {
	switch dir {
	case "up":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = false
			} else {
				b.bubbles.primaryPaginator.LineUp()
			}
		case internal.LogView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = false
			} else {
				b.bubbles.splashPaginator.LineUp()
			}
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineUp(1)
		case internal.SettingsView:
			b.nav.menuCursor--
		}
	case "down":
		switch b.activeBox {
		case internal.ModListView, internal.SearchView, internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = false
			} else {
				b.bubbles.primaryPaginator.LineDown()
			}
		case internal.LogView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = false
			} else {
				b.bubbles.splashPaginator.LineDown()
			}
		case internal.ModInfoView:
			b.bubbles.secondaryViewport.LineDown(1)
		case internal.SettingsView:
			b.nav.menuCursor++
		}
	case "left":
		switch b.activeBox {
		case internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.boolCursor = !b.nav.boolCursor
			} else {
				b.nav.boolCursor = false
				b.bubbles.primaryPaginator.PrevPage()
			}
		case internal.ModListView, internal.SearchView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = false
			} else {
				b.bubbles.primaryPaginator.PrevPage()
			}
		case internal.LogView:
			b.bubbles.splashPaginator.PrevPage()
		}
	case "right":
		switch b.activeBox {
		case internal.QueueView:
			if b.nav.listCursorHide {
				b.nav.boolCursor = !b.nav.boolCursor
			} else {
				b.nav.boolCursor = false
				b.bubbles.primaryPaginator.NextPage()
			}
		case internal.ModListView, internal.SearchView:
			if b.nav.listCursorHide {
				b.nav.listCursorHide = !b.nav.listCursorHide
			} else {
				b.bubbles.primaryPaginator.NextPage()
			}
		case internal.LogView:
			b.bubbles.splashPaginator.NextPage()
		}
	default:
		log.Panic("Invalid scroll direction: " + dir)
	}
}

func (b *Bubble) updateActiveMod() {
	if !b.nav.listCursorHide && len(b.registry.ModMapIndex) > 0 {
		cursor := b.bubbles.primaryPaginator.GetCursorIndex()
		id := b.registry.ModMapIndex[cursor]
		b.nav.activeMod = b.registry.UnsortedModMap[id.Key]
	}
}
