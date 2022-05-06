package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/theme"
)

var (
	SecondaryVP = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.AppTheme.InactiveBoxBorderColor)

	CommandVP = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.AppTheme.InactiveBoxBorderColor)

	SplashVP = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.AppTheme.ActiveBoxBorderColor).
			Align(lipgloss.Center)
)

var (
	PrimaryTitle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Height(3).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 0)

	SecondaryTitle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Height(3).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 0)
)

var (
	Spinner = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	TextInput = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	PaginatorActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")

	PaginatorInactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
)

var (
	ListSelected = lipgloss.NewStyle().
		Background(theme.AppTheme.SelectedListItemColor).
		Foreground(theme.AppTheme.UnselectedListItemColor)
)

var (
	AbstractInfo = lipgloss.NewStyle().
		Bold(false).
		Align(lipgloss.Center).
		Height(3).
		Padding(1, 2)
)

var (
	Installed = lipgloss.NewStyle().Foreground(theme.AppTheme.InstalledColor)

	Incompatible = lipgloss.NewStyle().Foreground(theme.AppTheme.Orange)
)

var (
	KeyStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Bold(true).
			Padding(0, 2)

	ValueStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Padding(0, 2)
)

var (
	StatusBar = lipgloss.NewStyle().
		Height(internal.StatusBarHeight)
)

var (
	LogLineNumber = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Width(5).
			Padding(0, 2, 0, 1)

	LogTimestamp = lipgloss.NewStyle().
			Foreground(theme.AppTheme.Green).
			MarginRight(1)

	LogFile = lipgloss.NewStyle().
		Foreground(theme.AppTheme.Blue).
		Width(17).
		MarginRight(1)
)

var (
	QueueTitle = lipgloss.NewStyle().
		Padding(2, 0, 1, 2).
		Bold(true)
)
