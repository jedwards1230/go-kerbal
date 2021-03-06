package internal

const (
	LogPath = "./logs/debug.log"
	DBPath  = "./data.db"
)

const (
	MinWidth  = 115
	MinHeight = 20
)

const (
	CommandView     = 0
	ModListView     = 1
	ModInfoView     = 2
	LogView         = 3
	EnterKspDirView = 4
	SearchView      = 5
	SettingsView    = 6
	QueueView       = 7
)

const (
	// StatusBarHeight represents the height of the status bar.
	StatusBarHeight = 1

	// BoxPadding represents the padding of the boxes.
	BoxPadding = 1

	// EllipsisStyle represents the characters displayed when overflowing.
	EllipsisStyle = "..."

	// FileSizeLoadingStyle represents the characters displayed when file sizes are loading.
	FileSizeLoadingStyle = "---"
)

const (
	MenuInputs     = 4
	MenuSortOrder  = 0
	MenuSortTag    = 1
	MenuCompatible = 2
	MenuKspDir     = 3
)
