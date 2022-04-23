package bubbles

import "github.com/charmbracelet/bubbles/key"

const spacebar = " "

// KeyMap defines the keybindings for the viewport. Note that you don't
// necessary need to use keybindings at all; the viewport can be controlled
// programmatically with methods like Model.LineDown(1). See the GoDocs for
// details.
type KeyMap struct {
	Quit     key.Binding
	Down     key.Binding
	Up       key.Binding
	Space    key.Binding
	Enter    key.Binding
	Esc      key.Binding
	SwapView key.Binding
	ShowLogs key.Binding

	RefreshList      key.Binding
	HideIncompatible key.Binding
	SwapSortOrder    key.Binding
	EnterKspDir      key.Binding
	Download         key.Binding
	Search           key.Binding
	Settings         key.Binding

	PageDown     key.Binding
	PageUp       key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
}

// GetKeyMap returns a set of default keybindings.
func GetKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("spacebar", "select"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Enter selection"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("escape", "Exit menu"),
		),
		SwapView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "select"),
		),
		ShowLogs: key.NewBinding(
			key.WithKeys("O"),
		),

		RefreshList: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "refresh mod list"),
		),
		HideIncompatible: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "toggle incompatible mods view"),
		),
		SwapSortOrder: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "swap sort order (asc/desc)"),
		),
		EnterKspDir: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "swap sort order (asc/desc)"),
		),
		Download: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "download selected mod"),
		),
		Search: key.NewBinding(
			key.WithKeys("6"),
			key.WithHelp("6", "search mods"),
		),
		Settings: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "open settings"),
		),

		PageDown: key.NewBinding(
			key.WithKeys("pgdown", spacebar, "f"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
		),
	}
}
