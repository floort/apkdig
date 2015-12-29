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

import ()

type GenericBlock interface {
	GetType() uint32
	GetSize() uint32
	GetOffset() int64
	UnmarshalBinary([]byte) error
}

type AxmlBlock struct {
	Type   uint32
	Size   uint32
	Offset int64
}

func (ab AxmlBlock) GetType() uint32 {
	return ab.Type
}

func (ab AxmlBlock) GetSize() uint32 {
	return ab.Size
}

func (ab AxmlBlock) GetOffset() int64 {
	return ab.Offset
}
