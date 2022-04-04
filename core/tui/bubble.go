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
		"Welcome to FM!",
		[]help.HelpEntry{
			{Key: "ctrl+c", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
			{Key: "h/left", Description: "Go back a directory"},
			{Key: "l/right", Description: "Read file or enter directory"},
			{Key: "p", Description: "Preview directory"},
			{Key: "gg", Description: "Go to top of filetree or box"},
			{Key: "G", Description: "Go to bottom of filetree or box"},
			{Key: "~", Description: "Go to home directory"},
			{Key: "/", Description: "Go to root directory"},
			{Key: ".", Description: "Toggle hidden files"},
			{Key: "S", Description: "Only show directories"},
			{Key: "s", Description: "Only show files"},
			{Key: "y", Description: "Copy file path to clipboard"},
			{Key: "Z", Description: "Zip currently selected tree item"},
			{Key: "U", Description: "Unzip currently selected tree item"},
			{Key: "n", Description: "Create new file"},
			{Key: "N", Description: "Create new directory"},
			{Key: "ctrl+d", Description: "Delete currently selected tree item"},
			{Key: "M", Description: "Move currently selected tree item"},
			{Key: "enter", Description: "Process command"},
			{Key: "E", Description: "Edit currently selected tree item"},
			{Key: "C", Description: "Copy currently selected tree item"},
			{Key: "esc", Description: "Reset FM to initial state"},
			{Key: "O", Description: "Show logs if debugging enabled"},
			{Key: "tab", Description: "Toggle between boxes"},
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
