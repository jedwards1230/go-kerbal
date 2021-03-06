package registry

import (
	"log"
	"os"
	"testing"

	"github.com/jedwards1230/go-kerbal/internal/ckan"
	"github.com/jedwards1230/go-kerbal/internal/config"
	"github.com/jedwards1230/go-kerbal/internal/database"
	"github.com/jedwards1230/go-kerbal/internal/queue"
)

var reg *Registry
var db *database.CkanDB
var logPath = "../../logs/registry_test.log"

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
	log.Println("Testing Registry")
	log.Println("*****************")

	config.LoadConfig("../../")
	log.Printf("Initializing test-db")
	db = database.GetDB("../../test-data.db")

	sortOpts := SortOptions{
		SortTag:   "name",
		SortOrder: "ascend",
	}

	reg = &Registry{
		DB:               db,
		SortOptions:      sortOpts,
		Queue:            queue.New(),
		InstalledModList: make(map[string]ckan.Ckan, 0),
	}

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

func TestGetEntireModList(t *testing.T) {
	modMap := reg.GetEntireModList()
	if modMap == nil && len(modMap) > 0 {
		t.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modMap), modMap)
	}
	reg.TotalModMap = modMap
}
func BenchmarkGetEntireModList(b *testing.B) {
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
