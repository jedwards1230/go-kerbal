package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/jedwards1230/go-kerbal/core/theme"
	"github.com/jedwards1230/go-kerbal/help"
	"github.com/jedwards1230/go-kerbal/registry"
)

type Bubble struct {
	// Contains config for TUI
	appConfig config.Config
	// Contains style details for TUI
	theme theme.Theme
	// Primary view (left panel)
	primaryViewport viewport.Model
	// Seconday view (right panel)
	secondaryViewport viewport.Model
	// full page view (splash panel)
	splashViewport viewport.Model
	// Contains mod list and db
	registry    registry.Registry
	sortOptions registry.SortOptions
	help        help.Bubble
	keyMap      KeyMap
	ready       bool
	activeBox   int
	cursor      int
	selected    int
	width       int
	height      int
	logs        []string
}

func InitialModel() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.AppTheme)
	reg := registry.GetRegistry()
	sortOpts := registry.SortOptions{
		SortTag:   "name",
		SortOrder: "ascend",
	}

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
			{Key: "q | ctrl+c", Description: "Exit FM"},
			{Key: "j | up", Description: "Move up"},
			{Key: "k | down", Description: "Move down"},
			{Key: "spacebar", Description: "Select an entry"},
			{Key: "tab", Description: "Swap active views"},
			{},
			{Key: "O", Description: "Show logs if debugging enabled"},
			{},
			{Key: "1", Description: "Refresh mod list"},
			{Key: "2", Description: "Toggle hiding incompatible mods"},
			{Key: "3", Description: "Toggle sort order (ascend/descend)"},
		})

	return Bubble{
		appConfig:         cfg,
		theme:             theme,
		primaryViewport:   primaryVP,
		secondaryViewport: secondaryVP,
		splashViewport:    splashVP,
		registry:          reg,
		sortOptions:       sortOpts,
		help:              h,
		selected:          -1,
		activeBox:         constants.PrimaryBoxActive,
		logs:              []string{"Initializing"},
		keyMap:            DefaultKeyMap(),
	}
}

func (b Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmd := b.getAvailableModsCmd()
	cmds = append(cmds, cmd)

	b.splashViewport.SetContent(b.loadingView())

	return tea.Batch(cmds...)
}
