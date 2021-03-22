package client

import (
	"ip-fastclient-go/fast/domain"
	"testing"
)

var (
	ipv6GeoConf = domain.FastGeoConf{
		LicenseFilePath:      "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/license-ipv6.lic",
		DataFilePath:         "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/ipv6-inner-common-geo.dex",
		BlockedIfRateLimited: true,
	}
	ipv6FastIpClient = GetSingleton(&ipv6GeoConf)
)

func TestSearchIpv6(t *testing.T) {

	properties := make(map[string]bool, 16)
	properties["country_code"] = true
	properties["country"] = true
	properties["isp"] = true

	ret, err := ipv6FastIpClient.Search("240e:00e0:4fc5:0000:0000:0000:0000:0001")
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf(ret)
	expected := "{\"city\":\"上海市\",\"city_en\":\"Shanghai\",\"country\":\"中国\",\"country_code\":\"CN\",\"country_en\":\"China\",\"county\":\"\",\"isp\":\"阿里巴巴\",\"latitude\":\"121.4726600\",\"longitude\":\"31.2317600\",\"province\":\"上海市\",\"province_en\":\"Shanghai\"}\n"
	if ret != expected {
		panic("查询结果与预期不一致")
	}
}
