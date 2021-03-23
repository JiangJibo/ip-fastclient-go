package xprt

type IPBitCalc struct {
}

func ReadInt(data []byte, p int) uint32 {
	return ReadVInt4(data, p)
}

func ReadVInt4(data []byte, p int) uint32 {
	x := uint32(data[p]) << 24
	p++
	y := ReadVInt3(data, p)
	return x | y
}

// 读取3个字节组合成int
func ReadVInt3(data []byte, p int) uint32 {
	x := uint32(data[p]) << 16
	p++
	y := uint32(data[p]) << 8
	p++
	z := uint32(data[p])
	return x | y | z
}

// 将int的数据写入字节数组，3个字节
func WriteVInt3(data []byte, offset, i int) {
	data[offset] = uint8(i >> 16)
	offset++
	data[offset] = uint8(i >> 8)
	offset++
	data[offset] = uint8(i)
}

func CopyOfRange(data []byte, from, to int) []byte {
	ret := make([]byte, to-from)
	for i := from; i < to; i++ {
		ret[i-from] = data[i]
	}
	return ret
}

func ArrayCopy(data []byte, from int, dest []byte, destPos, length int) {
	j := 0
	for i := from; i < from+length; i++ {
		dest[destPos+j] = data[i]
		j++
	}
}
