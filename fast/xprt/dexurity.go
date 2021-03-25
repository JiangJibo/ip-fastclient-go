package xprt

import (
	"crypto/rc4"
	"github.com/JiangJibo/ip-fastclient-go/fast/consts"
)

// rc4解密
func ParseBytes(bytes []byte) []byte {
	realCipherSize := len(bytes) - consts.MgConfusedSize

	confusedBytes := CopyOfRange(bytes, 0, consts.MgConfusedSize)
	keyBytes := CopyOfRange(confusedBytes, consts.MgKeyStartIndex, consts.MgKeyStartIndex+consts.MgKeySize)
	realCipherBytes := CopyOfRange(bytes, consts.MgConfusedSize, consts.MgConfusedSize+realCipherSize)

	cipher, _ := rc4.NewCipher(keyBytes)
	dest := make([]byte, realCipherSize)
	cipher.XORKeyStream(dest, realCipherBytes)

	return dest
}
