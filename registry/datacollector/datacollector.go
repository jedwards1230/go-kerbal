package datacollector

import (
	"encoding/json"
	"fmt"
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

func pullRepo() {
	dir := filepath.Join(".", "ckan_database")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Open repo from dir
	fmt.Println("Checking repo")
	r, err := git.PlainOpen(dir)
	if err != nil {
		fmt.Println("Cloning repo")
		// Clones the repository if not already downloaded
		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: "https://github.com/KSP-CKAN/CKAN-meta.git",
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Updating repo")
		// Get the working directory
		w, err := r.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		// Pull from origin
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			fmt.Println("No changes detected")
		}
	}
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
