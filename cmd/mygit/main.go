package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}
		
		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}
		
		fmt.Println("Initialized git directory")

	case "cat-file": 
		sch := os.Args[3]
		
		pathToFile := path.Join(".git/objects", sch[0:2], sch[2:])
		file, err := os.Open(pathToFile)
		if(err != nil) {
			log.Fatalln("Error opening file:", err)
		}

		r, err := zlib.NewReader(io.Reader(file))
		if( err != nil){
			log.Fatalln("Error reading file:", err)
		}
		defer r.Close()

		s, err := io.ReadAll(r)
		if(err != nil){
			log.Fatalln("Error reading file:", err)
		}
		parts := strings.Split(string(s), "\x00")
		fmt.Print(parts[1])

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
