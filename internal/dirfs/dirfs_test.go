package dirfs

import (
	"log"
	"os"
	"testing"

	"github.com/jedwards1230/go-kerbal/internal/config"
)

func TestMain(m *testing.M) {
	// Create log dir
	err := os.MkdirAll("../logs", os.ModePerm)
	if err != nil {
		log.Fatalf("error creating tmp dir: %v", err)
	}

	// clear previous logs
	if _, err := os.Stat("../logs/dirfs_test.log"); err == nil {
		if err := os.Truncate("../logs/dirfs_test.log", 0); err != nil {
			log.Printf("Failed to clear ../logs/dirfs_test.log: %v", err)
		}
	}

	// write new logs to file
	f, err := os.OpenFile("../logs/dirfs_test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	log.SetOutput(f)

	config.LoadConfig("../")

	log.Println("*****************")
	log.Println("Testing DirFS")
	log.Println("*****************")

	_ = m.Run()

	//os.Exit(code)
}

func TestFindKspPath(t *testing.T) {
	s, err := FindKspPath("")
	if err != nil {
		t.Errorf("Error finding KSP path: %v", s)
	}
	if s == "" {
		log.Printf("No KSP path found")
	}

}

func BenchmarkFindKspPath(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s, _ := FindKspPath("")
		if s == "" {
			b.Errorf("Error finding KSP path: %v", s)
		}
	}
}

func TestCheckInstalledMods(t *testing.T) {
	_, err := CheckInstalledMods()
	if err != nil {
		t.Errorf("error checking mods: %v", err)
	}
}

func BenchmarkCheckInstalledMods(b *testing.B) {
	for n := 0; n < b.N; n++ {
		installedMods, err := CheckInstalledMods()
		if err != nil || installedMods == nil {
			b.Errorf("error checking mods: %v", err)
		}
	}
}
