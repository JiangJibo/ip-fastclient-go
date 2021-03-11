package LicenseDomain

import (
	"ip-fastclient-go/license/enums"
	"ip-fastclient-go/license/errors"
	"testing"
)

func TestReturnNilStruct(t *testing.T) {
	err := returnNilStruct()
	t.Log(err)
}

func returnNilStruct() LicenseErrors.LicenseError {
	return LicenseErrors.LicenseError{}
}

func TestDataTypeOrdinal(t *testing.T) {
	t.Log(enums.IPV4)
}

func TestCalCipherSign(t *testing.T) {
	entity := CipherEntity{
		ExpireAt:  1622610855000,
		DelayAt:   1622610855000,
		ApplyAt:   1593312943800,
		RateLimit: "-1",
		DataType:  "ipv4",
		Token:     "C3MOx1T8",
		Id:        "00000000000000012345",
	}
	json, _ := entity.MakeCipherJson()
	t.Log(json)
	string, _ := entity.CalCipherSign()
	t.Log(string)
}
