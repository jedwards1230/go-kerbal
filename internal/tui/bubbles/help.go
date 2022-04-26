package bubbles

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpEntry represents a single entry in the help bubble.
type HelpEntry struct {
	Key         string
	Description string
}

// HelpBubble represents a help bubble.
type HelpBubble struct {
	Width     int
	Height    int
	Entries   []HelpEntry
	TextColor lipgloss.AdaptiveColor
}

// NewHelpBubble creates a new help bubble.
func NewHelpBubble(textColor lipgloss.AdaptiveColor) HelpBubble {
	return HelpBubble{
		TextColor: textColor,
		Entries: []HelpEntry{
			{Key: "up", Description: "Move up"},
			{Key: "down", Description: "Move down"},
			{Key: "spacebar", Description: "Toggle mod info"},
			{Key: "enter", Description: "Select mod for install"},
			{Key: "tab", Description: "Swap active views"},
			{},
			{Key: "shift + o", Description: "Show logs"},
			{},
			{Key: "1", Description: "Refresh mod list"},
			{Key: "2", Description: "Search mods"},
			{Key: "3", Description: "Apply mods"},
			{},
			{Key: "0", Description: "View settings"},
			{Key: "ctrl+c", Description: "Exit"},
		},
	}
}

// SetSize sets the size of the bubble.
func (b *HelpBubble) SetSize(w, h int) {
	b.Width = w
	b.Height = h
}

// Update handles updating the help bubble.
func (b HelpBubble) Update(msg tea.Msg) (HelpBubble, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
	}

	return b, nil
}

// View returns a string representation of the help bubble.
func (b HelpBubble) View() string {
	helpScreen := ""

	for _, content := range b.Entries {
		keyText := lipgloss.NewStyle().
			Bold(true).
			Foreground(b.TextColor).
			Width(20).
			Render(content.Key)
		descriptionText := lipgloss.NewStyle().
			Foreground(b.TextColor).
			Render(content.Description)
		row := lipgloss.JoinHorizontal(lipgloss.Top, keyText, descriptionText)
		helpScreen += fmt.Sprintf("%s\n", row)
	}

	helpScreen = lipgloss.NewStyle().
		Padding(1, 2).
		Render(helpScreen)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().
			Width(b.Width).
			Height(b.Height).
			Render(helpScreen),
	)
}
