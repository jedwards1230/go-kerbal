package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"

	"github.com/jedwards1230/go-kerbal/dirfs"
)

// SettingsConfig struct represents the config for the settings.
type (
	SettingsConfig struct {
		KerbalDir           string `mapstructure:"kerbal_dir"`
		StartDir            string `mapstructure:"start_dir"`
		ShowIcons           bool   `mapstructure:"show_icons"`
		EnableLogging       bool   `mapstructure:"enable_logging"`
		EnableMouseWheel    bool   `mapstructure:"enable_mousewheel"`
		PrettyMarkdown      bool   `mapstructure:"pretty_markdown"`
		Borderless          bool   `mapstructure:"borderless"`
		SimpleMode          bool   `mapstructure:"simple_mode"`
		CalculatedFileSizes bool   `mapstructure:"calculated_file_sizes"`
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

	if runtime.GOOS != "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		err = dirfs.CreateDirectory(filepath.Join(homeDir, ".config", "go-kerbal"))
		if err != nil {
			log.Fatal(err)
		}

		// TODO: Adjust where config is stored
		//viper.AddConfigPath("$HOME/.config/go-kerbal")
		viper.AddConfigPath("./")
	} else {
		//viper.AddConfigPath("$HOME")
		viper.AddConfigPath("./")
	}

	viper.SetConfigName("go-kerbal")
	viper.SetConfigType("yml")

	// Setup config defaults.
	viper.SetDefault("settings.kerbal_dir", dirfs.FindKspPath())
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("settings.enable_logging", true)
	viper.SetDefault("settings.enable_mousewheel", true)
	viper.SetDefault("settings.pretty_markdown", true)
	viper.SetDefault("settings.borderless", false)
	viper.SetDefault("settings.simple_mode", false)
	viper.SetDefault("settings.calculated_file_sizes", false)
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

	// Setup flag defaults.
	viper.SetDefault("start-dir", "")
	viper.SetDefault("selection-path", "")
}

// GetConfig returns the users config.
func GetConfig() (config Config) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}
