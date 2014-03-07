package apkdig

import (
	"encoding/binary"
	"io"
)

type ClassDefItem struct {
	ClassIdx        uint32
	AccessFlags     uint32
	SuperclassIdx   uint32
	InterfacesOff   uint32
	SourceFileIdx   uint32
	AnnotationsOff  uint32
	ClassDataOff    uint32
	StaticValuesOff uint32
}

type DexClass struct {
	ClassIdx    uint32
	Name        string
	AccessFlags uint32
	SourceFile  string
}

func (dex *DEX) readClassDefs(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.ClassDefsOff), 0)
	dex.ClassDefs = make([]ClassDefItem, dex.Header.ClassDefsSize)
	return binary.Read(file, binary.LittleEndian, &dex.ClassDefs)
}

func (dex *DEX) parseDexClasses() {
	dex.Classes = make([]DexClass, len(dex.ClassDefs))
	for n, c := range dex.ClassDefs {
		dex.Classes[n].ClassIdx = uint32(n)
		//dex.Classes[n].Name = dex.Strings[c.NameIdx]
		dex.Classes[n].SourceFile = dex.Strings[c.SourceFileIdx]
	}
}
