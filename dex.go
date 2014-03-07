package apkdig

import (
	"io"
)

type DEX struct {
	Header    *DEXHeader
	StringIds []StringIdItem
	Strings   []string
	TypeIds   []TypeIdItem
	ProtoIds  []ProtoIdItem
	FieldIds  []FieldIdItem
	MethodIds []MethodIdItem
	ClassDefs []ClassDefItem
	Classes   []DexClass
}

func ReadDex(file io.ReadSeeker) (*DEX, error) {
	dex := new(DEX)
	err := dex.readHeader(file)
	if err != nil {
		return dex, err
	}
	err = dex.readStringIds(file)
	if err != nil {
		return dex, err
	}
	err = dex.readStrings(file)
	if err != nil {
		return dex, err
	}
	err = dex.readTypeIds(file)
	if err != nil {
		return nil, err
	}
	err = dex.readProtoIds(file)
	if err != nil {
		return nil, err
	}
	err = dex.readFieldIds(file)
	if err != nil {
		return nil, err
	}
	err = dex.readMethodIds(file)
	if err != nil {
		return nil, err
	}
	err = dex.readClassDefs(file)
	if err != nil {
		return nil, err
	}
	dex.parseDexClasses()
	return dex, nil
}
