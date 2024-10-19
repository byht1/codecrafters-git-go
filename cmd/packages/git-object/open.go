package gitObject

import (
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/packages/constants"
	"github.com/codecrafters-io/git-starter-go/cmd/packages/helpers"
)

func OpenObject(objectName string) ([]byte, string, error) {
	path := constants.ObjectPathBuilder(objectName[:2], objectName[2:])
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	reader, err := zlib.NewReader(file)
	if err != nil {
		return nil, "", err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	idx := helpers.FindNull(data)
	var (
		objectType string
		size       int
	)
	fmt.Sscanf(string(data[:idx]), "%s %d", &objectType, &size)

	if idx+size+1 != len(data) {
		return nil, "", errors.New("bad object size")
	}

	return data[idx+1:], objectType, nil
}
