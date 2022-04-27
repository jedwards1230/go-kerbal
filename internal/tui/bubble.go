package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/theme"
	"github.com/jedwards1230/go-kerbal/internal/tui/bubbles"
	"github.com/jedwards1230/go-kerbal/registry"
)

type Bubble struct {
	appConfig      config.Config
	theme          theme.Theme
	bubbles        Bubbles
	inputRequested bool
	searchInput    bool
	registry       registry.Registry
	keyMap         bubbles.KeyMap
	nav            Nav
	ready          bool
	activeBox      int
	lastActiveBox  int
	width          int
	height         int
}

type Bubbles struct {
	primaryViewport   bubbles.Viewport
	secondaryViewport bubbles.Viewport
	splashViewport    bubbles.Viewport
	commandViewport   bubbles.Viewport
	help              bubbles.HelpBubble
	spinner           spinner.Model
	textInput         textinput.Model
}

type Nav struct {
	activeMod       registry.Ckan
	installSelected map[string]registry.Ckan
	listSelected    int
	listCursor      int
	menuCursor      int
	boolCursor      bool
}

func InitialModel() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.AppTheme)
	reg := registry.GetRegistry()
	reg.SetTheme(theme)

	iRequested := false
	if cfg.Settings.KerbalDir == "" {
		iRequested = true
	}

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	primaryBoxBorderColor := theme.ActiveBoxBorderColor
	secondaryBoxBorderColor := theme.InactiveBoxBorderColor
	splashBoxBorderColor := theme.InactiveBoxBorderColor

	t := textinput.New()
	t.Prompt = "❯ "
	t.CharLimit = -1
	t.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	primaryVP := bubbles.NewViewport(0, 0)
	primaryVP.Style = lipgloss.NewStyle().
		PaddingLeft(internal.BoxPadding).
		PaddingRight(internal.BoxPadding).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryBoxBorderColor)

	secondaryVP := bubbles.NewViewport(0, 0)
	secondaryVP.Style = lipgloss.NewStyle().
		PaddingLeft(internal.BoxPadding).
		PaddingRight(internal.BoxPadding).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryBoxBorderColor)

	commandVP := bubbles.NewViewport(0, 0)
	commandVP.Style = lipgloss.NewStyle().
		PaddingLeft(internal.BoxPadding).
		PaddingRight(internal.BoxPadding).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryBoxBorderColor)

	splashVP := bubbles.NewViewport(0, 0)
	splashVP.Style = lipgloss.NewStyle().
		PaddingLeft(internal.BoxPadding).
		PaddingRight(internal.BoxPadding).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(splashBoxBorderColor)

	h := bubbles.NewHelpBubble(theme.DefaultTextColor)

	nav := Nav{
		listSelected:    -1,
		menuCursor:      0,
		boolCursor:      false,
		installSelected: make(map[string]registry.Ckan),
	}

	bubs := Bubbles{
		primaryViewport:   primaryVP,
		secondaryViewport: secondaryVP,
		commandViewport:   commandVP,
		splashViewport:    splashVP,
		spinner:           spin,
		textInput:         t,
		help:              h,
	}

	return Bubble{
		appConfig:      cfg,
		theme:          theme,
		bubbles:        bubs,
		inputRequested: iRequested,
		searchInput:    false,
		ready:          false,
		registry:       reg,
		nav:            nav,
		activeBox:      internal.ModListView,
		keyMap:         bubbles.GetKeyMap(),
	}
}

func (b Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, b.getAvailableModsCmd(), b.bubbles.spinner.Tick, b.MyTickCmd())

	return tea.Batch(cmds...)
}
