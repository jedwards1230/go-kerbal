package tui

import (
	"bufio"
	"log"
	"os"
	"unicode/utf8"

	"github.com/jedwards1230/go-kerbal/internal"
)

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func (b *Bubble) checkLogs() []string {
	file, err := os.Open(internal.LogPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var fileList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileList = append(fileList, scanner.Text())

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return fileList
}
