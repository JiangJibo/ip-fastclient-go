package license_utils

import (
	"ip-fastclient-go/license/consts"
	"strconv"
	"strings"
	"time"
)

func echo(x string) string {
	timestamp := time.Now().Unix()

}

func decox(x string) bool {
	plainText := AesDecryptECB(x, []byte(consts.MgLsnEchoKey))
	parts := strings.SplitN(string(plainText), consts.MgLsnEchoSep, 2)
	timestamp := time.Now().Unix()
	decodedTimestamp, _ := strconv.Atoi(parts[0])
	return timestamp-int64(decodedTimestamp) <= 90
}
