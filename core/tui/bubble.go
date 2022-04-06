package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/jedwards1230/go-kerbal/core/theme"
	"github.com/jedwards1230/go-kerbal/registry/datacollector"
	"github.com/knipferrc/fm/help"
)

type Bubble struct {
	appConfig         config.Config
	theme             theme.Theme
	primaryViewport   viewport.Model
	secondaryViewport viewport.Model
	splashViewport    viewport.Model
	modList           []datacollector.Ckan
	help              help.Bubble
	sortFilter        string
	keyMap            KeyMap
	ready             bool
	activeBox         int
	cursor            int
	selected          int
	width             int
	height            int
	logs              []string
}

func InitialModel() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.Theme.AppTheme)

	primaryBoxBorder := lipgloss.NormalBorder()
	primaryBoxBorderColor := theme.ActiveBoxBorderColor
	secondaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorderColor := theme.InactiveBoxBorderColor
	splashBoxBorder := lipgloss.NormalBorder()
	splashBoxBorderColor := theme.InactiveBoxBorderColor

	primaryVP := viewport.New(0, 0)
	primaryVP.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder).
		BorderForeground(primaryBoxBorderColor)

	secondaryVP := viewport.New(0, 0)
	secondaryVP.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(secondaryBoxBorder).
		BorderForeground(secondaryBoxBorderColor)

	splashVP := viewport.New(0, 0)
	splashVP.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(splashBoxBorder).
		BorderForeground(splashBoxBorderColor)

	h := help.New(
		theme.DefaultTextColor,
		"go-kerbal help!",
		[]help.HelpEntry{
			{Key: "ctrl+c", Description: "Exit FM"},
			{Key: "j | up", Description: "Move up"},
			{Key: "k | down", Description: "Move down"},
			{Key: "spacebar", Description: "Select an entry"},
			{Key: "O", Description: "Show logs if debugging enabled"},
			{},
			{Key: "1", Description: "Refresh mod list"},
		})

	return Bubble{
		appConfig:         cfg,
		theme:             theme,
		primaryViewport:   primaryVP,
		secondaryViewport: secondaryVP,
		splashViewport:    splashVP,
		help:              h,
		selected:          -1,
		activeBox:         constants.PrimaryBoxActive,
		logs:              []string{"Initializing"},
		keyMap:            DefaultKeyMap(),
	}
}

func (b Bubble) Init() tea.Cmd {
	cmd := b.getAvailableModsCmd()

	b.splashViewport.SetContent(b.loadingView())

	return cmd
}
