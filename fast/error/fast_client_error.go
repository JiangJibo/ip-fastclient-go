package error

import LicenseErrors "github.com/jiangjibo/ip-fastclient-go/license/error"

type FastClientError struct {
	LicenseErrors.IpGeoErrorInterface
	msg  string
	code int
}

func (error *FastClientError) Error() string {
	return error.msg
}

func (error *FastClientError) Code() int {
	return error.code
}

var (
	// 成功
	SUCCESS = FastClientError{
		msg:  "success",
		code: 200,
	}

	//未知异常
	UNKNOWN = FastClientError{
		msg:  "未知异常",
		code: 1500,
	}

	//未知异常
	InvalidDat = FastClientError{
		msg:  "文件内容不合法",
		code: 1501,
	}

	//未知异常
	NotMatch = FastClientError{
		msg:  "文件不匹配",
		code: 1501,
	}
)
