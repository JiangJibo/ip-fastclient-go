package clientUtils

import (
	md52 "crypto/md5"
	"encoding/hex"
	json2 "encoding/json"
	"ip-fastclient-go/fast/consts"
	"ip-fastclient-go/fast/domain"
	"os"
	"strings"
)

// 对字节数组做md5加密
func Md5Hex(bytes []byte) (string, error) {
	md5 := md52.New()
	_, error := md5.Write(bytes)
	if error != nil {
		return "", error
	}
	return hex.EncodeToString(md5.Sum(nil)), nil
}

// 原始string数据转换成json
func RawToJson(version string, id string, rawContent string, storedProperties []string, loadProperties map[string]bool) string {
	splits := strings.Split(rawContent, consts.GEO_RAW_SEP)
	jsonObject := make(map[string]string, 10)
	property := os.Getenv(domain.IpSdkFilterEmptyPropertyConfKey)

	for i := 0; i < len(storedProperties); i++ {
		var value = consts.NOTFOUND_GEO_ITEM_VALUE
		if i < len(splits) {
			trimValue := strings.Trim(splits[i], " ")
			if len(trimValue) > 0 {
				value = trimValue
			}
		}
		if "true" == property && value == consts.NOTFOUND_GEO_ITEM_VALUE {
			continue
		}
		if _, ok := loadProperties[storedProperties[i]]; ok {
			jsonObject[storedProperties[i]] = value
		}
		_, okx := loadProperties[consts.GEO_X]
		_, oky := loadProperties[consts.GEO_Y]
		if okx && oky {
			makeGeo(version, id, jsonObject)
		}
	}
	json, _ := json2.Marshal(jsonObject)
	return string(json)
}

//思路： 把用户的uid（20位）分成10个chunk，每个chunk 为2个char的数字 对于水印ip的位置信息，归一化经纬度为6位小数
func makeGeo(version string, id string, jsonObject map[string]string) {
	latitude, ok := jsonObject[consts.GEO_X]
	longitude, _ := jsonObject[consts.GEO_Y]
	if !ok || !strings.Contains(latitude, "_") || latitude == "" || longitude == "" {
		return
	}
	splits := strings.SplitAfterN(longitude, "_", 2)
}
