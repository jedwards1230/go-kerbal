package registry

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/internal/theme"
)

type SortOptions struct {
	SortTag   string
	SortOrder string
}

type versions struct {
	Epoch  string
	Mod    string
	KspMin string
	KspMax string
	Spec   string
}

type download struct {
	Downloaded bool
	URL        string
	Path       string
}

type install struct {
	Installed bool
	FindRegex string
	Find      string
	File      string
	InstallTo string
}

type resource struct {
	Homepage    string
	Spacedock   string
	Repository  string
	XScreenshot string
}

type Entry struct {
	Key      string
	SearchBy string
}

func (entry ModIndex) Len() int           { return len(entry) }
func (entry ModIndex) Less(i, j int) bool { return entry[i].SearchBy < entry[j].SearchBy }
func (entry ModIndex) Swap(i, j int)      { entry[i], entry[j] = entry[j], entry[i] }

func (r *Registry) SetTheme(t theme.Theme) {
	r.theme = t
}

func (r *Registry) LogCommand(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(r.theme.Blue).Render(msg))
}

func (r *Registry) LogCommandf(format string, a ...interface{}) {
	r.LogCommand(fmt.Sprintf(format, a...))
}

func (r *Registry) LogSuccess(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(r.theme.Green).Render(msg))
}

func (r *Registry) LogSuccessf(format string, a ...interface{}) {
	r.LogSuccess(fmt.Sprintf(format, a...))
}

func (r *Registry) LogWarning(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(r.theme.Orange).Render(msg))
}

func (r *Registry) LogWarningf(format string, a ...interface{}) {
	r.LogWarning(fmt.Sprintf(format, a...))
}

func (r *Registry) LogError(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(r.theme.Red).Render(msg))
}

func (r *Registry) LogErrorf(format string, a ...interface{}) {
	r.LogError(fmt.Sprintf(format, a...))
}
