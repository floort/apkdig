package apkdig

import (
    "testing"
    "bytes"
)

func TestReadUleb(t *testing.T) {
    buf := bytes.NewBuffer([]byte{179, 173, 71})
    val, err := ULEB128Read(buf)
    if err != nil {
        t.Errorf("fgsfdfsd")
    }
    if val != 1169075 {
        t.Errorf("Wrong value. Got %d, expected 1169075.", val)
    }
}
