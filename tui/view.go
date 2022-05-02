package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/theme"
	"github.com/muesli/reflow/truncate"
)

func (b Bubble) View() string {
	var body string

	splashStyle := lipgloss.NewStyle().
		PaddingLeft(internal.BoxPadding).
		PaddingRight(internal.BoxPadding).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.AppTheme.ActiveBoxBorderColor).
		Align(lipgloss.Center).Render

	switch b.activeBox {
	case internal.LogView:
		body = connectVert(
			b.styleTitle("Logs"),
			splashStyle(b.bubbles.splashPaginator.GetContent()),
		)
	case internal.EnterKspDirView:
		body = connectVert(
			b.styleTitle("Enter Kerbal Space Program Directory"),
			splashStyle(b.bubbles.splashPaginator.GetContent()),
		)
	default:
		var primaryBox string
		var secondaryBox string

		// set colors
		primaryBoxBorderColor := theme.AppTheme.InactiveBoxBorderColor
		secondaryBoxBorderColor := theme.AppTheme.InactiveBoxBorderColor

		primaryTitle := b.styleTitle("Mod List")
		secondaryTitle := b.styleTitle("Go-Kerbal")
		switch b.activeBox {
		case internal.ModListView:
			primaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.ModInfoView:
			secondaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.SettingsView:
			secondaryTitle = b.styleTitle("Options")
			secondaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
		case internal.SearchView:
			primaryTitle = b.styleTitle("Search Mods")
			primaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleTitle(b.nav.activeMod.Name)
			}
		case internal.QueueView:
			primaryTitle = b.styleTitle("Queue")
			primaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
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
			pageStyle(b.bubbles.primaryPaginator.GetContent()),
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
		Padding(1, 0)

	switch b.activeBox {
	case internal.EnterKspDirView:
		s = truncate.StringWithTail(
			s,
			uint(b.bubbles.splashPaginator.Width-2),
			internal.EllipsisStyle)

		return titleStyle.
			Width(b.bubbles.splashPaginator.Width + 2).
			Render(s)
	case internal.LogView:
		s = truncate.StringWithTail(
			s,
			uint(b.bubbles.splashPaginator.Width-2),
			internal.EllipsisStyle)

		return titleStyle.
			Width(b.bubbles.splashPaginator.Width + 2).
			Render(s)
	default:
		s = truncate.StringWithTail(
			s,
			uint(b.bubbles.secondaryViewport.Width-2),
			internal.EllipsisStyle)

		return titleStyle.
			Width(b.bubbles.secondaryViewport.Width + 2).
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
				Foreground(theme.AppTheme.UnselectedListItemColor).
				Background(theme.AppTheme.SelectedListItemColor).
				Render(v))
	} else {
		return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
	}
}

func (b Bubble) drawHelpKV(k, v string) string {
	keyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Bold(true).
		Width(b.bubbles.commandViewport.Width / 5).
		PaddingLeft(2)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Faint(true).
		PaddingRight(4)

	return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
}

func connectHorz(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

func connectVert(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
