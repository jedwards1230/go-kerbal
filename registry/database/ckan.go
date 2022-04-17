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
	Abstract             string
	Download             string
	License              string
	Epoch                string
	Resources            resource
	SearchTags           map[string]interface{}
	ModDepends           map[string]interface{}
	ModConflicts         map[string]interface{}
	Installed            bool
	IsCompatible         bool
	SpecVersion          string
	Version              string
	VersionKspMax        string
	VersionKspMin        string
	SearchableName       string
	SearchableAbstract   string
	SearchableIdentifier string
}

type resource struct {
	Homepage    string `json:"homepage"`
	Spacedock   string `json:"spacedock"`
	Repository  string `json:"repository"`
	XScreenshot string `json:"x_screenshot"`
}

// Initialize struct values in-place
func (c *Ckan) Init(raw map[string]interface{}) error {
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
		//log.Printf("cleanVersions error: %v", err)
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

	return err
}

func (c *Ckan) cleanNames(raw map[string]interface{}) error {
	c.Name = strings.TrimSpace(raw["name"].(string))
	if c.Name == "" {
		return errors.New("invalid file name")
	}
	c.SearchableName = dirfs.Strip(c.Name)

	return nil
}

func (c *Ckan) cleanIdentifiers(raw map[string]interface{}) error {
	c.Identifier = strings.TrimSpace(raw["identifier"].(string))
	if c.Identifier == "" {
		return errors.New("invalid file identifier")
	}
	c.SearchableIdentifier = dirfs.Strip(c.Name)

	return nil
}

// Clean author name data.
//
//
func (c *Ckan) cleanAuthors(raw map[string]interface{}) error {
	switch author := raw["author"].(type) {
	case []interface{}:
		var list []string
		for _, v := range author {
			val := strings.TrimSpace(v.(string))
			list = append(list, val)
		}
		c.Author = strings.Join(list, ", ")
	case string:
		c.Author = strings.TrimSpace(author)
		if c.Author == "" {
			return errors.New("invalid author name")
		}
	default:
		return errors.New("type mismatch")
	}

	return nil
}

func (c *Ckan) cleanVersions(raw map[string]interface{}) error {
	if raw["version"] != nil && strings.TrimSpace(raw["version"].(string)) != "" {
		v := raw["version"]
		modVersion, epoch, err := c.cleanModVersion(v.(string))
		if err != nil {
			return fmt.Errorf("Error: %v, v: %v", err, raw["version"].(string))
		}
		c.Version = modVersion
		c.Epoch = epoch
		c.VersionKspMin = ""
		c.VersionKspMax = ""

		if raw["ksp_version_max"] != nil && strings.TrimSpace(raw["ksp_version_max"].(string)) != "" {
			v = raw["ksp_version_max"]
			vMax, _, _ := c.cleanModVersion(v.(string))
			c.VersionKspMax = vMax
		}

		if raw["ksp_version_min"] != nil && strings.TrimSpace(raw["ksp_version_min"].(string)) != "" {
			v = raw["ksp_version_min"]
			vMin, _, _ := c.cleanModVersion(v.(string))
			c.VersionKspMin = vMin
		}

		if raw["ksp_version"] != nil && strings.TrimSpace(raw["ksp_version"].(string)) != "" {
			v = raw["ksp_version"]
			cfg := config.GetConfig()

			if v == "any" {
				c.VersionKspMax = cfg.Settings.KerbalVer
				c.VersionKspMin = "0.0"
				return nil
			}

			vKsp, _, _ := c.cleanModVersion(v.(string))

			if c.VersionKspMax == "" {
				c.VersionKspMax = vKsp
			}
			if c.VersionKspMin == "" {
				c.VersionKspMin = vKsp
			}
		}

		if c.VersionKspMax == "" && c.VersionKspMin != "" {
			cfg := config.GetConfig()

			c.VersionKspMax = cfg.Settings.KerbalVer
		} else if c.VersionKspMin == "" && c.VersionKspMax != "" {
			c.VersionKspMin = "0.0"
		}

		if c.VersionKspMin == "" || c.VersionKspMax == "" {
			return fmt.Errorf("Error: ksp: %v, min: %v, max: %v", raw["ksp_version"], raw["ksp_version_min"], raw["ksp_version_max"])
		}

		return nil
	} else {
		return errors.New("Error: no version available")
	}
}

func (c *Ckan) cleanAbstract(raw map[string]interface{}) error {
	c.Abstract = strings.TrimSpace(raw["abstract"].(string))
	if c.Abstract == "" {
		c.SearchableAbstract = ""
		return errors.New("invalid abstract")
	}
	// TODO: strip common words (a, the, and, etc.)
	c.SearchableAbstract = dirfs.Strip(c.Abstract)

	return nil
}

func (c *Ckan) cleanLicense(raw map[string]interface{}) error {
	switch license := raw["license"].(type) {
	case []interface{}:
		for _, v := range license {
			c.License = v.(string)
			break
		}
	case string:
		c.License = strings.TrimSpace(raw["license"].(string))
		if c.License == "" {
			return errors.New("invalid license")
		}
	default:
		return errors.New("type mismatch")
	}

	return nil
}

func (c *Ckan) cleanDownload(raw map[string]interface{}) error {
	switch download := raw["download"].(type) {
	case string:
		c.Download = strings.TrimSpace(string(download))
	default:
		c.Download = ""
	}
	if c.Download == "" {
		return errors.New("invalid download path")
	}
	return nil
}

// Clean version string
//
// Returns Version, Epoch, and any errors
func (c *Ckan) cleanModVersion(rawVersion string) (string, string, error) {
	var epoch string
	var cleanVersion string

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

		newVersion, err = version.NewVersion(rawVersion)
		if err != nil {
			// TODO: Better version cleaning
			return "", "", err
		}
		cleanVersion = newVersion.String()
		return cleanVersion, epoch, nil
	}

	cleanVersion = newVersion.String()
	return cleanVersion, epoch, nil
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

	if c.VersionKspMin != "" {
		min, err := version.NewVersion(c.VersionKspMin)
		if err != nil {
			log.Printf("Error with kerbal min version: %v", err)
		}
		if min.GreaterThan(kerbalVer) {
			return false
		}
	}
	if c.VersionKspMax != "" {
		max, err := version.NewVersion(c.VersionKspMax)
		if err != nil {
			log.Printf("Error with kerbal max version: %v", err)
		}
		if max.LessThan(kerbalVer) {
			return false
		}
	}
	return true
}
