package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/jedwards1230/go-kerbal/internal/help"
	"github.com/jedwards1230/go-kerbal/internal/theme"
	"github.com/jedwards1230/go-kerbal/registry"
)

type Bubble struct {
	appConfig         config.Config
	theme             theme.Theme
	primaryViewport   viewport.Model
	secondaryViewport viewport.Model
	splashViewport    viewport.Model
	spinner           spinner.Model
	textInput         textinput.Model
	inputRequested    bool
	searchInput       bool
	registry          registry.Registry
	help              help.Bubble
	keyMap            KeyMap
	nav               Nav
	ready             bool
	activeBox         int
	width             int
	height            int
	logs              []string
}

type Nav struct {
	listCursor        int
	listSelected      int
	installSelected   map[string]registry.Ckan
	mainButtonsCursor int
}

func InitialModel() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.AppTheme)
	reg := registry.GetRegistry()

	iRequested := false
	if cfg.Settings.KerbalDir == "" {
		iRequested = true
	}

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	primaryBoxBorder := lipgloss.NormalBorder()
	primaryBoxBorderColor := theme.ActiveBoxBorderColor
	secondaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorderColor := theme.InactiveBoxBorderColor
	splashBoxBorder := lipgloss.NormalBorder()
	splashBoxBorderColor := theme.InactiveBoxBorderColor

	t := textinput.New()
	t.Prompt = "❯ "
	t.CharLimit = -1
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
		[]help.HelpEntry{
			{Key: "up", Description: "Move up"},
			{Key: "down", Description: "Move down"},
			{Key: "spacebar", Description: "Toggle mod info"},
			{Key: "enter", Description: "Select mod for install"},
			{Key: "tab", Description: "Swap active views"},
			{},
			{Key: "shift + o", Description: "Show logs"},
			{},
			{Key: "1", Description: "Refresh mod list"},
			{Key: "2", Description: "Toggle hiding incompatible mods"},
			{Key: "3", Description: "Toggle sort order (ascend/descend)"},
			{Key: "4", Description: "Update KSP directory"},
			{Key: "5", Description: "Download selected mod"},
			{Key: "6", Description: "Search mods"},
			{},
			{Key: "0", Description: "View settings"},
			{Key: "ctrl+c", Description: "Exit"},
		})

	nav := Nav{
		listSelected:      -1,
		mainButtonsCursor: 0,
		installSelected:   make(map[string]registry.Ckan),
	}

	return Bubble{
		appConfig:         cfg,
		theme:             theme,
		primaryViewport:   primaryVP,
		secondaryViewport: secondaryVP,
		splashViewport:    splashVP,
		spinner:           spin,
		textInput:         t,
		inputRequested:    iRequested,
		searchInput:       false,
		ready:             false,
		registry:          reg,
		help:              h,
		nav:               nav,
		activeBox:         constants.PrimaryBoxActive,
		logs:              []string{"Initializing"},
		keyMap:            getKeyMap(),
	}
}

func (b Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, b.getAvailableModsCmd(), b.spinner.Tick, b.MyTickCmd())

	return tea.Batch(cmds...)
}
