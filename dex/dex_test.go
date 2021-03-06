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
	//"fmt"
	"testing"

	"github.com/floort/apkdig/apk"
	//"github.com/floort/apkdig/axml"

	"github.com/davecgh/go-spew/spew"
)

func TestParseDex(t *testing.T) {
	testfile := "../tests/Orbot-release-12.0.3.apk"
	a, err := apk.OpenAPK(testfile)
	if err != nil {
		t.Errorf("Could not open %v: %v", testfile, err)
	}
	file, err := a.OpenFile("classes.dex")
	if err != nil {
		t.Errorf("Could not open classes.dex: %v", err)
	}
	d, err := ReadDex(file)
	if err != nil {
		t.Errorf("Could not parse dex: %v", err)
	}
	spew.Dump(d.Classes)

	/*manifestfile, err := a.OpenFile("AndroidManifest.xml")
	if err != nil {
		t.Errorf("Could not open manifest: %v", err)
	}
	manifest, err := axml.ReadAxml(manifestfile)
	spew.Dump(manifest)
	if err != nil {
		t.Errorf("Error parsing AXML file: %v", err)
	}
	fmt.Printf("%#v\n", manifest)*/
	a.Close()
}
