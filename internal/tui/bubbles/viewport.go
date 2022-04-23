package bubbles

import (
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// New returns a new model with the given width and height as well as default
// keymappings.
func NewViewport(width, height int) (m Viewport) {
	m.Width = width
	m.Height = height
	m.setInitialValues()
	return m
}

// Viewport is the Bubble Tea model for this viewport element.
type Viewport struct {
	Width  int
	Height int
	KeyMap KeyMap

	// Whether or not to respond to the mouse. The mouse must be enabled in
	// Bubble Tea for this to work. For details, see the Bubble Tea docs.
	MouseWheelEnabled bool

	// The number of lines the mouse wheel will scroll. By default, this is 3.
	MouseWheelDelta int

	// YOffset is the vertical scroll position.
	YOffset int

	// YPosition is the position of the viewport in relation to the terminal
	// window. It's used in high performance rendering only.
	YPosition int

	// Style applies a lipgloss style to the viewport. Realistically, it's most
	// useful for setting borders, margins and padding.
	Style lipgloss.Style

	initialized bool
	lines       []string
}

func (m *Viewport) setInitialValues() {
	m.KeyMap = GetKeyMap()
	m.MouseWheelEnabled = true
	m.MouseWheelDelta = 3
	m.initialized = true
}

// Init exists to satisfy the tea.Model interface for composability purposes.
func (m Viewport) Init() tea.Cmd {
	return nil
}

// AtTop returns whether or not the viewport is in the very top position.
func (m Viewport) AtTop() bool {
	return m.YOffset <= 0
}

// AtBottom returns whether or not the viewport is at or past the very bottom
// position.
func (m Viewport) AtBottom() bool {
	return m.YOffset >= m.maxYOffset()
}

// PastBottom returns whether or not the viewport is scrolled beyond the last
// line. This can happen when adjusting the viewport height.
func (m Viewport) PastBottom() bool {
	return m.YOffset > m.maxYOffset()
}

// ScrollPercent returns the amount scrolled as a float between 0 and 1.
func (m Viewport) ScrollPercent() float64 {
	if m.Height >= len(m.lines) {
		return 1.0
	}
	y := float64(m.YOffset)
	h := float64(m.Height)
	t := float64(len(m.lines) - 1)
	v := y / (t - h)
	return math.Max(0.0, math.Min(1.0, v))
}

// SetContent set the pager's text content. For high performance rendering the
// Sync command should also be called.
func (m *Viewport) SetContent(s string) {
	s = strings.ReplaceAll(s, "\r\n", "\n") // normalize line endings
	m.lines = strings.Split(s, "\n")

	if m.YOffset > len(m.lines)-1 {
		m.GotoBottom()
	}
}

// maxYOffset returns the maximum possible value of the y-offset based on the
// viewport's content and set height.
func (m Viewport) maxYOffset() int {
	return max(0, len(m.lines)-m.Height)
}

// visibleLines returns the lines that should currently be visible in the
// viewport.
func (m Viewport) visibleLines() (lines []string) {
	if len(m.lines) > 0 {
		top := max(0, m.YOffset)
		bottom := clamp(m.YOffset+m.Height, top, len(m.lines))
		lines = m.lines[top:bottom]
	}
	return lines
}

// scrollArea returns the scrollable boundaries for high performance rendering.
func (m Viewport) scrollArea() (top, bottom int) {
	top = max(0, m.YPosition)
	bottom = max(top, top+m.Height)
	if top > 0 && bottom > top {
		bottom--
	}
	return top, bottom
}

// SetYOffset sets the Y offset.
func (m *Viewport) SetYOffset(n int) {
	m.YOffset = clamp(n, 0, m.maxYOffset())
}

// ViewDown moves the view down by the number of lines in the viewport.
// Basically, "page down".
func (m *Viewport) ViewDown() []string {
	if m.AtBottom() {
		return nil
	}

	m.SetYOffset(m.YOffset + m.Height)
	return m.visibleLines()
}

// ViewUp moves the view up by one height of the viewport. Basically, "page up".
func (m *Viewport) ViewUp() []string {
	if m.AtTop() {
		return nil
	}

	m.SetYOffset(m.YOffset - m.Height)
	return m.visibleLines()
}

// HalfViewDown moves the view down by half the height of the viewport.
func (m *Viewport) HalfViewDown() (lines []string) {
	if m.AtBottom() {
		return nil
	}

	m.SetYOffset(m.YOffset + m.Height/2)
	return m.visibleLines()
}

// HalfViewUp moves the view up by half the height of the viewport.
func (m *Viewport) HalfViewUp() (lines []string) {
	if m.AtTop() {
		return nil
	}

	m.SetYOffset(m.YOffset - m.Height/2)
	return m.visibleLines()
}

// LineDown moves the view down by the given number of lines.
func (m *Viewport) LineDown(n int) (lines []string) {
	if m.AtBottom() || n == 0 {
		return nil
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we actually have left before we reach
	// the bottom.
	m.SetYOffset(m.YOffset + n)
	return m.visibleLines()
}

// LineUp moves the view down by the given number of lines. Returns the new
// lines to show.
func (m *Viewport) LineUp(n int) (lines []string) {
	if m.AtTop() || n == 0 {
		return nil
	}

	// Make sure the number of lines by which we're going to scroll isn't
	// greater than the number of lines we are from the top.
	m.SetYOffset(m.YOffset - n)
	return m.visibleLines()
}

// GotoTop sets the viewport to the top position.
func (m *Viewport) GotoTop() (lines []string) {
	if m.AtTop() {
		return nil
	}

	m.SetYOffset(0)
	return m.visibleLines()
}

// GotoBottom sets the viewport to the bottom position.
func (m *Viewport) GotoBottom() (lines []string) {
	m.SetYOffset(m.maxYOffset())
	return m.visibleLines()
}

// ViewDown is a high performance command that moves the viewport up by a given
// numer of lines. Use Model.ViewDown to get the lines that should be rendered.
// For example:
//
//     lines := model.ViewDown(1)
//     cmd := ViewDown(m, lines)
//
func ViewDown(m Viewport, lines []string) tea.Cmd {
	if len(lines) == 0 {
		return nil
	}
	top, bottom := m.scrollArea()
	return tea.ScrollDown(lines, top, bottom)
}

// ViewUp is a high performance command the moves the viewport down by a given
// number of lines height. Use Model.ViewUp to get the lines that should be
// rendered.
func ViewUp(m Viewport, lines []string) tea.Cmd {
	if len(lines) == 0 {
		return nil
	}
	top, bottom := m.scrollArea()
	return tea.ScrollUp(lines, top, bottom)
}

// Update handles standard message-based viewport updates.
func (m Viewport) Update(msg tea.Msg) (Viewport, tea.Cmd) {
	var cmd tea.Cmd
	m, cmd = m.updateAsModel(msg)
	return m, cmd
}

// Author's note: this method has been broken out to make it easier to
// potentially transition Update to satisfy tea.Model.
func (m Viewport) updateAsModel(msg tea.Msg) (Viewport, tea.Cmd) {
	if !m.initialized {
		m.setInitialValues()
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.PageDown):
			m.ViewDown()
		case key.Matches(msg, m.KeyMap.PageUp):
			m.ViewUp()
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.HalfViewDown()
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.HalfViewUp()
		case key.Matches(msg, m.KeyMap.Down):
			m.LineDown(1)

		case key.Matches(msg, m.KeyMap.Up):
			m.LineUp(1)
		}
	case tea.MouseMsg:
		if !m.MouseWheelEnabled {
			break
		}
		switch msg.Type {
		case tea.MouseWheelUp:
			m.LineUp(m.MouseWheelDelta)
		case tea.MouseWheelDown:
			m.LineDown(m.MouseWheelDelta)
		}
	}

	return m, cmd
}

// View renders the viewport into a string.
func (m Viewport) View() string {

	lines := m.visibleLines()

	// Fill empty space with newlines
	extraLines := ""
	if len(lines) < m.Height {
		extraLines = strings.Repeat("\n", max(0, m.Height-len(lines)))
	}

	return m.Style.Copy().
		UnsetWidth().
		UnsetHeight().
		Render(strings.Join(lines, "\n") + extraLines)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
