package main

import "path"

type AvailableCommandsStruct struct {
	Init       string
	CatFile    string
	HashObject string
	LsTree     string
}

const (
	ROOT_DIR string = ".git"

	TREE_MODE string = "040000"
	BLOB_TYPE string = "blob"
	TREE_TYPE string = "tree"

	RAW_TREE_MODE string = "40000"
)

var (
	AvailableCommands AvailableCommandsStruct
	ObjectDir         string = path.Join(ROOT_DIR, "objects")
	RefsDir           string = path.Join(ROOT_DIR, "refs")
)

func NewAvailableCommands() AvailableCommandsStruct {
	AvailableCommands = AvailableCommandsStruct{
		Init:       "init",
		CatFile:    "cat-file",
		HashObject: "hash-object",
		LsTree:     "ls-tree",
	}

	return AvailableCommands
}

func RootPathBuilder(args ...string) string {
	allArgs := append([]string{ROOT_DIR}, args...)
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
