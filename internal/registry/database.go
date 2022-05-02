package registry

import (
	// Using standard json encoder here because benchmarks showed segmentio to be slightly slower
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jedwards1230/go-kerbal/internal/ckan"
	"github.com/jedwards1230/go-kerbal/internal/common"
	"github.com/jedwards1230/go-kerbal/internal/config"
	"github.com/jedwards1230/go-kerbal/internal/dirfs"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
)

// Wrapper for buntDB
type CkanDB struct {
	*buntdb.DB
}

// Open database file
func GetDB() *CkanDB {
	database, _ := buntdb.Open("./data.db")
	db := &CkanDB{DB: database}
	return db
}

// Update the database by checking the repo and applying any new changes
func (r *Registry) UpdateDB(force_update bool) error {
	log.Printf("Updating DB. Force Update: %v", force_update)
	// Check if update is required
	if !force_update {
		changes := checkRepoChanges()
		if !changes {
			log.Printf("No repo changes detected")
			return nil
		}
	}

	// Clone repo
	fs, err := cloneRepo()
	if err != nil {
		log.Printf("Error cloning repo: %v", err)
		return err
	}

	// Get currently downloaded .ckan files
	log.Printf("Searching for .ckan files")
	var filesToScan []string
	filesToScan = append(filesToScan, dirfs.FindFilePaths(fs, ".ckan")...)

	err = r.updateDB(&fs, filesToScan)

	return err
}

func (r *Registry) updateDB(fs *billy.Filesystem, filesToScan []string) error {
	var mods []ckan.Ckan

	goodCount := 0
	ignoredCount := 0
	errCount := 0
	log.Print("Cleaning mod files")
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	wg.Add(len(filesToScan))
	for i := range filesToScan {
		// Parse .ckan from repo into JSON
		go func(i int) {
			defer wg.Done()

			mod, err := parseCKAN(*fs, filesToScan[i])
			if err != nil || !mod.Valid {
				if r.viewParseErrors(mod) {
					ignoredCount++
				} else {
					errCount++
				}
			} else {
				mu.Lock()
				mods = append(mods, mod)
				mu.Unlock()
				goodCount++
			}
		}(i)
	}
	wg.Wait()
	log.Printf("Scanned mod files | %d good | %d errors | %d missing info", goodCount, errCount, ignoredCount)

	err := r.DB.Update(func(tx *buntdb.Tx) error {
		for i := range mods {
			byteValue, err := json.Marshal(mods[i])
			if err != nil {
				log.Printf("Error: %s", err)
				return err
			}
			tx.Set("mod:"+strconv.Itoa(i), string(byteValue), nil)
		}
		log.Printf("Database updated with %d mods", len(mods))
		return nil
	})
	return err
}

// Return true if errors but ignored
//
// Errors would be ignored if they don't satisfy the required fields
func (r *Registry) viewParseErrors(mod ckan.Ckan) bool {
	if len(mod.Errors) > 0 {
		if mod.Errors["ignored"] == true {
			return true
		}
		if mod.Errors["raw"] != nil && false {
			raw := mod.Errors["raw"].(map[string]interface{})
			log.Print("***** RAW *****")
			for k, v := range raw {
				log.Printf("%v: %v", k, v)
			}
			log.Print("\n")
			for k, v := range mod.Errors {
				if k != "raw" {
					common.LogErrorf("%v: %v", k, v)
				}
			}
			log.Print("\n")
		}
	}
	return false
}

// Parse .ckan file into Ckan struct
func parseCKAN(repo billy.Filesystem, filePath string) (ckan.Ckan, error) {
	var mod ckan.Ckan

	// Read .ckan from filesystem
	file, err := repo.Open(filePath)
	if err != nil {
		return mod, err
	}
	defer file.Close()

	// read ckan data
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return mod, err
	}

	// Store .ckan in struct and interface
	var raw map[string]interface{}
	err = json.Unmarshal(byteValue, &raw)
	if err != nil {
		return mod, err
	}

	// clean extra data into struct
	mod = ckan.New(raw)
	if !mod.Valid {
		//log.Printf("Initialization error: %v", err)
		return mod, errors.New("invalid mod file")
	}

	return mod, err
}

// Checks for changes to the repo by comparing commit hashes
//
// Returns true if changes were detected
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

func cloneRepo() (billy.Filesystem, error) {
	cfg := config.GetConfig()
	log.Println("Cloning database repo")
	fs := memfs.New()
	storer := memory.NewStorage()
	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:   cfg.Settings.MetaRepo,
		Depth: 1,
	})
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	viper.Set("settings.last_repo_hash", ref.Hash().String())
	viper.WriteConfigAs(viper.ConfigFileUsed())

	return fs, nil
}
