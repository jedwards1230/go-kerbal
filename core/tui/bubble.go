package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/jedwards1230/go-kerbal/core/help"
	"github.com/jedwards1230/go-kerbal/core/theme"
	"github.com/jedwards1230/go-kerbal/registry"
)

type Bubble struct {
	appConfig         config.Config
	theme             theme.Theme
	primaryViewport   viewport.Model
	secondaryViewport viewport.Model
	splashViewport    viewport.Model
	textInput         textinput.Model
	inputRequested    bool
	registry          registry.Registry
	sortOptions       registry.SortOptions
	help              help.Bubble
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
	theme := theme.GetTheme(cfg.AppTheme)
	reg := registry.GetRegistry()
	sortOpts := registry.SortOptions{
		SortTag:   "name",
		SortOrder: "ascend",
	}

	iRequested := true
	if cfg.Settings.KerbalDir == "" {
		iRequested = true
	}

	primaryBoxBorder := lipgloss.NormalBorder()
	primaryBoxBorderColor := theme.ActiveBoxBorderColor
	secondaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorderColor := theme.InactiveBoxBorderColor
	splashBoxBorder := lipgloss.NormalBorder()
	splashBoxBorderColor := theme.InactiveBoxBorderColor

	t := textinput.New()
	t.Prompt = "❯ "
	t.CharLimit = 0
	t.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

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
			{Key: "ctrl+c", Description: "Exit"},
			{Key: "up", Description: "Move up"},
			{Key: "down", Description: "Move down"},
			{Key: "spacebar", Description: "Select an entry"},
			{Key: "tab", Description: "Swap active views"},
			{},
			{Key: "O", Description: "Show logs if debugging enabled"},
			{},
			{Key: "1", Description: "Refresh mod list"},
			{Key: "2", Description: "Toggle hiding incompatible mods"},
			{Key: "3", Description: "Toggle sort order (ascend/descend)"},
			{Key: "4", Description: "Update KSP directory"},
			{Key: "5", Description: "Download selected mod"},
		})

	return Bubble{
		appConfig:         cfg,
		theme:             theme,
		primaryViewport:   primaryVP,
		secondaryViewport: secondaryVP,
		splashViewport:    splashVP,
		textInput:         t,
		inputRequested:    iRequested,
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
	cmds = append(cmds, b.getAvailableModsCmd())

	return tea.Batch(cmds...)
}
