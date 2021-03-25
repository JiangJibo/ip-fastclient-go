package context

import (
	"github.com/jiangjibo/ip-fastclient-go/fast/domain"
	lsnClient "github.com/jiangjibo/ip-fastclient-go/license/client"
)

type FastIPGeoContext struct {
	LicenseClient *lsnClient.LicenseClient
	GeoConf       *domain.FastGeoConf
	Metric        *domain.FastMetric
	MetaInfo      *domain.FastMetaInfo
	Data          []byte
}
