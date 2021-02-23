package xprt

import "testing"

func TestReadVInt3(t *testing.T) {
	data := make([]byte, 3)
	t.Log(data)
	IPBitCalc{}.writeVInt3(data, 0, 32)
	t.Log(data)
}
