package database

import (
	"log"
	"os"
	"testing"

	"github.com/jedwards1230/go-kerbal/cmd/config"
)

func TestUpdateDB(t *testing.T) {
	config.LoadConfig()
	db := GetDB()
	err := db.UpdateDB(true)
	if err != nil {
		t.Errorf("Error updating database %v", err)
	}
}

func TestGetModList(t *testing.T) {
	config.LoadConfig()
	db := GetDB()
	modlist := db.GetModList()
	if modlist == nil && len(modlist) > 0 {
		t.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
	}
}

func BenchmarkGetModList(b *testing.B) {
	f, _ := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	log.SetOutput(f)
	config.LoadConfig()
	db := GetDB()
	for n := 0; n < b.N; n++ {
		modlist := db.GetModList()
		if modlist == nil && len(modlist) > 0 {
			b.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
		}
	}
}

func TestCheckRepoChanges(t *testing.T) {
	config.LoadConfig()

	_ = CheckRepoChanges()
}

func BenchmarkCheckRepoChanges(b *testing.B) {
	f, _ := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	log.SetOutput(f)
	config.LoadConfig()
	for n := 0; n < b.N; n++ {
		_ = CheckRepoChanges()
	}
}
