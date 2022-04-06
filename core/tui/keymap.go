package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the viewport. Note that you don't
// necessary need to use keybindings at all; the viewport can be controlled
// programmatically with methods like Model.LineDown(1). See the GoDocs for
// details.
type KeyMap struct {
	Quit     key.Binding
	Down     key.Binding
	Up       key.Binding
	Space    key.Binding
	SwapView key.Binding
	ShowLogs key.Binding
	One      key.Binding
}

// DefaultKeyMap returns a set of pager-like default keybindings.
// DefaultKeyMap returns a set of default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q", "quit"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("spacebar", "select"),
		),
		SwapView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "select"),
		),
		ShowLogs: key.NewBinding(
			key.WithKeys("O"),
		),
		One: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "refresh mod list"),
		),
	}
}
