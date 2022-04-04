package datacollector

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jedwards1230/go-kerbal/registry/module"
)

func GetAvailableMods() []module.ModuleVersion {
	modList := []module.ModuleVersion{}
	idList := make(map[string]bool)
	var filesToScan []string

	pullRepo()

	// get currently downloaded ckans
	filesToScan = append(filesToScan, findFilePaths("ckan_database", ".ckan")...)

	for i := range filesToScan {
		mod := parseCKAN(filesToScan[i])
		// check if mod ID has been tracked already
		if idList[mod.Identifier] {
			// TODO: handle storing older versions of mod
		} else {
			modList = append(modList, mod)
			idList[mod.Identifier] = true
		}
	}

	return modList
}

func findFilePaths(root, ext string) []string {
	var pathList []string
	filepath.WalkDir(root, func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(dir.Name()) == ext {
			pathList = append(pathList, s)
		}
		return nil
	})
	return pathList
}

func parseCKAN(filePath string) module.ModuleVersion {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// parse ckan data
	var result module.ModuleVersion
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return result
	} else {
		json.Unmarshal(byteValue, &result)
	}

	return result

}
