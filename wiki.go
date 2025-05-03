package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// updateWikiIfProjectNewer checks if the project files have been modified more recently than the wiki files.
// If they have, it triggers a wiki update (the update logic is a placeholder for now).
func updateWikiIfProjectNewer() {
	projectLatestModTime, err := getLatestModTimeInDir(".")
	if err != nil {
		fmt.Printf("Error getting latest modification time for project files: %v\n", err)
		return
	}

	wikiLatestModTime, err := getLatestModTimeInDir("./wiki") // Assuming wiki files are in a 'wiki' directory
	if err != nil {
		fmt.Printf("Error getting latest modification time for wiki files: %v\n", err)
		// If the wiki directory doesn't exist or there's an error reading it,
		// we might assume the wiki needs to be created or updated.
		fmt.Println("Could not read wiki directory, assuming wiki needs update.")
		triggerWikiUpdate() // Trigger update if wiki mod time can't be determined
		return
	}

	if projectLatestModTime.After(wikiLatestModTime) {
		fmt.Println("Project files are newer than wiki files. Triggering wiki update.")
		triggerWikiUpdate()
	} else {
		fmt.Println("Wiki files are up to date with project files.")
	}
}

// getLatestModTimeInDir recursively finds the latest modification time of files in a directory.
func getLatestModTimeInDir(dirPath string) (time.Time, error) {
	var latestModTime time.Time

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore directories themselves, only consider files
		if !info.IsDir() {
			if info.ModTime().After(latestModTime) {
				latestModTime = info.ModTime()
			}
		}
		return nil
	})

	if err != nil {
		return time.Time{}, err
	}

	return latestModTime, nil
}

// triggerWikiUpdate is a placeholder function for the actual wiki update logic.
func triggerWikiUpdate() {
	// TODO: Implement the actual wiki update logic here.
	// This might involve regenerating documentation, committing changes, etc.
	fmt.Println("Placeholder for wiki update logic.")
}
