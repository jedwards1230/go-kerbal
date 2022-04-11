package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
)

// Wrapper for buntDB
type CkanDB struct {
	*buntdb.DB
}

// Open database file
func GetDB() *CkanDB {
	database, _ := buntdb.Open(dirfs.RootDir() + "/data.db")
	db := &CkanDB{DB: database}
	return db
}

// Update the database by checking the repo and applying any new changes
//
// TODO: compare speeds between downloading to memory vs storage. currently uses <= 7GB of memory on git clones.
func (db *CkanDB) UpdateDB(force_update bool) error {
	log.Printf("Updating DB. Force Update: %v", force_update)
	if !force_update {
		changes := CheckRepoChanges()
		if !changes {
			log.Printf("No repo changes detected")
			return nil
		}
	}

	fs, err := cloneRepo()
	if err != nil {
		log.Printf("Error cloning repo: %v", err)
		return err
	}

	// get currently downloaded ckans
	log.Printf("Searching for .ckan files")
	var filesToScan []string
	filesToScan = append(filesToScan, dirfs.FindFilePaths(fs, ".ckan")...)

	log.Printf("Cleaning .ckan files and adding to database")
	var ckan Ckan
	err = db.Update(func(tx *buntdb.Tx) error {
		for i := range filesToScan {
			// Parse .ckan from repo into JSON
			ckan, _ = parseCKAN(fs, filesToScan[i])

			// Ckan to JSON
			b, err := json.Marshal(ckan)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return err
			}

			// Store in DB
			key := strconv.Itoa(i)
			tx.Set(key, string(b), nil)
		}
		return nil
	})
	log.Printf("Database updated with %d entries", len(filesToScan))
	return err
}

// Parse .ckan file into JSON string
func parseCKAN(repo billy.Filesystem, filePath string) (Ckan, error) {
	var ckan Ckan

	// Read .ckan from filesystem
	file, err := repo.Open(filePath)
	if err != nil {
		return ckan, err
	}
	defer file.Close()

	// parse ckan data
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return ckan, err
	}

	// Store .ckan in struct and interface
	var raw map[string]interface{}
	err = json.Unmarshal(byteValue, &raw)
	if err != nil {
		return ckan, err
	}

	// clean extra data into struct
	ckan.Init(raw)

	return ckan, nil
}

// Checks for changes to the repo by comparing commit hashes
//
// Returns true if changes were detected
func CheckRepoChanges() bool {
	log.Println("Checking repo for changes")

	// Load metadata repo
	cfg := config.GetConfig()
	storer := memory.NewStorage()
	rem := git.NewRemote(storer, &gitConfig.RemoteConfig{
		Name: "master",
		URLs: []string{cfg.Settings.MetaRepo},
	})

	// Gather reference list
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		log.Printf("Error loading remote list: %v", err)
	}

	// Finds last hash in master
	for _, ref := range refs {
		if ref.Name().IsBranch() && ref.Name() == "refs/heads/master" {
			log.Printf("Loading: %s %v", cfg.Settings.MetaRepo, ref.Name())
			log.Printf("Latest commit: %v", ref.Hash().String())
			// if hashes match, return false to show no changes
			if cfg.Settings.LastRepoHash == ref.Hash().String() {
				return false
			}
		}
	}
	return true
}

func cloneRepo() (billy.Filesystem, error) {
	cfg := config.GetConfig()
	log.Println("Cloning database repo")
	// Pull metadata repo
	fs := memfs.New()
	storer := memory.NewStorage()
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:   cfg.Settings.MetaRepo,
		Depth: 1,
	})
	if err != nil {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	viper.Set("settings.last_repo_hash", ref.Hash().String())
	viper.WriteConfigAs(viper.ConfigFileUsed())

	return fs, nil
}

/* func compareVersions(stored, mod Ckan, i int) {
	if stored.Version.LessThan(mod.Version) {
		log.Printf("%d | %s is less than %s", i, stored.Version, mod.Version)
	} else if stored.Version.GreaterThan(mod.Version) {
		log.Printf("%d | %s is greater than %s", i, stored.Version, mod.Version)
	} else {
		log.Printf("%d | %s is equal to %s", i, stored.Version, mod.Version)
	}
} */
