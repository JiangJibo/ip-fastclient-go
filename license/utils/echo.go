package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	LicenseConsts "github.com/jiangjibo/ip-fastclient-go/license/consts"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func Echo(x string) string {
	timestamp := time.Now().UnixNano() / 1e6
	plainText := fmt.Sprintf("%d%s%s", timestamp, LicenseConsts.MgLsnEchoSep, CreateRandomNumber(32))
	raw := []byte(LicenseConsts.MgLsnEchoKey)
	encrypted := AesEncryptECB([]byte(plainText), raw)
	return hex.EncodeToString(encrypted)
}

func Decox(x string) bool {
	plainText := AesDecryptECB(x, []byte(LicenseConsts.MgLsnEchoKey))
	parts := strings.SplitN(string(plainText), LicenseConsts.MgLsnEchoSep, 2)
	timestamp := time.Now().UnixNano() / 1e6
	decodedTimestamp, _ := strconv.Atoi(parts[0])
	return timestamp-int64(decodedTimestamp) <= 90
}

// 随机生成指定长度的字符串
func CreateRandomNumber(len int) string {
	var numbers = []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}
