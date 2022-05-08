package dirfs

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"howett.net/plist"

	"github.com/go-git/go-billy/v5"
	"github.com/hashicorp/go-version"
	"github.com/jedwards1230/go-kerbal/internal/common"
	"github.com/jedwards1230/go-kerbal/internal/config"
)

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
func FindKspPath(home string) (string, error) {
	if home == "" {
		switch runtime.GOOS {
		case "darwin":
			log.Printf("MacOS detected")
			home, _ = os.UserHomeDir()
			home += "/Library/Application Support/Steam/steamapps"

		case "windows":
			log.Printf("Windows OS detected")
			home = "C:\\Program Files (x86)\\steam\\SteamApps\\common"

		case "linux":
			// TODO: add paths for linux
			log.Printf("Linux OS detected")
			return "", nil
		}
	}

	var path string
	log.Printf("Searching directory: %s", home)
	//re := regexp.MustCompile("Kerbal Space Program")
	err := filepath.WalkDir(home, func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() && strings.Contains(dir.Name(), "Kerbal Space Program") {
			path = s
			return nil
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if path == "" {
		return "", errors.New("unable to find KSP path")
	}
	return path, nil
}

// Find version of installed KSP.
func FindKspVersion(fp string) *version.Version {
	var result *version.Version
	var rawVersion string

	switch runtime.GOOS {
	case "darwin":
		kspPlist, err := os.Open(fp + "/KSP.app/Contents/info.plist")
		if err != nil {
			common.LogErrorf("Error reading: %v", err)
			return result
		}
		defer kspPlist.Close()

		var data map[string]interface{}
		decoder := plist.NewDecoder(kspPlist)
		err = decoder.Decode(&data)
		if err != nil {
			log.Println(err)
		}
		rawVersion = data["CFBundleShortVersionString"].(string)

	case "windows", "linux":
		// parse readme.txt
		file, err := os.Open(fp + "/readme.txt")
		if err != nil {
			return result
		}
		defer file.Close()

		// version is one of the first lines in the file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "Version") {
				rawVersion = scanner.Text()
				break
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// regex the version
		re := regexp.MustCompile(`\d+(\.\d+)+`)
		rawVersion = fmt.Sprint(re.FindAllString(rawVersion, -1))
		if strings.Contains(rawVersion, "[") {
			rawVersion = strings.ReplaceAll(rawVersion, "[", "")
		}
		if strings.Contains(rawVersion, "]") {
			rawVersion = strings.ReplaceAll(rawVersion, "]", "")
		}
	}

	// create proper version
	result, err := version.NewVersion(rawVersion)
	if err != nil {
		log.Printf("Error writing KSP Version: %v", err)
		return nil
	}

	return result
}

// Unzip file to specified directory
func UnzipFile(f *zip.File, filePath string) error {
	// create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// create a destination file for unzipped data
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// unzip data and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func CheckInstalledMods() (map[string]bool, error) {
	cfg := config.GetConfig()
	installedMods := make(map[string]bool, 0)

	// get Kerbal folder
	destination, err := filepath.Abs(cfg.Settings.KerbalDir + "/GameData")
	if err != nil {
		return installedMods, fmt.Errorf("error getting KSP dir: %v", err)
	}

	log.Print("Checking installed mods")
	files, err := ioutil.ReadDir(destination)
	if err != nil {
		return installedMods, err
	}

	for _, f := range files {
		modName := f.Name()
		if modName != "Squad" && modName != "SquadExpansion" && f.IsDir() {
			installedMods[modName] = true
		} else if filepath.Ext(modName) == ".dll" {
			installedMods[modName] = true
		}
	}

	return installedMods, nil
}
