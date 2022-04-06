package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
)

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case getAvailableModsMsg:
		b.modList = msg
		b.logs = append(b.logs, "Mod list updated")
		b.checkPrimaryViewportBounds()
		b.splashScreenActive = false
		b.primaryViewport.GotoTop()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		b.help.Width = msg.Width

		b.splashViewport.Width = msg.Width - b.primaryViewport.Style.GetHorizontalFrameSize()
		b.splashViewport.Height = msg.Height - constants.StatusBarHeight - b.primaryViewport.Style.GetVerticalFrameSize()
		b.splashViewport.SetContent(b.loadingView())

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
		cmd = b.handleKeys(msg)
		cmds = append(cmds, cmd)

		return b, tea.Batch(cmds...)
	}

	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// handleKeys handles all keypresses.
func (b *Bubble) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, b.keyMap.Quit):
		b.logs = append(b.logs, "Quitting")
		return tea.Quit
	case key.Matches(msg, b.keyMap.Down):
		b.cursor++
		b.checkPrimaryViewportBounds()
		b.primaryViewport.SetContent(b.modListView())
	case key.Matches(msg, b.keyMap.Up):
		b.cursor--
		b.checkPrimaryViewportBounds()
		b.primaryViewport.SetContent(b.modListView())
	case key.Matches(msg, b.keyMap.Space):
		if b.selected == b.cursor {
			b.selected = -1
		} else {
			b.selected = b.cursor
		}
		b.checkPrimaryViewportBounds()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case key.Matches(msg, b.keyMap.Tab):
		b.checkPrimaryViewportBounds()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case key.Matches(msg, b.keyMap.ShowLogs):
		log.Println("show logs")
		if b.splashScreenActive {
			b.splashScreenActive = false
			b.primaryViewport.SetContent(b.modListView())
			b.secondaryViewport.SetContent(b.modInfoView())
		} else {
			b.splashScreenActive = true
			b.splashViewport.SetContent(b.logView())
		}
	case key.Matches(msg, b.keyMap.One):
		b.logs = append(b.logs, "Getting mod list")
		b.splashScreenActive = true
		b.splashViewport.SetContent(b.loadingView())
		cmd = b.getAvailableModsCmd()
		cmds = append(cmds, cmd)
	}
	b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)

	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// checkPrimaryViewportBounds handles wrapping of the filetree and
// scrolling of the viewport.
func (b *Bubble) checkPrimaryViewportBounds() {
	top := b.primaryViewport.YOffset - 3
	bottom := b.primaryViewport.Height + b.primaryViewport.YOffset - 4

	if b.cursor < top {
		b.primaryViewport.LineUp(1)
	} else if b.cursor > bottom {
		b.primaryViewport.LineDown(1)
	}

	if b.cursor > len(b.modList)-1 {
		b.primaryViewport.GotoTop()
		b.cursor = 0
	} else if b.cursor < 0 {
		b.primaryViewport.GotoBottom()
		b.cursor = len(b.modList) - 1
	}
}
