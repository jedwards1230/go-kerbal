package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents the properties that make up a theme.
type Theme struct {
	// Light: black | Dark: white
	SelectedListItemColor lipgloss.AdaptiveColor
	// Light: white | Dark: black
	UnselectedListItemColor lipgloss.AdaptiveColor
	// Light: green | Dark: green
	InstalledListItemColor lipgloss.AdaptiveColor
	ActiveBoxBorderColor   lipgloss.AdaptiveColor
	InactiveBoxBorderColor lipgloss.AdaptiveColor

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
	white string
	black string
	green string

	darkGray           string
	red                string
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

// Colors contains the different kinds of colors and their values.
var colors = appColors{
	white: "#FFFDF5",
	black: "#000000",
	green: "#00ff00",

	darkGray:           "#3c3836",
	red:                "#cc241d",
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
		SelectedListItemColor:   lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		UnselectedListItemColor: lipgloss.AdaptiveColor{Dark: colors.black, Light: colors.white},
		InstalledListItemColor:  lipgloss.AdaptiveColor{Dark: colors.green, Light: colors.green},
		ActiveBoxBorderColor:    lipgloss.AdaptiveColor{Dark: colors.holidayGreen, Light: colors.holidayGreen},
		InactiveBoxBorderColor:  lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},

		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.darkGray, Light: colors.darkGray},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.defaultLightPurple, Light: colors.defaultLightPurple},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.defaultDarkPurple, Light: colors.defaultDarkPurple},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
	},
}

// GetTheme returns a theme based on the given name.
func GetTheme(theme string) Theme {
	return themeMap["default"]
}
