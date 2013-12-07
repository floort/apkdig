package apkdig

import (
	"encoding/binary"
	"io"
)

type TypeIdItem uint32

func (dex *DEX) readTypeIds(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.TypeIdsOff), 0)
	dex.TypeIds = make([]TypeIdItem, dex.Header.TypeIdsSize)
	return binary.Read(file, binary.LittleEndian, &dex.TypeIds)
}
