package context

import (
	"github.com/JiangJibo/ip-fastclient-go/fast/domain"
	lsnClient "github.com/JiangJibo/ip-fastclient-go/license/client"
)

type FastIPGeoContext struct {
	LicenseClient *lsnClient.LicenseClient
	GeoConf       *domain.FastGeoConf
	Metric        *domain.FastMetric
	MetaInfo      *domain.FastMetaInfo
	Data          []byte
}
