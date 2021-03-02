package clientUtils

import (
	"ip-fastclient-go/fast/consts"
	"testing"
)

func TestMd5Hex(t *testing.T) {
	ret, err := Md5Hex([]byte{1, 2, 3})
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}

func TestSplitInFixedLength(t *testing.T) {
	splits := splitInFixedLength("abcdefg", 3)
	t.Log(splits)
}

func TestMakeGeo(t *testing.T) {
	jsonObject := make(map[string]string, 0)
	jsonObject[consts.GEO_X] = "128.123456_1"
	jsonObject[consts.GEO_Y] = "12.12345"
	makeGeo("abcdefghijkmn", "1234567890", jsonObject)
	t.Log(jsonObject)
}
