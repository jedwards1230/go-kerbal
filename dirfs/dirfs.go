package dirfs

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/hashicorp/go-version"
	"github.com/jedwards1230/go-kerbal/cmd/config"
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
			return "", nil
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

func DownloadMod(url string) error {
	cfg := config.GetConfig()
	// get response from url
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error reaching server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response status from server")
	}

	// convert zip to bytevalue
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading download: %v", err)
	}

	// create zip reader from download
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return fmt.Errorf("error converting download to zip: %v", err)
	}

	// get Kerbal folder
	destination, err := filepath.Abs(cfg.Settings.KerbalDir)
	if err != nil {
		return fmt.Errorf("error getting KSP dir: %v", err)
	}

	// unzip all into GameData folder
	for _, f := range zipReader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return fmt.Errorf("error unzipping file to filesystem: %v", err)
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// check file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

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
