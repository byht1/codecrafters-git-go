package gitObject

import "errors"

func ReadObjectHeader(packFile []byte) (size uint64, objectType ObjectCode, used int, err error) {
	data := packFile[used]
	used++

	objectType = ObjectCode((data >> 4) & 0x7)
	size = uint64(data & 0xF)
	shift := 4

	for data&0x80 != 0 {
		if len(packFile) <= used || 64 <= shift {
			return 0, ObjectCode(0), 0, errors.New("bad object header")
		}

		data = packFile[used]
		used++

		size += uint64(data&0x7F) << shift
		shift += 7
	}
	return size, objectType, used, nil
}
