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
	IsCompatible         bool
	SpecVersion          string
	Version              string
	VersionKspMax        string
	VersionKspMin        string
	searchableAuthor     []string
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
		log.Printf("cleanVersions error: %v", err)
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

// TODO: organize into one author field
func (c *Ckan) cleanAuthors(raw map[string]interface{}) error {
	switch author := raw["author"].(type) {
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

func (c *Ckan) cleanVersions(raw map[string]interface{}) error {
	v := raw["version"]
	if v != nil {
		modVersion, epoch, err := c.cleanModVersion(v.(string))
		if err != nil {
			// TODO: Only minor errors come through but could be fixed with better filtering
			return nil
		}
		c.Version = modVersion
		c.Epoch = epoch

		v = raw["ksp_version_max"]
		if v != nil {
			vMax, _, err := c.cleanModVersion(v.(string))
			if err != nil {
				return err
			}
			c.VersionKspMax = vMax
		}

		v = raw["ksp_version_min"]
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
			log.Printf("Error with kerbal min version: %v", err)
		}
		if max.LessThan(kerbalVer) {
			return false
		}
	}
	return true
}
