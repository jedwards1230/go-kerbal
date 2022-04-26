package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal"
)

func (b Bubble) View() string {
	var body string

	switch b.activeBox {
	case internal.LogView:
		b.bubbles.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(b.theme.ActiveBoxBorderColor)
		body = lipgloss.JoinVertical(lipgloss.Top,
			b.styleTitle("Logs"),
			b.bubbles.splashViewport.View(),
		)
	case internal.EnterKspDirView:
		b.bubbles.splashViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(b.theme.ActiveBoxBorderColor)
		body = lipgloss.JoinVertical(lipgloss.Top,
			b.styleTitle("Enter Kerbal Space Program Directory"),
			b.bubbles.splashViewport.View(),
		)
	default:
		var primaryBox string
		var secondaryBox string

		primaryTitle := b.styleTitle("Mod List")
		secondaryTitle := b.styleTitle("Help Menu")
		switch b.activeBox {
		case internal.ModListView, internal.ModInfoView:
			if b.nav.listSelected >= 0 && b.nav.listSelected < len(b.registry.ModMapIndex) {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.SettingsView:
			secondaryTitle = b.styleTitle("Options")
		case internal.SearchView:
			primaryTitle = b.styleTitle("Search Mods")
			if b.nav.listSelected >= 0 && b.nav.listSelected < len(b.registry.ModMapIndex) {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		}

		// set colors
		primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
		secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor

		// format active box
		switch b.activeBox {
		case internal.ModListView, internal.SearchView:
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		case internal.ModInfoView, internal.EnterKspDirView, internal.SettingsView:
			secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		}

		// format views
		b.bubbles.primaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(primaryBoxBorderColor).
			Align(lipgloss.Center)
		primaryBox = lipgloss.JoinVertical(lipgloss.Top,
			primaryTitle,
			b.bubbles.primaryViewport.View(),
		)

		b.bubbles.secondaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.NormalBorder()).
			BorderForeground(secondaryBoxBorderColor)
		secondaryBox = lipgloss.JoinVertical(lipgloss.Top,
			secondaryTitle,
			b.bubbles.secondaryViewport.View(),
		)

		// organize views
		body = connectSides(
			primaryBox,
			secondaryBox,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		b.getMainButtonsView(),
		body,
		b.statusBarView(),
	)
}

func (b Bubble) styleTitle(s string) string {
	switch b.activeBox {
	case internal.LogView, internal.EnterKspDirView:
		return lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Height(3).
			Border(lipgloss.NormalBorder()).
			Padding(1).
			Width(b.bubbles.splashViewport.Width).
			Render(s)
	default:
		return lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Height(3).
			Border(lipgloss.NormalBorder()).
			Padding(1).
			Width(b.bubbles.primaryViewport.Width + 2).
			Render(s)
	}
}

func (b Bubble) drawKV(k, v string, color bool) string {
	keyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Width((b.bubbles.secondaryViewport.Width/4)+3).
		Padding(0, 2)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(0, 2)

	if color {
		return connectSides(
			keyStyle.Copy().
				Render(k),
			valueStyle.Copy().
				Foreground(b.theme.UnselectedListItemColor).
				Background(b.theme.SelectedListItemColor).
				Render(v))
	} else {
		return connectSides(keyStyle.Render(k), valueStyle.Render(v))
	}
}

func connectSides(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}
