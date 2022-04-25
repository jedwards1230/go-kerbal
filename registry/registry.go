package registry

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/segmentio/encoding/json"
	"golang.org/x/sync/errgroup"

	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/tidwall/buntdb"
)

type Registry struct {
	TotalModMap            map[string][]Ckan
	CompatibleModMap       map[string][]Ckan
	SortedCompatibleMap    map[string]Ckan
	SortedNonCompatibleMap map[string]Ckan
	ModMapIndex            ModIndex
	InstalledModList       map[string]bool
	DB                     *CkanDB
	SortOptions            SortOptions
	InstallQueue           []Ckan
}

// Initializes database and registry
func GetRegistry() Registry {
	db := GetDB()
	installedList, err := dirfs.CheckInstalledMods()
	if err != nil {
		log.Printf("Error checking installed mods: %v", err)
	}

	sortOpts := SortOptions{
		SortTag:   "name",
		SortOrder: "ascend",
	}

	return Registry{
		DB:               db,
		InstalledModList: installedList,
		SortOptions:      sortOpts,
	}
}

func (r *Registry) SortModList() error {
	log.Printf("Sorting %d mods. Order: %s by %s", len(r.TotalModMap), r.SortOptions.SortOrder, r.SortOptions.SortTag)
	cfg := config.GetConfig()

	// Get map with most compatible mod
	compatMap, err := getLatestVersionMap(getCompatibleModMap(r.TotalModMap))
	if err != nil {
		return err
	}
	incompatMap, err := getLatestVersionMap(r.TotalModMap)
	if err != nil {
		return err
	}

	if cfg.Settings.HideIncompatibleMods {
		r.buildModIndex(compatMap)
	} else {
		r.buildModIndex(incompatMap)
	}

	r.SortedCompatibleMap = compatMap
	r.SortedNonCompatibleMap = incompatMap

	log.Printf("Sort result: %d/%d", len(r.ModMapIndex), len(r.TotalModMap))
	return nil
}

// Get list of Ckan objects from database
func (r *Registry) GetEntireModList() map[string][]Ckan {
	log.Println("Gathering mod list from database")

	var mod Ckan
	newMap := make(map[string][]Ckan)
	total := 0
	r.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(_, value string) bool {
			err := json.Unmarshal([]byte(value), &mod)
			if err != nil {
				log.Printf("Error loading into Ckan struct: %v", err)
			}

			// check if mod is installed
			if len(r.InstalledModList) > 0 {
				if r.InstalledModList[mod.Install.Find] || r.InstalledModList[mod.Install.File] {
					mod.Install.Installed = true
				} else if mod.Install.FindRegex != "" {
					re := regexp.MustCompile(mod.Install.FindRegex)
					for k, v := range r.InstalledModList {
						if re.MatchString(k) {
							mod.Install.Installed = v
							break
						}
					}
				} else {
					mod.Install.Installed = false
				}
			}

			// add to list
			newMap[mod.Identifier] = append(newMap[mod.Identifier], mod)
			total += 1
			return true
		})
		return nil
	})

	log.Printf("Loaded %v mod files from database", total)
	return newMap
}

func (r Registry) GetActiveModList() map[string]Ckan {
	cfg := config.GetConfig()
	var modMap map[string]Ckan
	if cfg.Settings.HideIncompatibleMods {
		modMap = r.SortedCompatibleMap
	} else {
		modMap = r.SortedNonCompatibleMap
	}
	return modMap
}

func (r *Registry) BuildSearchIndex(s string) (ModIndex, error) {
	modMap := r.GetActiveModList()
	s = strings.ToLower(s)
	re := regexp.MustCompile("(?i)" + s)

	searchMapIndex := make(ModIndex, 0)
	for id, mod := range modMap {
		if re.MatchString(mod.SearchSpace) {
			searchMapIndex = append(searchMapIndex, Entry{id, mod.Name})
		}
	}

	switch r.SortOptions.SortOrder {
	case "ascend":
		sort.Sort(searchMapIndex)
	case "descend":
		sort.Sort(sort.Reverse(searchMapIndex))
	}

	log.Printf("Found %d mods for \"%s\"", len(searchMapIndex), s)
	return searchMapIndex, nil
}

// Download selected mods
func (r *Registry) DownloadMods(toDownload map[string]Ckan) error {
	var mods []Ckan
	var err error

	// collect all mods and dependencies
	log.Print("Checking dependencies")
	if len(toDownload) > 0 {
		mods, err = r.checkDependencies(toDownload)
		if err != nil {
			return err
		}
	} else {
		return errors.New("no mods provided")
	}

	// check for any conflicts
	log.Print("Checking conflicts")
	err = r.checkConflicts(mods)
	if err != nil {
		return err
	}

	if len(mods) > 0 {
		log.Printf("Downloading %d mods", len(mods))

		// Create tmp dir
		err := os.MkdirAll("./tmp", os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating tmp dir: %v", err)
		}

		// download mods
		g := new(errgroup.Group)
		for i := range mods {
			mod := mods[i]
			g.Go(func() error {
				err := downloadMod(mod)
				if err != nil {
					return fmt.Errorf("%s: %v", mod.Name, err)
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
		r.InstallQueue = mods
		return nil
	}
	return errors.New("no URLS provided")
}

// Download a mod
func downloadMod(mod Ckan) error {
	resp, err := http.Get(mod.Download.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server code: %v", resp.StatusCode)
	}

	// Create zip file
	out, err := os.Create(mod.Download.Path)
	if err != nil {
		return fmt.Errorf("creating file: %v", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy contents to file: %v", err)
	}

	log.Printf("Downloaded: %v", mod.Name)

	return nil
}

// Install mods in the registry install queue
func (r *Registry) InstallMods() error {
	if len(r.InstallQueue) > 0 {

		log.Printf("Installing %d mods", len(r.InstallQueue))
		g := new(errgroup.Group)
		for i := range r.InstallQueue {
			mod := r.InstallQueue[i]
			g.Go(func() error {
				err := installMod(mod)
				if err != nil {
					return fmt.Errorf("%s: %v", mod.Name, err)
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
		log.Printf("Installed %v mods", len(r.InstallQueue))
		return nil
	}
	return errors.New("install queue empty")
}

// Install a mod
func installMod(mod Ckan) error {
	// open zip
	zipReader, err := zip.OpenReader(mod.Download.Path)
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

	re := regexp.MustCompile("(?i)" + mod.Install.InstallTo)
	// unzip all into GameData folder
	for _, f := range zipReader.File {
		if f.Name != "" {
			// verify mod being installed to folder location in metadata
			if re.MatchString(f.Name) {
				err := dirfs.UnzipFile(f, gameDataDir)
				if err != nil {
					return fmt.Errorf("unzipping file to filesystem: %v", err)
				}
			} else if re.MatchString("GameData") {
				// todo: sometimes doubles down on the GameData dir
				err := dirfs.UnzipFile(f, gameDataDir+"/GameData/")
				if err != nil {
					return fmt.Errorf("unzipping file to filesystem: %v", err)
				}
			} else {
				log.Printf("installing in separate dir: \"%v\"", f.Name)
				err := dirfs.UnzipFile(f, gameDataDir+"/GameData/")
				if err != nil {
					return fmt.Errorf("unzipping file to filesystem: %v", err)
				}
			}
		}
	}
	log.Printf("Installed: %v", mod.Name)
	return nil
}

// Gather list of mods and dependencies for download
func (r *Registry) checkDependencies(toDownload map[string]Ckan) ([]Ckan, error) {
	var mods []Ckan
	dependencies := make(map[string]bool)
	count := 0
	for _, mod := range toDownload {
		if !mod.IsCompatible {
			log.Printf("Warning: %v is not compatible with your current configuration", mod.Name)
		}
		if len(mod.ModDepends) > 0 {
			for i := range mod.ModDepends {
				dependent := r.SortedNonCompatibleMap[mod.ModDepends[i]]
				if dependent.Identifier != "" && !dependencies[dependent.Identifier] {
					if !dependent.IsCompatible {
						log.Printf("Warning: %v depends on %s which is not compatible with your current configuration", mod.Name, dependent.Name)
					}
					mods = append(mods, dependent)
					dependencies[dependent.Identifier] = true
					count++
				} else {
					return mods, fmt.Errorf("could not find dependency: %v for %v", mod.ModDepends[i], mod.Name)
				}
			}
		}
		if !dependencies[mod.Identifier] {
			mods = append(mods, mod)
			dependencies[mod.Identifier] = true
		}
	}
	if count > 0 {
		log.Printf("Found %d dependencies", count)
	}
	return mods, nil
}

// check for conflicts
//
// todo: could probably be a lot faster
func (r *Registry) checkConflicts(mods []Ckan) error {
	for i := range mods {
		if len(mods[i].ModConflicts) > 0 {
			for j, conflict := range mods[i].ModConflicts {
				// todo: link conflicts to install folder so filesystem can be checked for conflicts
				if r.InstalledModList[conflict] {
					return fmt.Errorf("%v requires %v which conflicts with %v", mods[i].Name, mods[j].Name, conflict)
				}

				for j := range mods {
					if mods[j].Identifier == conflict {
						return fmt.Errorf("%v conflicts with queued mod %v", mods[i].Name, mods[j].Name)
					}
				}
			}
		}
	}
	log.Printf("No conflicts found")
	return nil
}

// Create r.ModMapIndex from given modMap
//
// Sorts by order and tags saved to registry
func (r *Registry) buildModIndex(modMap map[string]Ckan) {
	r.ModMapIndex = make(ModIndex, 0)
	for k, v := range modMap {
		r.ModMapIndex = append(r.ModMapIndex, Entry{k, v.SearchableName})
	}

	switch r.SortOptions.SortOrder {
	case "ascend":
		sort.Sort(r.ModMapIndex)
	case "descend":
		sort.Sort(sort.Reverse(r.ModMapIndex))
	}
}

// Filter out incompatible mods
func getCompatibleModMap(incompatibleModMap map[string][]Ckan) map[string][]Ckan {
	countGood := 0
	countBad := 0
	compatibleModMap := make(map[string][]Ckan, len(incompatibleModMap))
	for id, modList := range incompatibleModMap {
		for i := range modList {
			if modList[i].IsCompatible {
				countGood += 1
				compatibleModMap[id] = append(compatibleModMap[id], modList[i])
			} else {
				countBad += 1
			}
		}
	}

	log.Printf("Total Compatible: %d | Incompatible: %d", countGood, countBad)
	return compatibleModMap
}

// Filters list by unique identifiers to ensure duplicate mods are not displayed
func getLatestVersionMap(modMapBuckets map[string][]Ckan) (map[string]Ckan, error) {
	modMap := make(map[string]Ckan)
	countGood := 0
	countBad := 0
	for id, modList := range modMapBuckets {
		for _, mod := range modList {
			// convert to proper version type for comparison
			foundVersion, err := version.NewVersion(mod.Versions.Mod)
			if err != nil {
				return modMap, fmt.Errorf("error creating version: %v", err)
			}

			// check if mod is stored already
			if modMap[id].Identifier != "" {
				// convert to proper version type for comparison
				storedVersion, err := version.NewVersion(modMap[id].Versions.Mod)
				if err != nil {
					return modMap, fmt.Errorf("error creating version: %v", err)
				}

				// compare versions and store most recent
				if foundVersion.GreaterThan(storedVersion) {
					// replace old mod
					modMap[id] = mod
				}
				countBad += 1
			} else {
				// store mod if slot is empty
				modMap[id] = mod
				countGood += 1
			}
		}
	}

	//log.Printf("Total filtered by identifier: Unique: %d | Extra: %d", countGood, countBad)
	return modMap, nil
}
