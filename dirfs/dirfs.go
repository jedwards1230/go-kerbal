package dirfs

import (
	"errors"
	"io/fs"
	"io/ioutil"
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

// Collect list of file paths
func FindFilePaths(root, ext string) []string {
	var pathList []string
	filepath.WalkDir(root, func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(dir.Name()) == ext {
			pathList = append(pathList, s)
		}
		return nil
	})
	return pathList
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

// Parse .ckan file into JSON string
func ParseCKAN(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// parse ckan data
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(byteValue), nil
}
