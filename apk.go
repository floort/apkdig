// Package apk provides functions for handling Android(tm) application packages
// (.apk files).
package apkdig

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
)

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
