package apkdig

import (
	"testing"
	"fmt"
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
	for n, s := range dex.Strings {
	    fmt.Printf("%5d: %#v\n", n, s)
	}
	apk.Close()
}
