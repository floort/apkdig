package apkdig

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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

/*       XML_START_TAG
 * +-----------------------------------+
 * | lineNumber uint32
 * | skip       uint32 = SKIP_BLOCK
 * | nsIdx      uint32
 * | 
 */


type StringsMeta struct {
	Nstrings         uint32
	StyleOffsetCount uint32
	Flags            uint32
	StringDataOffset uint32
	Stylesoffset     uint32
	DataOffset		[]uint32
}

type AXML struct {
	Header      uint32
	size        uint32
	stringsmeta StringsMeta
	Strings     []string
}

func ReadAXML(reader io.ReadSeeker) (AXML, error) {
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
		case CHUNK_STRINGS:
			fmt.Printf("@%04X[%04X]:\tCHUNK_STRINGS\n", offset, size)
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
		case CHUNK_XML_END_NAMESPACE:
			fmt.Printf("@%04X[%04X]:\tCHUNK_XML_END_NAMESPACE\n", offset, size)
		case CHUNK_XML_END_TAG:
			fmt.Printf("@%04X[%04X]:\tCHUNK_XML_END_TAG\n", offset, size)
		case CHUNK_XML_START_NAMESPACE:
			fmt.Printf("@%04X[%04X]:\tCHUNK_XML_START_NAMESPACE\n", offset, size)
		case CHUNK_XML_START_TAG:
			/* +------------------------------------+
			 * | lineNumber     uint32              |
			 * | skip           uint32 = SKIP_BLOCK |
			 * | nsIdx          uint32              |
			 * | nameIdx        uint32              |
			 * | flag           uint32 = 0x00140014 |
			 * | attributeCount uint16              |
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

			fmt.Printf("@%04X[%04X]:\tCHUNK_XML_START_TAG\n", offset, size)
			var lineNumber, skip, nsIdx, nameIdx, flag uint32
			var attributeCount uint
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
		case CHUNK_XML_TEXT:
			fmt.Printf("@%04X[%04X]:\tCHUNK_XML_TEXT\n", offset, size)
		}
		offset += size
		reader.Seek(int64(offset), 0)
	}
	return axml, nil
}
