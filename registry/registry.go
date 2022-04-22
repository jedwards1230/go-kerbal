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
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/segmentio/encoding/json"

	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/tidwall/buntdb"
)

type Registry struct {
	// Contains every mod in database
	TotalModMap map[string][]Ckan
	// Contains every compatible mod in database
	CompatibleModMap map[string][]Ckan
	// Contains every unique mod, sorted
	SortedCompatibleMap map[string]Ckan
	// Contains every unique compatible mod, sorted
	SortedNonCompatibleMap map[string]Ckan
	// Index for traversing mod map.
	ModMapIndex      ModIndex
	InstalledModList map[string]bool
	DB               *CkanDB
	SortOptions      SortOptions
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

func (r *Registry) SortModMap() error {
	log.Printf("Sorting %d mods. Order: %s by %s", len(r.TotalModMap), r.SortOptions.SortOrder, r.SortOptions.SortTag)
	cfg := config.GetConfig()

	// Get map with most compatible mod
	r.SortedCompatibleMap = getLatestVersionMap(getCompatibleModMap(r.TotalModMap))
	log.Printf("Total Incompatible: %d", len(r.TotalModMap))
	r.SortedNonCompatibleMap = getLatestVersionMap(r.TotalModMap)

	if cfg.Settings.HideIncompatibleMods {
		r.buildModMapIndex(r.SortedCompatibleMap)
	} else {
		r.buildModMapIndex(r.SortedNonCompatibleMap)
	}

	log.Printf("Sort result: %d/%d", len(r.ModMapIndex), len(r.TotalModMap))
	return nil
}

// Get list of Ckan objects from database
func (r *Registry) GetTotalModMap() map[string][]Ckan {
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

			// todo: better check for installed mods. this is not accurate at all.
			// check if mod is installed
			if r.InstalledModList[mod.Install.Find] {
				mod.Install.Installed = true
			} else {
				mod.Install.Installed = false
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

func (r *Registry) GetActiveModMap() map[string]Ckan {
	cfg := config.GetConfig()
	var modMap map[string]Ckan
	if cfg.Settings.HideIncompatibleMods {
		modMap = r.SortedCompatibleMap
	} else {
		modMap = r.SortedNonCompatibleMap
	}
	return modMap
}

func (r *Registry) BuildSearchMapIndex(s string) (ModIndex, error) {
	modMap := r.GetActiveModMap()
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

func (r *Registry) DownloadMods(toDownload map[string]Ckan) ([]Ckan, error) {
	var mods []Ckan
	var err error

	// collect all mods and dependencies
	if len(toDownload) > 0 {
		mods, err = r.checkDependencies(toDownload)
		if err != nil {
			return mods, err
		}
	} else {
		return mods, errors.New("no mods provided")
	}

	err = r.checkConflicts(mods)
	if err != nil {
		return mods, err
	}

	// download mods
	if len(mods) > 0 {
		log.Printf("Downloading %d mods (after checking dependencies)", len(mods))

		// Create tmp dir
		err := os.MkdirAll("./tmp", os.ModePerm)
		if err != nil {
			return mods, fmt.Errorf("error creating tmp dir: %v", err)
		}

		var wg sync.WaitGroup
		n := len(mods)
		wg.Add(n)
		for i := range mods {
			go func(i int) {
				err := downloadMod(mods[i])
				if err != nil {
					log.Printf("Error downloading %s: %v", mods[i].Name, err)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		log.Printf("Downloaded %d mods", n)
	} else {
		return mods, errors.New("no URLS provided")
	}
	return mods, nil
}

// Gather list of mods and dependencies for download
func (r *Registry) checkDependencies(toDownload map[string]Ckan) ([]Ckan, error) {
	var mods []Ckan

	log.Printf("Checking %d mods", len(toDownload))
	for _, mod := range toDownload {
		if len(mod.ModDepends) > 0 {
			for i := range mod.ModDepends {
				dependent := r.SortedNonCompatibleMap[mod.ModDepends[i]]
				if dependent.Identifier != "" {
					if dependent.IsCompatible {
						mods = append(mods, dependent)
					} else {
						log.Printf("Warning: Dependent mod %s is not compatible with your current mods", dependent.Name)
						mods = append(mods, dependent)
					}
				} else {
					return mods, fmt.Errorf("could not find dependency: %v %v", mod.ModDepends[i], dependent.Identifier)
				}
			}
		}
		mods = append(mods, mod)
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
					return fmt.Errorf("%v conflicts with installed mod %v", mods[i].Name, mods[j].Name)
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
func (r *Registry) buildModMapIndex(modMap map[string]Ckan) {
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

func downloadMod(mod Ckan) error {
	resp, err := http.Get(mod.Install.Download)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response from server: %v", resp.StatusCode)
	}

	// Create zip file
	out, err := os.Create("./tmp/" + mod.Identifier + ".zip")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
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

func InstallMods(mods []Ckan) error {
	var wg sync.WaitGroup
	wg.Add(len(mods))
	for i := range mods {
		go func(i int) {
			err := installMod(mods[i])
			if err != nil {
				log.Printf("Error installing mod: %v", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	return nil
}

// todo: more reliable installing to directory. really gotta validate the paths and compare whats in the zip.
func installMod(mod Ckan) error {
	cfg := config.GetConfig()
	// open zip
	zipReader, err := zip.OpenReader("./tmp/" + mod.Identifier + ".zip")
	if err != nil {
		return fmt.Errorf("error opening zip file: %v", err)
	}
	defer zipReader.Close()

	// get Kerbal folder
	destination, err := filepath.Abs(cfg.Settings.KerbalDir)
	if err != nil {
		return fmt.Errorf("error getting KSP dir: %v", err)
	}

	// unzip all into GameData folder
	for _, f := range zipReader.File {
		// verify mod being installed to folder location in metadata
		if strings.Contains(f.Name, mod.Install.InstallTo) {
			err := dirfs.UnzipFile(f, destination)
			if err != nil {
				return fmt.Errorf("error unzipping file to filesystem: %v", err)
			}
		} else {
			log.Printf("error unzipping: %v", f.Name)
		}
	}

	log.Printf("Installed: %v", mod.Name)
	return nil
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
func getLatestVersionMap(modMapBuckets map[string][]Ckan) map[string]Ckan {
	modMap := make(map[string]Ckan)
	countGood := 0
	countBad := 0
	for id, modList := range modMapBuckets {
		for _, mod := range modList {
			// convert to proper version type for comparison
			foundVersion, err := version.NewVersion(mod.Versions.Mod)
			if err != nil {
				log.Printf("Error creating version: %v", err)
			}

			// check if mod is stored already
			if modMap[id].Identifier != "" {
				// convert to proper version type for comparison
				storedVersion, err := version.NewVersion(modMap[id].Versions.Mod)
				if err != nil {
					log.Printf("Error creating version: %v", err)
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

	log.Printf("Total filtered by identifier: Unique: %d | Extra: %d", countGood, countBad)
	return modMap
}
