package dirfs

import (
	"log"
	"os"
	"testing"
)

func TestFindKspPath(t *testing.T) {
	f, _ := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	log.SetOutput(f)
	s := FindKspPath()
	if s == "" {
		t.Errorf("Error finding KSP path: %v", s)
	}

}

func BenchmarkFindKspPath(b *testing.B) {
	f, _ := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	log.SetOutput(f)
	for n := 0; n < b.N; n++ {
		s := FindKspPath()
		if s == "" {
			b.Errorf("Error finding KSP path: %v", s)
		}
	}
}
