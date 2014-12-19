// Package apkdig provides functions for handling Android(tm) application packages
// (.apk files).
package apk

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
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
)

// The APK type handles opening Android application packages.
type APK struct {
	ZipFile   *zip.ReadCloser
	FileNames []string
}

func OpenAPK(name string) (*APK, error) {
	apk := new(APK)
	zf, err := zip.OpenReader(name)
	apk.ZipFile = zf
	for _, file := range apk.ZipFile.File {
		apk.FileNames = append(apk.FileNames, file.Name)
	}
	return apk, err
}

func (apk *APK) OpenFile(name string) (io.ReadSeeker, error) {
	for _, f := range apk.ZipFile.File {
		if f.Name == name {
			zipedfile, err := f.Open()
			if err != nil {
				return nil, err
			}
			raw, err := ioutil.ReadAll(zipedfile)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(raw), nil
		}
	}
	return nil, errors.New("Apk file does not contain specified file.")
}

func (apk *APK) OpenDEX() (io.ReadSeeker, error) {
	return apk.OpenFile("classes.dex")
}

func (apk *APK) Close() {
	apk.ZipFile.Close()
}
