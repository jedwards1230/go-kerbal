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
	Valid          bool
	SearchTags     map[string]interface{}
	ModConflicts   []string
	ModDepends     []string
	IsCompatible   bool
	Versions       versions
	Install        install
	Download       download
	Resources      resource
	SearchSpace    string
	SearchableName string
	Errors         map[string]interface{}
}

// Initialize struct values
func CreateCkan(raw map[string]interface{}) Ckan {
	var mod Ckan
	var validMod = true
	/* for k, v := range raw {
		log.Printf("%v: %v", k, v)
	}
	log.Panic() */

	mod.Errors = make(map[string]interface{})
	valid, err := mod.checkValid(mod.cleanNames(raw))
	if !valid {
		mod.Errors["cleanNames"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanIdentifiers(raw))
	if !valid {
		mod.Errors["cleanIdentifiers"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanAuthors(raw))
	if !valid {
		mod.Errors["cleanAuthors"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanVersions(raw))
	if !valid {
		mod.Errors["cleanVersions"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanAbstract(raw))
	if !valid {
		mod.Errors["cleanAbstract"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanDescription(raw))
	if !valid {
		mod.Errors["cleanDescription"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanLicense(raw))
	if !valid {
		mod.Errors["cleanLicense"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanInstall(raw))
	if !valid {
		switch err.Error() {
		case "no install path":
			mod.Errors["ignored"] = true
		default:
			mod.Errors["cleanInstall"] = err
		}
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanDownload(raw))
	if !valid {
		mod.Errors["cleanDownload"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanDependencies(raw))
	if !valid {
		mod.Errors["cleanDependencies"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanConflicts(raw))
	if !valid {
		mod.Errors["cleanConflicts"] = err
		validMod = false
	}

	valid, err = mod.checkValid(mod.cleanSearchSpace(raw))
	if !valid {
		mod.Errors["cleanSearchSpace"] = err
		validMod = false
	}

	mod.Valid = validMod

	if !validMod || err != nil {
		mod.Errors["raw"] = raw
		return mod
	}

	return mod
}

func (c *Ckan) checkValid(err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return true, nil
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

func (c *Ckan) markDownloaded() {
	c.Download.Downloaded = true
}

func (c *Ckan) markInstalled() {
	c.Install.Installed = true
}
