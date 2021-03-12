package LicenseDomain

import (
	"ip-fastclient-go/license/errors"
	"testing"
)

var (
	entity = CipherEntity{
		ApplyAt:   1593312943800,
		DataType:  "ipv4",
		DelayAt:   1622610855000,
		ExpireAt:  1622610855000,
		Id:        "00000000000000012345",
		RateLimit: "-1",
		Token:     "C3MOx1T8",
	}
)

func TestReturnNilStruct(t *testing.T) {
	err := returnNilStruct()
	t.Log(err)
}

func returnNilStruct() LicenseErrors.LicenseError {
	return LicenseErrors.LicenseError{}
}

func TestCalCipherSign(t *testing.T) {

	json, _ := entity.MakeCipherJson()
	t.Log(json)
	string, _ := entity.CalCipherSign()
	t.Log(string)
}

func TestPointer(t *testing.T) {
	var x = &entity
	t.Log(x)
	t.Log(&entity)
}

func TestIntPointer(t *testing.T) {
	var i int = 10
	t.Log("i的地址=", &i)
}
