package database

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/dirfs"
)

// CKAN Spec: https://github.com/KSP-CKAN/CKAN/blob/master/Spec.md

type Ckan struct {
	Identifier           string
	Name                 string
	Author               string
	Authors              []string
	Abstract             string
	Download             string
	License              string
	Epoch                string
	Resources            resource
	SearchTags           map[string]interface{}
	ModDepends           map[string]interface{}
	ModConflicts         map[string]interface{}
	raw                  map[string]interface{}
	SpecVersion          *version.Version
	Version              *version.Version
	VersionKspMax        *version.Version
	VersionKspMin        *version.Version
	searchableAuthor     []string
	SearchableName       string
	SearchableIdentifier string
}

type resource struct {
	Homepage    string `json:"homepage"`
	Spacedock   string `json:"spacedock"`
	Repository  string `json:"repository"`
	XScreenshot string `json:"x_screenshot"`
}

// Initialize struct values in-place
func (c *Ckan) init() error {
	err := c.cleanNames()
	if err != nil {
		return err
	}

	err = c.cleanIdentifiers()
	if err != nil {
		return err
	}

	err = c.cleanAuthors()
	if err != nil {
		return err
	}

	err = c.cleanVersions()
	if err != nil {
		return err
	}

	return err
}

func (c *Ckan) cleanNames() error {
	c.Name = strings.TrimSpace(c.raw["name"].(string))
	if c.Name == "" {
		return errors.New("invalid file name")
	}
	c.SearchableName = dirfs.Strip(c.Name)

	return nil
}

func (c *Ckan) cleanIdentifiers() error {
	c.Identifier = strings.TrimSpace(c.raw["identifier"].(string))
	if c.Identifier == "" {
		return errors.New("invalid file identifier")
	}
	c.SearchableIdentifier = dirfs.Strip(c.Name)

	return nil
}

// TODO: organize into one author field
func (c *Ckan) cleanAuthors() error {
	switch author := c.raw["author"].(type) {
	case []interface{}:
		for i, v := range author {
			c.Authors = append(c.Authors, v.(string))
			c.searchableAuthor = append(c.searchableAuthor, dirfs.Strip(c.Authors[i]))
		}
	case string:
		c.Author = strings.TrimSpace(author)
		if c.Author == "" {
			return errors.New("invalid author name")
		}
		c.searchableAuthor = append(c.searchableAuthor, dirfs.Strip(c.Author))
	default:
		return errors.New("type mismatch")
	}

	return nil
}

func (c *Ckan) cleanVersions() error {
	v := c.raw["version"]
	if v != nil {
		modVersion, epoch, err := c.cleanModVersion(v.(string))
		if err != nil {
			// TODO: Only minor errors come through but could be fixed with better filtering
			return nil
		}
		c.Version = modVersion
		c.Epoch = epoch

		v = c.raw["ksp_version_max"]
		if v != nil {
			vMax, _, err := c.cleanModVersion(v.(string))
			if err != nil {
				return err
			}
			c.VersionKspMax = vMax
		}

		v = c.raw["ksp_version_min"]
		if v != nil {
			vMin, _, err := c.cleanModVersion(v.(string))
			if err != nil {
				return err
			}
			c.VersionKspMin = vMin
		}
		return nil
	} else {
		return errors.New("no version available")
	}
}

// Clean version string
//
// Returns Version, Epoch, and any errors
func (c *Ckan) cleanModVersion(rawVersion string) (*version.Version, string, error) {
	var v *version.Version
	var epoch string

	if strings.Contains(rawVersion, ":") {
		s := strings.Split(rawVersion, ":")
		epoch = s[0]
		rawVersion = s[1]
	}

	newVersion, err := version.NewVersion(rawVersion)
	if err != nil {
		re := regexp.MustCompile(`\d+(\.\d+)+`)
		rawVersion = fmt.Sprint(re.FindAllString(rawVersion, -1))

		if strings.Contains(rawVersion, "[") {
			rawVersion = strings.ReplaceAll(rawVersion, "[", "")
		}
		if strings.Contains(rawVersion, "]") {
			rawVersion = strings.ReplaceAll(rawVersion, "]", "")
		}

		fixedVersion, err := version.NewVersion(rawVersion)
		if err != nil {
			// TODO: Better version cleaning
			/* log.Printf("BAD VERSION: %v", err)
			log.Printf("raw: %v", c.raw["version"].(string))
			log.Printf("final: %s\n", rawVersion) */
			return v, epoch, err
		}
		return fixedVersion, epoch, nil
	}
	return newVersion, epoch, nil
}

// Compares installed KSP version to min/max compatible for the mod.
//
// Returns true if compatible
func (c *Ckan) CheckCompatible() bool {
	cfg := config.GetConfig()
	configVer := cfg.Settings.KerbalVer
	kerbalVer, err := version.NewVersion(configVer)
	if err != nil {
		log.Printf("Error with kerbal version: %v", err)
	}

	if c.VersionKspMin != nil {
		if !c.VersionKspMin.LessThanOrEqual(kerbalVer) {
			/* 			log.Printf("True!")
			   			log.Printf("%v <= %v <= %v\n\n", c.VersionKspMin, configVer, c.VersionKspMin)
			*/return false
		}

	}
	if c.VersionKspMax != nil {
		if !c.VersionKspMax.GreaterThanOrEqual(kerbalVer) {
			/* 			log.Printf("True!")
			   			log.Printf("%v <= %v <= %v\n\n", c.VersionKspMin, configVer, c.VersionKspMin)
			*/return false
		}
	}

	/* 	log.Printf("True!")
	   	log.Printf("%v <= %v <= %v\n\n", c.VersionKspMin, configVer, c.VersionKspMin)
	*/return true
}
