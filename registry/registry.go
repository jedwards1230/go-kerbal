package registry

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/registry/database"
	"github.com/tidwall/buntdb"
)

type Registry struct {
	// List of mods in the database
	ModList []database.Ckan
	// Database for handling CKAN files
	DB *database.CkanDB
}

type SortOptions struct {
	SortTag   string
	SortOrder string
}

// Initializes database and registry
func GetRegistry() Registry {
	db := database.GetDB()

	return Registry{
		DB: db,
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
	r.ModList = sortedModList
}

// Get list of Ckan objects from database
func (r *Registry) GetModList() []database.Ckan {
	log.Println("Gathering mod list from database")

	var ckan database.Ckan
	var newList []database.Ckan
	r.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(_, value string) bool {
			err := json.Unmarshal([]byte(value), &ckan.Raw)
			if err != nil {
				log.Printf("Error loading Ckan.raw struct: %v", err)
			}

			// initialize struct values
			err = ckan.Init()
			if err != nil {
				log.Printf("Error initializing ckan: %v", err)
			}

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
	var sortedModList []database.Ckan
	idList := make(map[string]bool)
	countGood := 0
	countBad := 0
	for _, mod := range modList {
		if idList[mod.Identifier] {
			// TODO: Compare versions
			countBad += 1
		} else {
			sortedModList = append(sortedModList, mod)
			idList[mod.Identifier] = true
			countGood += 1
		}
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

/* func (r *Registry) removeMod(i int) {
	r.ModList[i] = r.ModList[len(r.ModList)-1]
	list := r.ModList[:len(r.ModList)-1]
	r.ModList = list
}
*/
