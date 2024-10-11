package main

import (
	"fmt"
	"os"
)

func InitCmd() error {
	for _, dir := range []string{ROOT_DIR, ObjectDir, RefsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(RootPathBuilder("HEAD"), headFileContents, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println("Initialized git directory")

	return nil
}
