package apkdig

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

type StringIdItem uint32

func (dex *DEX) readStringIds(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.StringIdsOff), 0)
	dex.StringIds = make([]StringIdItem, dex.Header.StringIdsSize)
	return binary.Read(file, binary.LittleEndian, &dex.StringIds)
}

func (dex *DEX) readStrings(file io.ReadSeeker) error {
	dex.Strings = make([]string, len(dex.StringIds))
	for i, idx := range dex.StringIds {
		if (uint32(idx) < dex.Header.DataOff) || (uint32(idx) > dex.Header.DataOff+dex.Header.DataSize) {
			return errors.New("String offset outside of data block")
		}
		file.Seek(int64(idx), 0)
		size, err := ULEB128Read(file)
		if err != nil {
			return err
		}
		buf := make([]byte, size)
		n, err := file.Read(buf)
		if err != nil {
			return err
		}
		if uint32(n) != size {
			return errors.New("XXX")
		}
		dex.Strings[i] = string(buf)
	}
	return nil
}
