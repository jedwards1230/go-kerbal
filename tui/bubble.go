package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/ckan"
	"github.com/jedwards1230/go-kerbal/internal/config"
	"github.com/jedwards1230/go-kerbal/internal/keymap"
	"github.com/jedwards1230/go-kerbal/internal/paginator"
	"github.com/jedwards1230/go-kerbal/internal/registry"
	"github.com/jedwards1230/go-kerbal/internal/style"
	"github.com/jedwards1230/go-kerbal/internal/theme"
	"github.com/jedwards1230/go-kerbal/internal/viewport"
)

type Bubble struct {
	appConfig      config.Config
	bubbles        Bubbles
	inputRequested bool
	searchInput    bool
	registry       registry.Registry
	keyMap         keymap.KeyMap
	logs           []string
	nav            Nav
	ready          bool
	activeBox      int
	lastActiveBox  int
	width          int
	height         int
}

type Bubbles struct {
	primaryPaginator  paginator.Paginator
	splashPaginator   paginator.Paginator
	secondaryViewport viewport.Viewport
	commandViewport   viewport.Viewport
	spinner           spinner.Model
	textInput         textinput.Model
}

type Nav struct {
	activeMod      ckan.Ckan
	listCursorHide bool
	listCursor     int
	menuCursor     int
	boolCursor     bool
}

func InitialModel() Bubble {
	cfg := config.GetConfig()
	theme.SetTheme(cfg.AppTheme)
	reg := registry.New()

	iRequested := false
	if cfg.Settings.KerbalDir == "" {
		iRequested = true
	}

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = style.Spinner

	t := textinput.New()
	t.Prompt = "❯ "
	t.CharLimit = -1
	t.PlaceholderStyle = style.TextInput

	secondaryVP := viewport.NewViewport(0, 0)
	secondaryVP.Style = style.SecondaryVP

	commandVP := viewport.NewViewport(0, 0)
	commandVP.Style = style.CommandVP

	pages := paginator.New()
	pages.SetTotalPages(1)

	splashPages := paginator.New()
	splashPages.SetTotalPages(1)

	nav := Nav{
		listCursorHide: true,
		menuCursor:     0,
		boolCursor:     false,
	}

	bubs := Bubbles{
		secondaryViewport: secondaryVP,
		commandViewport:   commandVP,
		spinner:           spin,
		textInput:         t,
		primaryPaginator:  pages,
		splashPaginator:   splashPages,
	}

	return Bubble{
		appConfig:      cfg,
		bubbles:        bubs,
		inputRequested: iRequested,
		searchInput:    false,
		ready:          false,
		registry:       reg,
		nav:            nav,
		activeBox:      internal.ModListView,
		keyMap:         keymap.New(),
	}
}

func (b Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, b.getAvailableModsCmd(), b.bubbles.spinner.Tick, b.TickCmd())

	return tea.Batch(cmds...)
}
