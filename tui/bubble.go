package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/jedwards1230/go-kerbal/registry/datacollector"
)

type Bubble struct {
	//appConfig         config.Config
	primaryViewport   viewport.Model
	secondaryViewport viewport.Model
	loadingViewport   viewport.Model
	modList           []datacollector.Ckan
	keyMap            KeyMap
	loadingMsg        string
	ready             bool
	loading           bool
	cursor            int
	width             int
	height            int
	selected          int
}

func InitialModel() Bubble {
	primaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorder := lipgloss.NormalBorder()
	loadingBoxBorder := lipgloss.NormalBorder()

	pvp := viewport.New(0, 0)
	pvp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder)
		//BorderForeground(primaryBoxBorderColor)

	svp := viewport.New(0, 0)
	svp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(secondaryBoxBorder)
		//BorderForeground(secondaryBoxBorderColor)

	lvp := viewport.New(0, 0)
	lvp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(loadingBoxBorder)
		//BorderForeground(loadingBoxBorderColor)

	return Bubble{
		primaryViewport:   pvp,
		secondaryViewport: svp,
		loadingViewport:   lvp,
		modList:           registry.BuildRegistry(),
		selected:          -1,
		keyMap:            DefaultKeyMap(),
	}
}

func (b Bubble) Init() tea.Cmd {
	return nil
}
