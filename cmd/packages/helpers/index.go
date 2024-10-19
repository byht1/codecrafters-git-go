package helpers

import (
	"fmt"
	"log"
)

func ReadTreeEntry(data []byte) (used int, mode string, name string, hash [20]byte) {
	idx := FindNull(data)

	entry := string(data[:idx])
	fmt.Sscanf(entry, "%s %s", &mode, &name)

	copy(hash[:], data[idx+1:idx+21])

	return idx + 21, mode, name, hash
}

func FindNull(bytes []byte) int {
	for i, val := range bytes {
		if val == 0 {
			return i
		}
	}
	return len(bytes)
}

func ProcessCmdFunc(fn func() error) {
	err := fn()
	if err != nil {
		log.Fatalln(err)
	}
}
