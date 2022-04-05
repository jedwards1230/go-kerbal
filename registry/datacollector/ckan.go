package datacollector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
)

// CKAN Spec: https://github.com/KSP-CKAN/CKAN/blob/master/Spec.md

type Ckan struct {
	SpecVersion   string                 `json:"spec_version" binding:"required"`
	Identifier    string                 `json:"identifier" binding:"required"`
	Name          string                 `json:"name" binding:"required"`
	Abstract      string                 `json:"abstract" binding:"required"`
	Author        string                 `json:"author" binding:"required"`
	Download      string                 `json:"download" binding:"required"`
	License       string                 `json:"license" binding:"required"`
	Epoch         string                 `json:"epoch"`
	VersionKspMax string                 `json:"ksp_version_max"`
	VersionKspMin string                 `json:"ksp_version_min"`
	Resources     resource               `json:"resources"`
	Tags          map[string]interface{} `json:"tags"`
	Depends       map[string]interface{} `json:"depends"`
	Conflicts     map[string]interface{} `json:"conflicts"`
	raw           map[string]interface{}
	Version       *version.Version
	Authors       []string
}

type resource struct {
	Homepage    string `json:"homepage"`
	Spacedock   string `json:"spacedock"`
	Repository  string `json:"repository"`
	XScreenshot string `json:"x_screenshot"`
}

func (c *Ckan) init() error {
	err := c.cleanVersion()

	return err
}

func (c *Ckan) cleanVersion() error {
	rawVersion := c.raw["version"].(string)
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
