package client

import "testing"

func TestLicenseClientInit(t *testing.T) {
	lc := LicenseClient{
		LicenseFilePath: "/Users/JiangJibo/applications/ip-explorer/ip-geo-fastclient/src/test/resources/license-ipv4-15w.lic",
	}
	lc.Init()
}
