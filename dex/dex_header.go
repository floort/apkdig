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
	"errors"
	"io"
)

var DEX_FILE_MAGIC = [8]byte{100, 101, 120, 10, 48, 51, 53, 0}

const ENDIAN_CONSTANT = 0x12345678
const REVERSE_ENDIAN_CONSTANT = 0x78563412

type DEXHeader struct {
	Magic         [8]byte  // magic value. Should be equal to DEX_FILE_MAGIC
	Checksum      uint32   // adler32 checksum of whole file after Checksum
	Signature     [20]byte // Sha-1 hash of whole file after Signature
	FileSize      uint32   // size of the entire file (including the header), in bytes
	HeaderSize    uint32   // size of the header (this entire section), in bytes.
	EndianTag     uint32   // endianness tag. Should be equal to ENDIAN_CONSTANT
	LinkSize      uint32   // size of the link section, or 0 if this file isn't statically linked
	LinkOff       uint32   // offset from the start of the file to the link section, or 0 if link_size == 0. The offset, if non-zero, should be to an offset into the link_data section. The format of the data pointed at is left unspecified by this document; this header field (and the previous) are left as hooks for use by runtime implementations.
	MapOff        uint32   // offset from the start of the file to the map item, or 0 if this file has no map. The offset, if non-zero, should be to an offset into the data section, and the data should be in the format specified by "map_list" below.
	StringIdsSize uint32   // count of strings in the string identifiers list
	StringIdsOff  uint32   // offset from the start of the file to the string identifiers list, or 0 if string_ids_size == 0 (admittedly a strange edge case). The offset, if non-zero, should be to the start of the string_ids section.
	TypeIdsSize   uint32   // count of elements in the type identifiers list
	TypeIdsOff    uint32   // offset from the start of the file to the type identifiers list, or 0 if type_ids_size == 0 (admittedly a strange edge case). The offset, if non-zero, should be to the start of the type_ids section.
	ProtoIdsSize  uint32   // count of elements in the prototype identifiers list
	ProtoIdsOff   uint32   // offset from the start of the file to the prototype identifiers list, or 0 if proto_ids_size == 0 (admittedly a strange edge case). The offset, if non-zero, should be to the start of the proto_ids section.
	FieldIdsSize  uint32   // count of elements in the field identifiers list
	FieldIdsOff   uint32   // offset from the start of the file to the field identifiers list, or 0 if field_ids_size == 0. The offset, if non-zero, should be to the start of the field_ids section.
	MethodIdsSize uint32   // count of elements in the method identifiers list
	MethodIdsOff  uint32   // offset from the start of the file to the method identifiers list, or 0 if method_ids_size == 0. The offset, if non-zero, should be to the start of the method_ids section.
	ClassDefsSize uint32   // count of elements in the class definitions list
	ClassDefsOff  uint32   // offset from the start of the file to the class definitions list, or 0 if class_defs_size == 0 (admittedly a strange edge case). The offset, if non-zero, should be to the start of the class_defs section.
	DataSize      uint32   // Size of data section in bytes. Must be an even multiple of sizeof(uint).
	DataOff       uint32   // offset from the start of the file to the start of the data section.
}

func (dex *DEX) readHeader(file io.ReadSeeker) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	dex.Header = new(DEXHeader)
	err = binary.Read(file, binary.LittleEndian, dex.Header)
	if err != nil {
		return err
	}
	// Check magic marker
	if dex.Header.Magic != DEX_FILE_MAGIC {
		return errors.New("Magic header does not match.")
	}
	// Check endianness
	if dex.Header.EndianTag != ENDIAN_CONSTANT {
		return errors.New("File endianness does not match specifications.")
	}
	return nil
}
