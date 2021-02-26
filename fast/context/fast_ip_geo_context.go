package context

import (
	"ip-fastclient-go/fast/domain"
	lsnClient "ip-fastclient-go/license/client"
)

type FastIPGeoContext struct {
	LicenseClient *lsnClient.LicenseClient
	GeoConf       *domain.FastGeoConf
	Metric        *domain.FastMetric
	MetaInfo      *domain.FastMetaInfo
	Data          []byte
}
