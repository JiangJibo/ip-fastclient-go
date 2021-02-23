package context

import (
	"ip-fastclient-go/fast/domain"
	lsnClient "ip-fastclient-go/license/client"
)

type FastIPGeoContext struct {
	LicenseClient *lsnClient.LicenseClient
	geoConf       *domain.FastGeoConf
	metric        *domain.FastMetric
	metaInfo      *domain.FastMetaInfo
	data          []byte
}
