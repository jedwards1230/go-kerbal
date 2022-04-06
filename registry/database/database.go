package database

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/tidwall/buntdb"
)

// Wrapper for buntDB
type CkanDB struct {
	*buntdb.DB
}

// Open database file
func GetDB() *CkanDB {
	database, _ := buntdb.Open(dirfs.RootDir() + "/data.db")
	db := &CkanDB{database}
	db.CreateIndex("name", "*", buntdb.IndexJSON("name"))
	db.CreateIndex("age", "*", buntdb.IndexJSON("age"))

	return db
}

// Get list of Ckan objects from database
//
// TODO: Filtering
func (db *CkanDB) GetModList() []Ckan {
	log.Println("Gathering mod list from database")

	var ckan Ckan
	var modList []Ckan
	idList := make(map[string]bool)
	db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(_, value string) bool {
			err := json.Unmarshal([]byte(value), &ckan.raw)
			if err != nil {
				log.Printf("Error loading Ckan.raw struct: %v", err)
			}

			// initialize struct values
			err = ckan.init()
			if err != nil {
				log.Printf("Error initializing ckan: %v", err)
			}

			if idList[ckan.Identifier] {
				// TODO: handle multiple versions
			} else {
				modList = append(modList, ckan)
				idList[ckan.Identifier] = true
			}
			return true
		})
		return nil
	})
	log.Printf("Loaded %v mods from database", len(modList))
	return modList
}

// Update the database by checking the repo and applying any new changes
//
// TODO: a lot
func (db *CkanDB) UpdateDB(force_update bool) error {
	changes := checkChanges()
	if !changes && !force_update {
		return nil
	}

	log.Println("Updating database entries")

	// get currently downloaded ckans
	var filesToScan []string
	filesToScan = append(filesToScan, findFilePaths(dirfs.RootDir()+"/ckan_database", ".ckan")...)

	err := db.Update(func(tx *buntdb.Tx) error {
		for i := range filesToScan {
			modJSON, _ := parseCKAN(filesToScan[i])
			tx.Set(strconv.Itoa(i), modJSON, nil)
		}
		return nil
	})
	log.Println("Database updated")
	return err
}

// Check for any changes in metadata repo
func checkChanges() bool {
	dir := dirfs.RootDir() + "/ckan_database"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Open repo from dir
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
			return false
		}
	}
	return true
}

// Parse .ckan file into JSON string
func parseCKAN(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// parse ckan data
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(byteValue), nil

}

// Collect list of file paths
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

/* if stored.Version.LessThan(mod.Version) {
	log.Printf("%d | %s is less than %s", i, stored.Version, mod.Version)
} else if stored.Version.GreaterThan(mod.Version) {
	log.Printf("%d | %s is greater than %s", i, stored.Version, mod.Version)
} else {
	log.Printf("%d | %s is equal to %s", i, stored.Version, mod.Version)
} */
