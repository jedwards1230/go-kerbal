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
	ModList          []database.Ckan
	SortedModList    []database.Ckan
	InstalledModList map[string]bool
	DB               *database.CkanDB
}

type SortOptions struct {
	SortTag   string
	SortOrder string
}

// Initializes database and registry
func GetRegistry() Registry {
	db := database.GetDB()
	installedList, err := dirfs.CheckInstalledMods()
	if err != nil {
		log.Printf("Error checking installed mods: %v", err)
	}

	return Registry{
		DB:               db,
		InstalledModList: installedList,
	}
}

func (r *Registry) SortModList(opts SortOptions) {
	log.Printf("Sorting %d mods. Order: %s by %s", len(r.ModList), opts.SortOrder, opts.SortTag)
	cfg := config.GetConfig()
	var sortedModList []database.Ckan

	// Check compatible
	if cfg.Settings.HideIncompatibleMods {
		sortedModList = getCompatibleModList(r.ModList)
	} else {
		sortedModList = r.ModList
	}

	// Get list by unique identifiers
	sortedModList = getUniqueModList(sortedModList)

	// TODO: Filter by tag

	// Sort by order
	sortedModList = getSortedModList(sortedModList, opts.SortTag, opts.SortOrder)

	log.Printf("Sort result: %d/%d", len(sortedModList), len(r.ModList))
	r.SortedModList = sortedModList
}

// Get list of Ckan objects from database
func (r *Registry) GetModList() []database.Ckan {
	log.Println("Gathering mod list from database")

	var ckan database.Ckan
	var newList []database.Ckan
	r.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(_, value string) bool {
			// extract mod
			err := json.Unmarshal([]byte(value), &ckan)
			if err != nil {
				log.Printf("Error loading into Ckan struct: %v", err)
			}

			// TODO: better check for installed mods. this is not accurate at all.
			// check if mod is installed
			if r.InstalledModList[ckan.InstallInfo.Find] {
				ckan.Installed = true
			} else {
				ckan.Installed = false
			}

			// add to list
			newList = append(newList, ckan)
			return true
		})
		return nil
	})

	log.Printf("Loaded %v mods from database", len(newList))
	return newList
}

// Filter out incompatible mods if config is set
func getCompatibleModList(modList []database.Ckan) []database.Ckan {
	countGood := 0
	countBad := 0
	var compatibleModList []database.Ckan
	for i := range modList {
		if modList[i].IsCompatible {
			countGood += 1
			compatibleModList = append(compatibleModList, modList[i])
		} else {
			countBad += 1
		}
	}
	log.Printf("Total filtered by compatibility: Compatible: %d | Incompatible: %d", countGood, countBad)
	return compatibleModList
}

// Filters list by unique identifiers to ensure duplicate mods are not displayed
func getUniqueModList(modList []database.Ckan) []database.Ckan {
	sortedModMap := make(map[string]database.Ckan)
	countGood := 0
	countBad := 0
	for _, mod := range modList {
		// convert to proper version type for comparison
		foundVersion, err := version.NewVersion(mod.Version)
		if err != nil {
			log.Printf("Error creating version: %v", err)
		}

		// check if mod is stored already
		if sortedModMap[mod.Identifier].Identifier != "" {
			// convert to proper version type for comparison
			storedVersion, err := version.NewVersion(sortedModMap[mod.Identifier].Version)
			if err != nil {
				log.Printf("Error creating version: %v", err)
			}

			// compare versions and store most recent
			if foundVersion.GreaterThan(storedVersion) {
				// replace old mod
				sortedModMap[mod.Identifier] = mod
			}
			countBad += 1
		} else {
			// store mod if slot is empty
			sortedModMap[mod.Identifier] = mod
			countGood += 1
		}
	}

	// map to slice
	//
	// TODO: this is only done because i originally had a slice for this. check if keeping it as a map is better
	sortedModList := make([]database.Ckan, 0, countGood)
	for _, v := range sortedModMap {
		sortedModList = append(sortedModList, v)
	}

	log.Printf("Total filtered by identifier: Unique: %d | Extra: %d", countGood, countBad)
	return sortedModList
}

// Sort mods by order
func getSortedModList(modList []database.Ckan, tag, order string) []database.Ckan {
	sortedModList := modList
	switch tag {
	case "name":
		switch order {
		case "ascend":
			sort.Slice(sortedModList, func(i, j int) bool {
				return sortedModList[i].SearchableName < sortedModList[j].SearchableName
			})
		case "descend":
			sort.Slice(sortedModList, func(i, j int) bool {
				return sortedModList[i].SearchableName > sortedModList[j].SearchableName
			})
		}
	}
	return sortedModList
}

func (r *Registry) DownloadMods(toDownload map[string]bool) ([]database.Ckan, error) {
	var mods []database.Ckan
	if len(toDownload) > 0 {
		// mods and dependencies to download
		for i := range r.SortedModList {
			mod := r.SortedModList[i]
			if toDownload[mod.Identifier] {
				log.Printf("Mod download requested: %s", mod.Name)

				// TODO: find links for dependencies.
				if len(mod.ModDepends) > 0 {
					for i := range mod.ModDepends {
						log.Printf("Depends on: %v", mod.ModDepends[i])
					}
				} else {
					log.Print("No dependencies detected")
				}

				mods = append(mods, mod)
			}
		}
	} else {
		return mods, errors.New("no mods provided")
	}

	// download mods
	if len(mods) > 0 {
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
	log.Printf("Downloading mod: %s", mod.Name)
	resp, err := http.Get(mod.Download)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response from server: %v", resp.StatusCode)
	}

	// Create tmp dir
	err = os.MkdirAll("./tmp", os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating tmp dir: %v", err)
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

	log.Printf("Successfully downloaded: %v", mod.Name)

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
		// verify mod being isntalled to folder location in metadata
		if strings.Contains(f.Name, mod.InstallInfo.InstallTo) {
			err := dirfs.UnzipFile(f, destination)
			if err != nil {
				return fmt.Errorf("error unzipping file to filesystem: %v", err)
			}
		} else {
			return fmt.Errorf("error unzipping: %v", f.Name)
		}
	}

	log.Printf("Installed: %v", mod.Name)
	return nil
}

/* func (r *Registry) removeMod(i int) {
	r.ModList[i] = r.ModList[len(r.ModList)-1]
	list := r.ModList[:len(r.ModList)-1]
	r.ModList = list
}
*/
