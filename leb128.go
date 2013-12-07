package apkdig

import (
	"io"
)

// https://en.wikipedia.org/wiki/LEB128
func ULEB128Read(file io.Reader) (uint32, error) {
	var result uint32 = 0
	var shift uint32 = 0
	buf := make([]byte, 1, 1)
	for {
		_, err := file.Read(buf)
		if err != nil {
			return result, err
		}
		result |= ((uint32(buf[0]) & 127) << shift)
		if (uint32(buf[0]) & 128) == 0 {
			return result, nil
		}
		shift += 7
	}
}
