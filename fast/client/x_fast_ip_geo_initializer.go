package client

import (
	"encoding/json"
	"errors"
	"github.com/jiangjibo/ip-fastclient-go/fast/context"
	"github.com/jiangjibo/ip-fastclient-go/fast/domain"
	fasterror "github.com/jiangjibo/ip-fastclient-go/fast/error"
	"github.com/jiangjibo/ip-fastclient-go/fast/xprt"
	"github.com/jiangjibo/ip-fastclient-go/license/client"
	error2 "github.com/jiangjibo/ip-fastclient-go/license/error"
	"github.com/jiangjibo/ip-fastclient-go/license/utils"
)

// 初始化
func Init(geoConf *domain.FastGeoConf) (*context.FastIPGeoContext, error) {
	licenseClient := client.GetInstance(geoConf.CalculateLicenseKey(), geoConf.GetLicenseData())
	word := licenseClient.XTry()
	if !utils.Decox(word) {
		return &context.FastIPGeoContext{}, errors.New(error2.LicenseInvalid.Error())
	}

	// 解密dat文件
	encryptedBytes := geoConf.GetDexData()
	data := xprt.ParseBytes(encryptedBytes)

	//计算交集，得到实际要加载的字段
	metaLength := xprt.ReadInt(data, 0)
	mInfo := string(xprt.CopyOfRange(data, 4, 4+int(metaLength)))
	jsonObject := make(map[string]interface{}, 8)
	json.Unmarshal([]byte(mInfo), &jsonObject)

	//获取数据版本, 一般是日期作为版本
	version := jsonObject["version"].(string)
	// 所有属性
	storedProperties := jsonObject["storedProperties"].([]interface{})
	datDataType := jsonObject["dataType"].(string)

	dataType := licenseClient.GetDataType()
	if datDataType != dataType {
		return nil, errors.New(fasterror.NotMatch.Error())
	}
	// 移除meta信息, 方便计算md5
	for i := 0; i < 4+int(metaLength); i++ {
		data[i] = 0
	}
	// 校验 MD5
	md5 := xprt.Md5Hex(data)
	if md5 == "" || md5 != jsonObject["checksum"].(string) {
		return nil, errors.New(fasterror.InvalidDat.Error())
	}

	loadedProperties := geoConf.Properties
	if loadedProperties == nil || len(loadedProperties) == 0 {
		loadedProperties = sliceToMap(storedProperties)
	} else {
		retainAll(loadedProperties, storedProperties)
	}

	geoConf.Properties = loadedProperties

	metaInfo := domain.FastMetaInfo{
		StoredProperties: toStringSlice(storedProperties),
		Version:          version,
	}

	ipGeoContext := context.FastIPGeoContext{
		LicenseClient: licenseClient,
		GeoConf:       geoConf,
		Metric:        &domain.FastMetric{},
		MetaInfo:      &metaInfo,
		Data:          data,
	}

	return &ipGeoContext, nil
}

// 切片转换到map
func sliceToMap(slice []interface{}) map[string]bool {
	m := make(map[string]bool, len(slice))
	for _, i := range slice {
		m[i.(string)] = true
	}
	return m
}

// 清除loadProperties里不存在storedProperties的元素
func retainAll(loadedProperties map[string]bool, storedProperties []interface{}) {
	storedMap := sliceToMap(storedProperties)
	for key := range loadedProperties {
		if _, ok := storedMap[key]; !ok {
			delete(loadedProperties, key)
		}
	}
}

func toStringSlice(slice []interface{}) []string {
	strings := make([]string, len(slice))
	for i := 0; i < len(slice); i++ {
		strings[i] = slice[i].(string)
	}
	return strings
}
