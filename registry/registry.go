package registry

import (
	"github.com/jedwards1230/go-kerbal/registry/datacollector"
)

func BuildRegistry() []datacollector.Ckan {
	return datacollector.GetAvailableMods()
}
