package registry

import (
	"log"

	"github.com/hashicorp/go-version"
	"github.com/jedwards1230/go-kerbal/cmd/config"
)

// CKAN Spec: https://github.com/KSP-CKAN/CKAN/blob/master/Spec.md

type Ckan struct {
	Identifier     string
	Name           string
	Author         string
	Abstract       string
	Description    string
	License        string
	SearchTags     map[string]interface{}
	ModConflicts   []string
	ModDepends     []string
	IsCompatible   bool
	Versions       versions
	Install        install
	Resources      resource
	SearchSpace    string
	SearchableName string
}

// Initialize struct values
func CreateCkan(raw map[string]interface{}) (Ckan, error) {
	var mod Ckan
	/* for k, v := range raw {
		log.Printf("%v: %v", k, v)
	}
	log.Panic() */

	if err := mod.cleanNames(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanIdentifiers(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanAuthors(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanVersions(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanAbstract(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanDescription(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanLicense(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanInstall(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanDownload(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanDependencies(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanConflicts(raw); err != nil {
		return mod, err
	}

	if err := mod.cleanSearchSpace(raw); err != nil {
		return mod, err
	}

	return mod, nil
}

// Compares installed KSP version to min/max compatible for the mod.
//
// Returns true if compatible
func (c Ckan) checkCompatible() bool {
	cfg := config.GetConfig()
	configVer := cfg.Settings.KerbalVer
	kerbalVer, err := version.NewVersion(configVer)
	if err != nil {
		log.Printf("Error with kerbal version: %v", err)
	}

	if c.Versions.KspMin != "" {
		min, err := version.NewVersion(c.Versions.KspMin)
		if err != nil {
			log.Printf("Error with kerbal min version: %v", err)
		}
		if min.GreaterThan(kerbalVer) {
			return false
		}
	}
	if c.Versions.KspMax != "" {
		max, err := version.NewVersion(c.Versions.KspMax)
		if err != nil {
			log.Printf("Error with kerbal max version: %v", err)
		}
		if max.LessThan(kerbalVer) {
			return false
		}
	}
	return true
}
