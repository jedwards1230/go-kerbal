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

func (c *Ckan) cleanInstall(raw map[string]interface{}) error {
	var installInfo install
	rawI := raw["install"].([]interface{})
	if len(rawI) > 0 {
		rawInstall := rawI[0].(map[string]interface{})
		pathFound := false

		switch {
		case rawInstall["find"] != nil:
			installInfo.Find = rawInstall["find"].(string)
			pathFound = true
		case rawInstall["file"] != nil:
			installInfo.File = rawInstall["file"].(string)
			pathFound = true
		case rawInstall["find_regexp"] != nil:
			installInfo.FindRegex = rawInstall["find_regexp"].(string)
			pathFound = true
		}

		if rawInstall["install_to"] != nil {
			installInfo.InstallTo = rawInstall["install_to"].(string)
		}

		if pathFound && installInfo.InstallTo != "" {
			c.InstallInfo = installInfo
			return nil
		}
	}
	return fmt.Errorf("install interface empty: %v", raw["install"])
}

func (c *Ckan) cleanDependencies(raw map[string]interface{}) error {
	if raw["depends"] != nil {
		dependsInfo := make([]string, 0)
		rawI := raw["depends"].([]interface{})
		if len(rawI) > 0 {
			rawDepends := rawI[0].(map[string]interface{})
			if rawDepends["name"] != nil {
				dependsInfo = append(dependsInfo, rawDepends["name"].(string))

				if len(dependsInfo) > 0 {
					c.ModDepends = dependsInfo
					return nil
				}
				return fmt.Errorf("error proccessing install dependencies: %v", raw["depends"])
			}
		}
	}
	return nil
}

// Clean author name data.
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
	var vMod, vMin, vMax *version.Version
	var epoch string
	var err error
	if raw["version"] != nil && strings.TrimSpace(raw["version"].(string)) != "" {
		v := raw["version"]
		vMod, epoch, err = c.cleanModVersion(v.(string))
		if err != nil {
			return fmt.Errorf("error: %v, v: %v", err, raw["version"].(string))
		}
		c.Epoch = epoch

		if raw["ksp_version_max"] != nil && strings.TrimSpace(raw["ksp_version_max"].(string)) != "" {
			v = raw["ksp_version_max"]
			vMax, _, _ = c.cleanModVersion(v.(string))
		}

		if raw["ksp_version_min"] != nil && strings.TrimSpace(raw["ksp_version_min"].(string)) != "" {
			v = raw["ksp_version_min"]
			vMin, _, _ = c.cleanModVersion(v.(string))
		}

		if raw["ksp_version"] != nil && strings.TrimSpace(raw["ksp_version"].(string)) != "" {
			v = raw["ksp_version"]
			cfg := config.GetConfig()

			if v == "any" {
				vMax, _ = version.NewVersion(cfg.Settings.KerbalVer)
				vMin, _ = version.NewVersion("0.0")
			}

			newVKsp, _, err := c.cleanModVersion(v.(string))
			if err != nil {
				log.Printf("Error cleaning mod version: %v", err)
			}

			if vMax == nil {
				vMax = newVKsp
			}
			if vMin == nil {
				vMin = newVKsp
			}
		}

		if vMax == nil && vMin != nil {
			cfg := config.GetConfig()
			vMax, _ = version.NewVersion(cfg.Settings.KerbalVer)
		} else if vMin == nil && vMax != nil {
			vMin, _ = version.NewVersion("0.0")
		}

		if vMin == nil || vMax == nil {
			return fmt.Errorf("error: ksp: %v, min: %v, max: %v", raw["ksp_version"], raw["ksp_version_min"], raw["ksp_version_max"])
		}

		c.Versions.VersionMod = vMod.String()
		c.Versions.VersionKspMin = vMin.String()
		c.Versions.VersionKspMax = vMax.String()

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
		return fmt.Errorf("invalid download path: %v", raw["download"])
	}
	return nil
}

// Clean version string
//
// Returns Version, Epoch, and any errors
func (c *Ckan) cleanModVersion(rawVersion string) (*version.Version, string, error) {
	var epoch string

	// check if epoch is stored in version string
	if strings.Contains(rawVersion, ":") {
		s := strings.Split(rawVersion, ":")
		epoch = s[0]
		rawVersion = s[1]
	}

	// attempt to parse version
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
			return newVersion, "", err
		}
		return newVersion, epoch, nil
	}

	return newVersion, epoch, nil
}
