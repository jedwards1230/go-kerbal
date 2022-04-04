package datacollector

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
)

func pullRepo() {
	dir := filepath.Join(".", "ckan_database")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Open already cloned repo in path
	fmt.Println("Checking repo")
	r, err := git.PlainOpen(dir)
	if err != nil {
		fmt.Println("Cloning repo")
		// Clones the repository if not already downloaded
		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: "https://github.com/KSP-CKAN/CKAN-meta.git",
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Get the working directory
		w, err := r.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		// Pull from origin
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			fmt.Println("No changes detected")
		}
	}
}
