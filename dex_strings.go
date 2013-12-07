package apkdig

import (
	"encoding/binary"
	"io"
	"errors"
)

type StringIdItem uint32


func (dex *DEX) readStringIds(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.StringIdsOff), 0)
	dex.StringIds = make([]StringIdItem, dex.Header.StringIdsSize)
	return binary.Read(file, binary.LittleEndian, &dex.StringIds)
}

func (dex *DEX) readStrings(file io.ReadSeeker) error {
    dex.Strings = make([]string, len(dex.StringIds))
    for i, idx := range dex.StringIds {
        if (uint32(idx) < dex.Header.DataOff) || (uint32(idx) > dex.Header.DataOff + dex.Header.DataSize) {
            return errors.New("String offset outside of data block")
        }
        file.Seek(int64(idx), 0)
        size, err := ULEB128Read(file)
        if err != nil {
            return err
        }
        buf := make([]byte, size)
        n, err := file.Read(buf)
        if err != nil {
            return err
        }
        if uint32(n) != size {
            return errors.New("XXX")
        }
        dex.Strings[i] = string(buf)
    }
    return nil
}
