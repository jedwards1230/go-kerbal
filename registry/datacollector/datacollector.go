package datacollector

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetAvailableMods() []Ckan {
	modList := []Ckan{}
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

func parseCKAN(filePath string) Ckan {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// parse ckan data
	var Ckan Ckan
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return Ckan
	} else {
		json.Unmarshal(byteValue, &Ckan)
	}

	return Ckan

}
