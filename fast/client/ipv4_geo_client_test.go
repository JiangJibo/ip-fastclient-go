package client

import (
	"ip-fastclient-go/fast/domain"
	"testing"
	"time"
)

var (
	geoConf = domain.FastGeoConf{
		LicenseFilePath:      "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/license-ipv4.lic",
		DataFilePath:         "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/ipv4-inner-common-geo.dex",
		BlockedIfRateLimited: true,
	}
	fastIpClient *FastIPGeoClient
)

func init() {
	fastIpClient = GetSingleton(&geoConf)
}

func TestSearchIpv4(t *testing.T) {

	properties := make(map[string]bool, 16)
	properties["country_code"] = true
	properties["country"] = true
	properties["isp"] = true

	ret, err := fastIpClient.Search("47.116.2.4")
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

func TestMultiSearchIpv4(t *testing.T) {
	startTime := time.Now().UnixNano() / 1e6
	t.Log(startTime)
	for i := 0; i < 1000*1000*1000; i++ {
		fastIpClient.Search("47.116.2.4")
	}
	endTime := time.Now().UnixNano() / 1e6
	t.Log(endTime)

	t.Logf("检索100W次耗时:%d毫秒", endTime-startTime)

}

type User struct {
	name string
	age  int
}

func TestForRange(t *testing.T) {
	u1 := &User{
		name: "jiangjibo",
		age:  30,
	}
	u2 := &User{
		name: "jiangjibo",
		age:  11,
	}
	u3 := &User{
		name: "xiao",
		age:  20,
	}
	users := []*User{u1, u2, u3}

	for i := 0; i < len(users); i++ {
		t.Logf("%v\n", &users[i])
		t.Logf("%v\n", &users[i].name)
	}

	for _, u := range users {
		t.Logf("%v\n", &u)
	}

}
