package dirfs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/hashicorp/go-version"
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

// Walks dir in billy.Filesystem
func FindFilePaths(repo billy.Filesystem, ext string) []string {
	var pathList []string
	WalkDir(repo, "", func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(dir.Name()) == ext {
			pathList = append(pathList, "/"+s)
		}
		return nil
	})
	return pathList
}

// Find root directory of KSP
//
// TODO: add paths for linux
func FindKspPath(home string) (string, error) {
	if home == "" {
		if runtime.GOOS == "darwin" {
			log.Printf("MacOS detected")
			home, _ = os.UserHomeDir()
			home += "/Library/Application Support/Steam/steamapps"

		} else if runtime.GOOS == "windows" {
			log.Printf("Windows OS detected")
			home = "C:\\Program Files (x86)\\steam\\SteamApps\\common"

		} else if runtime.GOOS == "linux" {
			log.Printf("Linux OS detected")
			return "/FIXME", nil
		}
	}

	path := ""
	log.Printf("Searching directory: %s", home)
	err := filepath.WalkDir(home, func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() && strings.Contains(dir.Name(), "Kerbal Space Program") {
			path = s
			return io.EOF
		}
		return nil
	})
	if err == io.EOF {
		err = nil
	} else if err != nil {
		return path, err
	}

	if path == "" {
		return path, errors.New("unable to find KSP version")
	}
	return path, nil
}

// Find version of installed KSP.
// Refers to readme.txt in game directory
//
// TODO: Find more reliable source for version number
func FindKspVersion(filePath string) *version.Version {
	var result *version.Version

	// parse readme.txt
	file, err := os.Open(filePath + "/readme.txt")
	if err != nil {
		return result
	}
	defer file.Close()

	// version is one of the first lines in the file
	v := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Version") {
			v = scanner.Text()
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// regex the version
	re := regexp.MustCompile(`\d+(\.\d+)+`)
	v = fmt.Sprint(re.FindAllString(v, -1))
	if strings.Contains(v, "[") {
		v = strings.ReplaceAll(v, "[", "")
	}
	if strings.Contains(v, "]") {
		v = strings.ReplaceAll(v, "]", "")
	}

	// create proper version
	result, err = version.NewVersion(v)
	if err != nil {
		log.Printf("Error writing KSP Version: %v", err)
		return nil
	}

	return result
}

// Strip string of non-alphanumeric characters
func Strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}
