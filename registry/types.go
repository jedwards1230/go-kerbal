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

type install struct {
	Installed bool
	Download  string
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
	Key   string
	Value string
}

type (
	ModIndex []Entry
)

func (entry ModIndex) Len() int           { return len(entry) }
func (entry ModIndex) Less(i, j int) bool { return entry[i].Value < entry[j].Value }
func (entry ModIndex) Swap(i, j int)      { entry[i], entry[j] = entry[j], entry[i] }
