package axml

/*
 * Copyright (c) 2015 Floor Terra <floort@gmail.com>
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
)

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

func (b *XmlEndTagBlock) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &b.Type); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Size); err != nil {
		return err
	}
	if b.Type != CHUNK_XML_END_TAG {
		return fmt.Errorf("Expected type=%X, got type=%X", CHUNK_XML_END_TAG, b.Type)
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.NsIdx); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.NameIdx); err != nil {
		return err
	}
	if b.Size != 16 {
		return fmt.Errorf("Expected size=%d, got size=%d", 16, b.Size)
	}
	return nil
}

func (b XmlEndTagBlock) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, &b.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.Size); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.NsIdx); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.NameIdx); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
