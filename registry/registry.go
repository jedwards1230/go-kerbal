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
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/segmentio/encoding/json"

	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/jedwards1230/go-kerbal/registry/database"
	"github.com/tidwall/buntdb"
)

type Registry struct {
	TotalModMap      map[string][]database.Ckan
	SortedMap        map[string]database.Ckan
	ModMapIndex      ModIndex
	InstalledModList map[string]bool
	DB               *database.CkanDB
	SortOptions      SortOptions
}

type SortOptions struct {
	SortTag   string
	SortOrder string
}

type Entry struct {
	Key   string
	Value string
}

type ModIndex []Entry

func (entry ModIndex) Len() int           { return len(entry) }
func (entry ModIndex) Less(i, j int) bool { return entry[i].Value < entry[j].Value }
func (entry ModIndex) Swap(i, j int)      { entry[i], entry[j] = entry[j], entry[i] }

// Initializes database and registry
func GetRegistry() Registry {
	db := database.GetDB()
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

	modMapBuckets := r.TotalModMap

	// Check compatible
	if cfg.Settings.HideIncompatibleMods {
		modMapBuckets = getCompatibleModMap(r.TotalModMap)
	}

	// Get map with most compatible mod
	r.SortedMap = getLatestVersionMap(modMapBuckets)

	// Create r.ModMapIndex
	r.buildModMapIndex(r.SortedMap)

	log.Printf("Sort result: %d/%d", len(r.ModMapIndex), len(r.TotalModMap))
	return nil
}

// Get list of Ckan objects from database
func (r *Registry) GetTotalModMap() map[string][]database.Ckan {
	log.Println("Gathering mod list from database")

	var mod database.Ckan
	newMap := make(map[string][]database.Ckan)
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

// Filter out incompatible mods if config is set
func getCompatibleModMap(incompatibleModMap map[string][]database.Ckan) map[string][]database.Ckan {
	countGood := 0
	countBad := 0
	compatibleModMap := make(map[string][]database.Ckan, len(incompatibleModMap))
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

	log.Printf("Total filtered by compatibility: Compatible: %d | Incompatible: %d", countGood, countBad)
	return compatibleModMap
}

// Filters list by unique identifiers to ensure duplicate mods are not displayed
func getLatestVersionMap(modMapBuckets map[string][]database.Ckan) map[string]database.Ckan {
	modMap := make(map[string]database.Ckan)
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

// Create r.ModMapIndex from given modMap
//
// Sorts by order and tags saved to registry
func (r *Registry) buildModMapIndex(modMap map[string]database.Ckan) {
	r.ModMapIndex = make(ModIndex, len(modMap))
	i := 0
	for k, v := range modMap {
		switch r.SortOptions.SortTag {
		case "name":
			r.ModMapIndex[i] = Entry{k, v.Name}
		}
		i++
	}

	switch r.SortOptions.SortOrder {
	case "ascend":
		sort.Sort(r.ModMapIndex)
	case "descend":
		sort.Sort(sort.Reverse(r.ModMapIndex))
	}
}

func (r *Registry) DownloadMods(toDownload map[string]bool) ([]database.Ckan, error) {
	var mods []database.Ckan
	dependencies := make(map[string]bool)
	// collect all mods and dependencies
	if len(toDownload) > 0 {
		log.Printf("Checking %d mods", len(toDownload))
		for _, id := range r.ModMapIndex {
			mod := r.SortedMap[id.Key]

			if toDownload[mod.Identifier] {
				// todo: find links for dependencies.
				if len(mod.ModDepends) > 0 {
					for i := range mod.ModDepends {
						dependencies[mod.ModDepends[i]] = true
					}
				}

				mods = append(mods, mod)
			}
		}
		if len(dependencies) > 0 {
			for _, id := range r.ModMapIndex {
				if dependencies[r.SortedMap[id.Key].Identifier] {
					mods = append(mods, r.SortedMap[id.Key])
				}
			}
		}
	} else {
		return mods, errors.New("no mods provided")
	}

	// check for conflicts
	//
	// todo: could probably be a lot faster
	for i := range mods {
		if len(mods[i].ModConflicts) > 0 {
			for _, conflict := range mods[i].ModConflicts {
				// todo: link conflicts to install folder so filesystem can be checked for conflicts
				/* if r.InstalledModList[conflict.Install.Find] {
					return mods, fmt.Errorf("%v conflicts with %v", mods[i].Name, mods[j].Name)
				} */

				for j := range mods {
					if mods[j].Identifier == conflict {
						return mods, fmt.Errorf("%v conflicts with %v", mods[i].Name, mods[j].Name)
					}
				}
			}
		}
	}
	log.Printf("No conflicts found")

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

func downloadMod(mod database.Ckan) error {
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

func InstallMods(mods []database.Ckan) error {
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
func installMod(mod database.Ckan) error {
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
