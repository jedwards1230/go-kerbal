package dirfs

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	f, err := os.OpenFile(RootDir()+"/test-debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println()
	log.Println("*****************")
	log.Println("Testing DirFS")
	log.Println("*****************")

	code := m.Run()

	os.Exit(code)
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
