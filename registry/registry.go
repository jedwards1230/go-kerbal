package registry

import (
	"github.com/jedwards1230/go-kerbal/registry/database"
)

type Registry struct {
	ModList []database.Ckan
	DB      *database.CkanDB
}

func GetRegistry() Registry {
	db := database.GetDB()

	return Registry{
		DB: db,
	}
}
