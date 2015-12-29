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
	"testing"
)

func TestMarshallUnmarshallXmlStartTagBlock(t *testing.T) {
	a := XmlStartTagBlock{}
	a.Type = CHUNK_XML_START_TAG
	a.Size = 32
	bytes, err := a.MarshalBinary()
	if err != nil {
		t.Errorf("Error marshaling block: %v", err)
	}
	b := &XmlStartTagBlock{}
	err = b.UnmarshalBinary(bytes)
	if err != nil {
		t.Errorf("Error unmarshaling block: %v", err)
	}
	if a.Type != b.Type || a.Size != b.Size || a.LineNumber != b.LineNumber ||
		a.Skip != b.Skip || a.NsIdx != b.NsIdx || a.NameIdx != b.NameIdx ||
		a.Flag != b.Flag || a.AttributeCount != b.AttributeCount {
		t.Errorf("Struct changed during marshal+unmarshal")
	}
	if len(a.Attributes) != len(b.Attributes) {
		t.Errorf("Struct changed during marshal+unmarshal")
	}
	for i := range a.Attributes {
		if a.Attributes[i] != b.Attributes[i] {
			t.Errorf("Struct changed during marshal+unmarshal")
		}
	}
}
