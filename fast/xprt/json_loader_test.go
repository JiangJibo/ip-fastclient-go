package clientUtils

import "testing"

func TestMd5Hex(t *testing.T) {
	ret, err := Md5Hex([]byte{1, 2, 3})
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
