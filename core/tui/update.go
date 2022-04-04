package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
)

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if b.loading {
		b.primaryViewport.SetContent(b.loadingView())
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		b.help.Width = msg.Width

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
		return tea.Quit
	case key.Matches(msg, b.keyMap.Down):
		b.cursor++
		b.checkPrimaryViewportBounds()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case key.Matches(msg, b.keyMap.Up):
		b.cursor--
		b.checkPrimaryViewportBounds()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	case key.Matches(msg, b.keyMap.Space):
		if b.selected == b.cursor {
			b.selected = -1
		} else {
			b.selected = b.cursor
		}
		b.checkPrimaryViewportBounds()
		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modInfoView())
	}
	b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)

	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// checkPrimaryViewportBounds handles wrapping of the filetree and
// scrolling of the viewport.
func (b *Bubble) checkPrimaryViewportBounds() {
	top := b.primaryViewport.YOffset
	bottom := b.primaryViewport.Height + b.primaryViewport.YOffset - 1

	if b.cursor < top {
		b.primaryViewport.LineUp(1)
	} else if b.cursor > bottom {
		b.primaryViewport.LineDown(1)
	}

	if b.cursor > len(b.modList)-1 {
		b.primaryViewport.GotoTop()
		b.cursor = 0
	} else if b.cursor < top {
		b.primaryViewport.GotoBottom()
		b.cursor = len(b.modList) - 1
	}
}