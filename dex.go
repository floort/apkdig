package apkdig

/*
 * Copyright (c) 2014 Floor Terra <floort@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

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
