package apkdig

import (
	"encoding/binary"
	"io"
)

type FieldIdItem struct {
	ClassIdx uint16
	TypeIdx  uint16
	NameIdx  uint32
}

func (dex *DEX) readFieldIds(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.FieldIdsOff), 0)
	dex.FieldIds = make([]FieldIdItem, dex.Header.FieldIdsSize)
	return binary.Read(file, binary.LittleEndian, &dex.FieldIds)
}
