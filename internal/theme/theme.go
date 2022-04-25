package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents the properties that make up a theme.
type Theme struct {
	Green  lipgloss.Color
	Orange lipgloss.Color
	Red    lipgloss.Color
	Blue   lipgloss.Color

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
	white  string
	black  string
	green  string
	orange string
	red    string
	blue   string
}

// Colors contains the different kinds of Colors and their values.
var Colors = appColors{
	white:  "#FFFDF5",
	black:  "#000000",
	green:  "#00aa00",
	orange: "#cf8611",
	red:    "#cc241d",
	blue:   "#0040ff",
}

// themeMap represents the mapping of different themes.
var themeMap = map[string]Theme{
	"default": {
		Green:  lipgloss.Color(Colors.green),
		Orange: lipgloss.Color(Colors.orange),
		Red:    lipgloss.Color(Colors.red),
		Blue:   lipgloss.Color(Colors.blue),

		DefaultTextColor: lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
		ErrorColor:       lipgloss.AdaptiveColor{Dark: Colors.red, Light: Colors.red},

		InstalledColor:    lipgloss.AdaptiveColor{Dark: Colors.green, Light: Colors.green},
		IncompatibleColor: lipgloss.AdaptiveColor{Dark: Colors.orange, Light: Colors.orange},

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
