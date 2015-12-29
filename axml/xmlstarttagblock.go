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

func (b *XmlStartTagBlock) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &b.Type); err != nil {
		return err
	}
	if b.Type != CHUNK_XML_START_TAG {
		return fmt.Errorf("Expected type=%X, got type=%X", CHUNK_XML_START_TAG, b.Type)
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Size); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.LineNumber); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Skip); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.NsIdx); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.NameIdx); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.Flag); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &b.AttributeCount); err != nil {
		return err
	}
	if b.Size != 32+b.AttributeCount*20 {
		return fmt.Errorf("Expected size=%d, got size=%d", 32+b.AttributeCount*20, b.Size)
	}
	b.Attributes = make([]XmlTagAttribute, b.AttributeCount)
	for i := uint32(0); i < b.AttributeCount; i++ {
		if err := binary.Read(reader, binary.LittleEndian, &b.Attributes[i].NsIdx); err != nil {
			return err
		}
		if err := binary.Read(reader, binary.LittleEndian, &b.Attributes[i].NameIdx); err != nil {
			return err
		}
		if err := binary.Read(reader, binary.LittleEndian, &b.Attributes[i].ValueString); err != nil {
			return err
		}
		if err := binary.Read(reader, binary.LittleEndian, &b.Attributes[i].AValueType); err != nil {
			return err
		}
		if err := binary.Read(reader, binary.LittleEndian, &b.Attributes[i].AValue); err != nil {
			return err
		}
	}
	return nil
}

func (b XmlStartTagBlock) MarshalBinary() (data []byte, err error) {
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
	if err := binary.Write(buf, binary.LittleEndian, &b.NsIdx); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.NameIdx); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.Flag); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, &b.AttributeCount); err != nil {
		return nil, err
	}
	for i := range b.Attributes {
		if err := binary.Write(buf, binary.LittleEndian, &b.Attributes[i].NsIdx); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, &b.Attributes[i].NameIdx); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, &b.Attributes[i].ValueString); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, &b.Attributes[i].AValueType); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, &b.Attributes[i].AValue); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
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
