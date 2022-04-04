package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the viewport. Note that you don't
// necessary need to use keybindings at all; the viewport can be controlled
// programmatically with methods like Model.LineDown(1). See the GoDocs for
// details.
type KeyMap struct {
	Quit                key.Binding
	Down                key.Binding
	Up                  key.Binding
	Space               key.Binding
	Left                key.Binding
	Right               key.Binding
	Preview             key.Binding
	GotoBottom          key.Binding
	HomeShortcut        key.Binding
	RootShortcut        key.Binding
	ToggleHidden        key.Binding
	ShowDirectoriesOnly key.Binding
	ShowFilesOnly       key.Binding
	CopyPathToClipboard key.Binding
	Zip                 key.Binding
	Unzip               key.Binding
	NewFile             key.Binding
	NewDirectory        key.Binding
	Delete              key.Binding
	Move                key.Binding
	Enter               key.Binding
	Edit                key.Binding
	Copy                key.Binding
	Find                key.Binding
	Rename              key.Binding
	Escape              key.Binding
	ShowLogs            key.Binding
	ToggleBox           key.Binding
}

// DefaultKeyMap returns a set of pager-like default keybindings.
// DefaultKeyMap returns a set of default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
		),
	}
}
