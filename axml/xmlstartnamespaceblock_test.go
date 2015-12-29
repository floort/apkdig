package axml

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
	"testing"
)

func TestMarshallUnmarshallXmlStartNamespaceBlock(t *testing.T) {
	b := XmlStartNamespaceBlock{}
	b.Type = CHUNK_XML_START_NAMESPACE
	b.Size = 24
	b.LineNumber = 1
	b.Skip = SKIP_BLOCK
	b.Prefix = 1
	b.Uri = 1
	bytes, err := b.MarshalBinary()
	if err != nil {
		t.Errorf("Error marshaling block: %v", err)
	}
	b2 := &XmlStartNamespaceBlock{}
	err = b2.UnmarshalBinary(bytes)
	if err != nil {
		t.Errorf("Error unmarshaling block: %v", err)
	}
	if b != *b2 {
		t.Errorf("Struct changed during marshal+unmarshal")
	}
}
