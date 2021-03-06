package ckan

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/jedwards1230/go-kerbal/internal/config"
)

func (c *Ckan) cleanSearchSpace(raw map[string]interface{}) error {
	space := []string{
		c.Name,
		c.SearchableName,
		clean(c.Identifier),
		clean(c.Author),
		clean(c.Abstract),
	}
	c.SearchSpace = strings.Join(space, " ")
	return nil
}

func (c *Ckan) cleanNames(raw map[string]interface{}) error {
	c.Name = strings.TrimSpace(raw["name"].(string))
	if c.Name == "" {
		return errors.New("invalid file name")
	}
	c.SearchableName = clean(c.Name)

	return nil
}

func (c *Ckan) cleanIdentifiers(raw map[string]interface{}) error {
	c.Identifier = strings.TrimSpace(raw["identifier"].(string))
	if c.Identifier == "" {
		return errors.New("invalid file identifier")
	}
	return nil
}

func (c *Ckan) cleanInstall(raw map[string]interface{}) error {
	if raw["install"] != nil {
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
				find := rawInstall["find_regexp"].(string)
				if find == "." {
					find = ""
				}
				installInfo.FindRegex = find
				pathFound = true
			}

			if rawInstall["install_to"] != nil {
				installInfo.InstallTo = rawInstall["install_to"].(string)
			}

			if pathFound && installInfo.InstallTo != "" {
				c.Install = installInfo
				return nil
			}
		}
		return errors.New("empty install path")
	}
	return errors.New("no install path")
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

func (c *Ckan) cleanConflicts(raw map[string]interface{}) error {
	if raw["conflicts"] != nil {
		//log.Printf("Conflicts: %v | %T", raw["conflicts"], raw["conflicts"])
		conflictsInfo := make([]string, 0)
		rawI := raw["conflicts"].([]interface{})
		if len(rawI) > 0 {
			rawConflicts := rawI[0].(map[string]interface{})
			if rawConflicts["name"] != nil {
				conflictsInfo = append(conflictsInfo, rawConflicts["name"].(string))

				if len(conflictsInfo) > 0 {
					c.ModConflicts = conflictsInfo
					return nil
				}
				return fmt.Errorf("error proccessing install conflictions: %v", raw["depends"])
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
		c.Versions.Epoch = epoch

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
			} else {
				newVKsp, _, err := c.cleanModVersion(v.(string))
				if err != nil {
					return fmt.Errorf("invalid mod version: %v", v.(string))
				}

				if vMax == nil {
					vMax = newVKsp
				}
				if vMin == nil {
					vMin = newVKsp
				}
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

		c.Versions.Mod = vMod.String()
		c.Versions.KspMin = vMin.String()
		c.Versions.KspMax = vMax.String()

		c.IsCompatible = c.checkCompatible()

		return nil
	} else {
		return errors.New("no version available")
	}
}

func (c *Ckan) cleanAbstract(raw map[string]interface{}) error {
	c.Abstract = strings.TrimSpace(raw["abstract"].(string))
	if c.Abstract == "" {
		return errors.New("invalid abstract")
	}

	return nil
}

func (c *Ckan) cleanDescription(raw map[string]interface{}) error {
	if raw["description"] != nil {
		c.Description = strings.TrimSpace(raw["description"].(string))
		if c.Description == "" {
			return errors.New("invalid description")
		}
		return nil
	}
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
	if raw["download"] != nil {
		c.Download.URL = strings.TrimSpace(raw["download"].(string))
	} else {
		//log.Printf("cannot read download field: %v", raw["download"])
		return fmt.Errorf("cannot read download field: %v", raw["download"])
	}

	if c.Download.URL == "" {
		//log.Printf("invalid download path: %v", raw["download"])
		return fmt.Errorf("invalid download path: %v", raw["download"])
	}

	c.Download.Path = "/" + c.Identifier + ".zip"
	c.Download.Downloaded = false

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

func clean(dirty string) string {
	s := []byte(dirty)
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == ' ' {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}
