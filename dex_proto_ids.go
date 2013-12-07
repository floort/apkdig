package apkdig

import (
	"encoding/binary"
	"io"
)

type ProtoIdItem struct {
	ShortyIdx     uint32
	ReturnTypeIds uint32
	ParametersOff uint32
}

func (dex *DEX) readProtoIds(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.ProtoIdsOff), 0)
	dex.ProtoIds = make([]ProtoIdItem, dex.Header.ProtoIdsSize)
	return binary.Read(file, binary.LittleEndian, &dex.ProtoIds)
}
