package apkdig

import "testing"

func TestOpenAPK(t *testing.T) {
	testfile := "tests/Orbot-release-12.0.3.apk"
	apk, err := OpenAPK(testfile)
	if err != nil {
		t.Errorf("Could not open %v: %v", testfile, err)
	}
	apk.Close()
}
