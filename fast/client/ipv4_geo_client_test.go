package client

import (
	"github.com/jiangjibo/ip-fastclient-go/fast/domain"
	"runtime"
	"testing"
	"time"
)

var (
	ipv4GeoConf = domain.FastGeoConf{
		LicenseFilePath:      "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/license-ipv4.lic",
		DataFilePath:         "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/ipv4-inner-common-geo.dex",
		BlockedIfRateLimited: true,
	}
	ipv4FastIpClient *FastIPGeoClient
)

//func init() {
//	ipv4FastIpClient = GetSingleton(&ipv4GeoConf)
//}

func TestSearchIpv4(t *testing.T) {

	properties := make(map[string]bool, 16)
	properties["country_code"] = true
	properties["country"] = true
	properties["isp"] = true

	ret, err := ipv4FastIpClient.Search("47.116.2.4")
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf(ret)
}

func TestMultiSearchIpv4(t *testing.T) {
	startTime := time.Now().UnixNano() / 1e6
	t.Log(startTime)
	for i := 0; i < 1000*1000*1000; i++ {
		ipv4FastIpClient.Search("47.116.2.4")
	}
	endTime := time.Now().UnixNano() / 1e6
	t.Log(endTime)

	t.Logf("检索100W次耗时:%d毫秒", endTime-startTime)

}

func BenchmarkSearchIpv4InBench(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ipv4FastIpClient.Search("47.116.2.4")
	}
}
