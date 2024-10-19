package zlibHelper

import (
	"bytes"
	"compress/zlib"
	"io"
)

func ReadObject(data []byte) (int, []byte, error) {
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return 0, nil, err
	}
	defer r.Close()

	object, err := io.ReadAll(r)
	if err != nil {
		return 0, nil, err
	}

	bytesRead := int(b.Size()) - b.Len()
	return bytesRead, object, nil
}
