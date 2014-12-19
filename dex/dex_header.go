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

type DEXHeader struct {
	Magic         [8]byte
	Checksum      uint32
	Signature     [20]byte
	FileSize      uint32
	HeaderSize    uint32
	EndianTag     uint32
	LinkSize      uint32
	LinkOff       uint32
	MapOff        uint32
	StringIdsSize uint32
	StringIdsOff  uint32
	TypeIdsSize   uint32
	TypeIdsOff    uint32
	ProtoIdsSize  uint32
	ProtoIdsOff   uint32
	FieldIdsSize  uint32
	FieldIdsOff   uint32
	MethodIdsSize uint32
	MethodIdsOff  uint32
	ClassDefsSize uint32
	ClassDefsOff  uint32
	DataSize      uint32
	DataOff       uint32
}

var DEX_FILE_MAGIC = [8]byte{100, 101, 120, 10, 48, 51, 53, 0}

const ENDIAN_CONSTANT = 0x12345678
const REVERSE_ENDIAN_CONSTANT = 0x78563412

func (dex *DEX) readHeader(file io.ReadSeeker) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	dex.Header = new(DEXHeader)
	err = binary.Read(file, binary.LittleEndian, &dex.Header.Magic)
	if err != nil {
		return err
	}
	if dex.Header.Magic != DEX_FILE_MAGIC {
		return errors.New("Magic header does not match.")
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.Checksum)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.Signature)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.FileSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.HeaderSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.EndianTag)
	if err != nil {
		return err
	}
	if dex.Header.EndianTag != ENDIAN_CONSTANT {
		return errors.New("File endianness does not match specifications.")
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.LinkSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.LinkOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.MapOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.StringIdsSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.StringIdsOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.TypeIdsSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.TypeIdsOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.ProtoIdsSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.ProtoIdsOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.FieldIdsSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.FieldIdsOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.MethodIdsSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.MethodIdsOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.ClassDefsSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.ClassDefsOff)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.DataSize)
	if err != nil {
		return err
	}
	err = binary.Read(file, binary.LittleEndian, &dex.Header.DataOff)
	if err != nil {
		return err
	}
	return nil
}
