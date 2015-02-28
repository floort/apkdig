package axml

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
)

/* +------------------------------------+
 * | Type           uint32              |
 * | Size           uint32              |
 * | LineNumber     uint32              |
 * | Skip           uint32 = SKIP_BLOCK |
 * | Prefix         uint32              |
 * | Uri            uint32              |
 * +------------------------------------+
 */

type XmlStartNamespaceBlock struct {
	AxmlBlock
	LineNumber uint32
	Skip       uint32
	Prefix     uint32
	Uri        uint32
}

func ReadXmlStartNamespaceBlock(reader io.ReadSeeker, size uint32, offset int64) (b XmlStartNamespaceBlock, err error) {
	b.Type = CHUNK_RESOURCEIDS
	b.Size = size
	b.Offset = offset
	reader.Seek(offset+8, 0) // Skip Type and Size
	binary.Read(reader, binary.LittleEndian, &b.LineNumber)
	binary.Read(reader, binary.LittleEndian, &b.Skip)
	binary.Read(reader, binary.LittleEndian, &b.Prefix)
	binary.Read(reader, binary.LittleEndian, &b.Uri)
	return b, nil
}

/* +------------------------------------+
 * | Type           uint32              |
 * | Size           uint32              |
 * | lineNumber     uint32              |
 * | skip           uint32 = SKIP_BLOCK |
 * | nsIdx          int32               |
 * | nameIdx        uint32              |
 * | flag           uint32 = 0x00140014 |
 * | attributeCount uint16              |
 * +------------------------------------+
 * | +--------------------------------+ |
 * | | nsIdx       uint32             | |
 * | | nameIdx     uint32             | |
 * | | valueString uint32 // Skipped  | |
 * | | aValueType  uint32             | |
 * | | aValue      uint32             | |
 * | +--------------------------------+ |
 * |   Repeat attributeCount times      |
 * +------------------------------------+
 */

type XmlStartTagBlock struct {
	AxmlBlock
	LineNumber     uint32
	Skip           uint32
	NsIdx          uint32
	NameIdx        uint32
	Flag           uint32
	AttributeCount uint32
	Attributes     []XmlTagAttribute
}

type XmlTagAttribute struct {
	NsIdx       uint32
	NameIdx     uint32
	ValueString uint32
	AValueType  uint32
	AValue      uint32
}

func ReadXmlStartTagBlock(reader io.ReadSeeker, size uint32, offset int64) (b XmlStartTagBlock, err error) {
	b.Type = CHUNK_RESOURCEIDS
	b.Size = size
	b.Offset = offset
	reader.Seek(offset+8, 0) // Skip Type and Size
	binary.Read(reader, binary.LittleEndian, &b.LineNumber)
	binary.Read(reader, binary.LittleEndian, &b.Skip)
	binary.Read(reader, binary.LittleEndian, &b.NsIdx)
	binary.Read(reader, binary.LittleEndian, &b.NameIdx)
	binary.Read(reader, binary.LittleEndian, &b.Flag)
	binary.Read(reader, binary.LittleEndian, &b.AttributeCount)
	b.Attributes = make([]XmlTagAttribute, b.AttributeCount)
	for i := uint32(0); i < b.AttributeCount; i++ {
		binary.Read(reader, binary.LittleEndian, &b.Attributes[i].NsIdx)
		binary.Read(reader, binary.LittleEndian, &b.Attributes[i].NameIdx)
		binary.Read(reader, binary.LittleEndian, &b.Attributes[i].ValueString)
		binary.Read(reader, binary.LittleEndian, &b.Attributes[i].AValueType)
		binary.Read(reader, binary.LittleEndian, &b.Attributes[i].AValue)
	}
	return b, nil
}

/* +------------------------------------+
 * | Type           uint32              |
 * | Size           uint32              |
 * | NsIdx          uint32              |
 * | NameIdx        uint32              |
 * +------------------------------------+
 */

type XmlEndTagBlock struct {
	AxmlBlock
	NsIdx   uint32
	NameIdx uint32
}

func ReadXmlEndTagBlock(reader io.ReadSeeker, size uint32, offset int64) (b XmlEndTagBlock, err error) {
	b.Type = CHUNK_RESOURCEIDS
	b.Size = size
	b.Offset = offset
	reader.Seek(offset+8, 0) // Skip Type and Size
	binary.Read(reader, binary.LittleEndian, &b.NsIdx)
	binary.Read(reader, binary.LittleEndian, &b.NameIdx)
	return b, nil
}

/* +------------------------------------+
 * | Type           uint32              |
 * | Size           uint32              |
 * | NsIdx          uint32              |
 * | NameIdx        uint32              |
 * +------------------------------------+
 */

type XmlEndNamespaceBlock struct {
        AxmlBlock
}

func ReadXmlEndNamespaceBlock(reader io.ReadSeeker, size uint32, offset int64) (b XmlEndNamespaceBlock, err error) {
        b.Type = CHUNK_RESOURCEIDS
        b.Size = size
        b.Offset = offset
        reader.Seek(offset+8, 0) // Skip Type and Size
        return b, nil
}



