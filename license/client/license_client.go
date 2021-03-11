package license_client

import "golang.org/x/time/rate"

type LicenseClient struct {

	// 证书数据
	licenseData []byte

	// 限流
	rateLimiter rate.Limiter

	//是否不限速
	isUnlimited bool
}
