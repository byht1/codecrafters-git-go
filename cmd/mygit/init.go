package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/packages/constants"
)

func InitCmd() error {
	for _, dir := range []string{constants.ROOT_DIR, constants.ObjectDir, constants.RefsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(constants.RootPathBuilder("HEAD"), headFileContents, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println("Initialized git directory")

	return nil
}
