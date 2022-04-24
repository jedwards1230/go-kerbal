package tui

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/jedwards1230/go-kerbal/registry"
	"github.com/spf13/viper"
)

type (
	UpdatedModMapMsg    map[string][]registry.Ckan
	InstalledModListMsg map[string]bool
	UpdateKspDirMsg     bool
	ErrorMsg            error
	MyTickMsg           bool
	SearchMsg           registry.ModIndex
	SortedMsg           bool
)

// Request the mod list from the database
func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		b.LogCommand("Checking available mods")
		b.registry.DB.UpdateDB(false)
		updatedModMap := b.registry.GetTotalModMap()
		if len(updatedModMap) == 0 {
			b.registry.DB.UpdateDB(true)
			updatedModMap = b.registry.GetTotalModMap()
		}
		return UpdatedModMapMsg(updatedModMap)
	}
}

func (b *Bubble) sortModMapCmd() tea.Cmd {
	return func() tea.Msg {
		return SortedMsg(true)
	}
}

// Manually input KSP directory
func (b Bubble) updateKspDirCmd(s string) tea.Cmd {
	return func() tea.Msg {
		b.LogCommandf("Checking dir: %s", s)
		kerbalDir, err := dirfs.FindKspPath(s)
		if err != nil || kerbalDir == "" {
			b.LogErrorf("Error finding KSP directory: %v, %s", err, s)
			return UpdateKspDirMsg(false)
		}
		kerbalVer := dirfs.FindKspVersion(kerbalDir)
		viper.Set("settings.kerbal_dir", kerbalDir)
		viper.Set("settings.kerbal_ver", kerbalVer.String())
		viper.WriteConfigAs(viper.ConfigFileUsed())
		log.Printf("Kerbal dir: " + kerbalDir + "/")
		log.Printf("Kerbal Version: %v", kerbalVer)
		return UpdateKspDirMsg(true)
	}
}

// Download selected mods
func (b Bubble) downloadModCmd() tea.Cmd {
	return func() tea.Msg {
		mods, err := b.registry.DownloadMods(b.nav.installSelected)
		if err != nil {
			b.LogErrorf("Error downloading: %v", err)
			return ErrorMsg(err)
		}

		err = registry.InstallMods(mods)
		if err != nil {
			b.LogErrorf("Error downloading: %v", err)
			return ErrorMsg(err)
		}

		installedModList, err := dirfs.CheckInstalledMods()
		if err != nil {
			return ErrorMsg(err)
		}
		return InstalledModListMsg(installedModList)
	}
}

// Download selected mods
func (b Bubble) searchCmd(s string) tea.Cmd {
	return func() tea.Msg {
		searchMapIndex, err := b.registry.BuildSearchIndex(s)
		if err != nil {
			b.LogErrorf("Error building search index: %v", err)
			return SearchMsg(searchMapIndex)
		}
		return SearchMsg(searchMapIndex)
	}
}

func (b Bubble) MyTickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return MyTickMsg(true)
	}
}

func (b *Bubble) LogCommand(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(b.theme.Green).Render(msg))
}

func (b *Bubble) LogCommandf(format string, a ...interface{}) {
	b.LogCommand(fmt.Sprintf(format, a...))
}

func (b *Bubble) LogError(msg string) {
	log.Print(lipgloss.NewStyle().Foreground(b.theme.Red).Render(msg))
}

func (b *Bubble) LogErrorf(format string, a ...interface{}) {
	b.LogError(fmt.Sprintf(format, a...))
}