package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/style"
	"github.com/jedwards1230/go-kerbal/internal/theme"
	"github.com/muesli/reflow/truncate"
)

func (b Bubble) View() string {
	var body string

	splashStyle := style.SplashVP.Render

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
		secondaryTitle := b.styleSecondaryTitle("Go-Kerbal")
		switch b.activeBox {
		case internal.ModListView:
			primaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleSecondaryTitle(b.nav.activeMod.Name)
			}
		case internal.ModInfoView:
			secondaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleSecondaryTitle(b.nav.activeMod.Name)
			}
		case internal.SettingsView:
			secondaryTitle = b.styleSecondaryTitle("Options")
			secondaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
		case internal.SearchView:
			primaryTitle = b.styleTitle("Search Mods")
			primaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleSecondaryTitle(b.nav.activeMod.Name)
			}
		case internal.QueueView:
			primaryTitle = b.styleTitle("Queue")
			primaryBoxBorderColor = theme.AppTheme.ActiveBoxBorderColor
			if !b.nav.listCursorHide {
				secondaryTitle = b.styleSecondaryTitle(b.nav.activeMod.Name)
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
	switch b.activeBox {
	case internal.EnterKspDirView:
		s = trunc(
			s,
			b.bubbles.splashPaginator.Width-2,
		)

		return style.PrimaryTitle.
			Width(b.bubbles.splashPaginator.Width + 2).
			Render(s)
	case internal.LogView:
		s = trunc(
			s,
			b.bubbles.splashPaginator.Width-2,
		)

		return style.PrimaryTitle.
			Width(b.bubbles.splashPaginator.Width + 2).
			Render(s)
	default:
		s = trunc(
			s,
			b.bubbles.primaryPaginator.Width-2,
		)

		return style.PrimaryTitle.
			Width(b.bubbles.primaryPaginator.Width + 2).
			Render(s)
	}
}

func (b Bubble) styleSecondaryTitle(s string) string {
	s = trunc(
		s,
		b.bubbles.secondaryViewport.Width-2,
	)

	return style.SecondaryTitle.
		Width(b.bubbles.secondaryViewport.Width + 2).
		Render(s)
}

func (b Bubble) drawKV(k, v string, color bool) string {
	keyStyle := style.KeyStyle.Width((b.bubbles.secondaryViewport.Width / 3))

	valueStyle := style.ValueStyle

	if color {
		return connectHorz(
			keyStyle.Copy().
				Render(k),
			style.ValueStyle.Copy().
				Foreground(theme.AppTheme.UnselectedListItemColor).
				Background(theme.AppTheme.SelectedListItemColor).
				Render(v))
	} else {
		return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
	}
}

func (b Bubble) drawHelpKV(k, v string) string {
	keyStyle := style.KeyStyle.Copy().Width(11)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Faint(true).
		PaddingRight(2)

	return connectHorz(keyStyle.Render(k), valueStyle.Render(v))
}

func connectHorz(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

func connectVert(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}

func trunc(s string, i int) string {
	return truncate.StringWithTail(
		s,
		uint(i),
		internal.EllipsisStyle)
}

func styleWidth(i int) lipgloss.Style {
	return lipgloss.NewStyle().Width(i)
}
