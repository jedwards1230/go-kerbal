package database

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
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
	c.SearchableName = strip(c.Name)

	return nil
}

func (c *Ckan) cleanIdentifiers() error {
	c.Identifier = strings.TrimSpace(c.raw["identifier"].(string))
	if c.Identifier == "" {
		return errors.New("invalid file identifier")
	}
	c.SearchableIdentifier = strip(c.Name)

	return nil
}

func (c *Ckan) cleanVersions() error {
	err := c.cleanModVersion()
	if err != nil {
		return err
	}

	kspMax := c.raw["ksp_version_max"]
	if kspMax != nil {
		c.VersionKspMax, err = cleanKspVersion(kspMax.(string))
	}

	kspMin := c.raw["ksp_version_min"]
	if kspMin != nil {
		c.VersionKspMax, err = cleanKspVersion(kspMin.(string))
	}

	return err
}

func (c *Ckan) cleanModVersion() error {
	var rawVersion string

	v := c.raw["version"]
	if v != nil {
		rawVersion = v.(string)
	} else {
		return errors.New("no version available")

	}

	re := regexp.MustCompile(`\d+(\.\d+)+`)

	if strings.Contains(rawVersion, ":") {
		s := strings.Split(rawVersion, ":")
		c.Epoch = s[0]
		rawVersion = s[1]
	}

	newVersion, err := version.NewVersion(rawVersion)
	if err != nil {
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
			return nil
		}
		c.Version = fixedVersion
	}

	c.Version = newVersion
	return nil
}

func cleanKspVersion(rawVersion string) (*version.Version, error) {
	var v *version.Version

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
			return v, err
		}
		return fixedVersion, nil
	}
	return newVersion, nil
}

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}
