package registry

import (
	"github.com/jedwards1230/go-kerbal/registry/database"
)

type Registry struct {
	// List of mods in the database
	ModList []database.Ckan
	// Database for handling CKAN files
	DB *database.CkanDB
}

// Initializes database and registry
func GetRegistry() Registry {
	db := database.GetDB()

	return Registry{
		DB: db,
	}
}
