package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal"
)

func (b Bubble) View() string {
	var body string

	b.bubbles.splashViewport.Style = lipgloss.NewStyle().
		PaddingLeft(internal.BoxPadding).
		PaddingRight(internal.BoxPadding).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(b.theme.ActiveBoxBorderColor)

	switch b.activeBox {
	case internal.LogView:
		body = connectVert(
			b.styleTitle("Logs"),
			b.bubbles.splashViewport.View(),
		)
	case internal.EnterKspDirView:
		body = connectVert(
			b.styleTitle("Enter Kerbal Space Program Directory"),
			b.bubbles.splashViewport.View(),
		)
	default:
		var primaryBox string
		var secondaryBox string

		// set colors
		primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
		secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor

		primaryTitle := b.styleTitle("Mod List")
		secondaryTitle := b.styleTitle("Go-Kerbal")
		switch b.activeBox {
		case internal.ModListView:
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.ModInfoView:
			secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.SettingsView:
			secondaryTitle = b.styleTitle("Options")
			secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor
		case internal.SearchView:
			primaryTitle = b.styleTitle("Search Mods")
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.QueueView:
			primaryTitle = b.styleTitle("Queue")
			primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		}

		pageStyle := lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryBoxBorderColor).
			Align(lipgloss.Center).Render

		primaryBox = connectVert(
			primaryTitle,
			pageStyle(b.bubbles.paginator.GetContent()),
		)

		b.bubbles.secondaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(internal.BoxPadding).
			PaddingRight(internal.BoxPadding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryBoxBorderColor)
		secondaryBox = connectVert(
			secondaryTitle,
			b.bubbles.secondaryViewport.View(),
			b.bubbles.commandViewport.View(),
		)

		// organize views
		body = connectHorz(
			primaryBox,
			secondaryBox,
		)
	}

	return connectVert(
		//b.getMainButtonsView(),
		body,
		b.statusBarView(),
	)
}

func (b Bubble) styleTitle(s string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Height(3).
		Border(lipgloss.RoundedBorder()).
		Padding(1)
		//Margin(1, 0)

	switch b.activeBox {
	case internal.LogView, internal.EnterKspDirView:
		return titleStyle.
			Width(b.bubbles.splashViewport.Width + 2).
			Render(s)
	default:
		return titleStyle.
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
		return connectHorz(
			keyStyle.Copy().
				Render(k),
			valueStyle.Copy().
				Foreground(b.theme.UnselectedListItemColor).
				Background(b.theme.SelectedListItemColor).
				Render(v))
	} else {
		return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
	}
}

func (b Bubble) drawHelpKV(k, v string) string {
	keyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Width(b.bubbles.secondaryViewport.Width/6).
		Padding(0, 1, 0, 2)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Faint(true).
		PaddingRight(3)

	return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
}

func connectHorz(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

func connectVert(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
