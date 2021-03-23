package client

import (
	"ip-fastclient-go/fast/client/impl"
	"ip-fastclient-go/fast/domain"
	"testing"
	"time"
)

var (
	ipv6GeoConf = domain.FastGeoConf{
		LicenseFilePath:      "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/license-ipv6.lic",
		DataFilePath:         "/Users/jiangjibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/ipv6-inner-common-geo.dex",
		BlockedIfRateLimited: true,
	}
	//ipv6FastIpClient = GetSingleton(&ipv6GeoConf)
)

func TestSearchIpv6(t *testing.T) {

	properties := make(map[string]bool, 16)
	properties["country_code"] = true
	properties["country"] = true
	properties["isp"] = true

	//ret, err := ipv6FastIpClient.Search("::4fc5:0000:0000:0000:0000:0001")
	//if err != nil {
	//	t.Log(err)
	//	return
	//}
	//t.Logf(ret)
}

func TestMultiSearchIpv6(t *testing.T) {
	startTime := time.Now().UnixNano() / 1e6
	t.Log(startTime)
	num := 1000 * 1000
	//for i := 0; i < num; i++ {
	//	ipv6FastIpClient.Search("240e:00e0:4fc5:0000:0000:0000:0000:0001")
	//}
	endTime := time.Now().UnixNano() / 1e6
	t.Log(endTime)

	t.Logf("检索%dW次耗时:%d毫秒", num/10000, endTime-startTime)
}

// go test -v ip-fastclient-go/fast/client -test.benchmem  -test.bench SearchIpv6InBench -test.run SearchIpv6InBench -benchtime 1s -memprofile ipv6_mem_profile.out
func BenchmarkSearchIpv6InBench(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		impl.ToByteArray("240e:00e0:4fc5:0000:0000:0000:0000:0001")
		//ipv6FastIpClient.Search("240e:00e0:4fc5:0000:0000:0000:0000:0001")
	}
}
