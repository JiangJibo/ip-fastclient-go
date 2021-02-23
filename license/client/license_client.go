package license_client

type LicenseClient struct {
	// 证书数据
	licenseData []byte

	//是否不限速
	isUnlimited bool
}
