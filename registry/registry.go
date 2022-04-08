package registry

import (
	"log"
	"sort"

	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/registry/database"
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
	var sortedModList []database.Ckan

	// Check compatible
	sortedModList = getCompatibleModList(r.ModList)

	// TODO: Filter by tag

	// Sort by order
	sortedModList = getSortedModList(sortedModList, opts.SortTag, opts.SortOrder)

	log.Printf("Sort result: %d/%d", len(sortedModList), len(r.ModList))
	r.ModList = sortedModList
}

// Filter out incompatible mods if config is set
func getCompatibleModList(modList []database.Ckan) []database.Ckan {
	cfg := config.GetConfig()
	countGood := 0
	countBad := 0
	var compatibleModList []database.Ckan
	for _, mod := range modList {
		if cfg.Settings.HideIncompatibleMods {
			if mod.CheckCompatible() {
				countGood += 1
				compatibleModList = append(compatibleModList, mod)
			} else {
				countBad += 1
			}
		} else {
			compatibleModList = append(compatibleModList, mod)
		}
	}
	if cfg.Settings.HideIncompatibleMods {
		log.Printf("Total filtered by compatibility: Compatible: %d | Incompatible: %d", countGood, countBad)
	}
	return compatibleModList
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
