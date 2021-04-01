package domain

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/JiangJibo/ip-fastclient-go/license/consts"
	LicenseErrors "github.com/JiangJibo/ip-fastclient-go/license/error"
	"github.com/JiangJibo/ip-fastclient-go/license/utils"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type CipherEntity struct {
	// license颁发时间，毫秒时间戳
	// sdk拿到这个时间后会和服务器当前时间做一个比较，如果当前时间小于这个时间（超过1天, 全球各个时间不可能超过24个小时的时差），
	// 我们认为时间不合法
	ApplyAt int64 `json:"applyAt"`

	// v4/v6: v4表示ip4, v6表示ip6
	DataType string `json:"dataType"`

	//到期后，我们通常愿意留给用户一个缓冲的时间，毫秒时间戳
	DelayAt int64 `json:"delayAt"`

	//过期时间判断, 毫秒时间戳
	ExpireAt int64 `json:"expireAt"`

	//阿里云用户的id，用来标志身份,是一个整数，最长20位
	Id string `json:"id"`

	//qps限速
	RateLimit string `json:"rateLimit"`

	// 随机token
	Token string `json:"token"`
}

// 对于不合法的证书，抛出异常，方便使用放日志提醒用户
func (entity *CipherEntity) IsValidate() (string, LicenseErrors.LicenseError) {
	if entity.ApplyAt == 0 || entity.DelayAt == 0 || entity.ExpireAt == 0 || entity.RateLimit == "" {
		return "", LicenseErrors.LicenseInvalid
	}
	//id不合法
	if entity.Id == "" || len(entity.Id) != 20 {
		return "", LicenseErrors.LicenseInvalid
	}

	//不允许系统时间小于申请时间超过${MAX_DELTA_SECONDS}秒
	deltaSeconds := (time.Now().UnixNano()/1e6 - entity.ApplyAt) / 1000
	if deltaSeconds < -consts.MaxDeltaSeconds {
		return "", LicenseErrors.SystemTimeErr
	}

	nowMillis := time.Now().UnixNano() / 1e6
	diffExpireSeconds := (entity.ExpireAt - nowMillis) / 1000
	diffDelaySeconds := (entity.DelayAt - nowMillis) / 1000

	if diffExpireSeconds <= consts.MaxNotifyBeforeSeconds && diffDelaySeconds >= 0 {
		log.Printf("[fastip2geo] | 您的服务使用到期时间为 %s，请尽快续费并更新授权文件，以免服务暂停使用影响业务运转", fmt.Sprint(time.Unix(entity.ExpireAt, 0)))
	}

	if diffDelaySeconds < 0 {
		return "", LicenseErrors.LicenseExpire
	}

	return utils.Echo(""), LicenseErrors.SUCCESS
}

func (entity *CipherEntity) MakeCipherJson() string {
	bytes, err := json.Marshal(entity)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func ReturnPartSize(part int, totalSize int) int {
	return part ^ (part + consts.MagicNum)
}

func (entity *CipherEntity) ReturnChaosParts() int64 {
	var index int64
	if entity.DataType == "ipv4" {
		index = 0
	} else if entity.DataType == "ipv6" {
		index = 1
	} else {
		panic("data type must be in ipv4 or ipv6")
	}

	limit, err := strconv.Atoi(entity.RateLimit)
	if err != nil {
		panic(err)
	}
	x := entity.ExpireAt ^ (entity.ApplyAt + int64(math.Abs(float64(limit))) ^ index)
	return x % int64(consts.MaxChaosParts)
}

func (entity *CipherEntity) CalCipherSign() string {
	plainJson := entity.MakeCipherJson()
	parts := entity.ReturnChaosParts()
	offset := 0

	chars := make([]rune, len(plainJson))

	for i := 0; i < int(parts); i++ {
		size := ReturnPartSize(i, len(plainJson))
		if offset+size > len(plainJson) {
			size = len(plainJson) - offset
		}
		// 牧茗写的java版本有问题，但为了统一格式兼容
		partString := plainJson[offset:size]
		tmp := []rune(strings.ToUpper(partString))
		for j := 0; j < len(tmp); j++ {
			tmp[j] = tmp[j] ^ int32(size)
			chars[offset] = tmp[j]
			offset++
		}
		if partString == "" || len(partString) == 0 {
			continue
		}
	}
	plainText := fmt.Sprintf("%d-:-%s-:-%s", entity.DelayAt, string(chars), entity.Token)
	// md5生成数组，非切片
	arr := md5.Sum([]byte(plainText))
	return hex.EncodeToString(arr[:])

}
