package ip_geo_client

import "ip-fastclient-go/fast/context"

type IpGeoClient interface {

	// 检索
	Search(ip string) (string, error)

	// 加载
	Load(ctx context.FastIPGeoContext) (bool, error)
}
