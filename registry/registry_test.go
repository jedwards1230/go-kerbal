package registry

import (
	"log"
	"os"
	"testing"

	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/jedwards1230/go-kerbal/registry/database"
	"github.com/tidwall/buntdb"
)

var reg Registry

func TestMain(m *testing.M) {
	// Create log dir
	err := os.MkdirAll("../logs", os.ModePerm)
	if err != nil {
		log.Fatalf("error creating tmp dir: %v", err)
	}

	// clear previous logs
	if _, err := os.Stat("../logs/registry_test.log"); err == nil {
		if err := os.Truncate("../logs/registry_test.log", 0); err != nil {
			log.Printf("Failed to clear ../logs/registry_test.log: %v", err)
		}
	}

	// write new logs to file
	f, err := os.OpenFile("../logs/registry_test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("*****************")
	log.Println("Testing Registry")
	log.Println("*****************")

	config.LoadConfig("../")
	log.Printf("Initializing test-db")
	db, err := GetTestDB()
	if err != nil {
		log.Print(err)
	}

	err = db.UpdateDB(true)
	if err != nil {
		log.Print(err)
	}

	reg = Registry{DB: db}

	code := m.Run()
	err = os.Remove("../test-data.db")
	if err != nil {
		log.Print(err)
	}

	os.Exit(code)
}

// Open database file
func GetTestDB() (*database.CkanDB, error) {
	var db *database.CkanDB
	data, err := buntdb.Open("../test-data.db")
	if err != nil {
		return db, err
	}
	db = &database.CkanDB{DB: data}
	return db, err
}

func TestGetModList(t *testing.T) {
	modlist := reg.GetModList()
	if modlist == nil && len(modlist) > 0 {
		t.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
	}
	reg.ModList = modlist
}

func BenchmarkGetModList(b *testing.B) {
	for n := 0; n < b.N; n++ {
		modlist := reg.GetModList()
		if modlist == nil && len(modlist) > 0 {
			b.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
		}
		reg.ModList = modlist
	}
}

func TestGetCompatibleModList(t *testing.T) {
	modlist := getCompatibleModList(reg.ModList)
	if modlist == nil && len(modlist) > 0 {
		t.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
	}
}

func BenchmarkGetCompatibleModList(b *testing.B) {
	for n := 0; n < b.N; n++ {
		modlist := getCompatibleModList(reg.ModList)
		if modlist == nil && len(modlist) > 0 {
			b.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
		}
	}
}
