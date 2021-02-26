package clientUtils

type IPBitCalc struct {
}

func ReadInt(data []byte, p int) uint32 {
	return ReadInt(data, p)
}

func ReadVInt4(data []byte, p int) uint32 {
	x := uint32(data[p]) << 24 & 0xFF000000
	p++
	y := ReadVInt3(data, p)
	return x | y
}

// 读取3个字节组合成int
func ReadVInt3(data []byte, p int) uint32 {
	x := uint32(data[p]) << 16 & 0xFF0000
	p++
	y := uint32(data[p]<<8) & 0xFF00
	p++
	z := uint32(data[p]) & 0xFF
	return x | y | z
}

// 将int的数据写入字节数组，3个字节
func WriteVInt3(data []byte, offset int, i int) {
	data[offset] = uint8(i >> 16)
	offset++
	data[offset] = uint8(i >> 8)
	offset++
	data[offset] = uint8(i)
}

func CopyOfRange(data []byte, from int, to int) []byte {
	ret := make([]byte, to-from+1)
	for i := from; i < to; i++ {
		ret[i] = data[i]
	}
	return ret
}
