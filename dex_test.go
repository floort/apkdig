package apkdig

import (
	"fmt"
	"testing"
)

func TestParseDex(t *testing.T) {
	testfile := "tests/Orbot-release-12.0.3.apk"
	apk, err := OpenAPK(testfile)
	if err != nil {
		t.Errorf("Could not open %v: %v", testfile, err)
	}
	file, err := apk.OpenFile("classes.dex")
	if err != nil {
		t.Errorf("Could not open classes.dex: %v", err)
	}
	dex, err := ReadDex(file)
	if err != nil {
		t.Errorf("Could not parse dex: %v", err)
	}
	fmt.Printf("%d strings in dex\n", len(dex.Strings))
	manifestfile, err := apk.OpenFile("AndroidManifest.xml")
	if err != nil {
		t.Errorf("Could not open manifest: %v", err)
	}
	axml, err := ReadAXML(manifestfile)
	if err != nil {
		t.Errorf("Error parsing AXML file: %v", err)
	}
	fmt.Printf("%+v\n", axml)
	apk.Close()
}
