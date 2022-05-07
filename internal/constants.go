package internal

const LogPath = "./logs/debug.log"
const DBPath = "./data.db"

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
