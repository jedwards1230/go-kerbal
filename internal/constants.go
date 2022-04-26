package internal

import "github.com/charmbracelet/lipgloss"

const (
	ModListView     = 1
	ModInfoView     = 2
	LogView         = 3
	EnterKspDirView = 4
	SearchView      = 5
	SettingsView    = 6
	QueueView       = 7
)

const (
	// StatusBarHeight represents the height of the status bar.
	StatusBarHeight = 1

	// BoxPadding represents the padding of the boxes.
	BoxPadding = 1

	// EllipsisStyle represents the characters displayed when overflowing.
	EllipsisStyle = "..."

	// FileSizeLoadingStyle represents the characters displayed when file sizes are loading.
	FileSizeLoadingStyle = "---"
)

const (
	MenuInputs     = 4
	MenuSortOrder  = 0
	MenuSortTag    = 1
	MenuCompatible = 2
	MenuKspDir     = 3
)

// todo: remove
// BoldTextStyle is the style used for bold text.
var BoldTextStyle = lipgloss.NewStyle().Bold(true)
