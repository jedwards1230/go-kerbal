package tui

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedwards1230/go-kerbal/internal/ckan"
	"github.com/jedwards1230/go-kerbal/internal/dirfs"
	"github.com/jedwards1230/go-kerbal/internal/registry"
	"github.com/spf13/viper"
)

type (
	UpdatedModMapMsg    map[string][]ckan.Ckan
	InstalledModListMsg map[string]interface{}
	UpdateKspDirMsg     bool
	ErrorMsg            error
	SearchMsg           registry.ModIndex
	SortedMsg           map[string]interface{}
	MyTickMsg           map[string]interface{}
)

// Request the mod list from the database
func (b Bubble) getAvailableModsCmd() tea.Cmd {
	return func() tea.Msg {
		b.LogCommand("Checking available mods")
		b.registry.UpdateDB(false)
		updatedModMap := b.registry.GetEntireModList()
		if len(updatedModMap) == 0 {
			b.registry.UpdateDB(true)
			updatedModMap = b.registry.GetEntireModList()
		}
		return UpdatedModMapMsg(updatedModMap)
	}
}

func (b *Bubble) sortModMapCmd() tea.Cmd {
	return func() tea.Msg {
		return SortedMsg{}
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
func (b *Bubble) applyModsCmd() tea.Cmd {
	return func() tea.Msg {

		// Remove Mods
		if b.registry.Queue.RemoveLen() > 0 {
			b.LogCommandf("Removing %d mods", b.registry.Queue.RemoveLen())
			err := b.registry.RemoveMods()
			if err != nil {
				return fmt.Errorf("error removing: %v", err)
			}
		}

		// Install Mods
		if b.registry.Queue.InstallLen() > 0 {
			b.LogCommandf("Downloading %d mods", b.registry.Queue.InstallLen())
			err := b.registry.DownloadMods()
			if err != nil {
				return fmt.Errorf("error downloading: %v", err)
			}

			b.LogCommandf("Installing %d mods", b.registry.Queue.InstallLen())
			err = b.registry.InstallMods()
			if err != nil {
				return fmt.Errorf("error installing: %v", err)
			}
		}
		return InstalledModListMsg{}
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
		time.Sleep(500 * time.Millisecond)
		return MyTickMsg{}
	}
}
