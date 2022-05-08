package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal"
	"github.com/jedwards1230/go-kerbal/internal/theme"
)

func (b Bubble) homeView() string {
	contentStyle := styleWidth(b.bubbles.commandViewport.Width).
		PaddingLeft(1).
		Render

	var content string
	switch b.activeBox {
	case internal.QueueView:
		content = "" +
			fmt.Sprintf("Installing %d mods \n", b.registry.Queue.InstallLen()) +
			fmt.Sprintf("Removing %d mods \n", b.registry.Queue.RemoveLen()) +
			"\n" +
			"Press up/down to scroll the list \n" +
			"Press enter to remove the selected mod \n" +
			"\n" +
			"Press tab to get back to the confirmation window \n"
		content = styleWidth(b.bubbles.secondaryViewport.Width).
			Align(lipgloss.Left).
			Render(content)
	default:
		content = "" +
			"To do:\n " +
			"- New layout, probably\n " +
			"- Split options screen with a more button \n " +
			"- Display error messages\n " +
			"- Better KSP version check on Windows/Linux \n " +
			"- Map dependencies for uninstalls\n " +
			"- Window resizing on Windows\n " +
			"- Better mod info formatting\n " +
			"- Remember last page when swapping views\n " +
			"- Info screens per view\n " +
			"- Mouse clicking\n " +
			"- Ensure mods install/uninstall properly\n " +
			"- Dynamic command view\n "
	}

	return contentStyle(content)
}

func (b Bubble) inputKspView() string {
	question := styleWidth(b.width).
		Align(lipgloss.Left).
		Padding(1).
		Render("Please enter the path to your Kerbal Space Program directory:")

	inText := ""
	if b.inputRequested {
		inText = b.bubbles.textInput.View()
	} else {
		inText = b.bubbles.textInput.Value()
		inText += "\n\nPress Esc to close"
	}

	inText = styleWidth(b.width).
		Align(lipgloss.Left).
		Padding(1).
		Render(inText)

	content := connectVert(
		question,
		inText,
	)

	return styleWidth(b.bubbles.splashPaginator.Width).
		Height(b.bubbles.splashPaginator.Height + 1).
		Render(content)
}

// todo: make this easier to use between different views with different inputs
func (b Bubble) helpView() string {
	leftColumn := []string{
		b.drawHelpKV("up", "Move up"),
		b.drawHelpKV("down", "Move down"),
		b.drawHelpKV("space", "Toggle mod info"),
		b.drawHelpKV("enter", "Add to queue"),
		b.drawHelpKV("tab", "Swap windows"),
	}

	rightColumn := []string{
		b.drawHelpKV("1", "Refresh"),
		b.drawHelpKV("2", "Search"),
		b.drawHelpKV("3", "Apply"),
		b.drawHelpKV("0", "Settings"),
		b.drawHelpKV("shift+o", "Logs"),
	}

	var content string
	if b.bubbles.commandViewport.Width >= 49 {
		content = connectHorz(connectVert(leftColumn...), connectVert(rightColumn...))
	} else {
		content = connectVert(connectVert(leftColumn...), connectVert(rightColumn...))
	}

	content = lipgloss.NewStyle().
		Margin(1, 0).
		Render(content)

	return styleWidth(b.bubbles.commandViewport.Width).
		Align(lipgloss.Left).
		Padding(1, 1).
		Render(content)
}

func (b Bubble) getBoolOptionsView() string {
	titleStyle := styleWidth(b.bubbles.commandViewport.Width-4).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 2, 0)

	optionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		BorderForeground(theme.AppTheme.InactiveBoxBorderColor).
		Padding(0, 4).
		Faint(true).
		Margin(1, 1)

	cancel := optionStyle.Render("Cancel")
	confirm := optionStyle.Render("Confirm")

	if b.nav.listCursorHide {
		if b.nav.boolCursor {
			confirm = optionStyle.Copy().
				Border(lipgloss.RoundedBorder()).
				Faint(false).
				Render("Confirm")
			cancel = optionStyle.Copy().
				Render("Cancel")
		} else {
			cancel = optionStyle.Copy().
				Border(lipgloss.RoundedBorder()).
				Faint(false).
				Render("Cancel")
			confirm = optionStyle.Copy().
				Render("Confirm")
		}
	}

	options := connectHorz(cancel, "  ", confirm)
	options = styleWidth(b.bubbles.commandViewport.Width - 4).
		Align(lipgloss.Center).
		Render(options)

	content := connectVert(
		titleStyle.Render("Apply?"),
		options,
	)

	return styleWidth(b.bubbles.commandViewport.Width).
		Height(b.bubbles.commandViewport.Height).
		Render(content)
}
