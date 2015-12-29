package axml

/*
 * Copyright (c) 2014, 2015 Floor Terra <floort@gmail.com>
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
	"fmt"
	"io"
	"unicode/utf16"
)

/* +------------------------------------+
 * | Type             uint32            |
 * | Size             uint32            |
 * | NStrings         uint32            |
 * | StyleOffsetCount uint32            |
 * | Flags            uint32            |
 * | StringDataOffset uint32            |
 * | StylesOffset     uint32            |
 * +------------------------------------+
 * | +--------------------------------+ |
 * | | DataOffset uint32              | |
 * | +--------------------------------+ |
 * |       Repeat NStrings times        |
 * +------------------------------------+
 * | +--------------------------------+ |
 * | | Size uint16                    | |
 * | | Data [Size]byte                | |
 * | +--------------------------------+ |
 * |       Repeat NStrings times        |
 * +------------------------------------+
 */
type StringsBlock struct {
	AxmlBlock
	NStrings         uint32
	StyleOffsetCount uint32
	Flags            uint32
	StringDataOffset uint32
	StylesOffset     uint32
	DataOffset       []uint32
	Strings          []string
}

func (b *StringsBlock) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &b.Type); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Size); err != nil {
		return err
	}
	if b.Type != CHUNK_STRINGS {
		return fmt.Errorf("Expected type=%X, got type=%X", CHUNK_STRINGS, b.Type)
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.NStrings); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.StyleOffsetCount); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Flags); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.StringDataOffset); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.StylesOffset); err != nil {
		return err
	}
	b.DataOffset = make([]uint32, b.NStrings)
	for i := uint32(0); i < b.NStrings; i++ {
		if err := binary.Read(reader, binary.LittleEndian, &b.DataOffset[i]); err != nil {
			return err
		}
	}
	nbytes := 28 + b.NStrings*4
	if 0 != (b.Flags & UTF8_FLAG) {
		// String will be in UTF-8 encoding
		return fmt.Errorf("Strings are encoded in UTF-8: not implemented")
	} else {
		// String will be in UTF-16LE encoding
		for i := uint32(0); i < b.NStrings; i++ {
			var size uint16
			binary.Read(reader, binary.LittleEndian, &size)
			stringbytes := make([]uint16, size)
			binary.Read(reader, binary.LittleEndian, &stringbytes)
			b.Strings = append(b.Strings, string(utf16.Decode(stringbytes)))
			if i != b.NStrings-1 {
				reader.Seek(2, 1) // Skip 0x0000 on all but the last string
			}
			nbytes += 2 + uint32(size)
		}
	}
	if b.Size != nbytes {
		return fmt.Errorf("Expected size=%d, got size=%d", nbytes, b.Size)
	}
	return nil
}

func ReadStringsBlock(reader io.ReadSeeker, size uint32, offset int64) (b StringsBlock, err error) {
	b.Type = CHUNK_STRINGS
	b.Size = size
	b.Offset = offset
	reader.Seek(offset+8, 0) // Skip Type and Size
	binary.Read(reader, binary.LittleEndian, &b.NStrings)
	if (b.NStrings*4)+(5*4) > size {
		return b, fmt.Errorf("NStrings = %ud, max: %d", b.NStrings, (size-(5*4))/4)
	}
	binary.Read(reader, binary.LittleEndian, &b.StyleOffsetCount)
	binary.Read(reader, binary.LittleEndian, &b.Flags)
	binary.Read(reader, binary.LittleEndian, &b.StringDataOffset)
	binary.Read(reader, binary.LittleEndian, &b.StylesOffset)
	b.DataOffset = make([]uint32, b.NStrings)
	for i := uint32(0); i < b.NStrings; i++ {
		binary.Read(reader, binary.LittleEndian, &b.DataOffset[i])
	}
	if 0 != (b.Flags & UTF8_FLAG) {
		// String will be in UTF-8 encoding
		return b, fmt.Errorf("Strings are encoded in UTF-8: not implemented")
	} else {
		// String will be in UTF-16LE encoding
		for i := uint32(0); i < b.NStrings; i++ {
			var size uint16
			binary.Read(reader, binary.LittleEndian, &size)
			stringbytes := make([]uint16, size)
			binary.Read(reader, binary.LittleEndian, &stringbytes)
			b.Strings = append(b.Strings, string(utf16.Decode(stringbytes)))
			if i != b.NStrings-1 {
				reader.Seek(2, 1) // Skip 0x0000 on all but the last string
			}
		}
	}
	return b, nil
}
