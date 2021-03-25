package client

import (
	ip_geo_client "github.com/jiangjibo/ip-fastclient-go/fast/client/impl"
	"github.com/jiangjibo/ip-fastclient-go/fast/context"
	"github.com/jiangjibo/ip-fastclient-go/fast/domain"
	"github.com/jiangjibo/ip-fastclient-go/fast/xprt"
	"sync"
	"time"
)

type FastIPGeoClient struct {
	fastIPGeoContext *context.FastIPGeoContext
	ipGeoClient      IpGeoClient
}

var (
	mutex                            = new(sync.Mutex)
	fastIpGeoClientConcurrentHashMap = make(map[string]*FastIPGeoClient, 4)
)

//单例模式, 文件路径一致只会实例化一次
func GetSingleton(geoConf *domain.FastGeoConf) *FastIPGeoClient {
	key := xprt.Md5Hex(geoConf.GetDexData()) + xprt.Md5Hex(geoConf.GetLicenseData())
	if _, ok := fastIpGeoClientConcurrentHashMap[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if _, ok := fastIpGeoClientConcurrentHashMap[key]; !ok {
			fastIpGeoClientConcurrentHashMap[key] = newInstance(geoConf)
		}
	}
	return fastIpGeoClientConcurrentHashMap[key]
}

//非单例模式，保持兼容性
func newInstance(geoConf *domain.FastGeoConf) *FastIPGeoClient {
	this := FastIPGeoClient{}
	startTime := time.Now().UnixNano() / 1e6
	ctx, err := Init(geoConf)
	if err != nil || ctx == nil {
		panic(err)
	}
	this.fastIPGeoContext = ctx
	fastMetric := ctx.Metric
	fastMetric.DecryptCost = fastMetric.TimeCost(startTime)

	this.ipGeoClient = this.getIPGeoClient()

	startTime = time.Now().UnixNano() / 1e6
	this.ipGeoClient.Load(ctx)

	fastMetric.LoadCost = fastMetric.TimeCost(startTime)

	geoConf.ReleaseResources()
	fastMetric.Echo()

	return &this
}

func (client *FastIPGeoClient) getIPGeoClient() IpGeoClient {
	dataType := client.fastIPGeoContext.LicenseClient.GetDataType()
	if dataType == "ipv4" {
		return &ip_geo_client.Ipv4GeoClient{}
	} else if dataType == "ipv6" {
		return &ip_geo_client.IPv6GeoClient{}
	}
	panic("invalid dataType " + dataType)
}

func (client FastIPGeoClient) Search(ip string) (string, error) {
	return client.ipGeoClient.Search(ip)
}
