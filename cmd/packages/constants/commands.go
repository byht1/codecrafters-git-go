package constants

type AvailableCommandsStruct struct {
	Init       string
	CatFile    string
	HashObject string
	LsTree     string
	WriteTree  string
	CommitTree string
	Clone      string
}

var AvailableCommands AvailableCommandsStruct

func NewAvailableCommands() AvailableCommandsStruct {
	AvailableCommands = AvailableCommandsStruct{
		Init:       "init",
		CatFile:    "cat-file",
		HashObject: "hash-object",
		LsTree:     "ls-tree",
		WriteTree:  "write-tree",
		CommitTree: "commit-tree",
		Clone:      "clone",
	}

	return AvailableCommands
}
