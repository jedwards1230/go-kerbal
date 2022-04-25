package registry

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/segmentio/encoding/json"

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
	InstalledModList       map[string]Ckan
	DB                     *CkanDB
	SortOptions            SortOptions
	InstallQueue           []Ckan
}

// Initializes database and registry
func GetRegistry() Registry {
	db := GetDB()

	installedList := make(map[string]Ckan, 0)

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

	installedMap, err := dirfs.CheckInstalledMods()
	if err != nil {
		log.Printf("Error checking installed mods: %v", err)
	}

	var mod Ckan
	newMap := make(map[string][]Ckan)
	total := 0
	err = r.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(_, value string) bool {
			err := json.Unmarshal([]byte(value), &mod)
			if err != nil {
				log.Printf("Error loading into Ckan struct: %v", err)
			}

			// check if mod is installed
			r.checkModInstalled(&mod, installedMap)

			// add to list
			newMap[mod.Identifier] = append(newMap[mod.Identifier], mod)
			total += 1
			return true
		})
		return nil
	})
	if err != nil {
		log.Fatalf("Error viewing db: %v", err)
	}

	log.Printf("Loaded %v mod files from database", total)
	if len(r.InstalledModList) > 0 {
		for _, mod := range r.InstalledModList {
			log.Printf("Found installed: %v", mod.Name)
		}
	} else {
		log.Printf("Found %d mods installed", len(r.InstalledModList))
	}

	return newMap
}

func (r *Registry) checkModInstalled(mod *Ckan, installedMap map[string]bool) {
	if len(installedMap) > 0 {
		if installedMap[mod.Install.Find] || installedMap[mod.Install.File] {
			mod.Install.Installed = true
			r.InstalledModList[mod.Identifier] = *mod
		} else if mod.Install.FindRegex != "" {
			re := regexp.MustCompile(mod.Install.FindRegex)
			for k, v := range installedMap {
				if v && re.MatchString(k) {
					mod.Install.Installed = true
					r.InstalledModList[mod.Identifier] = *mod
					break
				}
			}
		} else {
			mod.Install.Installed = false
		}
	}
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
			searchMapIndex = append(searchMapIndex, Entry{id, mod.SearchableName})
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

func (r *Registry) ApplyMods(modMap map[string]Ckan) error {
	if len(modMap) <= 0 {
		return errors.New("no mods to apply")
	}
	var toDownload, toRemove []Ckan
	for i := range modMap {
		if modMap[i].Install.Installed {
			toRemove = append(toRemove, modMap[i])
		} else {
			toDownload = append(toDownload, modMap[i])
		}
	}

	// Remove Mods
	if len(toRemove) > 0 {
		err := r.RemoveMods(toRemove)
		if err != nil {
			return fmt.Errorf("error removing: %v", err)
		}
	}

	// Install Mods
	if len(toDownload) > 0 {
		err := r.DownloadMods(toDownload)
		if err != nil {
			return fmt.Errorf("error downloading: %v", err)
		}

		if len(r.InstallQueue) > 0 {
			err = r.InstallMods()
			if err != nil {
				return fmt.Errorf("error installing: %v", err)
			}
		}
	}
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
