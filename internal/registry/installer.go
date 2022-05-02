package registry

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jedwards1230/go-kerbal/internal/ckan"
	"github.com/jedwards1230/go-kerbal/internal/common"
	"github.com/jedwards1230/go-kerbal/internal/config"
	"github.com/jedwards1230/go-kerbal/internal/dirfs"
	"golang.org/x/sync/errgroup"
)

func (r *Registry) AddToQueue(mod ckan.Ckan) error {
	if mod.Installed() {
		r.Queue.AddRemoval(mod)
	} else {
		r.Queue.AddSelection(mod)

		mods, err := r.CheckDependencies(mod)
		if err != nil {
			return err
		}
		if len(mods) > 0 {
			for _, mod := range mods {
				if r.Queue.CheckRemovals(mod.Identifier) {
					common.LogError("dependent mod being removed!!!!")
				}
				r.Queue.AddDependency(mod)
			}
		}
	}
	return nil
}

func (r *Registry) RemoveMods() error {
	for _, mod := range r.Queue.GetRemovals() {
		err := r.removeMod(mod)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) removeMod(mod ckan.Ckan) error {
	common.LogErrorf("Removing %v", mod.Name)

	// find path
	cfg := config.GetConfig()

	// get Kerbal folder
	destination, err := filepath.Abs(cfg.Settings.KerbalDir + "/GameData")
	if err != nil {
		return fmt.Errorf("cannot get KSP dir: %v", err)
	}

	files, err := ioutil.ReadDir(destination)
	if err != nil {
		return err
	}

	var removePath string
	for _, f := range files {
		modName := f.Name()
		if mod.Install.Find != "" {
			if strings.Contains(modName, mod.Install.Find) {
				removePath = modName
			}
		} else if mod.Install.File != "" {
			if strings.Contains(modName, mod.Install.File) {
				removePath = modName
			}
		} else if mod.Install.FindRegex != "" {
			re := regexp.MustCompile(mod.Install.FindRegex)
			if re.MatchString(mod.Install.FindRegex) {
				removePath = modName
			}
		} else {
			common.LogErrorf("Cannot find for %v", mod.Name)
		}
	}
	removePath = destination + "/" + removePath
	//r.LogErrorf("Deleting \"%v\"", removePath)
	err = os.RemoveAll(removePath)
	if err != nil {
		return fmt.Errorf("cannot remove mod %s: %v", mod.Name, err)
	}
	return nil
}

// Download selected mods
func (r *Registry) DownloadMods() error {
	var mods []ckan.Ckan
	var err error

	for _, mod := range r.Queue.GetSelections() {
		mods = append(mods, mod)
	}

	for _, mod := range r.Queue.GetDependencies() {
		mods = append(mods, mod)
	}

	// check for any conflicts
	log.Print("Checking conflicts")
	err = r.checkConflicts(mods)
	if err != nil {
		return err
	}

	if len(mods) > 0 {
		common.LogCommandf("Downloading %d mods", len(mods))

		// download mods
		g := new(errgroup.Group)
		for i := range mods {
			mod := mods[i]
			g.Go(func() error {
				err := r.downloadMod(mod)
				if err != nil {
					return fmt.Errorf("%s: %v", mod.Name, err)
				}
				mod.MarkDownloaded()
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no URLS provided")
}

// Download a mod
func (r *Registry) downloadMod(mod ckan.Ckan) error {
	resp, err := http.Get(mod.Download.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server code: %v", resp.StatusCode)
	}

	// Create zip file
	out, err := os.Create(r.GetTempDir() + mod.Download.Path)
	if err != nil {
		return fmt.Errorf("creating file: %v", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy contents to file: %v", err)
	}

	// Add mod to install queue if successfully downloaded
	log.Printf("Downloaded: %v", mod.Name)

	return nil
}

// Install mods in the registry install queue
func (r *Registry) InstallMods() error {
	if r.Queue.InstallLen() > 0 {

		// install dependencies
		for _, mod := range r.Queue.GetDependencies() {
			if !mod.Installed() {
				err := r.installMod(&mod)
				if err != nil {
					return fmt.Errorf("%s: %v", mod.Name, err)
				}
				mod.SetInstalled(true)
			}
		}

		// install the rest
		for _, mod := range r.Queue.GetSelections() {
			if !mod.Installed() {
				err := r.installMod(&mod)
				if err != nil {
					return fmt.Errorf("%s: %v", mod.Name, err)
				}
				mod.SetInstalled(true)
			}
		}

		common.LogSuccessf("Installed %v mods", r.Queue.InstallLen())
		return nil
	}
	return errors.New("install queue empty")
}

// Install a mod
func (r *Registry) installMod(mod *ckan.Ckan) error {
	// open zip
	zipReader, err := zip.OpenReader(r.GetTempDir() + mod.Download.Path)
	if err != nil {
		return fmt.Errorf("opening zip file: %v", mod.Download.Path)
	}
	defer zipReader.Close()

	// get Kerbal folder
	cfg := config.GetConfig()
	gameDataDir, err := filepath.Abs(cfg.Settings.KerbalDir)
	if err != nil {
		return fmt.Errorf("getting KSP dir: %v", err)
	}

	installTo := regexp.MustCompile("(?i)" + mod.Install.InstallTo)
	// unzip all into GameData folder
	for _, f := range zipReader.File {
		destination, err := r.getInstallDir(f.Name, gameDataDir, installTo)
		if err != nil {
			return err
		}

		err = dirfs.UnzipFile(f, destination)
		if err != nil {
			return fmt.Errorf("unzipping file to filesystem: %v", err)
		}

	}
	log.Printf("Installed: %v", mod.Name)
	return nil
}

func (r *Registry) getInstallDir(file, gameDataDir string, installTo *regexp.Regexp) (string, error) {
	if file != "" {
		// verify file path matches GameData directory
		if installTo.MatchString(file) {
			// trim filepath to ensure it aligns with GameData directory
			if !strings.HasPrefix(file, "GameData") {
				dest := file
				for _, char := range file {
					dest = strings.TrimPrefix(dest, string(char))
					if strings.HasPrefix(dest, "GameData") {
						file = dest
						break
					}
				}
				if !strings.HasPrefix(dest, "GameData") {
					// double check the trim worked
					return "", fmt.Errorf("bad trim: f.Name: \"%v\" -> \"%v\", prefix: %v", file, gameDataDir, strings.HasPrefix(file, "GameData"))
				}
			}
		} else {
			// dump extras into GameData
			// this is done for depenendencies like .dlls and merging mods
			// might be able to do this in a safer way
			gameDataDir = gameDataDir + "/GameData/"
		}

		// merge file paths
		filePath := filepath.Join(gameDataDir, file)
		if !strings.HasPrefix(filePath, filepath.Clean(gameDataDir)+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path: %s", filePath)
		}

		// warn if overwriting vanilla game data
		if strings.Contains(filePath, "GameData/Squad/") || strings.Contains(filePath, "GameData/SquadExpansion/") {
			common.LogWarningf("Warning: attempting to overwrite KSP data: %s", filePath)
		}

		return filePath, nil
	}
	return "", errors.New("empty file string")
}

// Gather list of mods and dependencies for download
func (r *Registry) CheckDependencies(mod ckan.Ckan) (map[string]ckan.Ckan, error) {
	mods := make(map[string]ckan.Ckan)
	count := 0
	if mod.Identifier == "" {
		return mods, errors.New("empty mod provided")
	}

	if !mod.IsCompatible {
		common.LogWarningf("Warning: %v is not compatible with your current configuration", mod.Name)
	}

	if len(mod.ModDepends) > 0 {
		for i := range mod.ModDepends {
			dependent := r.UnsortedModMap[mod.ModDepends[i]]
			if dependent.Identifier != "" {
				if mods[dependent.Identifier].Identifier == "" {
					if !dependent.IsCompatible {
						common.LogWarningf("Warning: %v depends on %s (incompatible with current configuration)", mod.Name, dependent.Name)
					}
					mods[dependent.Identifier] = dependent
					count++
				}
			} else {
				return mods, fmt.Errorf("could not find dependency: %v for %v", mod.ModDepends[i], mod.Name)
			}
		}
	}
	if count > 0 {
		log.Printf("Found %d dependencies", count)
	}

	return mods, nil
}

// check for conflicts
func (r *Registry) checkConflicts(mods []ckan.Ckan) error {
	// find conflicts for each queued mod
	for i := range mods {
		if len(mods[i].ModConflicts) > 0 {
			for j, conflict := range mods[i].ModConflicts {
				// check conflicts with installed mods
				if r.InstalledModList[conflict].Identifier != "" {
					return fmt.Errorf("%v requires %v which conflicts with installed %v", mods[i].Name, mods[j].Name, conflict)
				}

				// check conflicts with queued mods
				for j := range mods {
					if mods[j].Identifier == conflict {
						return fmt.Errorf("%v requires %v which conflicts with queued %v", mods[i].Name, mods[j].Name, conflict)
					}
				}
			}
		}
	}
	log.Printf("No conflicts found")
	return nil
}
