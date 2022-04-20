package database

import (
	"log"

	"github.com/hashicorp/go-version"
	"github.com/jedwards1230/go-kerbal/cmd/config"
)

// CKAN Spec: https://github.com/KSP-CKAN/CKAN/blob/master/Spec.md

type Ckan struct {
	Identifier           string
	Name                 string
	Author               string
	Abstract             string
	License              string
	SearchTags           map[string]interface{}
	ModConflicts         map[string]interface{}
	ModDepends           []string
	IsCompatible         bool
	Versions             versions
	Install              install
	Resources            resource
	SearchableName       string
	SearchableAbstract   string
	SearchableIdentifier string
}

type versions struct {
	Epoch  string
	Mod    string
	KspMin string
	KspMax string
	Spec   string
}

type install struct {
	Installed bool
	Download  string
	FindRegex string
	Find      string
	File      string
	InstallTo string
}

type resource struct {
	Homepage    string
	Spacedock   string
	Repository  string
	XScreenshot string
}

// Initialize struct values in-place
func CreateCkan(raw map[string]interface{}) (Ckan, error) {
	var mod Ckan
	/* for k, v := range raw {
		log.Printf("%v: %v", k, v)
	}
	log.Panic() */
	err := mod.cleanNames(raw)
	if err != nil {
		return mod, err
	}

	err = mod.cleanIdentifiers(raw)
	if err != nil {
		return mod, err
	}

	err = mod.cleanAuthors(raw)
	if err != nil {
		return mod, err
	}

	err = mod.cleanVersions(raw)
	if err != nil {
		return mod, err
	}

	mod.IsCompatible = mod.CheckCompatible()

	_ = mod.cleanAbstract(raw)

	err = mod.cleanLicense(raw)
	if err != nil {
		return mod, err
	}

	err = mod.cleanInstall(raw)
	if err != nil {
		return mod, err
	}

	err = mod.cleanDownload(raw)
	if err != nil {
		return mod, err
	}

	err = mod.cleanDependencies(raw)
	if err != nil {
		return mod, err
	}

	return mod, err
}

// Compares installed KSP version to min/max compatible for the mod.
//
// Returns true if compatible
func (c Ckan) CheckCompatible() bool {
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
