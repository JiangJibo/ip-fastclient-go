package domain

import (
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

func TestCalCipherSign(t *testing.T) {

	json := entity.MakeCipherJson()
	t.Log(json)
	string := entity.CalCipherSign()
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
