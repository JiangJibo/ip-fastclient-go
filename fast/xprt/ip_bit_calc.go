package xprt

import "fmt"

type IPBitCalc struct {
}

func (iPBitCalc IPBitCalc) readVInt3(data []byte, p int32) int32 {
	fmt.Println(data)
	x := (int32)(data[p]<<16) & 0xFF0000
	p++
	y := (int32)(data[p]<<8) & 0xFF00
	p++
	z := (int32)(data[p]) & 0xFF
	return x | y | z
}

func (iPBitCalc IPBitCalc) writeVInt3(data []byte, offset int32, i int) {
	data[offset] = (byte)(i >> 16)
	data[offset] = (byte)(i >> 8)
	data[offset] = (byte)(i)
	fmt.Printf("%p \n", data)
}
