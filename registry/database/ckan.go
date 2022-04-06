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
	SpecVersion    string                 `json:"spec_version" binding:"required"`
	Identifier     string                 `json:"identifier" binding:"required"`
	Name           string                 `json:"name" binding:"required"`
	Abstract       string                 `json:"abstract" binding:"required"`
	Author         string                 `json:"author" binding:"required"`
	Download       string                 `json:"download" binding:"required"`
	License        string                 `json:"license" binding:"required"`
	Epoch          string                 `json:"epoch"`
	Resources      resource               `json:"resources"`
	Tags           map[string]interface{} `json:"tags"`
	Depends        map[string]interface{} `json:"depends"`
	Conflicts      map[string]interface{} `json:"conflicts"`
	raw            map[string]interface{}
	Version        *version.Version
	VersionKspMax  *version.Version
	VersionKspMin  *version.Version
	Authors        []string
	SearchableName string
}

type resource struct {
	Homepage    string `json:"homepage"`
	Spacedock   string `json:"spacedock"`
	Repository  string `json:"repository"`
	XScreenshot string `json:"x_screenshot"`
}

func (c *Ckan) init() error {
	err := c.cleanVersions()
	if err != nil {
		return err
	}

	c.Name = strings.TrimSpace(c.Name)
	c.SearchableName = strip(c.Name)

	return err
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
			return err
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
