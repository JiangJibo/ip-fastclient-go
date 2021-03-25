package client

import (
	"github.com/JiangJibo/ip-fastclient-go/fast/domain"
	"testing"
)

func TestInit(t *testing.T) {
	geoConf := &domain.FastGeoConf{
		LicenseFilePath: "/Users/JiangJibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/license-ipv4-15w.lic",
		DataFilePath:    "/Users/JiangJibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/ipv4-outer-common-geo-testonly.dex",
	}

	context, err := Init(geoConf)
	if err != nil {
		t.Log(err)
	}
	t.Log(context)

}
