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
	"fmt"
)

/* +------------------------------------+
 * | Type             uint32            |
 * | Size             uint32            |
 * | Nstrings         uint32            |
 * | StyleOffsetCount uint32            |
 * | Flags            uint32            |
 * | StringDataOffset uint32            |
 * | StylesOffset     uint32            |
 * +------------------------------------+
 * | +--------------------------------+ |
 * | | DataOffset uint32              | |
 * | +--------------------------------+ |
 * |       Repeat Nstrings times        |
 * +------------------------------------+
 * |
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
}

func ReadStringsBlock(reader io.ReadSeeker, size uint32, offset int64) (b StringsBlock, err error) {
	b.Type = CHUNK_STRINGS
	b.Size = size
	b.Offset = offset
	reader.Seek(offset, 0)
	binary.Read(reader, binary.LittleEndian, &b.NStrings)
	fmt.Printf("%#v\n", b)
	if (b.NStrings*4) + (5*4) > size {
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
	return b, nil
}
