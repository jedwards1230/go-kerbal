package help

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

// Bubble represents a help bubble.
type Bubble struct {
	Width     int
	Height    int
	Entries   []HelpEntry
	TextColor lipgloss.AdaptiveColor
}

// New creates a new help bubble.
func New(textColor lipgloss.AdaptiveColor, entries []HelpEntry) Bubble {
	return Bubble{
		TextColor: textColor,
		Entries:   entries,
	}
}

// SetSize sets the size of the bubble.
func (b *Bubble) SetSize(w, h int) {
	b.Width = w
	b.Height = h
}

// Update handles updating the help bubble.
func (b Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
	}

	return b, nil
}

// View returns a string representation of the help bubble.
func (b Bubble) View() string {
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
