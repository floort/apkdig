package dex

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
	"encoding/binary"
	"io"
	"fmt"
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
	AccessFlags uint32
	SourceFile  string
	Super       string
	Methods     []Method
}

func (dex *DEX) readClassDefs(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.ClassDefsOff), 0)
	dex.ClassDefs = make([]ClassDefItem, dex.Header.ClassDefsSize)
	return binary.Read(file, binary.LittleEndian, &dex.ClassDefs)
}

func (dex *DEX) parseDexClasses() {
	dex.Classes = make(map[string]DexClass)
	for n, c := range dex.ClassDefs {
		name := dex.Strings[dex.TypeIds[c.ClassIdx]]
		fmt.Println(c.SourceFileIdx)
		class := DexClass{
			ClassIdx:    uint32(n),
			AccessFlags: c.AccessFlags,
			SourceFile:  dex.GetString(c.SourceFileIdx),
			Super:       dex.Strings[dex.TypeIds[c.SuperclassIdx]],
			Methods:     make([]Method, 0, 0),
		}
		dex.Classes[name] = class
	}
	for _, m := range dex.MethodIds {
		classname := dex.Strings[dex.TypeIds[m.ClassIdx]]
		class := dex.Classes[classname]
		name := dex.Strings[m.NameIdx]
		class.Methods = append(class.Methods, Method{Name: name})
		dex.Classes[classname] = class
	}
}
