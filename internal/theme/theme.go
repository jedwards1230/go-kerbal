package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents the properties that make up a theme.
type Theme struct {
	LightGray lipgloss.AdaptiveColor
	Blue      lipgloss.AdaptiveColor
	Green     lipgloss.AdaptiveColor
	Red       lipgloss.AdaptiveColor
	Orange    lipgloss.AdaptiveColor

	DefaultTextColor lipgloss.AdaptiveColor
	ErrorColor       lipgloss.AdaptiveColor

	InstalledColor    lipgloss.AdaptiveColor
	IncompatibleColor lipgloss.AdaptiveColor

	SelectedListItemColor   lipgloss.AdaptiveColor
	UnselectedListItemColor lipgloss.AdaptiveColor
	ActiveBoxBorderColor    lipgloss.AdaptiveColor
	InactiveBoxBorderColor  lipgloss.AdaptiveColor
}

// appColors contains the different types of colors.
type appColors struct {
	white                  string
	black                  string
	lightGray              string
	green                  string
	orange                 string
	slightlyBrighterOrange string
	red                    string
	slightlyBrighterRed    string
	blue                   string
}

// Colors contains the different kinds of Colors and their values.
var Colors = appColors{
	white:                  "#f2faff",
	black:                  "#000000",
	lightGray:              "#bcbfc2",
	green:                  "#00bf0d",
	orange:                 "#cf8611",
	slightlyBrighterOrange: "#ffff00",
	red:                    "#cc241d",
	slightlyBrighterRed:    "#ff0a00",
	blue:                   "#0091ff",
}

// themeMap represents the mapping of different themes.
var themeMap = map[string]Theme{
	"default": {
		LightGray: lipgloss.AdaptiveColor{Dark: Colors.lightGray, Light: Colors.lightGray},
		Blue:      lipgloss.AdaptiveColor{Dark: Colors.blue, Light: Colors.blue},
		Green:     lipgloss.AdaptiveColor{Dark: Colors.green, Light: Colors.green},
		Red:       lipgloss.AdaptiveColor{Dark: Colors.slightlyBrighterRed, Light: Colors.red},
		Orange:    lipgloss.AdaptiveColor{Dark: Colors.slightlyBrighterOrange, Light: Colors.orange},

		DefaultTextColor: lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
		ErrorColor:       lipgloss.AdaptiveColor{Dark: Colors.slightlyBrighterRed, Light: Colors.red},

		InstalledColor:    lipgloss.AdaptiveColor{Dark: Colors.green, Light: Colors.green},
		IncompatibleColor: lipgloss.AdaptiveColor{Dark: Colors.slightlyBrighterOrange, Light: Colors.orange},

		SelectedListItemColor:   lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
		UnselectedListItemColor: lipgloss.AdaptiveColor{Dark: Colors.black, Light: Colors.white},
		ActiveBoxBorderColor:    lipgloss.AdaptiveColor{Dark: Colors.green, Light: Colors.green},
		InactiveBoxBorderColor:  lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
	},
}

// GetTheme returns a theme based on the given name.
func GetTheme(theme string) Theme {
	return themeMap["default"]
}
