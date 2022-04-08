package database

import (
	"encoding/json"
	"log"
	"strconv"

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
// TODO: compare speeds between downloading to memory vs storage. currently uses <= 7GB of memory on git clones.
func (db *CkanDB) UpdateDB(force_update bool) error {
	cfg := config.GetConfig()

	if !force_update {
		changes := checkRepoChanges()
		if !changes {
			log.Printf("No repo changes detected")
			return nil
		}
	}

	log.Println("Cloning database repo")
	// Pull metadata repo
	fs := memfs.New()
	storer := memory.NewStorage()
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: cfg.Settings.MetaRepo,
	})
	if err != nil {
		log.Printf("Error cloning repo: %v", err)
	}

	ref, err := r.Head()
	if err != nil {
		log.Printf("Error collecting HEAD: %v", err)
	}

	viper.Set("settings.last_repo_hash", ref.Hash().String())
	viper.WriteConfigAs(viper.ConfigFileUsed())

	// get currently downloaded ckans
	log.Printf("Updating DB entries")
	var filesToScan []string
	filesToScan = append(filesToScan, dirfs.FindFilePaths(fs, ".ckan")...)

	err = db.Update(func(tx *buntdb.Tx) error {
		for i := range filesToScan {
			modJSON, _ := dirfs.ParseCKAN(fs, filesToScan[i])
			tx.Set(strconv.Itoa(i), modJSON, nil)
		}
		return nil
	})
	log.Println("Database updated")
	return err
}

// Checks for changes to the repo by comparing commit hashes
//
// Returns true if changes were detected
//
// TODO: look into using git.Repository instead of git.Remote. Which is better?
func checkRepoChanges() bool {
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

/* if stored.Version.LessThan(mod.Version) {
	log.Printf("%d | %s is less than %s", i, stored.Version, mod.Version)
} else if stored.Version.GreaterThan(mod.Version) {
	log.Printf("%d | %s is greater than %s", i, stored.Version, mod.Version)
} else {
	log.Printf("%d | %s is equal to %s", i, stored.Version, mod.Version)
} */
