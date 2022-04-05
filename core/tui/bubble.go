package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/jedwards1230/go-kerbal/core/theme"
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/jedwards1230/go-kerbal/registry/datacollector"
	"github.com/knipferrc/fm/help"
)

type Bubble struct {
	//appConfig         config.Config
	primaryViewport   viewport.Model
	secondaryViewport viewport.Model
	loadingViewport   viewport.Model
	modList           []datacollector.Ckan
	help              help.Bubble
	keyMap            KeyMap
	status            string
	ready             bool
	loading           bool
	cursor            int
	width             int
	height            int
	selected          int
	//loadingMsg        string
}

func InitialModel() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.Theme.AppTheme)

	primaryBoxBorder := lipgloss.NormalBorder()
	primaryBoxBorderColor := theme.ActiveBoxBorderColor
	secondaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorderColor := theme.InactiveBoxBorderColor
	loadingBoxBorder := lipgloss.NormalBorder()
	loadingBoxBorderColor := theme.InactiveBoxBorderColor

	pvp := viewport.New(0, 0)
	pvp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder).
		BorderForeground(primaryBoxBorderColor)

	svp := viewport.New(0, 0)
	svp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(secondaryBoxBorder).
		BorderForeground(secondaryBoxBorderColor)

	lvp := viewport.New(0, 0)
	lvp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(loadingBoxBorder).
		BorderForeground(loadingBoxBorderColor)

	h := help.New(
		theme.DefaultTextColor,
		"go-kerbal help!",
		[]help.HelpEntry{
			{Key: "ctrl+c", Description: "Exit FM"},
			{Key: "j | up", Description: "Move up"},
			{Key: "k | down", Description: "Move down"},
			{Key: "spacebar", Description: "Select an entry"},
			{Key: "1", Description: "Refresh mod list"},
		})

	return Bubble{
		primaryViewport:   pvp,
		secondaryViewport: svp,
		loadingViewport:   lvp,
		modList:           registry.BuildRegistry(),
		help:              h,
		selected:          -1,
		keyMap:            DefaultKeyMap(),
	}
}

func (b Bubble) Init() tea.Cmd {
	return nil
}
