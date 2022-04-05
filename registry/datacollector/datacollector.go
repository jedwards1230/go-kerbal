package datacollector

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
)

func GetAvailableMods() []Ckan {
	modList := []Ckan{}
	idList := make(map[string]bool)
	var filesToScan []string

	pullRepo()

	// get currently downloaded ckans
	filesToScan = append(filesToScan, findFilePaths("ckan_database", ".ckan")...)
	log.Printf("%v files to scan", len(filesToScan))

	for i := range filesToScan {
		mod, _ := parseCKAN(filesToScan[i])
		if mod.Version != nil {
			// check if mod ID has been tracked already
			if idList[mod.Identifier] {
				/* for i, stored := range modList {
					if stored.Identifier == mod.Identifier {
						if stored.Version.LessThan(mod.Version) {
							log.Printf("%d | %s is less than %s", i, stored.Version, mod.Version)
						} else if stored.Version.GreaterThan(mod.Version) {
							log.Printf("%d | %s is greater than %s", i, stored.Version, mod.Version)
						} else {
							log.Printf("%d | %s is equal to %s", i, stored.Version, mod.Version)
						}
					}
				} */
				// TODO: handle storing older versions of mod
			} else {
				modList = append(modList, mod)
				idList[mod.Identifier] = true
			}
		}
	}

	log.Printf("%v mods in list", len(modList))
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

func pullRepo() {
	dir := filepath.Join(".", "ckan_database")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Open repo from dir
	log.Println("Checking repo")
	r, err := git.PlainOpen(dir)
	if err != nil {
		log.Println("Cloning repo")
		// Clones the repository if not already downloaded
		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: "https://github.com/KSP-CKAN/CKAN-meta.git",
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Updating repo")
		// Get the working directory
		w, err := r.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		// Pull from origin
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			log.Println("No changes detected")
		}
	}
}

func parseCKAN(filePath string) (Ckan, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Ckan{}, err
	}
	defer file.Close()

	// parse ckan data
	var ckan Ckan
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return Ckan{}, err
	}

	// unmarshal to struct
	// store extra in raw field for cleaning later
	json.Unmarshal(byteValue, &ckan)
	json.Unmarshal(byteValue, &ckan.raw)

	// clean data and assign necessary struct values
	err = ckan.init()
	if err != nil {
		return Ckan{}, err
	}

	return ckan, nil

}
