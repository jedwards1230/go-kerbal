package registry

type SortOptions struct {
	SortTag   string
	SortOrder string
}

type versions struct {
	Epoch  string
	Mod    string
	KspMin string
	KspMax string
	Spec   string
}

type download struct {
	URL  string
	Path string
}

type install struct {
	Installed bool
	FindRegex string
	Find      string
	File      string
	InstallTo string
}

type resource struct {
	Homepage    string
	Spacedock   string
	Repository  string
	XScreenshot string
}

type Entry struct {
	Key      string
	SearchBy string
}

type ModIndex []Entry

func (entry ModIndex) Len() int           { return len(entry) }
func (entry ModIndex) Less(i, j int) bool { return entry[i].SearchBy < entry[j].SearchBy }
func (entry ModIndex) Swap(i, j int)      { entry[i], entry[j] = entry[j], entry[i] }
