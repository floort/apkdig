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
	"io"
)

type FieldIdItem struct {
	ClassIdx uint16
	TypeIdx  uint16
	NameIdx  uint32
}

func (dex *DEX) readFieldIds(file io.ReadSeeker) error {
	file.Seek(int64(dex.Header.FieldIdsOff), 0)
	dex.FieldIds = make([]FieldIdItem, dex.Header.FieldIdsSize)
	return binary.Read(file, binary.LittleEndian, &dex.FieldIds)
}
