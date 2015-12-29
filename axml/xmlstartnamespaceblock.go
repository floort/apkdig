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
	"fmt"
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

func (b *XmlStartNamespaceBlock) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &b.Type); err != nil {
		return err
	}
	if b.Type != CHUNK_XML_START_NAMESPACE {
		return fmt.Errorf("Expected type=%X, got type=%X", CHUNK_XML_START_NAMESPACE, b.Type)
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Size); err != nil {
		return err
	}
	if b.Size != 24 {
		return fmt.Errorf("Expected size=%d, got size=%d", 24, b.Size)
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.LineNumber); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Skip); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Prefix); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Uri); err != nil {
		return err
	}
	return nil
}

func (b XmlStartNamespaceBlock) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, &b.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.Size); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.LineNumber); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.Skip); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.Prefix); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.Uri); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
