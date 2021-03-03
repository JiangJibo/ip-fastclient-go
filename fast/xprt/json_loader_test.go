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
	jsonObject[consts.GEO_X] = "128.1112223_1"
	jsonObject[consts.GEO_Y] = "12.223432"
	makeGeo("20200505", "12345678901234567890", jsonObject)
	t.Log(jsonObject)
}

func TestRawToJson(t *testing.T) {
	storedProperties := make([]string, 10)
	storedProperties[0] = "country"
	storedProperties[1] = "country_code"
	storedProperties[2] = "province"
	storedProperties[3] = "city"
	storedProperties[4] = "isp"
	storedProperties[5] = "province_en"
	storedProperties[6] = "latitude"
	storedProperties[7] = "country_en"
	storedProperties[8] = "city_en"
	storedProperties[9] = "longitude"

	loadProperties := make(map[string]bool, len(storedProperties))

	for i := 0; i < len(storedProperties); i++ {
		loadProperties[storedProperties[i]] = true
	}

	rawContent := "中国\u0000CN\u0000上海市\u0000上海\u0000阿里巴巴\u0000Shanghai\u000031.2317000_2\u0000China\u0000Shanghai\u0000121.4726000"

	json := RawToJson("20200505", "12345678901234567890", rawContent, storedProperties, loadProperties)
	t.Log(json)
}
