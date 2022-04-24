package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents the properties that make up a theme.
type Theme struct {
	SelectedListItemColor   lipgloss.AdaptiveColor
	UnselectedListItemColor lipgloss.AdaptiveColor
	InstalledListItemColor  lipgloss.AdaptiveColor
	ActiveBoxBorderColor    lipgloss.AdaptiveColor
	InactiveBoxBorderColor  lipgloss.AdaptiveColor
	Green                   lipgloss.Color
	Orange                  lipgloss.Color
	Red                     lipgloss.Color
	Blue                    lipgloss.Color

	SelectedTreeItemColor                lipgloss.AdaptiveColor
	UnselectedTreeItemColor              lipgloss.AdaptiveColor
	SpinnerColor                         lipgloss.AdaptiveColor
	StatusBarSelectedFileForegroundColor lipgloss.AdaptiveColor
	StatusBarSelectedFileBackgroundColor lipgloss.AdaptiveColor
	StatusBarBarForegroundColor          lipgloss.AdaptiveColor
	StatusBarBarBackgroundColor          lipgloss.AdaptiveColor
	StatusBarTotalFilesForegroundColor   lipgloss.AdaptiveColor
	StatusBarTotalFilesBackgroundColor   lipgloss.AdaptiveColor
	StatusBarLogoForegroundColor         lipgloss.AdaptiveColor
	StatusBarLogoBackgroundColor         lipgloss.AdaptiveColor
	ErrorColor                           lipgloss.AdaptiveColor
	DefaultTextColor                     lipgloss.AdaptiveColor
}

// appColors contains the different types of colors.
type appColors struct {
	white  string
	black  string
	green  string
	orange string
	red    string
	blue   string

	darkGray           string
	defaultPink        string
	defaultLightPurple string
	defaultDarkPurple  string
	gruvGreen          string
	gruvBlue           string
	gruvYellow         string
	gruvOrange         string
	nordRed            string
	nordGreen          string
	nordBlue           string
	nordYellow         string
	nordWhite          string
	nordBlack          string
	nordGray           string
	nordOrange         string
	spookyPurple       string
	spookyOrange       string
	spookyYellow       string
	holidayRed         string
	holidayGreen       string
	holidayGold        string
}

// Colors contains the different kinds of Colors and their values.
var Colors = appColors{
	white:  "#FFFDF5",
	black:  "#000000",
	green:  "#00aa00",
	orange: "#cf8611",
	red:    "#cc241d",
	blue:   "#0040ff",

	darkGray:           "#3c3836",
	defaultPink:        "#F25D94",
	defaultLightPurple: "#A550DF",
	defaultDarkPurple:  "#6124DF",
	gruvGreen:          "#b8bb26",
	gruvBlue:           "#458588",
	gruvYellow:         "#d79921",
	gruvOrange:         "#d65d0e",
	nordRed:            "#bf616a",
	nordGreen:          "#a3be8c",
	nordBlue:           "#81a1c1",
	nordYellow:         "#ebcb8b",
	nordWhite:          "#e5e9f0",
	nordBlack:          "#3b4252",
	nordGray:           "#4c566a",
	nordOrange:         "#d08770",
	spookyPurple:       "#881EE4",
	spookyOrange:       "#F75F1C",
	spookyYellow:       "#FF9A00",
	holidayRed:         "#B70D00",
	holidayGreen:       "#005C01",
	holidayGold:        "#CC9901",
}

// themeMap represents the mapping of different themes.
var themeMap = map[string]Theme{
	"default": {
		SelectedListItemColor:   lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
		UnselectedListItemColor: lipgloss.AdaptiveColor{Dark: Colors.black, Light: Colors.white},
		InstalledListItemColor:  lipgloss.AdaptiveColor{Dark: Colors.green, Light: Colors.green},
		ActiveBoxBorderColor:    lipgloss.AdaptiveColor{Dark: Colors.green, Light: Colors.green},
		InactiveBoxBorderColor:  lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
		Green:                   lipgloss.Color(Colors.green),
		Orange:                  lipgloss.Color(Colors.orange),
		Red:                     lipgloss.Color(Colors.red),
		Blue:                    lipgloss.Color(Colors.blue),

		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: Colors.defaultPink, Light: Colors.defaultPink},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: Colors.defaultPink, Light: Colors.defaultPink},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.white},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: Colors.defaultPink, Light: Colors.defaultPink},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.white},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: Colors.darkGray, Light: Colors.darkGray},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.white},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: Colors.defaultLightPurple, Light: Colors.defaultLightPurple},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.white},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: Colors.defaultDarkPurple, Light: Colors.defaultDarkPurple},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: Colors.red, Light: Colors.red},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: Colors.white, Light: Colors.black},
	},
}

// GetTheme returns a theme based on the given name.
func GetTheme(theme string) Theme {
	return themeMap["default"]
}
