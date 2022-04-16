package database

import (
	"log"
	"os"
	"testing"

	"github.com/jedwards1230/go-kerbal/cmd/config"
	"github.com/tidwall/buntdb"
)

var db *CkanDB

func TestMain(m *testing.M) {
	f, err := os.OpenFile("../../test-debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println()
	log.Println("*****************")
	log.Println("Testing Database")
	log.Println("*****************")

	config.LoadConfig("../../")
	log.Printf("Initializing test-db")
	db, err = GetTestDB()
	if err != nil {
		log.Print(err)
	}

	/* err = db.UpdateDB(true)
	if err != nil {
		log.Print(err)
	} */

	code := m.Run()

	DeleteTestDB()

	os.Exit(code)
}

// Open database file
func GetTestDB() (*CkanDB, error) {
	var db *CkanDB
	database, err := buntdb.Open("../../test-data.db")
	if err != nil {
		return db, err
	}
	db = &CkanDB{database}
	return db, err
}

func DeleteTestDB() {
	err := os.Remove("../../test-data.db")
	if err != nil {
		log.Print(err)
	}
}
func TestUpdateDB(t *testing.T) {
	err := db.UpdateDB(false)
	if err != nil {
		t.Errorf("Error updating database %v", err)
	}
}

func BenchmarkUpdateDB(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := db.UpdateDB(true)
		if err != nil {
			b.Errorf("Error updating database %v", err)
		}
	}
}

func TestCheckRepoChanges(t *testing.T) {
	_ = CheckRepoChanges()
}

func BenchmarkCheckRepoChanges(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = CheckRepoChanges()
	}
}
