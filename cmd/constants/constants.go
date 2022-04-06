package constants

import "github.com/charmbracelet/lipgloss"

const (
	// PrimaryBoxActive represents when the left box is active.
	PrimaryBoxActive = 1

	// SecondaryBoxActive represents when the right box is active.
	SecondaryBoxActive = 2

	SplashBoxActive = 3
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
