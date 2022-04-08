package config

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/jedwards1230/go-kerbal/dirfs"
)

// SettingsConfig struct represents the config for the settings.
type (
	SettingsConfig struct {
		KerbalDir            string `mapstructure:"kerbal_dir"`
		KerbalVer            string `mapstructure:"kerbal_ver"`
		MetaRepo             string `mapstructure:"meta_repo"`
		LastRepoHash         string `mapstructure:"last_repo_hash"`
		EnableLogging        bool   `mapstructure:"enable_logging"`
		EnableMouseWheel     bool   `mapstructure:"enable_mousewheel"`
		HideIncompatibleMods bool   `mapstructure:"hide_incompatible"`
	}

	SyntaxThemeConfig struct {
		Light string `mapstructure:"light"`
		Dark  string `mapstructure:"dark"`
	}

	ThemeConfig struct {
		AppTheme    string            `mapstructure:"app_theme"`
		SyntaxTheme SyntaxThemeConfig `mapstructure:"syntax_theme"`
	}

	Config struct {
		Settings SettingsConfig `mapstructure:"settings"`
		Theme    ThemeConfig    `mapstructure:"theme"`
	}
)

// LoadConfig loads a users config and creates the config if it does not exist
// located at ~/.config/fm.yml.
func LoadConfig() {

	/* if runtime.GOOS != "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		err = dirfs.CreateDirectory(filepath.Join(homeDir, ".config", "go-kerbal"))
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath("$HOME/.config/go-kerbal")
	} else {
		viper.AddConfigPath("$HOME")
	} */

	viper.AddConfigPath(dirfs.RootDir())

	viper.SetConfigName("config")
	viper.SetConfigType("json")

	kerbalDir := dirfs.FindKspPath()
	kerbalVer := dirfs.FindKspVersion(kerbalDir)

	// Setup config defaults.
	viper.SetDefault("settings.kerbal_dir", kerbalDir)
	viper.SetDefault("settings.kerbal_ver", kerbalVer.String())
	viper.SetDefault("settings.meta_repo", "https://github.com/KSP-CKAN/CKAN-meta.git")
	viper.SetDefault("settings.last_repo_hash", "")
	viper.SetDefault("settings.enable_logging", true)
	viper.SetDefault("settings.enable_mousewheel", true)
	viper.SetDefault("settings.hide_incompatible", true)
	viper.SetDefault("theme.app_theme", "default")
	viper.SetDefault("theme.syntax_theme.light", "pygments")
	viper.SetDefault("theme.syntax_theme.dark", "dracula")

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfig()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %v", e.Name)
	})
	viper.WatchConfig()
}

// GetConfig returns the users config.
func GetConfig() (config Config) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}
