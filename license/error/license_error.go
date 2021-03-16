package error

type IpGeoErrorInterface interface {
	error
	Code() int
}

type LicenseError struct {
	IpGeoErrorInterface
	msg  string
	code int
}

func (error *LicenseError) Error() string {
	return error.msg
}

func (error *LicenseError) Code() int {
	return error.code
}

var (

	// 成功
	SUCCESS = LicenseError{
		msg:  "success",
		code: 200,
	}

	//未知异常
	UNKNOWN = LicenseError{
		msg:  "未知异常",
		code: 500,
	}

	//文件已经存在
	AlreadyExist = LicenseError{
		msg:  "文件已经存在",
		code: 100,
	}

	//不是绝对路径
	NotAbsPath = LicenseError{
		msg:  "不是绝对路径",
		code: 101,
	}

	//文件夹不存在
	DirNotExist = LicenseError{
		msg:  "文件夹不存在",
		code: 102,
	}

	//证书异常
	LicenseFileNotExists = LicenseError{
		msg:  "证书不存在",
		code: 501,
	}

	//证书异常
	LicenseInvalid = LicenseError{
		msg:  "证书异常",
		code: 501,
	}

	//系统时间异常
	SystemTimeErr = LicenseError{
		msg:  "系统时间异常",
		code: 502,
	}

	//证书即将过期，请提前申请，过期后将无法使用
	LicenseWillExpire = LicenseError{
		msg:  "证书即将过期，请提前申请，过期后将无法使用",
		code: 503,
	}

	//证书过期, 请重新申请
	LicenseExpire = LicenseError{
		msg:  "证书过期, 请重新申请",
		code: 504,
	}

	//证书延期异常
	LicenseDelayErr = LicenseError{
		msg:  "证书延期异常",
		code: 505,
	}

	//非法回应
	LicenceErrEcho = LicenseError{
		msg:  "非法回应",
		code: 506,
	}

	//超过购买的qps阈值
	LicenseErrRatelimit = LicenseError{
		msg:  "超过购买的qps阈值",
		code: 507,
	}

	//数据/SDK版本不兼容
	VersionNotCompatible = LicenseError{
		msg:  "数据/SDK版本不兼容",
		code: 508,
	}
)
