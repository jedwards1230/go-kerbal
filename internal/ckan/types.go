package ckan

type versions struct {
	Epoch  string
	Mod    string
	KspMin string
	KspMax string
	Spec   string
}

type download struct {
	Downloaded bool
	URL        string
	Path       string
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
