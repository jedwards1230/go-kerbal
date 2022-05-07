package database

import (
	"log"
	"os"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/jedwards1230/go-kerbal/internal/config"
	"github.com/jedwards1230/go-kerbal/internal/dirfs"
)

var db *CkanDB
var fs billy.Filesystem
var logPath = "../../logs/database_test.log"

func TestMain(m *testing.M) {
	// Create log dir
	err := os.MkdirAll("../../logs", os.ModePerm)
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
	log.Println("Testing Database")
	log.Println("*****************")

	config.LoadConfig("../../")
	log.Printf("Initializing test-db")
	db = GetDB("../../test-data.db")

	code := m.Run()

	DeleteTestDB()
	os.Exit(code)
}

func DeleteTestDB() {
	err := os.Remove("../../test-data.db")
	if err != nil {
		log.Print(err)
	}
}

func TestUpdateDB(t *testing.T) {
	err := db.UpdateDB(true)
	if err != nil {
		t.Errorf("Error updating database %v", err)
	}
}

/* func BenchmarkUpdateDB(b *testing.B) {
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

/* func TestCloneRepo(t *testing.T) {
	var err error
	repo, err := cloneRepo()
	if err != nil {
		t.Error(err)
	}
	fs = repo
} */

func BenchmarkCloneRepo(b *testing.B) {
	var err error
	for n := 0; n < b.N; n++ {
		fs, err = cloneRepo()
		if err != nil {
			b.Error(err)
		}
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
