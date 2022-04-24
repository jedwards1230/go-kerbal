package registry

import (
	"log"
	"os"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/dirfs"
	"github.com/tidwall/buntdb"
)

var reg Registry
var db *CkanDB
var fs billy.Filesystem
var logPath = "../logs/registry_test.log"

func TestMain(m *testing.M) {
	// Create log dir
	err := os.MkdirAll("../logs", os.ModePerm)
	if err != nil {
		log.Fatalf("Failed creating tmp dir: %v", err)
	}

	// clear previous logs
	if _, err := os.Stat(logPath); err == nil {
		if err := os.Truncate(logPath, 0); err != nil {
			log.Printf("Failed to clear %s: %v", logPath, err)
		}
	}

	// write new logs to file
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open %s", logPath)
	}
	defer f.Close()
	log.SetOutput(f)

	// Set log flags
	var LstdFlags = log.Lmsgprefix | log.Ltime | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(LstdFlags)

	log.Println("*****************")
	log.Println("Testing Registry")
	log.Println("*****************")

	config.LoadConfig("../")
	log.Printf("Initializing test-db")
	db, err = GetTestDB()
	if err != nil {
		log.Fatal(err)
	}

	installedList, err := dirfs.CheckInstalledMods()
	if err != nil {
		log.Printf("Error checking installed mods: %v", err)
	}

	sortOpts := SortOptions{
		SortTag:   "name",
		SortOrder: "ascend",
	}

	reg = Registry{
		DB:               db,
		InstalledModList: installedList,
		SortOptions:      sortOpts,
	}

	code := m.Run()

	DeleteTestDB()
	os.Exit(code)
}

// Open database file
func GetTestDB() (*CkanDB, error) {
	var db *CkanDB
	data, err := buntdb.Open("../test-data.db")
	if err != nil {
		return db, err
	}
	db = &CkanDB{DB: data}

	/* err = db.UpdateDB(true)
	if err != nil {
		return db, err
	} */

	return db, err
}

func DeleteTestDB() {
	err := os.Remove("../test-data.db")
	if err != nil {
		log.Print(err)
	}
}

/* func TestUpdateDB(t *testing.T) {
	err := db.UpdateDB(true)
	if err != nil {
		t.Errorf("Error updating database %v", err)
	}
}

func BenchmarkUpdateDB(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := db.UpdateDB(true)
		if err != nil {
			b.Error(err)
		}
	}
} */

func TestCheckRepoChanges(t *testing.T) {
	_ = checkRepoChanges()
}

func BenchmarkCheckRepoChanges(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = checkRepoChanges()
	}
}

func TestCloneRepo(t *testing.T) {
	var err error
	fs, err = cloneRepo()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkCloneRepo(b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		fs, err = cloneRepo()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestUpdateDBDirect(t *testing.T) {
	var filesToScan []string
	filesToScan = append(filesToScan, dirfs.FindFilePaths(fs, ".ckan")...)
	err := db.updateDB(&fs, filesToScan)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkUpdateDBDirect(b *testing.B) {
	var filesToScan []string
	filesToScan = append(filesToScan, dirfs.FindFilePaths(fs, ".ckan")...)
	for n := 0; n < b.N; n++ {
		err := db.updateDB(&fs, filesToScan)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestGetTotalModMap(t *testing.T) {
	modMap := reg.GetEntireModList()
	if modMap == nil && len(modMap) > 0 {
		t.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modMap), modMap)
	}
	reg.TotalModMap = modMap
}
func BenchmarkGetTotalModMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		modMap := reg.GetEntireModList()
		if modMap == nil && len(modMap) > 0 {
			b.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modMap), modMap)
		}
		reg.TotalModMap = modMap
	}
}

func TestGetCompatibleModMap(t *testing.T) {
	modlist := getCompatibleModMap(reg.TotalModMap)
	if modlist == nil && len(modlist) > 0 {
		t.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
	}
}

func BenchmarkGetCompatibleModMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		modlist := getCompatibleModMap(reg.TotalModMap)
		if modlist == nil && len(modlist) > 0 {
			b.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
		}
	}
}

func TestSortModMap(t *testing.T) {
	if err := reg.SortModList(); err != nil {
		t.Errorf("could not sort mod list: %v", err)
	}
}

func BenchmarkSortModMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if err := reg.SortModList(); err != nil {
			b.Errorf("could not sort mod list: %v", err)
		}
	}
}
