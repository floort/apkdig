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
	"io"
)

/* +------------------------------------+
 * | Type           uint32              |
 * | Size           uint32              |
 * | NsIdx          uint32              |
 * | NameIdx        uint32              |
 * +------------------------------------+
 */
// https://github.com/android/platform_frameworks_base/blob/master/tools/aapt/XMLNode.cpp
type XmlEndNamespaceBlock struct {
	AxmlBlock
}

func ReadXmlEndNamespaceBlock(reader io.ReadSeeker, size uint32, offset int64) (b XmlEndNamespaceBlock, err error) {
	b.Type = CHUNK_RESOURCEIDS
	b.Size = size
	b.Offset = offset
	reader.Seek(offset+8, 0) // Skip Type and Size
	return b, nil
}
