package tui

import (
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

func (b *Bubble) LogCommand(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(b.theme.Blue).Render(msg))
}

func (b *Bubble) LogCommandf(format string, a ...interface{}) {
	b.LogCommand(fmt.Sprintf(format, a...))
}

func (b *Bubble) LogSuccess(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(b.theme.Green).Render(msg))
}

func (b *Bubble) LogSuccessf(format string, a ...interface{}) {
	b.LogCommand(fmt.Sprintf(format, a...))
}

func (b *Bubble) LogError(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(b.theme.Red).Render(msg))
}

func (b *Bubble) LogErrorf(format string, a ...interface{}) {
	b.LogError(fmt.Sprintf(format, a...))
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
