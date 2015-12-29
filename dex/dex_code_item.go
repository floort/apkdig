package dex

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
	"encoding/binary"
	"io"
)

type CodeItem struct {
	RegistersSize uint16   // the number of registers used by this code
	InsSize       uint16   // the number of words of incoming arguments to the method that this code is for
	OutsSize      uint16   // the number of words of outgoing argument space required by this code for method invocation
	TriesSize     uint16   // the number of try_items for this instance. If non-zero, then these appear as the tries array just after the insns in this instance.
	DebugInfoOff  uint32   // offset from the start of the file to the debug info (line numbers + local variable info) sequence for this code, or 0 if there simply is no information. The offset, if non-zero, should be to a location in the data section. The format of the data is specified by "debug_info_item" below.
	InsnsSize     uint32   // offset from the start of the file to the debug info (line numbers + local variable info) sequence for this code, or 0 if there simply is no information. The offset, if non-zero, should be to a location in the data section. The format of the data is specified by "debug_info_item" below.
	Insns         []uint16 // actual array of bytecode. The format of code in an insns array is specified by the companion document Dalvik bytecode. Note that though this is defined as an array of ushort, there are some internal structures that prefer four-byte alignment. Also, if this happens to be in an endian-swapped file, then the swapping is only done on individual ushorts and not on the larger internal structures.
	// Optional padding to make next item 4-byte aligned
}

func (dex *DEX) readCodeItem(file io.ReadSeeker) error {
	dex.LinkSection = make([]byte, dex.Header.LinkSize)
	if dex.Header.LinkSize > 0 {
		_, err := file.Seek(int64(dex.Header.LinkOff), 0)
		if err != nil {
			return err
		}
		err = binary.Read(file, binary.LittleEndian, dex.LinkSection)
		if err != nil {
			return err
		}

	}
	return nil
}
