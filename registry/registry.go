package registry

import (
	"github.com/jedwards1230/go-kerbal/registry/datacollector"
	"github.com/jedwards1230/go-kerbal/registry/module"
)

func BuildRegistry() []module.ModuleVersion {
	return datacollector.GetAvailableMods()
}
