package dirfs

import (
	"log"
	"os"
	"testing"

	"github.com/jedwards1230/go-kerbal/cmd/config"
)

func TestMain(m *testing.M) {
	f, err := os.OpenFile("../test-debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	config.LoadConfig("../")
	log.SetOutput(f)
	log.Println()
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
