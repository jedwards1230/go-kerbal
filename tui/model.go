package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/cmd/constants"
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/jedwards1230/go-kerbal/registry/module"
)

type Bubble struct {
	//appConfig         config.Config
	primaryViewport     viewport.Model
	secondaryViewport   viewport.Model
	modList             []module.ModuleVersion
	keyMap              KeyMap
	cursor              int
	width               int
	height              int
	ready               bool
	loadingScreenActive bool
	selected            int
}

func InitialModel() Bubble {
	primaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorder := lipgloss.NormalBorder()

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

	return Bubble{
		primaryViewport:   pvp,
		secondaryViewport: svp,
		modList:           registry.BuildRegistry(),
		selected:          -1,
		keyMap:            DefaultKeyMap(),
	}
}

func (b Bubble) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
