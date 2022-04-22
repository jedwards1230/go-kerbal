package internal

import "github.com/charmbracelet/lipgloss"

const (
	PrimaryBoxActive = 1
	ModInfoBox       = 2
	SplashBoxActive  = 3
	LogView          = 4
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

// BoldTextStyle is the style used for bold text.
var BoldTextStyle = lipgloss.NewStyle().Bold(true)
