package xprt

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"ip-fastclient-go/fast/consts"
	"os"
	"strconv"
	"strings"
)

var (
	IpSdkFilterEmptyPropertyConfKey = "ip.sdk.filter.empty.property"
)

// 对字节数组做md5加密
func Md5Hex(bytes []byte) string {
	md5 := md5.New()
	md5.Write(bytes)
	return hex.EncodeToString(md5.Sum(nil))
}

// 原始string数据转换成json
func RawToJson(version, id, rawContent string, storedProperties []string, loadProperties map[string]bool) string {
	splits := strings.Split(rawContent, consts.GeoRawSep)
	jsonObject := make(map[string]string, 10)
	property := os.Getenv(IpSdkFilterEmptyPropertyConfKey)

	for i := 0; i < len(storedProperties); i++ {
		var value = consts.NotfoundGeoItemValue
		if i < len(splits) {
			trimValue := strings.Trim(splits[i], " ")
			if len(trimValue) > 0 {
				value = trimValue
			}
		}
		if "true" == property && value == consts.NotfoundGeoItemValue {
			continue
		}
		if _, ok := loadProperties[storedProperties[i]]; ok {
			jsonObject[storedProperties[i]] = value
		}
	}
	_, okx := loadProperties[consts.GeoX]
	_, oky := loadProperties[consts.GeoY]
	if okx && oky {
		makeGeo(version, id, jsonObject)
	}
	json, _ := json.Marshal(jsonObject)
	return string(json)
}

//思路： 把用户的uid（20位）分成10个chunk，每个chunk 为2个char的数字 对于水印ip的位置信息，归一化经纬度为6位小数
func makeGeo(version, id string, jsonObject map[string]string) {
	latitude, ok := jsonObject[consts.GeoY]
	longitude, _ := jsonObject[consts.GeoX]

	//水印位置,需要对ip经纬度做替换
	//129.618605_1
	if !ok || !strings.Contains(latitude, "_") || latitude == "" || longitude == "" {
		return
	}
	splits := strings.SplitN(latitude, "_", 2)
	realLatitude := formatString(splits[0])
	realLongitude := formatString(longitude)

	//水印ip固定为10个吧，不然太复杂
	//20200505 把version前面补0，10位
	versions := splitInFixedLength(version, 1)

	versions = append(versions, "0", "0")
	marks := splitInFixedLength(id, 2)

	order, _ := strconv.Atoi(splits[1])

	postfix := versions[order] + marks[order]
	finalLatitude := realLatitude[:len(realLatitude)-3] + postfix
	finalLongitude := realLongitude[:len(realLongitude)-3] + postfix
	jsonObject[consts.GeoX] = finalLatitude
	jsonObject[consts.GeoY] = finalLongitude
}

// 将经纬度点后的6位替换成7位
func formatString(str string) string {
	if !strings.Contains(str, ".") {
		return str
	}
	splits := strings.SplitAfterN(str, ".", 2)
	suffix := splits[1]
	// . 后是7位
	if len(suffix) == 7 {
		return str
	} else
	// .后小于7位
	if len(suffix) < 7 {
		for i := 0; i < 7-len(suffix); i++ {
			str = str + "0"
		}
		return str
	} else {
		//.后大于7位
		suffix = suffix[0:7]
		return splits[0] + "." + suffix
	}
}

// 将字符串拆分成固定段
func splitInFixedLength(text string, length int) []string {
	splits := make([]string, 0)
	x := 0
	var s string
	for i := 0; i < len(text); i++ {
		x++
		s = s + text[i:i+1]
		if x == length {
			splits = append(splits, s)
			s = ""
			x = 0
		}
		if i == len(text)-1 && len(s) > 0 {
			splits = append(splits, s)
		}
	}
	return splits
}
