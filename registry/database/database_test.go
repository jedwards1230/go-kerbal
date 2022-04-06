package database

import (
	"testing"
)

func TestUpdateDB(t *testing.T) {
	db := GetDB()
	err := db.UpdateDB(true)
	if err != nil {
		t.Errorf("Error updating database %v", err)
	}
}

func TestGetModList(t *testing.T) {
	db := GetDB()
	modlist := db.GetModList()
	if modlist == nil && len(modlist) > 0 {
		t.Errorf("Mod list came back nil. Length: %v | Type: %T", len(modlist), modlist)
	}
}
