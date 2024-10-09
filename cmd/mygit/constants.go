package main

import "path"

type AvailableCommandsStruct struct {
	Init       string
	CatFile    string
	HashObject string
}

const (
	RootDir string = ".git"
)

var (
	AvailableCommands AvailableCommandsStruct
	ObjectDir         string = path.Join(RootDir, "objects")
	RefsDir           string = path.Join(RootDir, "refs")
)

func NewAvailableCommands() AvailableCommandsStruct {
	AvailableCommands = AvailableCommandsStruct{
		Init:       "init",
		CatFile:    "cat-file",
		HashObject: "hash-object",
	}

	return AvailableCommands
}

func RootPathBuilder(args ...string) string {
	allArgs := append([]string{RootDir}, args...)
	return path.Join(allArgs...)
}

func ObjectPathBuilder(args ...string) string {
	allArgs := append([]string{ObjectDir}, args...)
	return path.Join(allArgs...)
}

func RefsPathBuilder(args ...string) string {
	allArgs := append([]string{ObjectDir}, args...)
	return path.Join(allArgs...)
}
