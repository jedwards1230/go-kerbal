package tui

import "github.com/charmbracelet/lipgloss"

func (b Bubble) homeView() string {
	contentStyle := lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width - 5).
		Padding(2).
		Render

	content := "" +
		"To do:\n " +
		"- Display error messages\n " +
		"- Make install queue editable\n " +
		"- Window resizing on Windows\n " +
		"- Better mod info formatting\n " +
		"- Remember last page when swapping views\n " +
		"- Info screens per view\n " +
		"- Mouse clicking\n " +
		"- Ensure mods install/uninstall properly\n " +
		"- Dynamic command view\n "

	return contentStyle(content)
}

func (b Bubble) inputKspView() string {
	question := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(b.width).
		Padding(1).
		Render("Please enter the path to your Kerbal Space Program directory:")

	inText := ""
	if b.inputRequested {
		inText = b.bubbles.textInput.View()
	} else {
		inText = b.bubbles.textInput.Value()
		inText += "\n\nPress Esc to close"
	}

	inText = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(b.width).
		Padding(1).
		Render(inText)

	return connectVert(
		question,
		inText,
	)
}

func (b Bubble) helpView() string {
	leftColumn := []string{
		b.drawHelpKV("up", "Move up"),
		b.drawHelpKV("down", "Move down"),
		b.drawHelpKV("space", "Toggle mod info"),
		b.drawHelpKV("enter", "Add to queue"),
		b.drawHelpKV("tab", "Swap active windows"),
	}

	rightColumn := []string{
		b.drawHelpKV("1", "Refresh"),
		b.drawHelpKV("2", "Search"),
		b.drawHelpKV("3", "Apply"),
		b.drawHelpKV("0", "Settings"),
		b.drawHelpKV("shift+o", "Logs"),
	}

	content := connectHorz(connectVert(leftColumn...), connectVert(rightColumn...))

	content = lipgloss.NewStyle().
		Padding(1).
		Margin(1, 0).
		Render(content)

	return lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width).
		Render(content)
}

func (b Bubble) getBoolOptionsView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(b.bubbles.commandViewport.Width-4).
		Padding(1, 2, 0)

	optionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		BorderForeground(b.theme.InactiveBoxBorderColor).
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
	options = lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width - 4).
		Align(lipgloss.Center).
		Render(options)

	content := connectVert(
		titleStyle.Render("Apply?"),
		options,
	)

	return lipgloss.NewStyle().
		Width(b.bubbles.commandViewport.Width).
		//Padding(3).
		Render(content)
}
