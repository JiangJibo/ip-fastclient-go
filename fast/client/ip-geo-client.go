package client

import FastClientContext "ip-fastclient-go/fast/context"

type IpGeoClient interface {

	// 检索
	Search(ip string) (string, error)

	// 加载
	Load(ctx FastClientContext.FastIPGeoContext) (bool, error)
}
