package apkdig

import (
	"encoding/binary"
	"io"
)

type MethodIdItem struct {
	ClassIdx uint16
	ProtoIdx uint16
	NameIdx  uint32
}

func (dex *DEX) readMethodIds(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.MethodIdsOff), 0)
	dex.MethodIds = make([]MethodIdItem, dex.Header.MethodIdsSize)
	return binary.Read(file, binary.LittleEndian, &dex.MethodIds)
}
