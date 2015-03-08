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
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"unicode/utf16"
)

const (
	CHUNK_AXML_FILE           = 0x00080003
	CHUNK_RESOURCEIDS         = 0x00080180
	CHUNK_STRINGS             = 0x001C0001
	CHUNK_XML_END_NAMESPACE   = 0x00100101
	CHUNK_XML_END_TAG         = 0x00100103
	CHUNK_XML_START_NAMESPACE = 0x00100100
	CHUNK_XML_START_TAG       = 0x00100102
	CHUNK_XML_TEXT            = 0x00100104
	UTF8_FLAG                 = 0x00000100
	SKIP_BLOCK                = 0xFFFFFFFF
	TYPE_FIRST_COLOR_INT      = 28
	TYPE_FIRST_INT            = 16
	TYPE_FRACTION             = 6
	TYPE_INT_BOOLEAN          = 18
	TYPE_INT_COLOR_ARGB4      = 30
	TYPE_INT_COLOR_ARGB8      = 28
	TYPE_INT_COLOR_RGB4       = 31
	TYPE_INT_COLOR_RGB8       = 29
	TYPE_INT_DEC              = 16
	TYPE_INT_HEX              = 17
	TYPE_LAST_COLOR_INT       = 31
	TYPE_LAST_INT             = 31
	TYPE_NULL                 = 0x00000000
	TYPE_REFERENCE            = 0x01000000
	TYPE_ATTRIBUTE            = 0x02000000
	TYPE_STRING               = 0x03000000
	TYPE_FLOAT                = 0x04000000
	TYPE_DIMENSION            = 0x05000000
)

/*          AXML Data structure
 * +-----------------------------------+
 * | Header   uint32 = CHUNK_AXML_FILE |
 * | FileSize uint32 // Filesize       |
 * +-----------------------------------+
 * | +-------------------------------+ |
 * | | Blocktype uint32              | |
 * | | Size      uint32              | |
 * | +-------------------------------+ |
 * | | Depends on Blocktype          | |
 * | +-------------------------------+ |
 * | +-------------------------------+ |
 * | | BlockType uint32              | |
 * | | Size      uint32              | |
 * | +-------------------------------+ |
 * | | Depends on Blocktype          | |
 * | +-------------------------------+ |
 * |      .         .         .        |
 * |      .         .         .        |
 * |      .         .         .        |
 * +-----------------------------------+
 */

type Axml struct {
	Header   uint32
	FileSize uint32
	Blocks   []GenericBlock
}

func ReadAxml(reader io.ReadSeeker) (axml Axml, err error) {
	binary.Read(reader, binary.LittleEndian, &axml.Header)
	if axml.Header != CHUNK_AXML_FILE {
		return axml, errors.New("AXML file has wrong header")
	}
	binary.Read(reader, binary.LittleEndian, &axml.FileSize)
	var blocktype, size uint32
	for offset := int64(8); offset < int64(axml.FileSize); {
		binary.Read(reader, binary.LittleEndian, &blocktype)
		binary.Read(reader, binary.LittleEndian, &size)
		var b GenericBlock
		switch blocktype {
		default:
			return axml, fmt.Errorf("Unkown Axml blocktype: %08X", blocktype)
		case CHUNK_RESOURCEIDS:
			b, err = ReadResourceIdsBlock(reader, size, offset)
		case CHUNK_STRINGS:
			b, err = ReadStringsBlock(reader, size, offset)
		case CHUNK_XML_START_NAMESPACE:
			b, err = ReadXmlStartNamespaceBlock(reader, size, offset)
		case CHUNK_XML_START_TAG:
			b, err = ReadXmlStartTagBlock(reader, size, offset)
		case CHUNK_XML_END_TAG:
			b, err = ReadXmlEndTagBlock(reader, size, offset)
		case CHUNK_XML_END_NAMESPACE:
			b, err = ReadXmlEndNamespaceBlock(reader, size, offset)
		}
		if err != nil {
			return axml, err
		}
		axml.Blocks = append(axml.Blocks, b)
		offset += int64(size)
		reader.Seek(offset, 0)
	}
	return axml, nil
}

/*func (axml Axml) String() (s string) {
    xmlbuffer := new(bytes.Buffer)
	xmlencoder := xml.NewEncoder(xmlbuffer)
	xmlencoder.Indent("", "  ")
	for b, i := range axml.Blocks {
	    if b.Type == CHUNK_XML_START_TAG {

	    }
	}
}*/

type StringsMeta struct {
	Nstrings         uint32
	StyleOffsetCount uint32
	Flags            uint32
	StringDataOffset uint32
	Stylesoffset     uint32
	DataOffset       []uint32
}

type Attribute struct {
	ansidx       uint32
	anameidx     uint32
	avaluestring uint32
	avaluetype   uint32
	avalue       int32
}

type AXML struct {
	Header      uint32
	size        uint32
	stringsmeta StringsMeta
	Strings     []string
	Resources   []uint32
	XML         string
}

func ReadAXML(reader io.ReadSeeker) (AXML, error) {
	xmlbuffer := new(bytes.Buffer)
	xmlencoder := xml.NewEncoder(xmlbuffer)
	xmlencoder.Indent("", "  ")
	namestack := make([]xml.Name, 0, 10)
	axml := AXML{}
	binary.Read(reader, binary.LittleEndian, &axml.Header)
	if axml.Header != CHUNK_AXML_FILE {
		return axml, errors.New("AXML file has wrong header")
	}
	binary.Read(reader, binary.LittleEndian, &axml.size)
	var blocktype, size uint32
	// Start offset at 8 bytes for header and size
	for offset := uint32(8); offset < axml.size; {
		binary.Read(reader, binary.LittleEndian, &blocktype)
		binary.Read(reader, binary.LittleEndian, &size)
		switch blocktype {
		default:
			return axml, fmt.Errorf("Unkown chunk type: %X", blocktype)
		case CHUNK_RESOURCEIDS:
			fmt.Printf("@%04X[%04X]:\tCHUNK_RESOURCEIDS\n", offset, size)
			var id uint32
			nids := uint32(size/4 - 2)
			for i := uint32(0); i < nids; i++ {
				binary.Read(reader, binary.LittleEndian, &id)
				axml.Resources = append(axml.Resources, id)
			}
			fmt.Printf("%#v\n", axml.Resources)
		case CHUNK_STRINGS:
			/* +------------------------------------+
			 * | Nstrings         uint32            |
			 * | StyleOffsetCount uint32            |
			 * | Flags            uint32            |
			 * | StringDataOffset uint32            |
			 * | flag             uint32            |
			 * | Stylesoffset     uint32            |
			 * +------------------------------------+
			 * | +--------------------------------+ |
			 * | | DataOffset uint32              | |
			 * | +--------------------------------+ |
			 * |       Repeat Nstrings times        |
			 * +------------------------------------+
			 * |
			 * +------------------------------------+
			 */
			binary.Read(reader, binary.LittleEndian, &axml.stringsmeta.Nstrings)
			binary.Read(reader, binary.LittleEndian, &axml.stringsmeta.StyleOffsetCount)
			binary.Read(reader, binary.LittleEndian, &axml.stringsmeta.Flags)
			binary.Read(reader, binary.LittleEndian, &axml.stringsmeta.StringDataOffset)
			binary.Read(reader, binary.LittleEndian, &axml.stringsmeta.Stylesoffset)
			for i := uint32(0); i < axml.stringsmeta.Nstrings; i++ {
				var offset uint32
				binary.Read(reader, binary.LittleEndian, &offset)
				axml.stringsmeta.DataOffset = append(axml.stringsmeta.DataOffset, offset)
			}
			if 0 != (axml.stringsmeta.Flags & UTF8_FLAG) {
				// String will be in UTF-8 encoding
				var s string
				binary.Read(reader, binary.LittleEndian, &s)
			} else {
				// String will be in UTF-16LE encoding
				for i := uint32(0); i < axml.stringsmeta.Nstrings; i++ {
					var size uint16
					binary.Read(reader, binary.LittleEndian, &size)
					stringbytes := make([]uint16, size)
					binary.Read(reader, binary.LittleEndian, &stringbytes)
					axml.Strings = append(axml.Strings, string(utf16.Decode(stringbytes)))
					if i != axml.stringsmeta.Nstrings-1 {
						reader.Seek(2, 1) // Skip 0x0000 on all but the last string
					}
				}
			}
		case CHUNK_XML_END_NAMESPACE:
			//fmt.Printf("@%04X[%04X]:\tCHUNK_XML_END_NAMESPACE\n", offset, size)
		case CHUNK_XML_END_TAG:
			var end xml.Name
			end, namestack = namestack[len(namestack)-1], namestack[:len(namestack)-1]
			xmlencoder.EncodeToken(xml.EndElement{Name: end})
		case CHUNK_XML_START_NAMESPACE:
			/* +--------------------------------------+
			 * | lineNumber     uint32
			 * | skip           uint32 = SKIP_BLOCK
			 * |
			 */
			fmt.Printf("@%04X[%04X]:\tCHUNK_XML_START_NAMESPACE\n", offset, size)
		case CHUNK_XML_START_TAG:
			/* +------------------------------------+
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

			var lineNumber, skip, nameIdx, flag uint32
			var nsIdx int32
			var attributeCount uint16
			binary.Read(reader, binary.LittleEndian, &lineNumber)
			binary.Read(reader, binary.LittleEndian, &skip)
			if skip != SKIP_BLOCK {
				return axml, errors.New("Error: Expected block 0xFFFFFFFF")
			}
			binary.Read(reader, binary.LittleEndian, &nsIdx)
			binary.Read(reader, binary.LittleEndian, &nameIdx)
			binary.Read(reader, binary.LittleEndian, &flag)
			// Check if flag is magick number
			// https://code.google.com/p/axml/source/browse/src/main/java/pxb/android/axml/AxmlReader.java?r=9bc9e64ef832736a93750998a9fa1d4406b858c3#102
			if flag != 0x00140014 {
				return axml, fmt.Errorf("Expected flag 0x00140014, found %08X at %08X\n", flag, offset+4*6)
			}
			binary.Read(reader, binary.LittleEndian, &attributeCount)
			name := xml.Name{Local: axml.Strings[nameIdx]}
			if nsIdx > -1 {
				name.Space = axml.Strings[nsIdx]
			}
			token := xml.StartElement{
				Name: name,
			}
			namestack = append(namestack, name)
			reader.Seek(6, 1)
			for i := uint16(0); i < attributeCount; i++ {
				var attr Attribute
				binary.Read(reader, binary.LittleEndian, &attr.ansidx)
				binary.Read(reader, binary.LittleEndian, &attr.anameidx)
				binary.Read(reader, binary.LittleEndian, &attr.avaluestring)
				binary.Read(reader, binary.LittleEndian, &attr.avaluetype)
				binary.Read(reader, binary.LittleEndian, &attr.avalue)
				fmt.Printf("%#v\n", attr)
				switch attr.avaluetype & 0x7FFF0000 {
				default:
					fmt.Printf("%08X\n", attr.avaluetype&0x7FFF0000)
				case TYPE_NULL:
					fmt.Printf("TYPE_NULL")
				case TYPE_STRING:
					fmt.Printf("%08X\n", attr.avaluetype&0x7FFF0000)
					name := xml.Name{Local: axml.Strings[attr.anameidx]}
					name.Space = axml.Strings[attr.ansidx]
					tag := xml.Attr{Name: name}
					tag.Value = axml.Strings[attr.avaluestring]
					token.Attr = append(token.Attr, tag)
				}
			}
			xmlencoder.EncodeToken(token)
		case CHUNK_XML_TEXT:
			fmt.Printf("@%04X[%04X]:\tCHUNK_XML_TEXT\n", offset, size)
		}
		offset += size
		reader.Seek(int64(offset), 0)
	}
	xmlencoder.Flush()
	axml.XML = xmlbuffer.String()
	return axml, nil
}
