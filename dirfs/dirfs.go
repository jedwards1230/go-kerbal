package dirfs

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// useeful for keeping unit tests on track with created directories
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

// GetHomeDirectory returns the users home directory.
func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

// CreateDirectory creates a new directory given a name.
func CreateDirectory(name string) error {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(name, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// Find root directory of KSP
func FindKspPath() string {
	path := ""
	home, _ := os.UserHomeDir()
	home += "/Library/Application Support/Steam/steamapps"
	filepath.WalkDir(home, func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() && strings.Contains(dir.Name(), "Kerbal Space Program") {
			path = s
		}
		return nil
	})

	return path
}
