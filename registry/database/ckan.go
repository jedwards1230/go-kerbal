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
	Download             string
	License              string
	Epoch                string
	Resources            resource
	SearchTags           map[string]interface{}
	ModDepends           []string
	ModConflicts         map[string]interface{}
	InstallInfo          install
	Installed            bool
	IsCompatible         bool
	SpecVersion          string
	Versions             Versions
	SearchableName       string
	SearchableAbstract   string
	SearchableIdentifier string
}

type Versions struct {
	VersionMod    string
	VersionKspMin string
	VersionKspMax string
}

type install struct {
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
func (c *Ckan) Init(raw map[string]interface{}) error {
	/* for k, v := range raw {
		log.Printf("%v: %v", k, v)
	}
	log.Panic() */
	err := c.cleanNames(raw)
	if err != nil {
		return err
	}

	err = c.cleanIdentifiers(raw)
	if err != nil {
		return err
	}

	err = c.cleanAuthors(raw)
	if err != nil {
		return err
	}

	err = c.cleanVersions(raw)
	if err != nil {
		return err
	}

	c.IsCompatible = c.CheckCompatible()

	_ = c.cleanAbstract(raw)

	err = c.cleanLicense(raw)
	if err != nil {
		return err
	}

	err = c.cleanDownload(raw)
	if err != nil {
		return err
	}

	if raw["install"] != nil {
		err = c.cleanInstall(raw)
		if err != nil {
			return err
		}
	}

	err = c.cleanDependencies(raw)
	if err != nil {
		return err
	}

	return err
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

	if c.Versions.VersionKspMin != "" {
		min, err := version.NewVersion(c.Versions.VersionKspMin)
		if err != nil {
			log.Printf("Error with kerbal min version: %v", err)
		}
		if min.GreaterThan(kerbalVer) {
			return false
		}
	}
	if c.Versions.VersionKspMax != "" {
		max, err := version.NewVersion(c.Versions.VersionKspMax)
		if err != nil {
			log.Printf("Error with kerbal max version: %v", err)
		}
		if max.LessThan(kerbalVer) {
			return false
		}
	}
	return true
}
