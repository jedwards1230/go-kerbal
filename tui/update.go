package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
)

func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return b, tea.Quit
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			b.cursor--
			b.checkPrimaryViewportBounds()
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			b.cursor++
			b.checkPrimaryViewportBounds()
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := b.selected[b.cursor]
			if ok {
				delete(b.selected, b.cursor)
			} else {
				b.selected[b.cursor] = struct{}{}
			}
		}
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height

		b.primaryViewport.Width = (msg.Width / 2) - b.primaryViewport.Style.GetHorizontalFrameSize()
		b.primaryViewport.Height = msg.Height - constants.StatusBarHeight - b.primaryViewport.Style.GetVerticalFrameSize()
		b.secondaryViewport.Width = (msg.Width / 2) - b.secondaryViewport.Style.GetHorizontalFrameSize()
		b.secondaryViewport.Height = msg.Height - constants.StatusBarHeight - b.secondaryViewport.Style.GetVerticalFrameSize()

		b.primaryViewport.SetContent(b.modListView())
		b.secondaryViewport.SetContent(b.modListView())

		if !b.ready {
			b.ready = true
		}

		return b, nil
	}

	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
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
